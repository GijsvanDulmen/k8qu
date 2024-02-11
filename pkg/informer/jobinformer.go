package informer

import (
	"k8qu/pkg/clientset/v1alpha1"
	"k8s.io/client-go/kubernetes"
)

type JobInformer struct {
	clientSet           v1alpha1.JobClientInterface
	coreClientSet       *kubernetes.Clientset
	jobReconcileRequest JobReconcileRequest
}

func NewJobInformer(clientSet v1alpha1.JobClientInterface, coreClientSet *kubernetes.Clientset, jobReconcileRequest JobReconcileRequest) (informer *JobInformer) {
	return &JobInformer{
		clientSet:           clientSet,
		coreClientSet:       coreClientSet,
		jobReconcileRequest: jobReconcileRequest,
	}
}
