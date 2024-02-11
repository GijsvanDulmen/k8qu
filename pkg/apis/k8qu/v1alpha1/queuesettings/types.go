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
	TtlAfterSuccesfullCompletion string `json:"ttlAfterSuccesfullCompletion"`
	TtlAfterFailedCompletion     string `json:"ttlAfterFailedCompletion"`
	Timeout                      string `json:"timeout"`
	DeadlineTimeout              string `json:"deadlineTimeout"`
}

type QueueSettingsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []QueueSettings `json:"items"`
}
