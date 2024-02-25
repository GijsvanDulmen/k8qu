package informers

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/markqueuejobcomplete"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (informer *Informers) WatchMarkQueueJobComplete() (cache.Store, cache.Controller) {
	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return informer.clientSet.MarkQueueJobComplete("").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return informer.clientSet.MarkQueueJobComplete("").Watch(lo)
			},
		},
		&markqueuejobcomplete.MarkQueueJobComplete{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				typed := obj.(*markqueuejobcomplete.MarkQueueJobComplete)
				LogForMarkQueueJobComplete(*typed, "added")
				informer.ReconcileMarkQueueJobComplete(typed)
			},
			UpdateFunc: func(old, new interface{}) {
				typed := new.(*markqueuejobcomplete.MarkQueueJobComplete)
				LogForMarkQueueJobComplete(*typed, "updated")
				informer.ReconcileMarkQueueJobComplete(typed)
			},
			DeleteFunc: func(obj interface{}) {
				typed := obj.(*markqueuejobcomplete.MarkQueueJobComplete)
				LogForMarkQueueJobComplete(*typed, "deleted")
			},
		},
	)

	go controller.Run(wait.NeverStop)
	return store, controller
}

func (informer *Informers) ReconcileMarkQueueJobComplete(mqjc *markqueuejobcomplete.MarkQueueJobComplete) {
	if mqjc.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(mqjc, finalizerName) {
			LogForMarkQueueJobComplete(*mqjc, "adding finalizer")
			controllerutil.AddFinalizer(mqjc, finalizerName)

			_, err := informer.clientSet.MarkQueueJobComplete(mqjc.ObjectMeta.Namespace).Update(mqjc, metav1.UpdateOptions{})
			if err != nil {
				LogForMarkQueueJobComplete(*mqjc, err.Error())
			}
			return
		}
		informer.jobReconcileRequest(mqjc.GetQueueName())
	} else {
		if controllerutil.ContainsFinalizer(mqjc, finalizerName) {
			controllerutil.RemoveFinalizer(mqjc, finalizerName)

			LogForMarkQueueJobComplete(*mqjc, "removing finalizer")

			_, err := informer.clientSet.MarkQueueJobComplete(mqjc.ObjectMeta.Namespace).Update(mqjc, metav1.UpdateOptions{})
			if err != nil {
				LogForMarkQueueJobComplete(*mqjc, "could not remove finalizer")
				LogForMarkQueueJobComplete(*mqjc, err.Error())
			}
		}
	}
}
