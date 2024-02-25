package informers

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/queuesettings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (informer *Informers) WatchQueueSettings() (cache.Store, cache.Controller) {
	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return informer.clientSet.QueueSettings("").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return informer.clientSet.QueueSettings("").Watch(lo)
			},
		},
		&queuesettings.QueueSettings{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				typed := obj.(*queuesettings.QueueSettings)
				LogForQueueSetttings(*typed, "added")
				informer.ReconcileQueueSetting(typed)
			},
			UpdateFunc: func(old, new interface{}) {
				typed := new.(*queuesettings.QueueSettings)
				LogForQueueSetttings(*typed, "updated")

				informer.ReconcileQueueSetting(typed)
			},
			DeleteFunc: func(obj interface{}) {
				typed := obj.(*queuesettings.QueueSettings)
				LogForQueueSetttings(*typed, "deleted")
			},
		},
	)

	go controller.Run(wait.NeverStop)
	return store, controller
}

func (informer *Informers) ReconcileQueueSetting(qs *queuesettings.QueueSettings) {
	if qs.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(qs, finalizerName) {
			LogForQueueSetttings(*qs, "adding finalizer")
			controllerutil.AddFinalizer(qs, finalizerName)

			_, err := informer.clientSet.QueueSettings(qs.ObjectMeta.Namespace).Update(qs, metav1.UpdateOptions{})
			if err != nil {
				LogForQueueSetttings(*qs, err.Error())
			}
			return
		}

		informer.jobReconcileRequest(qs.Name)
	} else {
		if controllerutil.ContainsFinalizer(qs, finalizerName) {
			LogForQueueSetttings(*qs, "removing finalizer")
			controllerutil.RemoveFinalizer(qs, finalizerName)

			_, err := informer.clientSet.QueueSettings(qs.ObjectMeta.Namespace).Update(qs, metav1.UpdateOptions{})
			if err != nil {
				LogForQueueSetttings(*qs, "could not remove finalizer")
				LogForQueueSetttings(*qs, err.Error())
			}
		}
	}
}
