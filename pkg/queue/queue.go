package queue

import (
	"github.com/rs/zerolog"
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
	"k8qu/pkg/clientset/v1alpha1"
	logger "k8qu/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"sort"
)

var log = logger.Logger()

type Queue struct {
	Name     string
	Jobs     []*job.Job
	Settings Settings
}

type Settings struct {
	Parallelism                  int64
	TtlAfterSuccesfullCompletion string
	TtlAfterFailedCompletion     string
	Timeout                      string
	DeadlineTimeout              string
}

func NewQueue(name string, settings Settings) *Queue {
	return &Queue{
		Name:     name,
		Jobs:     []*job.Job{},
		Settings: settings,
	}
}

func (q *Queue) IsEmpty() bool {
	return len(q.Jobs) == 0
}

func (q *Queue) Reconcile(client *v1alpha1.Client, c discovery.ServerResourcesInterface, dc dynamic.Interface) {
	log.WithLevel(zerolog.DebugLevel).Msgf("%s - total jobs %d", q.Name, len(q.Jobs))
	log.WithLevel(zerolog.DebugLevel).Msgf("%s - parallelism %d", q.Name, q.Settings.Parallelism)

	// remove all completed
	toBeDoneJobs, canNotProcessFurther := q.GetToBeDoneJobs(client, q.Settings.DeadlineTimeout)
	if canNotProcessFurther {
		return
	}

	log.WithLevel(zerolog.DebugLevel).Msgf("%s - jobs still running or need to run %d", q.Name, len(toBeDoneJobs))

	// get not running jobs
	notRunningJobs, numberOfRunning, canNotProcessFurther := q.GetNotRunningJobs(client, toBeDoneJobs, q.Settings.Timeout)
	if canNotProcessFurther {
		return
	}

	// 1 running - parallel 1 = OK
	// 0 running - parallel 1 = NOK
	log.WithLevel(zerolog.DebugLevel).Msgf("used parallism %d", q.Settings.Parallelism)
	log.WithLevel(zerolog.DebugLevel).Msgf("numberOfrunning %d", numberOfRunning)
	log.WithLevel(zerolog.DebugLevel).Msgf("number of not running %d", len(notRunningJobs))

	if numberOfRunning < q.Settings.Parallelism {

		// sort by creation timestamp
		// equal = original order
		sort.SliceStable(notRunningJobs, func(i, j int) bool {
			return notRunningJobs[i].ObjectMeta.CreationTimestamp.Time.After(notRunningJobs[j].ObjectMeta.CreationTimestamp.Time)
		})

		// get the first X
		numberToStart := q.Settings.Parallelism - numberOfRunning

		log.WithLevel(zerolog.DebugLevel).Msgf("have to start %d jobs", len(notRunningJobs))
		log.WithLevel(zerolog.DebugLevel).Msgf("numberToStart %d jobs", numberToStart)

		if len(notRunningJobs) > 0 {
			for i := int64(0); i < numberToStart; i++ {
				if int64(len(notRunningJobs)) > i {
					log.WithLevel(zerolog.DebugLevel).Msgf("start job %s", notRunningJobs[i].Name)

					if q.StartNextJob(client, c, dc, &notRunningJobs[i]) {
						return
					}
				} else {
					break // no need to process further
				}
			}
		}
	}
}

func (q *Queue) StartNextJob(client *v1alpha1.Client, c discovery.ServerResourcesInterface, dc dynamic.Interface, nextJob *job.Job) bool {
	nextJob.MarkRunning()
	_, err := client.Job(nextJob.Namespace).Update(nextJob, metav1.UpdateOptions{})
	if err != nil {
		log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job should be marked as running but could not mark it %s", q.Name, nextJob.Name)
		log.WithLevel(zerolog.ErrorLevel).Err(err)
		return true
	}

	q.CreateTemplates(nextJob, c, dc)
	return false
}

func (q *Queue) CreateTemplates(jb *job.Job, c discovery.ServerResourcesInterface, dc dynamic.Interface) {
	if jb.Spec.Templates != nil {
		for _, template := range jb.Spec.Templates {
			resource := *template
			log.WithLevel(zerolog.DebugLevel).Msgf("%s", string(resource.Raw))
			err := CreateResource(jb, resource, c, dc)
			if err != nil {
				log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job marked as running but could not create resources %s", q.Name, jb.Name)
				log.WithLevel(zerolog.ErrorLevel).Err(err)
				return
			}
		}
	}
}

func (q *Queue) GetNotRunningJobs(client *v1alpha1.Client, jobs []job.Job, timeout string) ([]job.Job, int64, bool) {
	var notRunningJobs []job.Job
	hasRunning := int64(0)

	for i := range jobs {
		if jobs[i].IsRunning() {
			// could be timed out
			err, isTimedOut := jobs[i].IsTimedOut(timeout)
			if err != nil {
				log.WithLevel(zerolog.ErrorLevel).Msgf("%s - could not check if job was timed out %s", q.Name, q.Jobs[i].Name)
				log.WithLevel(zerolog.ErrorLevel).Err(err)
				return nil, hasRunning, true
			} else if isTimedOut {
				jobs[i].MarkTimedOut()
				log.WithLevel(zerolog.DebugLevel).Msgf("%s - job timed out %s", q.Name, jobs[i].Name)

				_, err := client.Job(jobs[i].Namespace).Update(&jobs[i], metav1.UpdateOptions{})
				if err != nil {
					log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job timed out but could not mark it %s", q.Name, jobs[i].Name)
					log.WithLevel(zerolog.ErrorLevel).Err(err)
					return nil, hasRunning, true
				}

				notRunningJobs = append(notRunningJobs, jobs[i])
			} else {
				log.WithLevel(zerolog.DebugLevel).Msgf("%s - current job running %s", q.Name, jobs[i].Name)
				hasRunning += 1
			}
		} else {
			notRunningJobs = append(notRunningJobs, jobs[i])
		}
	}
	return notRunningJobs, hasRunning, false
}

func (q *Queue) GetToBeDoneJobs(client *v1alpha1.Client, deadlineTimeout string) ([]job.Job, bool) {
	var toBeDoneJobs []job.Job

	for i := range q.Jobs {
		if !q.Jobs[i].IsCompleted() {
			// check if deadlined
			err, isDeadlinedTimeout := q.Jobs[i].IsDeadlinedTimeout(deadlineTimeout)
			if err != nil {
				log.WithLevel(zerolog.ErrorLevel).Msgf("%s - could not check if job was deadlined timed out %s", q.Name, q.Jobs[i].Name)
				log.WithLevel(zerolog.ErrorLevel).Err(err)
				return nil, true
			} else if isDeadlinedTimeout {
				q.Jobs[i].MarkDeadlinedTimeout()
				log.WithLevel(zerolog.DebugLevel).Msgf("%s - job deadline timed out %s", q.Name, q.Jobs[i].Name)

				_, err := client.Job(q.Jobs[i].Namespace).Update(q.Jobs[i], metav1.UpdateOptions{})
				if err != nil {
					log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job deadlined timed out but could not mark it %s", q.Name, q.Jobs[i].Name)
					log.WithLevel(zerolog.ErrorLevel).Err(err)
					return nil, true
				}
			} else {
				toBeDoneJobs = append(toBeDoneJobs, *q.Jobs[i])
			}
		} else {
			// check if we need to delete
			err, isCompletedAndCanBeDeleted := q.Jobs[i].IsCompletedAndCanBeDeleted(q.Settings.TtlAfterSuccesfullCompletion, q.Settings.TtlAfterFailedCompletion)
			if err != nil {
				log.WithLevel(zerolog.ErrorLevel).Msgf("%s - could not check if job was completed and needs to be deleted %s", q.Name, q.Jobs[i].Name)
				log.WithLevel(zerolog.ErrorLevel).Err(err)
				return nil, true
			} else if isCompletedAndCanBeDeleted {
				log.WithLevel(zerolog.DebugLevel).Msgf("%s - job completed and can be deleted %s", q.Name, q.Jobs[i].Name)

				err := client.Job(q.Jobs[i].Namespace).Delete(q.Jobs[i], metav1.DeleteOptions{})
				if err != nil {
					log.WithLevel(zerolog.ErrorLevel).Msgf("%s - job completed and should delete but can not %s", q.Name, q.Jobs[i].Name)
					log.WithLevel(zerolog.ErrorLevel).Err(err)
					return nil, true
				}
			}
		}
	}
	return toBeDoneJobs, false
}

func (q *Queue) Add(addJob *job.Job) {
	q.Jobs = append(q.Jobs, addJob)
}
