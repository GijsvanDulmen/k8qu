package queue

import (
	"errors"
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
	"testing"
)

type JobUpdaterMock struct {
	Result error
}

func (j *JobUpdaterMock) UpdateJobForTimeout(jb *job.Job) error {
	return j.Result
}

func (j *JobUpdaterMock) UpdateJobForDeadlineTimeout(jb *job.Job) error {
	return j.Result
}

func (j *JobUpdaterMock) StartJob(jb *job.Job) bool {
	return true
}

func (j *JobUpdaterMock) DeleteJob(jb *job.Job) error {
	return j.Result
}

func (j *JobUpdaterMock) UpdateJob(jb *job.Job) error {
	return j.Result
}

func TestEmpty(t *testing.T) {
	queue := NewQueue("a", Settings{})
	if !queue.IsEmpty() {
		t.Failed()
	}

	queue.DebugErr(errors.New("test"))
	queue.DebugLog("abc")
}
