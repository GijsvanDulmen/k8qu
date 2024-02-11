package queuesettings

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func (qs *QueueSettings) DeepCopyInto(out *QueueSettings) {
	*out = *qs
	out.TypeMeta = qs.TypeMeta
	qs.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = qs.Spec
}

func (qs *QueueSettings) DeepCopy() *QueueSettings {
	if qs == nil {
		return nil
	}
	out := new(QueueSettings)
	qs.DeepCopyInto(out)
	return out
}

func (qs *QueueSettings) DeepCopyObject() runtime.Object {
	if c := qs.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (qsl *QueueSettingsList) DeepCopyInto(out *QueueSettingsList) {
	*out = *qsl
	out.TypeMeta = qsl.TypeMeta
	qsl.ListMeta.DeepCopyInto(&out.ListMeta)
	if qsl.Items != nil {
		in, out := &qsl.Items, &out.Items
		*out = make([]QueueSettings, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (qsl *QueueSettingsList) DeepCopy() *QueueSettingsList {
	if qsl == nil {
		return nil
	}
	out := new(QueueSettingsList)
	qsl.DeepCopyInto(out)
	return out
}

func (qsl *QueueSettingsList) DeepCopyObject() runtime.Object {
	if c := qsl.DeepCopy(); c != nil {
		return c
	}
	return nil
}
