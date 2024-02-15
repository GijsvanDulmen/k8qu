package queuejob

import (
	"fmt"
	"time"
)

var timedOutMessage = "Timed out"
var maxTimeInQueueMessage = "Too long in queue"
var markedCompletedMessage = "Marked completed"
var markedFailedMessage = "Marked failed"

func (jb *QueueJob) IsRunning() bool {
	return jb.Status.StartedAt != nil && jb.Status.CompletedAt == nil
}

func (jb *QueueJob) IsCompleted() bool {
	return jb.Status.CompletedAt != nil
}

func (jb *QueueJob) IsCompletedAndCanBeDeleted(ttsSuccess string, ttsFailed string) (error, bool) {
	ttlAfterSuccessfulCompletion := jb.Spec.TtlAfterSuccessfulCompletion
	if jb.Spec.TtlAfterSuccessfulCompletion != "" {
		ttlAfterSuccessfulCompletion = jb.Spec.TtlAfterSuccessfulCompletion
	}
	if ttsSuccess != "" {
		ttlAfterSuccessfulCompletion = ttsSuccess
	}

	ttlAfterFailedCompletion := jb.Spec.TtlAfterFailedCompletion
	if jb.Spec.TtlAfterFailedCompletion != "" {
		ttlAfterFailedCompletion = jb.Spec.TtlAfterFailedCompletion
	}
	if ttsFailed != "" {
		ttlAfterFailedCompletion = ttsFailed
	}

	if jb.Status.CompletedAt != nil && jb.Status.IsSuccessful != nil {
		if ttlAfterSuccessfulCompletion != "" &&
			(*jb.Status.IsSuccessful) {
			timeDuration, err := time.ParseDuration(ttlAfterSuccessfulCompletion)

			if err != nil {
				return err, false
			}
			return nil, time.Now().After(jb.Status.CompletedAt.Add(timeDuration))
		} else if ttlAfterFailedCompletion != "" &&
			!(*jb.Status.IsSuccessful) {
			timeDuration, err := time.ParseDuration(ttlAfterFailedCompletion)

			if err != nil {
				return err, false
			}
			return nil, time.Now().After(jb.Status.CompletedAt.Add(timeDuration))
		}
	}
	return nil, false
}

func (jb *QueueJob) IsExecutionTimedOut(globalExecutionTimeout string) (error, bool) {
	if jb.Status.StartedAt == nil {
		return nil, false
	}

	if jb.Status.CompletedAt != nil {
		return nil, false
	}

	executionTimeout := jb.Spec.ExecutionTimeout
	if globalExecutionTimeout != "" {
		executionTimeout = globalExecutionTimeout
	}

	if executionTimeout != "" {
		timeDuration, err := time.ParseDuration(executionTimeout)

		if err != nil {
			return err, false
		}

		return nil, time.Now().After(jb.Status.StartedAt.Add(timeDuration))
	}
	return nil, false
}

func (jb *QueueJob) IsTooLongInQueue(globalMaxTimeInQueue string) (error, bool) {
	if jb.Status.StartedAt != nil {
		return nil, false
	}

	maxTimeInQueue := jb.Spec.MaxTimeInQueue
	if globalMaxTimeInQueue != "" {
		maxTimeInQueue = globalMaxTimeInQueue
	}

	if maxTimeInQueue != "" {
		timeDuration, err := time.ParseDuration(maxTimeInQueue)

		if err != nil {
			return err, false
		}

		return nil, time.Now().After(jb.ObjectMeta.CreationTimestamp.Time.Add(timeDuration))
	}
	return nil, false
}

func (jb *QueueJob) MarkTimedOut() {
	now := time.Now()
	falseBool := false

	jb.Status.CompletedAt = &now
	jb.Status.IsSuccessful = &falseBool
	jb.Status.Outcome = &timedOutMessage
}

func (jb *QueueJob) MarkTooLongInQueue() {
	now := time.Now()
	falseBool := false

	jb.Status.CompletedAt = &now
	jb.Status.IsSuccessful = &falseBool
	jb.Status.Outcome = &maxTimeInQueueMessage
}

func (jb *QueueJob) MarkRunning() {
	now := time.Now()
	jb.Status.StartedAt = &now
}

func (jb *QueueJob) MarkCompleted() {
	now := time.Now()
	trueBool := true
	jb.Status.CompletedAt = &now
	jb.Status.IsSuccessful = &trueBool
	jb.Status.Outcome = &markedCompletedMessage
}

func (jb *QueueJob) MarkFailed() {
	now := time.Now()
	falseBool := false
	jb.Status.CompletedAt = &now
	jb.Status.IsSuccessful = &falseBool
	jb.Status.Outcome = &markedFailedMessage
}

func (jb *QueueJob) GetQueueName() string {
	return fmt.Sprintf("%s.%s", jb.Namespace, jb.Spec.Queue)
}
