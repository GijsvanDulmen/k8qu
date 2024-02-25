package markqueuejobcomplete

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func (mqjc *MarkQueueJobComplete) DeepCopyInto(out *MarkQueueJobComplete) {
	*out = *mqjc
	out.TypeMeta = mqjc.TypeMeta
	mqjc.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = mqjc.Spec
}

func (mqjc *MarkQueueJobComplete) DeepCopy() *MarkQueueJobComplete {
	if mqjc == nil {
		return nil
	}
	out := new(MarkQueueJobComplete)
	mqjc.DeepCopyInto(out)
	return out
}

func (mqjc *MarkQueueJobComplete) DeepCopyObject() runtime.Object {
	if c := mqjc.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (qsl *MarkQueueJobCompleteList) DeepCopyInto(out *MarkQueueJobCompleteList) {
	*out = *qsl
	out.TypeMeta = qsl.TypeMeta
	qsl.ListMeta.DeepCopyInto(&out.ListMeta)
	if qsl.Items != nil {
		in, out := &qsl.Items, &out.Items
		*out = make([]MarkQueueJobComplete, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (qsl *MarkQueueJobCompleteList) DeepCopy() *MarkQueueJobCompleteList {
	if qsl == nil {
		return nil
	}
	out := new(MarkQueueJobCompleteList)
	qsl.DeepCopyInto(out)
	return out
}

func (qsl *MarkQueueJobCompleteList) DeepCopyObject() runtime.Object {
	if c := qsl.DeepCopy(); c != nil {
		return c
	}
	return nil
}
