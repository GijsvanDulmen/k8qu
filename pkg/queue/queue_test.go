package queue

import (
	"errors"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"testing"
)

type JobUpdaterMock struct {
	Result error
}

func (j *JobUpdaterMock) UpdateJobForExecutionTimeout(jb *queuejob.QueueJob) error {
	return j.Result
}

func (j *JobUpdaterMock) UpdateJobForMaxTimeInQueue(jb *queuejob.QueueJob) error {
	return j.Result
}

func (j *JobUpdaterMock) StartJob(jb *queuejob.QueueJob) bool {
	return true
}

func (j *JobUpdaterMock) DeleteJob(jb *queuejob.QueueJob) error {
	return j.Result
}

func (j *JobUpdaterMock) UpdateJob(jb *queuejob.QueueJob) error {
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
