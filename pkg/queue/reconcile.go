package queue

import (
	"github.com/rs/zerolog"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"sort"
)

func (q *Queue) Reconcile(jobUpdater JobUpdater) {
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
	sort.SliceStable(notRunningJobs, func(i, j int) bool {
		return notRunningJobs[j].ObjectMeta.CreationTimestamp.Time.After(notRunningJobs[i].ObjectMeta.CreationTimestamp.Time)
	})

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
