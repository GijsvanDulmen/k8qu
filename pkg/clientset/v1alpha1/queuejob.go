package v1alpha1

import (
	"context"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const queueJobPlural = "queuejobs"

type QueueJobInterface interface {
	List(opts metav1.ListOptions) (*queuejob.QueueJobList, error)
	Get(name string, options metav1.GetOptions) (*queuejob.QueueJob, error)
	Create(*queuejob.QueueJob) (*queuejob.QueueJob, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Update(*queuejob.QueueJob, metav1.UpdateOptions) (*queuejob.QueueJob, error)
	Delete(job *queuejob.QueueJob, options metav1.DeleteOptions) error
}

type queueJobClient struct {
	restClient rest.Interface
	ns         string
}

func (c *queueJobClient) List(opts metav1.ListOptions) (*queuejob.QueueJobList, error) {
	result := queuejob.QueueJobList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(queueJobPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *queueJobClient) Get(name string, opts metav1.GetOptions) (*queuejob.QueueJob, error) {
	result := queuejob.QueueJob{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(queueJobPlural).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *queueJobClient) Delete(jb *queuejob.QueueJob, opts metav1.DeleteOptions) error {
	result := queuejob.QueueJob{}
	err := c.restClient.
		Delete().
		Namespace(c.ns).
		Resource(queueJobPlural).
		Name(jb.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return err
}

func (c *queueJobClient) Update(jb *queuejob.QueueJob, opts metav1.UpdateOptions) (*queuejob.QueueJob, error) {
	result := queuejob.QueueJob{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(queueJobPlural).
		Name(jb.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(jb).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *queueJobClient) Create(jb *queuejob.QueueJob) (*queuejob.QueueJob, error) {
	result := queuejob.QueueJob{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(queueJobPlural).
		Body(jb).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *queueJobClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(queueJobPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}
