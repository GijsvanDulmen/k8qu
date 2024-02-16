package informer

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

func (informer *QueueSettingsInformer) WatchQueueSettings() (cache.Store, cache.Controller) {
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
				var cs = obj
				typed := cs.(*queuesettings.QueueSettings)
				LogForQueueSetttings(*typed, "added")
				informer.ReconcileQueueSetting(typed)
			},
			UpdateFunc: func(old, new interface{}) {
				var cs = new
				typed := cs.(*queuesettings.QueueSettings)
				LogForQueueSetttings(*typed, "updated")

				informer.ReconcileQueueSetting(typed)
			},
			DeleteFunc: func(obj interface{}) {
				var cs = obj
				typed := cs.(*queuesettings.QueueSettings)
				LogForQueueSetttings(*typed, "deleted")
			},
		},
	)

	go controller.Run(wait.NeverStop)
	return store, controller
}

func (informer *QueueSettingsInformer) ReconcileQueueSetting(cs *queuesettings.QueueSettings) {
	if cs.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(cs, finalizerName) {
			LogForQueueSetttings(*cs, "adding finalizer")
			controllerutil.AddFinalizer(cs, finalizerName)

			_, err := informer.clientSet.QueueSettings(cs.ObjectMeta.Namespace).Update(cs, metav1.UpdateOptions{})
			if err != nil {
				LogForQueueSetttings(*cs, err.Error())
			}
			return
		}

		informer.jobReconcileRequest(cs.Name)
	} else {
		if controllerutil.ContainsFinalizer(cs, finalizerName) {
			controllerutil.RemoveFinalizer(cs, finalizerName)

			LogForQueueSetttings(*cs, "removing finalizer")

			_, err := informer.clientSet.QueueSettings(cs.ObjectMeta.Namespace).Update(cs, metav1.UpdateOptions{})
			if err != nil {
				LogForQueueSetttings(*cs, "could not remove finalizer")
				LogForQueueSetttings(*cs, err.Error())
			}
		}
	}
}

func LogForQueueSetttings(qs queuesettings.QueueSettings, s string) {
	log.Debug().Msgf("queue setting %s/%s log: %s", qs.Namespace, qs.Name, s)
}
