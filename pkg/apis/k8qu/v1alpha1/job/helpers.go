package job

import (
	"fmt"
	"time"
)

var timedOutMessage = "Timed out"
var deadlinedMessage = "Deadlined"
var markedCompletedMessage = "Marked completed"
var markedFailedMessage = "Marked failed"

func (jb *Job) IsRunning() bool {
	return jb.Status.StartedAt != nil && jb.Status.CompletedAt == nil
}

func (jb *Job) IsCompleted() bool {
	return jb.Status.CompletedAt != nil
}

func (jb *Job) IsCompletedAndCanBeDeleted(ttsSuccess string, ttsFailed string) (error, bool) {
	ttlAfterSuccesfullCompletion := jb.Spec.TtlAfterSuccesfullCompletion
	if jb.Spec.TtlAfterSuccesfullCompletion != "" {
		ttlAfterSuccesfullCompletion = jb.Spec.TtlAfterSuccesfullCompletion
	}
	if ttsSuccess != "" {
		ttlAfterSuccesfullCompletion = ttsSuccess
	}

	ttlAfterFailedCompletion := jb.Spec.TtlAfterFailedCompletion
	if jb.Spec.TtlAfterFailedCompletion != "" {
		ttlAfterFailedCompletion = jb.Spec.TtlAfterFailedCompletion
	}
	if ttsFailed != "" {
		ttlAfterFailedCompletion = ttsFailed
	}

	if jb.Status.CompletedAt != nil && jb.Status.IsSuccesfull != nil {
		if ttlAfterSuccesfullCompletion != "" &&
			(*jb.Status.IsSuccesfull) {
			timeDuration, err := time.ParseDuration(ttlAfterSuccesfullCompletion)

			if err != nil {
				return err, false
			}
			return nil, time.Now().After(jb.Status.CompletedAt.Add(timeDuration))
		} else if ttlAfterFailedCompletion != "" &&
			!(*jb.Status.IsSuccesfull) {
			timeDuration, err := time.ParseDuration(ttlAfterFailedCompletion)

			if err != nil {
				return err, false
			}
			return nil, time.Now().After(jb.Status.CompletedAt.Add(timeDuration))
		}
	}
	return nil, false
}

func (jb *Job) IsTimedOut(globalTimeout string) (error, bool) {
	if jb.Status.StartedAt == nil {
		return nil, false
	}

	if jb.Status.CompletedAt != nil {
		return nil, false
	}

	timeout := jb.Spec.Timeout
	if globalTimeout != "" {
		timeout = globalTimeout
	}

	if timeout != "" {
		timeDuration, err := time.ParseDuration(timeout)

		if err != nil {
			return err, false
		}

		return nil, time.Now().After(jb.Status.StartedAt.Add(timeDuration))
	}
	return nil, false
}

func (jb *Job) IsDeadlinedTimeout(globalDeadlineTimeout string) (error, bool) {
	if jb.Status.StartedAt != nil {
		return nil, false
	}

	deadlineTimeout := jb.Spec.DeadlineTimeout
	if globalDeadlineTimeout != "" {
		deadlineTimeout = globalDeadlineTimeout
	}

	if deadlineTimeout != "" {
		timeDuration, err := time.ParseDuration(deadlineTimeout)

		if err != nil {
			return err, false
		}

		return nil, time.Now().After(jb.ObjectMeta.CreationTimestamp.Time.Add(timeDuration))
	}
	return nil, false
}

func (jb *Job) MarkTimedOut() {
	now := time.Now()
	falseBool := false

	jb.Status.CompletedAt = &now
	jb.Status.IsSuccesfull = &falseBool
	jb.Status.Outcome = &timedOutMessage
}

func (jb *Job) MarkDeadlinedTimeout() {
	now := time.Now()
	falseBool := false

	jb.Status.CompletedAt = &now
	jb.Status.IsSuccesfull = &falseBool
	jb.Status.Outcome = &deadlinedMessage
}

func (jb *Job) MarkRunning() {
	now := time.Now()
	jb.Status.StartedAt = &now
}

func (jb *Job) MarkCompleted() {
	now := time.Now()
	trueBool := true
	jb.Status.CompletedAt = &now
	jb.Status.IsSuccesfull = &trueBool
	jb.Status.Outcome = &markedCompletedMessage
}

func (jb *Job) MarkFailed() {
	now := time.Now()
	falseBool := false
	jb.Status.CompletedAt = &now
	jb.Status.IsSuccesfull = &falseBool
	jb.Status.Outcome = &markedFailedMessage
}

func (jb *Job) GetQueueName() string {
	return fmt.Sprintf("%s.%s", jb.Namespace, jb.Spec.Queue)
}
