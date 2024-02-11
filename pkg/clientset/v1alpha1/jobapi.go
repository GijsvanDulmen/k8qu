package v1alpha1

import (
	k8qucontroller "k8qu/pkg/apis/k8qu"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type JobClientInterface interface {
	Job(namespace string) JobInterface
}

type Client struct {
	restClient rest.Interface
}

func NewForJob(c *rest.Config) (*Client, error) {
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

func (c *Client) Job(namespace string) JobInterface {
	return &jobClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
