package informers

import (
	"k8qu/pkg/apis/k8qu/v1alpha1/markqueuejobcomplete"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuesettings"
)

func LogForJob(j queuejob.QueueJob, s string) {
	log.Debug().Msgf("job %s/%s log: %s", j.Namespace, j.Name, s)
}

func LogForMarkQueueJobComplete(j markqueuejobcomplete.MarkQueueJobComplete, s string) {
	log.Debug().Msgf("mqjc %s/%s log: %s", j.Namespace, j.Name, s)
}

func LogForQueueSetttings(qs queuesettings.QueueSettings, s string) {
	log.Debug().Msgf("queue setting %s/%s log: %s", qs.Namespace, qs.Name, s)
}
