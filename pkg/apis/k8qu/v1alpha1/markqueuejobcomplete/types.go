package markqueuejobcomplete

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MarkQueueJobComplete struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec Spec `json:"spec"`
}

type Spec struct {
	Queue string `json:"queue"`

	Completed *bool `json:"completed,omitempty"`
	Failed    *bool `json:"failed,omitempty"`

	CompletedParts []string `json:"completedParts"`
}

type MarkQueueJobCompleteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []MarkQueueJobComplete `json:"items"`
}
