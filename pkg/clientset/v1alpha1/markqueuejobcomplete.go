package v1alpha1

import (
	"context"
	"k8qu/pkg/apis/k8qu/v1alpha1/markqueuejobcomplete"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const markqueuejobcompletesPlural = "markqueuejobcompletes"

type MarkQueueJobCompleteInterface interface {
	List(opts metav1.ListOptions) (*markqueuejobcomplete.MarkQueueJobCompleteList, error)
	Get(name string, options metav1.GetOptions) (*markqueuejobcomplete.MarkQueueJobComplete, error)
	Create(settings *markqueuejobcomplete.MarkQueueJobComplete) (*markqueuejobcomplete.MarkQueueJobComplete, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Update(*markqueuejobcomplete.MarkQueueJobComplete, metav1.UpdateOptions) (*markqueuejobcomplete.MarkQueueJobComplete, error)
	Delete(settings *markqueuejobcomplete.MarkQueueJobComplete, options metav1.DeleteOptions) error
}

type markQueueJobCompleteClient struct {
	restClient rest.Interface
	ns         string
}

func (c *markQueueJobCompleteClient) List(opts metav1.ListOptions) (*markqueuejobcomplete.MarkQueueJobCompleteList, error) {
	result := markqueuejobcomplete.MarkQueueJobCompleteList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(markqueuejobcompletesPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *markQueueJobCompleteClient) Get(name string, opts metav1.GetOptions) (*markqueuejobcomplete.MarkQueueJobComplete, error) {
	result := markqueuejobcomplete.MarkQueueJobComplete{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(markqueuejobcompletesPlural).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *markQueueJobCompleteClient) Delete(qs *markqueuejobcomplete.MarkQueueJobComplete, opts metav1.DeleteOptions) error {
	result := markqueuejobcomplete.MarkQueueJobComplete{}
	err := c.restClient.
		Delete().
		Namespace(c.ns).
		Resource(markqueuejobcompletesPlural).
		Name(qs.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return err
}

func (c *markQueueJobCompleteClient) Update(jb *markqueuejobcomplete.MarkQueueJobComplete, opts metav1.UpdateOptions) (*markqueuejobcomplete.MarkQueueJobComplete, error) {
	result := markqueuejobcomplete.MarkQueueJobComplete{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(markqueuejobcompletesPlural).
		Name(jb.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(jb).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *markQueueJobCompleteClient) Create(qs *markqueuejobcomplete.MarkQueueJobComplete) (*markqueuejobcomplete.MarkQueueJobComplete, error) {
	result := markqueuejobcomplete.MarkQueueJobComplete{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(markqueuejobcompletesPlural).
		Body(qs).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *markQueueJobCompleteClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(markqueuejobcompletesPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}
