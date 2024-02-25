package informers

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (informer *Informers) WatchJob() (cache.Store, cache.Controller) {
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
				typed := obj.(*queuejob.QueueJob)
				LogForJob(*typed, "added")
				informer.ReconcileJob(typed)
			},
			UpdateFunc: func(old, new interface{}) {
				typed := new.(*queuejob.QueueJob)
				LogForJob(*typed, "updated")
				informer.ReconcileJob(typed)
			},
			DeleteFunc: func(obj interface{}) {
				typed := obj.(*queuejob.QueueJob)
				LogForJob(*typed, "deleted")
			},
		},
	)

	go controller.Run(wait.NeverStop)
	return store, controller
}

func (informer *Informers) ReconcileJob(qj *queuejob.QueueJob) {
	if qj.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(qj, finalizerName) {
			LogForJob(*qj, "adding finalizer")
			controllerutil.AddFinalizer(qj, finalizerName)

			_, err := informer.clientSet.QueueJob(qj.ObjectMeta.Namespace).Update(qj, metav1.UpdateOptions{})
			if err != nil {
				LogForJob(*qj, err.Error())
			}
			return
		}
		informer.jobReconcileRequest(qj.GetQueueName())
	} else {
		if controllerutil.ContainsFinalizer(qj, finalizerName) {
			controllerutil.RemoveFinalizer(qj, finalizerName)

			LogForJob(*qj, "removing finalizer")

			_, err := informer.clientSet.QueueJob(qj.ObjectMeta.Namespace).Update(qj, metav1.UpdateOptions{})
			if err != nil {
				LogForJob(*qj, "could not remove finalizer")
				LogForJob(*qj, err.Error())
			}
		}
	}
}
