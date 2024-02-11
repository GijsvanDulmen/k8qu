package job

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func (jb *Job) DeepCopyInto(out *Job) {
	*out = *jb
	out.TypeMeta = jb.TypeMeta
	jb.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = jb.Spec
	out.Status = jb.Status
}

func (jb *Job) DeepCopy() *Job {
	if jb == nil {
		return nil
	}
	out := new(Job)
	jb.DeepCopyInto(out)
	return out
}

func (jb *Job) DeepCopyObject() runtime.Object {
	if c := jb.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (jbl *JobList) DeepCopyInto(out *JobList) {
	*out = *jbl
	out.TypeMeta = jbl.TypeMeta
	jbl.ListMeta.DeepCopyInto(&out.ListMeta)
	if jbl.Items != nil {
		in, out := &jbl.Items, &out.Items
		*out = make([]Job, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (jbl *JobList) DeepCopy() *JobList {
	if jbl == nil {
		return nil
	}
	out := new(JobList)
	jbl.DeepCopyInto(out)
	return out
}

func (jbl *JobList) DeepCopyObject() runtime.Object {
	if c := jbl.DeepCopy(); c != nil {
		return c
	}
	return nil
}
