package queue

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestProcessForCompletionToBeDeleted(t *testing.T) {
	jb := job.CreateMockJob()
	jb.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now())

	jb2 := job.CreateMockJob()
	now := time.Now()
	jb2.ObjectMeta.CreationTimestamp = metav1.NewTime(now)
	jb2.Status.CompletedAt = &now

	jobs := []*job.Job{&jb, &jb2}

	resultJobs, err := ProcessForCompletedJobs(jobs, &JobUpdaterMock{}, Settings{
		Parallelism:                  0,
		TtlAfterSuccesfullCompletion: "",
		TtlAfterFailedCompletion:     "",
		Timeout:                      "",
		DeadlineTimeout:              "",
	})

	assert.Equal(t, 1, len(resultJobs))
	assert.Nil(t, err)
}

func TestProcessForCompletionToBeDeletedErrored(t *testing.T) {
	jb2 := job.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)
	trueV := true

	jb2.ObjectMeta.CreationTimestamp = metav1.NewTime(past)
	jb2.Status.CompletedAt = &past
	jb2.Status.IsSuccesfull = &trueV
	jb2.Spec.TtlAfterSuccesfullCompletion = "5s"

	jobs := []*job.Job{&jb2}

	resultJobs, err := ProcessForCompletedJobs(jobs, &JobUpdaterMock{Result: errors.New("broken")}, Settings{
		Parallelism:                  0,
		TtlAfterSuccesfullCompletion: "",
		TtlAfterFailedCompletion:     "",
		Timeout:                      "",
		DeadlineTimeout:              "",
	})

	assert.Equal(t, 0, len(resultJobs))
	assert.NotNil(t, err)
}

func TestProcessForDeadlineTimeout(t *testing.T) {
	// past job
	jb := job.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)

	jb.Spec.DeadlineTimeout = "5s"
	jb.ObjectMeta.CreationTimestamp = metav1.NewTime(past)

	jobs := []*job.Job{&jb}

	_, returnV := jb.IsDeadlinedTimeout("5s")
	assert.True(t, returnV)

	retJobs, err := ProcessForDeadlineTimeout(jobs, &JobUpdaterMock{}, "5s")

	assert.Equal(t, 0, len(retJobs))
	assert.Nil(t, err)
}

func TestProcessForDeadlineTimeoutErrored(t *testing.T) {
	// past job
	jb := job.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)

	jb.Spec.DeadlineTimeout = "5s"
	jb.ObjectMeta.CreationTimestamp = metav1.NewTime(past)

	jobs := []*job.Job{&jb}

	_, returnV := jb.IsDeadlinedTimeout("5s")
	assert.True(t, returnV)

	retJobs, err := ProcessForDeadlineTimeout(jobs, &JobUpdaterMock{Result: errors.New("broken")}, "5s")

	assert.Equal(t, 0, len(retJobs))
	assert.NotNil(t, err)
}

func TestProcessForTimeouts(t *testing.T) {
	// timed out job
	timedOutJob := job.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)

	timedOutJob.Spec.Timeout = "5s"
	timedOutJob.ObjectMeta.CreationTimestamp = metav1.NewTime(past)
	timedOutJob.Status.StartedAt = &past

	// running job
	runningJob := job.CreateMockJob()
	runningJob.ObjectMeta.CreationTimestamp = metav1.NewTime(past)
	runningJob.Status.StartedAt = &past
	runningJob.Status.CompletedAt = nil

	jobs := []*job.Job{&timedOutJob, &runningJob}

	_, returnV := timedOutJob.IsTimedOut("")
	assert.True(t, returnV)

	notRunningJobs, runningJobs, err := ProcessForTimeouts(jobs, &JobUpdaterMock{}, "")

	assert.Equal(t, 0, len(notRunningJobs))
	assert.Equal(t, int64(1), runningJobs)
	assert.Nil(t, err)
}

func TestProcessForTimeoutsErrored(t *testing.T) {
	// timed out job
	timedOutJob := job.CreateMockJob()
	duration, _ := time.ParseDuration("15s")
	past := time.Now().Add(-duration)

	timedOutJob.Spec.Timeout = "5s"
	timedOutJob.ObjectMeta.CreationTimestamp = metav1.NewTime(past)
	timedOutJob.Status.StartedAt = &past

	jobs := []*job.Job{&timedOutJob}

	notRunningJobs, runningJobs, err := ProcessForTimeouts(jobs, &JobUpdaterMock{Result: errors.New("broken")}, "")

	assert.Equal(t, 0, len(notRunningJobs))
	assert.Equal(t, int64(0), runningJobs)
	assert.NotNil(t, err)
}
