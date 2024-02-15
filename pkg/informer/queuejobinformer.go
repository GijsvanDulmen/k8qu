package informer

import (
	"k8qu/pkg/clientset/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type QueueJobInformer struct {
	clientSet           v1alpha1.QueueJobClientInterface
	coreClientSet       *kubernetes.Clientset
	jobReconcileRequest JobReconcileRequest
}

func NewQueueJobInformer(clientSet v1alpha1.QueueJobClientInterface, coreClientSet *kubernetes.Clientset, jobReconcileRequest JobReconcileRequest) (informer *QueueJobInformer) {
	return &QueueJobInformer{
		clientSet:           clientSet,
		coreClientSet:       coreClientSet,
		jobReconcileRequest: jobReconcileRequest,
	}
}
