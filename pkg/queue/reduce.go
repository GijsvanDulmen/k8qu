package queue

import "k8qu/pkg/apis/k8qu/v1alpha1/queuejob"

type Reducer func(jb queuejob.QueueJob) (bool, error)

func ReduceJobs(jobs []*queuejob.QueueJob, reducer Reducer) ([]*queuejob.QueueJob, error) {
	var resultJobs []*queuejob.QueueJob

	for i := range jobs {
		shouldInclude, err := reducer(*jobs[i])
		if err != nil {
			return []*queuejob.QueueJob{}, err
		} else if shouldInclude {
			resultJobs = append(resultJobs, jobs[i])
		}
	}
	return resultJobs, nil
}
