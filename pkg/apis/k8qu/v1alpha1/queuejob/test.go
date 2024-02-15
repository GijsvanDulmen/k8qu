package queuejob

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func CreateMockJob() QueueJob {
	return QueueJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "mock",
			GenerateName:               "",
			Namespace:                  "ns",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "",
			Generation:                 0,
			CreationTimestamp:          metav1.NewTime(time.Now()),
			DeletionTimestamp:          nil,
			DeletionGracePeriodSeconds: nil,
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            nil,
			Finalizers:                 nil,
			ManagedFields:              nil,
		},
		Spec: Spec{
			Queue:                        "queue",
			ExecutionTimeout:             "",
			MaxTimeInQueue:               "",
			Completed:                    nil,
			Failed:                       nil,
			TtlAfterSuccessfulCompletion: "",
			TtlAfterFailedCompletion:     "",
			Templates:                    nil,
		},
		Status: Status{},
	}
}
