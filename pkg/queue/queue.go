package queue

import (
	"fmt"
	"github.com/rs/zerolog"
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
	logger "k8qu/pkg/log"
)

var log = logger.Logger()

type Queue struct {
	Name     string
	Jobs     []*job.Job
	Settings Settings
}

type Settings struct {
	Parallelism                  int64
	TtlAfterSuccesfullCompletion string
	TtlAfterFailedCompletion     string
	Timeout                      string
	DeadlineTimeout              string
}

func NewQueue(name string, settings Settings) *Queue {
	return &Queue{
		Name:     name,
		Jobs:     []*job.Job{},
		Settings: settings,
	}
}

type JobUpdater interface {
	DeleteJob(jb *job.Job) error
	UpdateJobForTimeout(jb *job.Job) error
	UpdateJobForDeadlineTimeout(jb *job.Job) error
	StartJob(jb *job.Job) bool
}

func (q *Queue) IsEmpty() bool {
	return len(q.Jobs) == 0
}

func (q *Queue) DebugLog(msg ...any) {
	log.WithLevel(zerolog.DebugLevel).Msgf("%s - %s", q.Name, fmt.Sprintf("%s", msg))
}

func (q *Queue) DebugErr(err error) {
	log.WithLevel(zerolog.DebugLevel).Err(err)
}

func (q *Queue) Add(addJob *job.Job) {
	q.Jobs = append(q.Jobs, addJob)
}
