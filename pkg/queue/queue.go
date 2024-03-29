package queue

import (
	"fmt"
	"k8qu/pkg/apis/k8qu/v1alpha1/markqueuejobcomplete"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	logger "k8qu/pkg/log"
)

var log = logger.Logger()

type Queue struct {
	Name                 string
	Jobs                 []*queuejob.QueueJob
	MarkQueueJobComplete []*markqueuejobcomplete.MarkQueueJobComplete
	Settings             Settings
}

type Settings interface {
	GetParallelism() int64
	GetTtlAfterSuccessfulCompletion() string
	GetTtlAfterFailedCompletion() string
	GetExecutionTimeout() string
	GetMaxTimeInQueue() string
}

func NewQueue(name string, settings Settings) *Queue {
	return &Queue{
		Name:     name,
		Jobs:     []*queuejob.QueueJob{},
		Settings: settings,
	}
}

type JobUpdater interface {
	DeleteJob(jb *queuejob.QueueJob) error
	UpdateJobForExecutionTimeout(jb *queuejob.QueueJob) error
	UpdateJobForMaxTimeInQueue(jb *queuejob.QueueJob) error
	StartJob(jb *queuejob.QueueJob) bool
	UpdateJobForCompletion(jb *queuejob.QueueJob) error
	UpdateJob(jb *queuejob.QueueJob) error
	UpdateJobForFailure(jb *queuejob.QueueJob) error
	DeleteMarkQueueJobComplete(jb *markqueuejobcomplete.MarkQueueJobComplete) error
}

func (q *Queue) IsEmpty() bool {
	return len(q.Jobs) == 0
}

func (q *Queue) DebugLog(msg ...any) {
	log.Debug().Msgf("%s - %s", q.Name, fmt.Sprintf("%s", msg))
}

func (q *Queue) DebugErr(err error) {
	log.Debug().Err(err)
}

func (q *Queue) Add(addJob *queuejob.QueueJob) {
	q.Jobs = append(q.Jobs, addJob)
}

func (q *Queue) AddMarkQueueJobComplete(mqjc *markqueuejobcomplete.MarkQueueJobComplete) {
	q.MarkQueueJobComplete = append(q.MarkQueueJobComplete, mqjc)
}
