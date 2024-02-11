package job

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

type Job struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Spec   `json:"spec"`
	Status Status `json:"status"`
}

type Spec struct {
	Queue                        string `json:"queue"`
	Timeout                      string `json:"timeout"`
	DeadlineTimeout              string `json:"deadlineTimeout"`
	Completed                    *bool  `json:"completed,omitempty"`
	Failed                       *bool  `json:"failed,omitempty"`
	TtlAfterSuccesfullCompletion string `json:"ttlAfterSuccesfullCompletion,omitempty"`
	TtlAfterFailedCompletion     string `json:"ttlAfterFailedCompletion,omitempty"`

	Templates []*runtime.RawExtension `json:"templates,omitempty"`
}

type Status struct {
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	StartedAt    *time.Time `json:"startedAt,omitempty"`
	IsSuccesfull *bool      `json:"isSuccessfull,omitempty"`
	Outcome      *string    `json:"outcome"`
}

type JobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Job `json:"items"`
}
