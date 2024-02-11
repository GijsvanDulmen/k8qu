package informer

import (
	"github.com/rs/zerolog"
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (informer *JobInformer) WatchJob() (cache.Store, cache.Controller) {
	store, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return informer.clientSet.Job("").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return informer.clientSet.Job("").Watch(lo)
			},
		},
		&job.Job{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				var cs = obj
				typed := cs.(*job.Job)
				LogForJob(*typed, "added", zerolog.DebugLevel)
				informer.ReconcileJob(typed)
			},
			UpdateFunc: func(old, new interface{}) {
				var cs = new
				typed := cs.(*job.Job)
				LogForJob(*typed, "updated", zerolog.DebugLevel)

				informer.ReconcileJob(typed)
			},
			DeleteFunc: func(obj interface{}) {
				var cs = obj
				typed := cs.(*job.Job)
				LogForJob(*typed, "deleted", zerolog.DebugLevel)
			},
		},
	)

	go controller.Run(wait.NeverStop)
	return store, controller
}

func (informer *JobInformer) ReconcileJob(cs *job.Job) {
	if cs.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(cs, finalizerName) {
			LogForJob(*cs, "adding finalizer", zerolog.DebugLevel)
			controllerutil.AddFinalizer(cs, finalizerName)

			_, err := informer.clientSet.Job(cs.ObjectMeta.Namespace).Update(cs, metav1.UpdateOptions{})
			if err != nil {
				LogForJob(*cs, err.Error(), zerolog.ErrorLevel)
			}
			return
		}

		// update completion data
		completed := cs.Spec.Completed
		if (completed != nil) && *completed && cs.Status.CompletedAt == nil {
			cs.MarkCompleted()
			_, err := informer.clientSet.Job(cs.ObjectMeta.Namespace).Update(cs, metav1.UpdateOptions{})
			if err != nil {
				LogForJob(*cs, err.Error(), zerolog.ErrorLevel)
			}
		}

		// update failed data
		failed := cs.Spec.Failed
		if (failed != nil) && *failed && cs.Status.CompletedAt == nil {
			cs.MarkFailed()
			_, err := informer.clientSet.Job(cs.ObjectMeta.Namespace).Update(cs, metav1.UpdateOptions{})
			if err != nil {
				LogForJob(*cs, err.Error(), zerolog.ErrorLevel)
			}
		}

		informer.jobReconcileRequest(cs.GetQueueName())
	} else {
		if controllerutil.ContainsFinalizer(cs, finalizerName) {
			controllerutil.RemoveFinalizer(cs, finalizerName)

			LogForJob(*cs, "removing finalizer", zerolog.DebugLevel)

			_, err := informer.clientSet.Job(cs.ObjectMeta.Namespace).Update(cs, metav1.UpdateOptions{})
			if err != nil {
				LogForJob(*cs, "could not remove finalizer", zerolog.ErrorLevel)
				LogForJob(*cs, err.Error(), zerolog.ErrorLevel)
			}
		}
	}
}

func LogForJob(j job.Job, s string, level zerolog.Level) {
	log.WithLevel(level).Msgf("job %s/%s log: %s", j.Namespace, j.Name, s)
}
