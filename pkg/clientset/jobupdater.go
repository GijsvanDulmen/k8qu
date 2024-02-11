package clientset

import (
	"github.com/rs/zerolog"
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
	"k8qu/pkg/clientset/v1alpha1"
	logger "k8qu/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
)

var log = logger.Logger()

type JobUpdater struct {
	Client          *v1alpha1.Client
	ServerResources discovery.ServerResourcesInterface
	Dynamic         dynamic.Interface
}

func (j *JobUpdater) StartJob(nextJob *job.Job) bool {
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
		log.WithLevel(zerolog.ErrorLevel).Err(err)
		return false
	}
	return false
}

func (j *JobUpdater) UpdateJobForDeadlineTimeout(jb *job.Job) error {
	err := j.UpdateJob(jb)
	if err != nil {
		return err
	}

	err = CreateTemplates(jb.Spec.OnDeadlineTimeoutTemplates, jb, j.ServerResources, j.Dynamic)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobUpdater) UpdateJobForTimeout(jb *job.Job) error {
	err := j.UpdateJob(jb)
	if err != nil {
		return err
	}

	err = CreateTemplates(jb.Spec.OnTimeoutTemplates, jb, j.ServerResources, j.Dynamic)
	if err != nil {
		return err
	}
	return nil
}

func (j *JobUpdater) DeleteJob(jb *job.Job) error {
	err := j.Client.Job(jb.Namespace).Delete(jb, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (j *JobUpdater) UpdateJob(jb *job.Job) error {
	_, err := j.Client.Job(jb.Namespace).Update(jb, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func CreateTemplates(resources []*runtime.RawExtension, jb *job.Job, c discovery.ServerResourcesInterface, dc dynamic.Interface) error {
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
