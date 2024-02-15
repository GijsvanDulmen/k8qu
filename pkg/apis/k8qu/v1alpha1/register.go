package v1alpha1

import (
	k8qucontroller "k8qu/pkg/apis/k8qu"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuesettings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: k8qucontroller.GroupName, Version: "v1alpha1"}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&queuejob.QueueJob{},
		&queuejob.QueueJobList{},
		&queuesettings.QueueSettings{},
		&queuesettings.QueueSettingsList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
