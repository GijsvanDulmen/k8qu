package queue

import (
	"github.com/rs/zerolog"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"sort"
	"strings"
)

func (q *Queue) Reconcile(jobUpdater JobUpdater) {
	for _, cs := range q.Jobs {
		completed := cs.Spec.Completed
		if (completed != nil) && *completed && cs.Status.CompletedAt == nil {
			cs.MarkCompleted()
			err := jobUpdater.UpdateJobForCompletion(cs)
			if err != nil {
				return
			}
		}

		// update failed data
		failed := cs.Spec.Failed
		if (failed != nil) && *failed && cs.Status.CompletedAt == nil {
			cs.MarkFailed()
			err := jobUpdater.UpdateJobForFailure(cs)
			if err != nil {
				return
			}
		}
	}

	log.WithLevel(zerolog.DebugLevel).Msgf("%s - total jobs %d", q.Name, len(q.Jobs))
	log.WithLevel(zerolog.DebugLevel).Msgf("%s - parallelism %d", q.Name, q.Settings.Parallelism)

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

	log.WithLevel(zerolog.DebugLevel).Msgf("%s - jobs still running or need to run %d", q.Name, len(toBeDoneJobs))

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
