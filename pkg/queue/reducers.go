package queue

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/markqueuejobcomplete"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
)

func ProcessForCompletedJobs(jobs []*queuejob.QueueJob, jobUpdater JobUpdater, settings Settings) ([]*queuejob.QueueJob, error) {
	return ReduceJobs(jobs, func(jb queuejob.QueueJob) (bool, error) {
		if jb.IsCompleted() {
			err, isCompletedAndCanBeDeleted := jb.IsCompletedAndCanBeDeleted(settings.GetTtlAfterSuccessfulCompletion(), settings.GetTtlAfterFailedCompletion())
			if err != nil {
				return false, err
			} else if isCompletedAndCanBeDeleted {
				err = jobUpdater.DeleteJob(&jb)
				if err != nil {
					return false, err
				}
			}
			return false, err
		}
		return true, nil
	})
}

func ProcessForMarkQueueJobCompletedJobs(jobs []*queuejob.QueueJob, jobUpdater JobUpdater, markQueueJobComplete *markqueuejobcomplete.MarkQueueJobComplete) ([]*queuejob.QueueJob, error) {
	markQueueProcessed := false

	return ReduceJobs(jobs, func(jb queuejob.QueueJob) (bool, error) {
		if markQueueProcessed {
			return false, nil
		}

		if jb.IsRunning() {
			if markQueueJobComplete.Spec.Completed != nil {
				jb.MarkCompleted()
				err := jobUpdater.UpdateJobForCompletion(&jb)
				if err != nil {
					return true, err
				}
				markQueueProcessed = true
			} else if markQueueJobComplete.Spec.Failed != nil {
				jb.MarkFailed()
				err := jobUpdater.UpdateJobForFailure(&jb)
				if err != nil {
					return true, err
				}
				markQueueProcessed = true
			} else if len(markQueueJobComplete.Spec.CompletedParts) > 0 {
				// mark with parts
				for _, v := range markQueueJobComplete.Spec.CompletedParts {
					if jb.Spec.CompletedParts == nil {
						jb.Spec.CompletedParts = map[string]bool{}
					}
					jb.Spec.CompletedParts[v] = true
				}

				if jb.HasAllCompletedParts() {
					jb.MarkCompleted()
					err := jobUpdater.UpdateJobForCompletion(&jb)
					if err != nil {
						return false, err
					}
				} else {
					err := jobUpdater.UpdateJob(&jb)
					if err != nil {
						return true, err
					}
				}
				markQueueProcessed = true
			}
		}
		return true, nil
	})
}

func ProcessForTooLongInQueue(jobs []*queuejob.QueueJob, jobUpdater JobUpdater, maxTimeInQueue string) ([]*queuejob.QueueJob, error) {
	return ReduceJobs(jobs, func(jb queuejob.QueueJob) (bool, error) {
		if jb.IsCompleted() {
			return true, nil
		}

		err, IsTooLongInQueue := jb.IsTooLongInQueue(maxTimeInQueue)
		if err != nil {
			return false, err
		}

		if IsTooLongInQueue {
			jb.MarkTooLongInQueue()
			err = jobUpdater.UpdateJobForMaxTimeInQueue(&jb)
			if err != nil {
				return false, err
			}
			return false, nil
		}

		return true, nil
	})
}

func ProcessForExecutionTimeouts(jobs []*queuejob.QueueJob, jobUpdater JobUpdater, timeout string) ([]*queuejob.QueueJob, int64, error) {
	runningJobs := int64(0)
	reduceJobs, err := ReduceJobs(jobs, func(jb queuejob.QueueJob) (bool, error) {
		if jb.IsRunning() {
			err, isExecutionTimedOut := jb.IsExecutionTimedOut(timeout)
			if err != nil {
				return false, err
			}

			if isExecutionTimedOut {
				jb.MarkTimedOut()

				err = jobUpdater.UpdateJobForExecutionTimeout(&jb)
				if err != nil {
					return false, err
				}
				return false, nil
			} else {
				runningJobs++
				return false, nil
			}
		}
		return true, nil
	})

	return reduceJobs, runningJobs, err
}
