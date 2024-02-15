package clientset

import (
	"github.com/rs/zerolog"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"k8qu/pkg/clientset/v1alpha1"
	logger "k8qu/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

var log = logger.Logger()

type QueueJobUpdater struct {
	Client          *v1alpha1.Client
	ServerResources discovery.ServerResourcesInterface
	Dynamic         dynamic.Interface
}

func (j *QueueJobUpdater) UpdateJobForCompletion(jb *queuejob.QueueJob) error {
	err := j.UpdateJob(jb)
	if err != nil {
		log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job could not be marked for completion %s", jb.GetQueueName(), jb.Name)
		log.WithLevel(zerolog.ErrorLevel).Err(err)
	} else {
		err = CreateTemplates(jb.Spec.OnCompletionTemplates, jb, j.ServerResources, j.Dynamic)
		if err != nil {
			log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job marked as completed but could not create resources %s", jb.GetQueueName(), jb.Name)
			log.WithLevel(zerolog.ErrorLevel).Msg(err.Error())
			log.WithLevel(zerolog.ErrorLevel).Err(err)
		}
	}
	return err
}

func (j *QueueJobUpdater) UpdateJobForFailure(jb *queuejob.QueueJob) error {
	err := j.UpdateJob(jb)
	if err != nil {
		log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job could not be marked for failure %s", jb.GetQueueName(), jb.Name)
		log.WithLevel(zerolog.ErrorLevel).Err(err)
	} else {
		err = CreateTemplates(jb.Spec.OnFailureTemplates, jb, j.ServerResources, j.Dynamic)
		if err != nil {
			log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job marked as failure but could not create resources %s", jb.GetQueueName(), jb.Name)
			log.WithLevel(zerolog.ErrorLevel).Msg(err.Error())
			log.WithLevel(zerolog.ErrorLevel).Err(err)
		}
	}
	return err
}

func (j *QueueJobUpdater) StartJob(nextJob *queuejob.QueueJob) bool {
	nextJob.MarkRunning()
	err := j.UpdateJob(nextJob)
	if err != nil {
		log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job should be marked as running but could not mark it %s", nextJob.GetQueueName(), nextJob.Name)
		log.WithLevel(zerolog.ErrorLevel).Err(err)
		return true
	}

	err = CreateTemplates(nextJob.Spec.Templates, nextJob, j.ServerResources, j.Dynamic)
	if err != nil {
		log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job marked as running but could not create resources %s", nextJob.GetQueueName(), nextJob.Name)
		log.WithLevel(zerolog.ErrorLevel).Msg(err.Error())
		log.WithLevel(zerolog.ErrorLevel).Err(err)
		return false
	}
	return false
}

func (j *QueueJobUpdater) UpdateJobForMaxTimeInQueue(jb *queuejob.QueueJob) error {
	err := j.UpdateJob(jb)
	if err != nil {
		return err
	}

	err = CreateTemplates(jb.Spec.OnTooLongInQueueTemplates, jb, j.ServerResources, j.Dynamic)
	if err != nil {
		return err
	}
	return nil
}

func (j *QueueJobUpdater) UpdateJobForExecutionTimeout(jb *queuejob.QueueJob) error {
	err := j.UpdateJob(jb)
	if err != nil {
		return err
	}

	err = CreateTemplates(jb.Spec.OnExecutionTimeoutTemplates, jb, j.ServerResources, j.Dynamic)
	if err != nil {
		return err
	}
	return nil
}

func (j *QueueJobUpdater) DeleteJob(jb *queuejob.QueueJob) error {
	err := j.Client.QueueJob(jb.Namespace).Delete(jb, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (j *QueueJobUpdater) UpdateJob(jb *queuejob.QueueJob) error {
	_, err := j.Client.QueueJob(jb.Namespace).Update(jb, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func CreateTemplates(resources []*runtime.RawExtension, jb *queuejob.QueueJob, c discovery.ServerResourcesInterface, dc dynamic.Interface) error {
	if resources != nil {
		for _, template := range resources {
			resource := *template
			err := CreateResource(jb, resource, c, dc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
