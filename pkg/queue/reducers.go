package queue

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
)

func ProcessForCompletedJobs(jobs []*job.Job, jobUpdater JobUpdater, settings Settings) ([]*job.Job, error) {
	return ReduceJobs(jobs, func(jb job.Job) (bool, error) {
		if jb.IsCompleted() {
			err, isCompletedAndCanBeDeleted := jb.IsCompletedAndCanBeDeleted(settings.TtlAfterSuccesfullCompletion, settings.TtlAfterFailedCompletion)
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

func ProcessForDeadlineTimeout(jobs []*job.Job, jobUpdater JobUpdater, deadlineTimeout string) ([]*job.Job, error) {
	return ReduceJobs(jobs, func(jb job.Job) (bool, error) {
		if jb.IsCompleted() {
			return true, nil
		}

		// check if deadlined
		err, isDeadlinedTimeout := jb.IsDeadlinedTimeout(deadlineTimeout)
		if err != nil {
			return false, err
		}

		if isDeadlinedTimeout {
			jb.MarkDeadlinedTimeout()
			err = jobUpdater.UpdateJobForDeadlineTimeout(&jb)
			if err != nil {
				return false, err
			}
			return false, nil
		}

		return true, nil
	})
}

func ProcessForTimeouts(jobs []*job.Job, jobUpdater JobUpdater, timeout string) ([]*job.Job, int64, error) {
	runningJobs := int64(0)
	reduceJobs, err := ReduceJobs(jobs, func(jb job.Job) (bool, error) {
		if jb.IsRunning() {
			err, isTimedOut := jb.IsTimedOut(timeout)
			if err != nil {
				return false, err
			}

			if isTimedOut {
				jb.MarkTimedOut()

				err = jobUpdater.UpdateJobForTimeout(&jb)
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
