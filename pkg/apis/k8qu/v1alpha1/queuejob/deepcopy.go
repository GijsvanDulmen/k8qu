package queuejob

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func (jb *QueueJob) DeepCopyInto(out *QueueJob) {
	*out = *jb
	out.TypeMeta = jb.TypeMeta
	jb.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = jb.Spec
	out.Status = jb.Status
}

func (jb *QueueJob) DeepCopy() *QueueJob {
	if jb == nil {
		return nil
	}
	out := new(QueueJob)
	jb.DeepCopyInto(out)
	return out
}

func (jb *QueueJob) DeepCopyObject() runtime.Object {
	if c := jb.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (jbl *QueueJobList) DeepCopyInto(out *QueueJobList) {
	*out = *jbl
	out.TypeMeta = jbl.TypeMeta
	jbl.ListMeta.DeepCopyInto(&out.ListMeta)
	if jbl.Items != nil {
		in, out := &jbl.Items, &out.Items
		*out = make([]QueueJob, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (jbl *QueueJobList) DeepCopy() *QueueJobList {
	if jbl == nil {
		return nil
	}
	out := new(QueueJobList)
	jbl.DeepCopyInto(out)
	return out
}

func (jbl *QueueJobList) DeepCopyObject() runtime.Object {
	if c := jbl.DeepCopy(); c != nil {
		return c
	}
	return nil
}
