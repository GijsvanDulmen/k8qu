package v1alpha1

import (
	k8qucontroller "k8qu/pkg/apis/k8qu"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type QueueSettingsClientInterface interface {
	QueueSettings(namespace string) QueueSettingsInterface
}

type QueueSettingsClient struct {
	restClient rest.Interface
}

func NewForQueueSettings(c *rest.Config) (*QueueSettingsClient, error) {
	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: k8qucontroller.GroupName, Version: "v1alpha1"}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &QueueSettingsClient{restClient: client}, nil
}

func (c *QueueSettingsClient) QueueSettings(namespace string) QueueSettingsInterface {
	return &queuesettingsClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
