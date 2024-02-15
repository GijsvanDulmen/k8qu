package informer

import (
	"github.com/rs/zerolog"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (informer *QueueJobInformer) WatchJob() (cache.Store, cache.Controller) {
	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return informer.clientSet.QueueJob("").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return informer.clientSet.QueueJob("").Watch(lo)
			},
		},
		&queuejob.QueueJob{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				var cs = obj
				typed := cs.(*queuejob.QueueJob)
				LogForJob(*typed, "added", zerolog.DebugLevel)
				informer.ReconcileJob(typed)
			},
			UpdateFunc: func(old, new interface{}) {
				var cs = new
				typed := cs.(*queuejob.QueueJob)
				LogForJob(*typed, "updated", zerolog.DebugLevel)

				informer.ReconcileJob(typed)
			},
			DeleteFunc: func(obj interface{}) {
				var cs = obj
				typed := cs.(*queuejob.QueueJob)
				LogForJob(*typed, "deleted", zerolog.DebugLevel)
			},
		},
	)

	go controller.Run(wait.NeverStop)
	return store, controller
}

func (informer *QueueJobInformer) ReconcileJob(cs *queuejob.QueueJob) {
	if cs.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(cs, finalizerName) {
			LogForJob(*cs, "adding finalizer", zerolog.DebugLevel)
			controllerutil.AddFinalizer(cs, finalizerName)

			_, err := informer.clientSet.QueueJob(cs.ObjectMeta.Namespace).Update(cs, metav1.UpdateOptions{})
			if err != nil {
				LogForJob(*cs, err.Error(), zerolog.ErrorLevel)
			}
			return
		}
		informer.jobReconcileRequest(cs.GetQueueName())
	} else {
		if controllerutil.ContainsFinalizer(cs, finalizerName) {
			controllerutil.RemoveFinalizer(cs, finalizerName)

			LogForJob(*cs, "removing finalizer", zerolog.DebugLevel)

			_, err := informer.clientSet.QueueJob(cs.ObjectMeta.Namespace).Update(cs, metav1.UpdateOptions{})
			if err != nil {
				LogForJob(*cs, "could not remove finalizer", zerolog.ErrorLevel)
				LogForJob(*cs, err.Error(), zerolog.ErrorLevel)
			}
		}
	}
}

func LogForJob(j queuejob.QueueJob, s string, level zerolog.Level) {
	log.WithLevel(level).Msgf("job %s/%s log: %s", j.Namespace, j.Name, s)
}
