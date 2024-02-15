package main

import (
	"flag"
	"github.com/rs/zerolog"
	v1alpha12 "k8qu/pkg/apis/k8qu/v1alpha1"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuesettings"
	"k8qu/pkg/clientset"
	"k8qu/pkg/clientset/v1alpha1"
	"k8qu/pkg/informer"
	logger "k8qu/pkg/log"
	"k8qu/pkg/queue"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

var kubeconfig string
var log = logger.Logger()

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "path to Kubernetes config file")
	flag.Parse()
}

func main() {
	var restConfig *rest.Config
	var err error

	if kubeconfig == "" {
		log.Info().Msg("using in-cluster configuration")
		restConfig, err = rest.InClusterConfig()
	} else {
		log.Info().Msgf("using configuration from '%s'", kubeconfig)
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	_ = v1alpha12.AddToScheme(scheme.Scheme)

	jobClientSet, err := v1alpha1.NewForQueueJob(restConfig)
	if err != nil {
		panic(err)
	}

	queueSettingsClientSet, err := v1alpha1.NewForQueueSettings(restConfig)
	if err != nil {
		panic(err)
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	coreClientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Error().Err(err)
		os.Exit(3)
		return
	}

	if err != nil {
		log.Error().Err(err)
		os.Exit(3)
		return
	}

	var jobStore *cache.Store
	var jobController *cache.Controller
	var qsStore *cache.Store
	var qsController *cache.Controller

	reconcileChannel := make(chan string)

	// empty string = all queues
	reconcileRequest := func(queue string) {
		go func() {
			reconcileChannel <- queue
		}()
	}

	// job informer
	jobInformer := informer.NewQueueJobInformer(jobClientSet, coreClientSet, func(queue string) {
		reconcileRequest(queue)
	})
	js, jc := jobInformer.WatchJob()
	jobStore = &js
	jobController = &jc

	// queue settings informer
	queueSettingsInformer := informer.NewQueueSettingsInformer(queueSettingsClientSet, coreClientSet, func(queue string) {
		reconcileRequest(queue)
	})
	qss, qsc := queueSettingsInformer.WatchQueueSettings()
	qsStore = &qss
	qsController = &qsc

	go func() {
		for {
			reconcileRequest("")
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		reconcileQueue := <-reconcileChannel

		if !(*qsController).HasSynced() {
			log.WithLevel(zerolog.DebugLevel).Msg("waiting for full sync of queue settings")
			continue
		}

		if !(*jobController).HasSynced() {
			log.WithLevel(zerolog.DebugLevel).Msg("waiting for full sync of job settings")
			continue
		}

		(*qsStore).List()

		log.WithLevel(zerolog.DebugLevel).Msg("checking if we are synced yet")
		log.WithLevel(zerolog.DebugLevel).Msgf("reconciling '%s'", reconcileQueue)

		queues := queue.NewQueues()

		qsFromStore := (*qsStore).List()
		for i := range qsFromStore {
			qs := qsFromStore[i].(*queuesettings.QueueSettings)
			queues.AddQueue(qs.GetQueueName(), queue.Settings{
				Parallelism:                  qs.Spec.Parallelism,
				TtlAfterSuccessfulCompletion: qs.Spec.TtlAfterSuccessfulCompletion,
				TtlAfterFailedCompletion:     qs.Spec.TtlAfterFailedCompletion,
				ExecutionTimeout:             qs.Spec.ExecutionTimeout,
				MaxTimeInQueue:               qs.Spec.MaxTimeInQueue,
			})
		}

		jobsFromStore := (*jobStore).List()
		for i := range jobsFromStore {
			castedJob := jobsFromStore[i].(*queuejob.QueueJob) // safe cast
			if reconcileQueue == "" {
				queues.AddJob(castedJob)
			} else if reconcileQueue == castedJob.GetQueueName() {
				queues.AddJob(castedJob)
			}
		}

		queues.Reconcile(&clientset.QueueJobUpdater{
			Client:          jobClientSet,
			ServerResources: discoveryClient,
			Dynamic:         dynamicClient,
		})
	}
}
