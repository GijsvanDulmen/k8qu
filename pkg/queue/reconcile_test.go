package queue

import (
	"github.com/stretchr/testify/assert"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func TestToStartJob(t *testing.T) {
	// first job
	firstJob := queuejob.CreateMockJob()
	duration, _ := time.ParseDuration("5s")
	firstJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-duration))
	firstJob.ObjectMeta.Name = "first"

	// second job
	secondJob := queuejob.CreateMockJob()
	secondDuration, _ := time.ParseDuration("10s")
	secondJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-secondDuration))
	secondJob.ObjectMeta.Name = "second"

	// wrong order
	jobs := []*queuejob.QueueJob{&firstJob, &secondJob}

	toStartJob := GetToStartJob(jobs, 1, 0)

	assert.Equal(t, len(toStartJob), 1)
	assert.Equal(t, toStartJob[0].Name, "second")
}

func TestToStartJobWithMultipleParallism(t *testing.T) {
	// first job
	firstJob := queuejob.CreateMockJob()
	duration, _ := time.ParseDuration("5s")
	firstJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-duration))
	firstJob.ObjectMeta.Name = "first"

	// second job
	secondJob := queuejob.CreateMockJob()
	secondDuration, _ := time.ParseDuration("10s")
	secondJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-secondDuration))
	secondJob.ObjectMeta.Name = "second"

	// wrong order
	jobs := []*queuejob.QueueJob{&firstJob, &secondJob}

	toStartJob := GetToStartJob(jobs, 2, 0)

	assert.Equal(t, len(toStartJob), 2)
	assert.Equal(t, toStartJob[0].Name, "second")
	assert.Equal(t, toStartJob[1].Name, "first")
}

func TestToStartJobWithMultipleParallismWithRunning(t *testing.T) {
	// first job
	firstJob := queuejob.CreateMockJob()
	duration, _ := time.ParseDuration("5s")
	firstJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-duration))
	firstJob.ObjectMeta.Name = "first"

	// second job
	secondJob := queuejob.CreateMockJob()
	secondDuration, _ := time.ParseDuration("10s")
	secondJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-secondDuration))
	secondJob.ObjectMeta.Name = "second"

	// wrong order
	jobs := []*queuejob.QueueJob{&firstJob, &secondJob}

	toStartJob := GetToStartJob(jobs, 2, 1)

	assert.Equal(t, len(toStartJob), 1)
	assert.Equal(t, toStartJob[0].Name, "second")
}

func TestSortQueueJobs(t *testing.T) {
	// first job
	firstJob := queuejob.CreateMockJob()
	duration, _ := time.ParseDuration("5s")
	firstJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-duration))
	firstJob.ObjectMeta.Name = "first"

	// second job
	secondJob := queuejob.CreateMockJob()
	secondDuration, _ := time.ParseDuration("10s")
	secondJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-secondDuration))
	secondJob.ObjectMeta.Name = "second"

	// wrong order
	jobs := []*queuejob.QueueJob{&firstJob, &secondJob}

	SortQueueJobs(jobs)

	assert.Equal(t, len(jobs), 2)
	assert.Equal(t, jobs[0].Name, "second")
	assert.Equal(t, jobs[1].Name, "first")
}

func TestSortQueueJobsWithSame(t *testing.T) {
	// first job
	duration, _ := time.ParseDuration("5s")

	firstJob := queuejob.CreateMockJob()
	firstJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-duration))
	firstJob.ObjectMeta.Name = "job-b"

	// second job
	secondJob := queuejob.CreateMockJob()
	secondJob.ObjectMeta.CreationTimestamp = metav1.NewTime(time.Now().Add(-duration))
	secondJob.ObjectMeta.Name = "job-a"

	// wrong order
	jobs := []*queuejob.QueueJob{&firstJob, &secondJob}

	SortQueueJobs(jobs)

	assert.Equal(t, len(jobs), 2)
	assert.Equal(t, jobs[0].Name, "job-a")
	assert.Equal(t, jobs[1].Name, "job-b")
}
