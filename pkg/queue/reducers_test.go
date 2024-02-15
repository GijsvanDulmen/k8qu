package queue

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestProcessForCompletionToBeDeleted(t *testing.T) {
	jb := queuejob.CreateMockJob()
	jb.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now())

	jb2 := queuejob.CreateMockJob()
	now := time.Now()
	jb2.ObjectMeta.CreationTimestamp = metav1.NewTime(now)
	jb2.Status.CompletedAt = &now

	jobs := []*queuejob.QueueJob{&jb, &jb2}

	resultJobs, err := ProcessForCompletedJobs(jobs, &JobUpdaterMock{}, Settings{
		Parallelism:                  0,
		TtlAfterSuccessfulCompletion: "",
		TtlAfterFailedCompletion:     "",
		ExecutionTimeout:             "",
		MaxTimeInQueue:               "",
	})

	assert.Equal(t, 1, len(resultJobs))
	assert.Nil(t, err)
}

func TestProcessForCompletionToBeDeletedErrored(t *testing.T) {
	jb2 := queuejob.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)
	trueV := true

	jb2.ObjectMeta.CreationTimestamp = metav1.NewTime(past)
	jb2.Status.CompletedAt = &past
	jb2.Status.IsSuccessful = &trueV
	jb2.Spec.TtlAfterSuccessfulCompletion = "5s"

	jobs := []*queuejob.QueueJob{&jb2}

	resultJobs, err := ProcessForCompletedJobs(jobs, &JobUpdaterMock{Result: errors.New("broken")}, Settings{
		Parallelism:                  0,
		TtlAfterSuccessfulCompletion: "",
		TtlAfterFailedCompletion:     "",
		ExecutionTimeout:             "",
		MaxTimeInQueue:               "",
	})

	assert.Equal(t, 0, len(resultJobs))
	assert.NotNil(t, err)
}

func TestProcessForTooLongInQueue(t *testing.T) {
	// past job
	jb := queuejob.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)

	jb.Spec.MaxTimeInQueue = "5s"
	jb.ObjectMeta.CreationTimestamp = metav1.NewTime(past)

	jobs := []*queuejob.QueueJob{&jb}

	_, returnV := jb.IsTooLongInQueue("5s")
	assert.True(t, returnV)

	retJobs, err := ProcessForTooLongInQueue(jobs, &JobUpdaterMock{}, "5s")

	assert.Equal(t, 0, len(retJobs))
	assert.Nil(t, err)
}

func TestProcessForMaxTimeInQueueErrored(t *testing.T) {
	// past job
	jb := queuejob.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)

	jb.Spec.MaxTimeInQueue = "5s"
	jb.ObjectMeta.CreationTimestamp = metav1.NewTime(past)

	jobs := []*queuejob.QueueJob{&jb}

	_, returnV := jb.IsTooLongInQueue("5s")
	assert.True(t, returnV)

	retJobs, err := ProcessForTooLongInQueue(jobs, &JobUpdaterMock{Result: errors.New("broken")}, "5s")

	assert.Equal(t, 0, len(retJobs))
	assert.NotNil(t, err)
}

func TestProcessForTimeouts(t *testing.T) {
	// timed out job
	timedOutJob := queuejob.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)

	timedOutJob.Spec.ExecutionTimeout = "5s"
	timedOutJob.ObjectMeta.CreationTimestamp = metav1.NewTime(past)
	timedOutJob.Status.StartedAt = &past

	// running job
	runningJob := queuejob.CreateMockJob()
	runningJob.ObjectMeta.CreationTimestamp = metav1.NewTime(past)
	runningJob.Status.StartedAt = &past
	runningJob.Status.CompletedAt = nil

	jobs := []*queuejob.QueueJob{&timedOutJob, &runningJob}

	_, returnV := timedOutJob.IsExecutionTimedOut("")
	assert.True(t, returnV)

	notRunningJobs, runningJobs, err := ProcessForExecutionTimeouts(jobs, &JobUpdaterMock{}, "")

	assert.Equal(t, 0, len(notRunningJobs))
	assert.Equal(t, int64(1), runningJobs)
	assert.Nil(t, err)
}

func TestProcessForTimeoutsErrored(t *testing.T) {
	// timed out job
	timedOutJob := queuejob.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)

	timedOutJob.Spec.ExecutionTimeout = "5s"
	timedOutJob.ObjectMeta.CreationTimestamp = metav1.NewTime(past)
	timedOutJob.Status.StartedAt = &past

	jobs := []*queuejob.QueueJob{&timedOutJob}

	notRunningJobs, runningJobs, err := ProcessForExecutionTimeouts(jobs, &JobUpdaterMock{Result: errors.New("broken")}, "")

	assert.Equal(t, 0, len(notRunningJobs))
	assert.Equal(t, int64(0), runningJobs)
	assert.NotNil(t, err)
}
