package queue

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
	"testing"
)

func TestAddQueue(t *testing.T) {
	queues := NewQueues()

	jb := job.CreateMockJob()
	jb.Spec.Queue = "a"
	jb.ObjectMeta.Namespace = "ns"

	queues.AddJob(&jb)

	if len(queues.Queues) != 1 {
		t.Fail()
	}

	if queue, ok := queues.Queues["ns.a"]; ok {
		if queue.Settings.Parallelism != 1 {
			t.Fail()
		}

		if queue.Settings.TtlAfterFailedCompletion != "" {
			t.Fail()
		}

		if queue.Settings.TtlAfterFailedCompletion != "" {
			t.Fail()
		}

		if queue.Settings.DeadlineTimeout != "" {
			t.Fail()
		}

		if queue.Settings.Timeout != "" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}

func TestAddQueueWithQueueSpec(t *testing.T) {
	queues := NewQueues()

	jb := job.CreateMockJob()
	jb.Spec.Queue = "a"
	jb.ObjectMeta.Namespace = "ns"

	queues.AddQueue("ns.a", Settings{
		Parallelism:                  2,
		TtlAfterSuccesfullCompletion: "2s",
		TtlAfterFailedCompletion:     "2s",
		Timeout:                      "2s",
		DeadlineTimeout:              "2s",
	})
	queues.AddJob(&jb)

	if len(queues.Queues) != 1 {
		t.Fail()
	}

	if queue, ok := queues.Queues["ns.a"]; ok {
		if queue.Settings.Parallelism != 2 {
			t.Fail()
		}

		if queue.Settings.TtlAfterFailedCompletion != "2s" {
			t.Fail()
		}

		if queue.Settings.TtlAfterFailedCompletion != "2s" {
			t.Fail()
		}

		if queue.Settings.DeadlineTimeout != "2s" {
			t.Fail()
		}

		if queue.Settings.Timeout != "2s" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
