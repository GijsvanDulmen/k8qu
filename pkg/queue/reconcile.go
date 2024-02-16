package queue

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"sort"
	"strings"
)

func (q *Queue) Reconcile(jobUpdater JobUpdater) {
	for _, jb := range q.Jobs {
		if jb.IsCompleted() {
			continue
		}
		// part based completion
		if len(jb.Spec.NeedsCompletedParts) > 0 {
			missesOnePart := false
			for _, completedPart := range jb.Spec.NeedsCompletedParts {
				if value, ok := jb.Spec.CompletedParts[completedPart]; ok {
					if !value {
						missesOnePart = true
					}
				} else {
					missesOnePart = true
				}
			}

			if !missesOnePart {
				jb.MarkCompleted()
				err := jobUpdater.UpdateJobForCompletion(jb)
				if err != nil {
					return
				}
			}
		} else {
			completed := jb.Spec.Completed
			if (completed != nil) && *completed && jb.Status.CompletedAt == nil {
				jb.MarkCompleted()
				err := jobUpdater.UpdateJobForCompletion(jb)
				if err != nil {
					return
				}
			}
		}

		// update failed data
		failed := jb.Spec.Failed
		if (failed != nil) && *failed && jb.Status.CompletedAt == nil {
			jb.MarkFailed()
			err := jobUpdater.UpdateJobForFailure(jb)
			if err != nil {
				return
			}
		}
	}

	log.Debug().Msgf("%s - total jobs %d", q.Name, len(q.Jobs))
	log.Debug().Msgf("%s - parallelism %d", q.Name, q.Settings.Parallelism)

	toBeDoneJobs, err := ProcessForTooLongInQueue(q.Jobs, jobUpdater, q.Settings.MaxTimeInQueue)
	if err != nil {
		q.DebugLog("could not process for max time in queue")
		q.DebugErr(err)
		return
	}

	toBeDoneJobs, err = ProcessForCompletedJobs(toBeDoneJobs, jobUpdater, q.Settings)
	if err != nil {
		q.DebugLog("could not process for completion to be deleted")
		q.DebugErr(err)
		return
	}

	log.Debug().Msgf("%s - jobs still running or need to run %d", q.Name, len(toBeDoneJobs))

	// get not running jobs
	notRunningJobs, numberOfRunning, err := ProcessForExecutionTimeouts(toBeDoneJobs, jobUpdater, q.Settings.ExecutionTimeout)
	if err != nil {
		q.DebugLog("could not process for execution timeouts")
		q.DebugErr(err)
		return
	}

	if numberOfRunning < q.Settings.Parallelism {

		toStartJobs := GetToStartJob(notRunningJobs, q.Settings.Parallelism, numberOfRunning)

		for i := range toStartJobs {
			if jobUpdater.StartJob(notRunningJobs[i]) {
				return // job failed to start
			}
		}
	}
}

func GetToStartJob(notRunningJobs []*queuejob.QueueJob, parallelism int64, numberOfRunning int64) []*queuejob.QueueJob {
	// sort by creation timestamp
	// equal = original order
	SortQueueJobs(notRunningJobs)

	numberToStart := parallelism - numberOfRunning

	jobsLength := len(notRunningJobs)

	var startJobs []*queuejob.QueueJob
	if jobsLength > 0 {
		for i := int64(0); i < numberToStart; i++ {
			if int64(jobsLength) > i {
				startJobs = append(startJobs, notRunningJobs[i])
			} else {
				break // no need to process further
			}
		}
	}
	return startJobs
}

func SortQueueJobs(notRunningJobs []*queuejob.QueueJob) {
	sort.SliceStable(notRunningJobs, func(i, j int) bool {
		if notRunningJobs[j].ObjectMeta.CreationTimestamp.Time.Unix() == notRunningJobs[i].ObjectMeta.CreationTimestamp.Time.Unix() {
			return strings.Compare(notRunningJobs[j].ObjectMeta.Name, notRunningJobs[i].ObjectMeta.Name) > 0
		}
		return notRunningJobs[j].ObjectMeta.CreationTimestamp.Time.After(notRunningJobs[i].ObjectMeta.CreationTimestamp.Time)
	})
}
