package markqueuejobcomplete

import (
	"fmt"
)

func (mqjc *MarkQueueJobComplete) GetQueueName() string {
	return fmt.Sprintf("%s.%s", mqjc.Namespace, mqjc.Spec.Queue)
}
