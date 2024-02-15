package v1alpha1

import (
	k8qucontroller "k8qu/pkg/apis/k8qu"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type QueueJobClientInterface interface {
	QueueJob(namespace string) QueueJobInterface
}

type Client struct {
	restClient rest.Interface
}

func NewForQueueJob(c *rest.Config) (*Client, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: k8qucontroller.GroupName, Version: "v1alpha1"}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{restClient: client}, nil
}

func (c *Client) QueueJob(namespace string) QueueJobInterface {
	return &queueJobClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
