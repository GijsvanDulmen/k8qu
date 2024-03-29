package main

import (
	"flag"
	v1alpha12 "k8qu/pkg/apis/k8qu/v1alpha1"
	"k8qu/pkg/apis/k8qu/v1alpha1/markqueuejobcomplete"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuesettings"
	"k8qu/pkg/clientset"
	"k8qu/pkg/clientset/v1alpha1"
	"k8qu/pkg/informers"
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
		log.Error().Err(err)
		os.Exit(3)
		return
	}

	_ = v1alpha12.AddToScheme(scheme.Scheme)

	jobClientSet, err := v1alpha1.NewForK8Qu(restConfig)
	if err != nil {
		log.Error().Err(err)
		os.Exit(3)
		return
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		log.Error().Err(err)
		os.Exit(3)
		return
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		log.Error().Err(err)
		os.Exit(3)
		return
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
	var mqjcStore *cache.Store
	var mqjcController *cache.Controller

	reconcileChannel := make(chan string)

	// empty string = all queues
	reconcileRequest := func(queue string) {
		go func() {
			reconcileChannel <- queue
		}()
	}

	newInformers := informers.NewInformers(jobClientSet, coreClientSet, func(queue string) {
		reconcileRequest(queue)
	})
	js, jc := newInformers.WatchJob()
	jobStore = &js
	jobController = &jc

	qss, qsc := newInformers.WatchQueueSettings()
	qsStore = &qss
	qsController = &qsc

	mqjcs, mqjcc := newInformers.WatchMarkQueueJobComplete()
	mqjcStore = &mqjcs
	mqjcController = &mqjcc

	go func() {
		for {
			reconcileRequest("")
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		reconcileQueue := <-reconcileChannel

		log.Debug().Msg("checking if we are synced yet")

		if !(*qsController).HasSynced() {
			log.Debug().Msg("waiting for full sync of queue settings")
			continue
		}

		if !(*jobController).HasSynced() {
			log.Debug().Msg("waiting for full sync of job settings")
			continue
		}

		if !(*mqjcController).HasSynced() {
			log.Debug().Msg("waiting for full sync of mark queue job complete")
			continue
		}

		log.Debug().Msgf("reconciling '%s'", reconcileQueue)

		queues := queue.NewQueues()

		qsFromStore := (*qsStore).List()
		for i := range qsFromStore {
			qs := qsFromStore[i].(*queuesettings.QueueSettings)
			queues.AddQueue(qs.GetQueueName(), qs.Spec)
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

		mqjcFromStore := (*mqjcStore).List()
		for i := range mqjcFromStore {
			castedMqjc := mqjcFromStore[i].(*markqueuejobcomplete.MarkQueueJobComplete) // safe cast
			if reconcileQueue == "" {
				queues.AddMarkQueueJobComplete(castedMqjc)
			} else if reconcileQueue == castedMqjc.GetQueueName() {
				queues.AddMarkQueueJobComplete(castedMqjc)
			}
		}

		queues.Reconcile(&clientset.QueueJobUpdater{
			Client:          jobClientSet,
			ServerResources: discoveryClient,
			Dynamic:         dynamicClient,
		})
	}
}
