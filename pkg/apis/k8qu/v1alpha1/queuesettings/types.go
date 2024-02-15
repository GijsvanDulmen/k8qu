package queuesettings

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type QueueSettings struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec Spec `json:"spec"`
}

type Spec struct {
	Parallelism                  int64  `json:"parallelism"`
	TtlAfterSuccessfulCompletion string `json:"ttlAfterSuccessfulCompletion"`
	TtlAfterFailedCompletion     string `json:"ttlAfterFailedCompletion"`
	ExecutionTimeout             string `json:"executionTimeout"`
	MaxTimeInQueue               string `json:"maxTimeInQueue"`
}

type QueueSettingsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []QueueSettings `json:"items"`
}
