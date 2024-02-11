package informer

import (
	"k8qu/pkg/clientset/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type QueueSettingsInformer struct {
	clientSet           v1alpha1.QueueSettingsClientInterface
	coreClientSet       *kubernetes.Clientset
	jobReconcileRequest JobReconcileRequest
}

func NewQueueSettingsInformer(clientSet v1alpha1.QueueSettingsClientInterface, coreClientSet *kubernetes.Clientset, jobReconcileRequest JobReconcileRequest) (informer *QueueSettingsInformer) {
	return &QueueSettingsInformer{
		clientSet:           clientSet,
		coreClientSet:       coreClientSet,
		jobReconcileRequest: jobReconcileRequest,
	}
}
