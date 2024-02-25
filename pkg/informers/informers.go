package informers

import (
	"k8qu/pkg/clientset/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type Informers struct {
	clientSet           v1alpha1.ClientInterface
	coreClientSet       *kubernetes.Clientset
	jobReconcileRequest JobReconcileRequest
}

func NewInformers(clientSet v1alpha1.ClientInterface, coreClientSet *kubernetes.Clientset, jobReconcileRequest JobReconcileRequest) (informer *Informers) {
	return &Informers{
		clientSet:           clientSet,
		coreClientSet:       coreClientSet,
		jobReconcileRequest: jobReconcileRequest,
	}
}
