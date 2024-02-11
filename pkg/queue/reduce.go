package queue

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/job"
)

type Reducer func(jb job.Job) (bool, error)

func ReduceJobs(jobs []*job.Job, reducer Reducer) ([]*job.Job, error) {
	var resultJobs []*job.Job

	for i := range jobs {
		shouldInclude, err := reducer(*jobs[i])
		if err != nil {
			return []*job.Job{}, err
		} else if shouldInclude {
			resultJobs = append(resultJobs, jobs[i])
		}
	}
	return resultJobs, nil
}
