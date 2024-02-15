package queuejob

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"time"
)

type QueueJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Spec   `json:"spec"`
	Status Status `json:"status"`
}

type Spec struct {
	Queue            string `json:"queue"`
	ExecutionTimeout string `json:"executionTimeout"`
	MaxTimeInQueue   string `json:"maxTimeInQueue"`

	Completed *bool `json:"completed,omitempty"`
	Failed    *bool `json:"failed,omitempty"`

	NeedsCompletedParts []string        `json:"needsCompletedParts"`
	CompletedParts      map[string]bool `json:"completedParts"`

	TtlAfterSuccessfulCompletion string `json:"ttlAfterSuccessfulCompletion,omitempty"`
	TtlAfterFailedCompletion     string `json:"ttlAfterFailedCompletion,omitempty"`

	Templates                   []*runtime.RawExtension `json:"templates,omitempty"`
	OnTooLongInQueueTemplates   []*runtime.RawExtension `json:"onTooLongInQueueTemplates,omitempty"`
	OnExecutionTimeoutTemplates []*runtime.RawExtension `json:"onExecutionTimeoutTemplates,omitempty"`
	OnCompletionTemplates       []*runtime.RawExtension `json:"onCompletionTemplates,omitempty"`
	OnFailureTemplates          []*runtime.RawExtension `json:"onFailureTemplates,omitempty"`
}

type Status struct {
	CompletedAt  *time.Time `json:"completedAt,omitempty"`
	StartedAt    *time.Time `json:"startedAt,omitempty"`
	IsSuccessful *bool      `json:"isSuccessful,omitempty"`
	Outcome      *string    `json:"outcome"`
}

type QueueJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []QueueJob `json:"items"`
}
