package v1alpha1

import (
	"context"
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const jobPlural = "jobs"

type JobInterface interface {
	List(opts metav1.ListOptions) (*job.JobList, error)
	Get(name string, options metav1.GetOptions) (*job.Job, error)
	Create(*job.Job) (*job.Job, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Update(*job.Job, metav1.UpdateOptions) (*job.Job, error)
	Delete(job *job.Job, options metav1.DeleteOptions) error
}

type jobClient struct {
	restClient rest.Interface
	ns         string
}

func (c *jobClient) List(opts metav1.ListOptions) (*job.JobList, error) {
	result := job.JobList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(jobPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *jobClient) Get(name string, opts metav1.GetOptions) (*job.Job, error) {
	result := job.Job{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(jobPlural).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *jobClient) Delete(jb *job.Job, opts metav1.DeleteOptions) error {
	result := job.Job{}
	err := c.restClient.
		Delete().
		Namespace(c.ns).
		Resource(jobPlural).
		Name(jb.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return err
}

func (c *jobClient) Update(jb *job.Job, opts metav1.UpdateOptions) (*job.Job, error) {
	result := job.Job{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(jobPlural).
		Name(jb.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(jb).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *jobClient) Create(jb *job.Job) (*job.Job, error) {
	result := job.Job{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(jobPlural).
		Body(jb).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *jobClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(jobPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}
