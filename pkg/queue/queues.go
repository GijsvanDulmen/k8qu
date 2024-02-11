package queue

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
)

type Queues struct {
	Queues map[string]*Queue
}

func NewQueues() *Queues {
	return &Queues{
		Queues: map[string]*Queue{},
	}
}

func (q *Queues) AddQueue(queue string, settings Settings) {
	q.Queues[queue] = NewQueue(queue, settings)
}

func (q *Queues) AddJob(job *job.Job) {
	queueName := job.GetQueueName()
	if _, ok := q.Queues[queueName]; !ok {
		q.Queues[queueName] = NewQueue(queueName, Settings{
			Parallelism:                  1,
			TtlAfterFailedCompletion:     "",
			TtlAfterSuccesfullCompletion: "",
			Timeout:                      "",
			DeadlineTimeout:              "",
		})
	}
	q.Queues[queueName].Add(job)
}

func (q *Queues) Reconcile(jobUpdater JobUpdater) {
	for _, q := range q.Queues {
		q.Reconcile(jobUpdater)
	}
}
