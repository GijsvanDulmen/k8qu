package queuejob

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestIsCompleted(t *testing.T) {
	job := CreateMockJob()
	now := time.Now()
	job.Status.CompletedAt = &now

	if !job.IsCompleted() {
		t.Fail()
	}

	job = CreateMockJob()
	if job.IsCompleted() {
		t.Fail()
	}
}

func TestIsStarted(t *testing.T) {
	job := CreateMockJob()
	now := time.Now()
	job.Status.StartedAt = &now

	if !job.IsRunning() {
		t.Fail()
	}

	job = CreateMockJob()
	if job.IsRunning() {
		t.Fail()
	}
}

func TestIsCompletedAndCanBeDeleted(t *testing.T) {
	job := CreateMockJob()
	duration, _ := time.ParseDuration("10s")
	past := time.Now().Add(-duration)

	boolV := true

	job.Spec.TtlAfterSuccessfulCompletion = "5s"
	job.Status.CompletedAt = &past
	job.Status.IsSuccessful = &boolV

	_, isCompletedAndCanBeDeleted := job.IsCompletedAndCanBeDeleted("", "")
	if !isCompletedAndCanBeDeleted {
		t.Fail()
	}
}

func TestIsCompletedAndCanBeDeletedFalse(t *testing.T) {
	job := CreateMockJob()
	duration, _ := time.ParseDuration("10s")
	past := time.Now().Add(-duration)

	boolV := true

	job.Spec.TtlAfterSuccessfulCompletion = "15s"
	job.Status.CompletedAt = &past
	job.Status.IsSuccessful = &boolV

	_, isCompletedAndCanBeDeleted := job.IsCompletedAndCanBeDeleted("", "")
	if isCompletedAndCanBeDeleted {
		t.Fail()
	}
}

func TestIsCompletedAndCanBeDeletedFailed(t *testing.T) {
	job := CreateMockJob()
	duration, _ := time.ParseDuration("10s")
	past := time.Now().Add(-duration)

	falseV := false

	job.Spec.TtlAfterFailedCompletion = "5s"
	job.Spec.TtlAfterSuccessfulCompletion = "20s"
	job.Status.CompletedAt = &past
	job.Status.IsSuccessful = &falseV

	_, isCompletedAndCanBeDeleted := job.IsCompletedAndCanBeDeleted("", "")
	if !isCompletedAndCanBeDeleted {
		t.Fail()
	}
}

func TestIsCompletedAndCanBeDeletedFalseFailed(t *testing.T) {
	job := CreateMockJob()
	duration, _ := time.ParseDuration("10s")
	past := time.Now().Add(-duration)

	falseV := false

	job.Spec.TtlAfterFailedCompletion = "15s"
	job.Spec.TtlAfterSuccessfulCompletion = "5s"
	job.Status.CompletedAt = &past
	job.Status.IsSuccessful = &falseV

	_, isCompletedAndCanBeDeleted := job.IsCompletedAndCanBeDeleted("", "")
	if isCompletedAndCanBeDeleted {
		t.Fail()
	}
}

func TestIsTimedOut(t *testing.T) {
	job := CreateMockJob()
	duration, _ := time.ParseDuration("10s")
	past := time.Now().Add(-duration)

	job.Spec.ExecutionTimeout = "5s"
	job.Status.StartedAt = &past

	_, isTimedOut := job.IsExecutionTimedOut("")
	if !isTimedOut {
		t.Fail()
	}
}

func TestIsTimedOutFalse(t *testing.T) {
	job := CreateMockJob()
	duration, _ := time.ParseDuration("4s")
	past := time.Now().Add(-duration)

	job.Spec.ExecutionTimeout = "5s"
	job.Status.StartedAt = &past

	_, isTimedOut := job.IsExecutionTimedOut("")
	if isTimedOut {
		t.Fail()
	}
}

func TestIsTooLongInQueue(t *testing.T) {
	job := CreateMockJob()
	duration, _ := time.ParseDuration("10s")
	past := time.Now().Add(-duration)

	job.Spec.MaxTimeInQueue = "5s"
	job.ObjectMeta.CreationTimestamp = metav1.NewTime(past)

	_, isTimedOut := job.IsTooLongInQueue("")
	if !isTimedOut {
		t.Fail()
	}
}

func TestIsTooLongInQueueFalse(t *testing.T) {
	job := CreateMockJob()
	duration, _ := time.ParseDuration("4s")
	past := time.Now().Add(-duration)

	job.Spec.MaxTimeInQueue = "5s"
	job.ObjectMeta.CreationTimestamp = metav1.NewTime(past)

	_, isTimedOut := job.IsTooLongInQueue("")
	if isTimedOut {
		t.Fail()
	}
}
