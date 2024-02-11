package queuesettings

import "fmt"

func (qs *QueueSettings) GetQueueName() string {
	return fmt.Sprintf("%s.%s", qs.Namespace, qs.Name)
}
