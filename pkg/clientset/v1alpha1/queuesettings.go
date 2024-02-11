package v1alpha1

import (
	"context"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuesettings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const queuesettingsPlural = "queuesettings"

type QueueSettingsInterface interface {
	List(opts metav1.ListOptions) (*queuesettings.QueueSettingsList, error)
	Get(name string, options metav1.GetOptions) (*queuesettings.QueueSettings, error)
	Create(settings *queuesettings.QueueSettings) (*queuesettings.QueueSettings, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Update(*queuesettings.QueueSettings, metav1.UpdateOptions) (*queuesettings.QueueSettings, error)
	Delete(settings *queuesettings.QueueSettings, options metav1.DeleteOptions) error
}

type queuesettingsClient struct {
	restClient rest.Interface
	ns         string
}

func (c *queuesettingsClient) List(opts metav1.ListOptions) (*queuesettings.QueueSettingsList, error) {
	result := queuesettings.QueueSettingsList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(queuesettingsPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *queuesettingsClient) Get(name string, opts metav1.GetOptions) (*queuesettings.QueueSettings, error) {
	result := queuesettings.QueueSettings{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource(queuesettingsPlural).
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *queuesettingsClient) Delete(qs *queuesettings.QueueSettings, opts metav1.DeleteOptions) error {
	result := queuesettings.QueueSettings{}
	err := c.restClient.
		Delete().
		Namespace(c.ns).
		Resource(queuesettingsPlural).
		Name(qs.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.Background()).
		Into(&result)

	return err
}

func (c *queuesettingsClient) Update(jb *queuesettings.QueueSettings, opts metav1.UpdateOptions) (*queuesettings.QueueSettings, error) {
	result := queuesettings.QueueSettings{}
	err := c.restClient.Put().
		Namespace(c.ns).
		Resource(queuesettingsPlural).
		Name(jb.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(jb).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *queuesettingsClient) Create(qs *queuesettings.QueueSettings) (*queuesettings.QueueSettings, error) {
	result := queuesettings.QueueSettings{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource(queuesettingsPlural).
		Body(qs).
		Do(context.Background()).
		Into(&result)

	return &result, err
}

func (c *queuesettingsClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource(queuesettingsPlural).
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.Background())
}
