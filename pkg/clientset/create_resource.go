package clientset

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"k8qu/pkg/apis/k8qu/v1alpha1/queuejob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"strings"
)

func addLabels(us *unstructured.Unstructured, labelsToAdd map[string]string) (*unstructured.Unstructured, error) {
	labels, _, err := unstructured.NestedStringMap(us.Object, "metadata", "labels")
	if err != nil {
		return nil, err
	}

	if labels == nil {
		labels = make(map[string]string)
	}
	for k, v := range labelsToAdd {
		l := fmt.Sprintf("%s/%s", "k8qu", strings.TrimLeft(k, "/"))
		labels[l] = v
	}

	us.SetLabels(labels)
	return us, nil
}

func CreateResource(jb *queuejob.QueueJob, raw runtime.RawExtension, c discovery.ServerResourcesInterface, dc dynamic.Interface) error {

	// replacements
	replaced := strings.ReplaceAll(string(raw.Raw), "[[JOBNAME]]", jb.ObjectMeta.Name)
	raw.Raw = []byte(replaced)

	data := new(unstructured.Unstructured)

	if err := data.UnmarshalJSON(raw.Raw); err != nil {
		return fmt.Errorf("couldn't unmarshal json from the job template: %v", err)
	}

	data, err := addLabels(data, map[string]string{
		"job": jb.ObjectMeta.Name,
	})

	if err != nil {
		return err
	}

	apiResource, err := findAPIResource(data.GetAPIVersion(), data.GetKind(), c)
	if err != nil {
		return fmt.Errorf("couldn't find API resource for template: %v", err)
	}

	name := data.GetName()
	if name == "" {
		name = data.GetGenerateName()
	}

	log.WithLevel(zerolog.InfoLevel).Msgf("generating resource: kind: %s, name: %s", apiResource, name)

	gvr := schema.GroupVersionResource{
		Group:    apiResource.Group,
		Version:  apiResource.Version,
		Resource: apiResource.Name,
	}

	log.WithLevel(zerolog.InfoLevel).Msgf("for job %q creating resource %v", jb.Name, gvr)

	if _, err := dc.Resource(gvr).Namespace(jb.Namespace).Create(context.Background(), data, metav1.CreateOptions{}); err != nil {
		return fmt.Errorf("couldn't create resource with group version kind %q: %v", gvr, err)
	}
	return nil
}

func findAPIResource(apiVersion, kind string, c discovery.ServerResourcesInterface) (*metav1.APIResource, error) {
	resourceList, err := c.ServerResourcesForGroupVersion(apiVersion)
	if err != nil {
		return nil, fmt.Errorf("error getting kubernetes server resources for apiVersion %s: %s", apiVersion, err)
	}
	for i := range resourceList.APIResources {
		r := &resourceList.APIResources[i]
		if r.Kind != kind {
			continue
		}

		if r.Version == "" || r.Group == "" {
			gv, err := schema.ParseGroupVersion(resourceList.GroupVersion)
			if err != nil {
				return nil, fmt.Errorf("error parsing parsing GroupVersion: %v", err)
			}
			r.Group = gv.Group
			r.Version = gv.Version
		}
		return r, nil
	}
	return nil, fmt.Errorf("error could not find resource with apiVersion %s and kind %s", apiVersion, kind)
}
