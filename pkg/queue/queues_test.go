package queue

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"testing"
)

func TestAddQueue(t *testing.T) {
	queues := NewQueues()

	jb := queuejob.CreateMockJob()
	jb.Spec.Queue = "a"
	jb.ObjectMeta.Namespace = "ns"

	queues.AddJob(&jb)

	if len(queues.Queues) != 1 {
		t.Fail()
	}

	if queue, ok := queues.Queues["ns.a"]; ok {
		if queue.Settings.GetParallelism() != 1 {
			t.Fail()
		}

		if queue.Settings.GetTtlAfterFailedCompletion() != "" {
			t.Fail()
		}

		if queue.Settings.GetTtlAfterFailedCompletion() != "" {
			t.Fail()
		}

		if queue.Settings.GetMaxTimeInQueue() != "" {
			t.Fail()
		}

		if queue.Settings.GetExecutionTimeout() != "" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}

func TestAddQueueWithQueueSpec(t *testing.T) {
	queues := NewQueues()

	jb := queuejob.CreateMockJob()
	jb.Spec.Queue = "a"
	jb.ObjectMeta.Namespace = "ns"

	queues.AddQueue("ns.a", DefaultSettings{
		Parallelism:                  2,
		TtlAfterSuccessfulCompletion: "2s",
		TtlAfterFailedCompletion:     "2s",
		ExecutionTimeout:             "2s",
		MaxTimeInQueue:               "2s",
	})
	queues.AddJob(&jb)

	if len(queues.Queues) != 1 {
		t.Fail()
	}

	if queue, ok := queues.Queues["ns.a"]; ok {
		if queue.Settings.GetParallelism() != 2 {
			t.Fail()
		}

		if queue.Settings.GetTtlAfterFailedCompletion() != "2s" {
			t.Fail()
		}

		if queue.Settings.GetTtlAfterFailedCompletion() != "2s" {
			t.Fail()
		}

		if queue.Settings.GetMaxTimeInQueue() != "2s" {
			t.Fail()
		}

		if queue.Settings.GetExecutionTimeout() != "2s" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
