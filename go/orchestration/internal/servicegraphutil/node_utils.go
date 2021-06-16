package servicegraphutil

import (
	"fmt"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	fogapps "k8s.rainbow-h2020.eu/rainbow/orchestration/apis/fogapps/v1"
	"k8s.rainbow-h2020.eu/rainbow/orchestration/internal/util"
)

const (
	RainbowGeneratedPodLabelName = "rainbow-generated-pod-label"
)

// CreatePodSpec creates a PodSpec from the specified node.
func CreatePodTemplate(node *fogapps.ServiceGraphNode, graph *fogapps.ServiceGraph) (*core.PodTemplateSpec, error) {
	podTemplate := core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{
			Labels: getPodLabels(node, graph),
		},
		Spec: core.PodSpec{
			InitContainers: node.InitContainers,
			Containers:     node.Containers,
			Volumes:        node.Volumes,
		},
	}

	if serviceAccountName := getServiceAccountName(node, graph); serviceAccountName != nil {
		podTemplate.Spec.ServiceAccountName = *serviceAccountName
	}

	return &podTemplate, nil
}

// CreateDeployment creates a Deployment form the specified node.
func CreateDeployment(node *fogapps.ServiceGraphNode, graph *fogapps.ServiceGraph) (*apps.Deployment, error) {
	podTemplate, err := CreatePodTemplate(node, graph)
	if err != nil {
		return nil, err
	}
	replicas := getInitialReplicas(node)

	deployment := apps.Deployment{
		ObjectMeta: *createObjectMeta(node, graph),
		Spec: apps.DeploymentSpec{
			Selector: createLabelSelector(node, graph),
			Template: *podTemplate,
			Replicas: &replicas,
		},
	}

	return &deployment, nil
}

// CreateStatefulSet creates a StatefulSet form the specified node.
func CreateStatefulSet(node *fogapps.ServiceGraphNode, graph *fogapps.ServiceGraph) (*apps.StatefulSet, error) {
	podTemplate, err := CreatePodTemplate(node, graph)
	if err != nil {
		return nil, err
	}
	replicas := getInitialReplicas(node)

	statefulSet := apps.StatefulSet{
		ObjectMeta: *createObjectMeta(node, graph),
		Spec: apps.StatefulSetSpec{
			Selector: createLabelSelector(node, graph),
			Template: *podTemplate,
			Replicas: &replicas,
		},
	}

	return &statefulSet, nil
}

func createObjectMeta(node *fogapps.ServiceGraphNode, graph *fogapps.ServiceGraph) *meta.ObjectMeta {
	return &meta.ObjectMeta{
		Name:        node.Name,
		Namespace:   graph.Namespace,
		Labels:      getPodLabels(node, graph),
		Annotations: make(map[string]string),
	}
}

func createLabelSelector(node *fogapps.ServiceGraphNode, graph *fogapps.ServiceGraph) *meta.LabelSelector {
	return &meta.LabelSelector{
		MatchLabels: getPodLabels(node, graph),
	}
}

func getInitialReplicas(node *fogapps.ServiceGraphNode) int32 {
	if node.Replicas.InitialCount != nil {
		return *node.Replicas.InitialCount
	}
	return node.Replicas.Min
}

func getServiceAccountName(node *fogapps.ServiceGraphNode, graph *fogapps.ServiceGraph) *string {
	if node.ServiceAccountName != nil {
		return node.ServiceAccountName
	}
	return graph.Spec.ServiceAccountName
}

func getPodLabels(node *fogapps.ServiceGraphNode, graph *fogapps.ServiceGraph) map[string]string {
	labels := util.DeepCopyStringMap(node.PodLabels)
	if len(labels) == 0 {
		labels[RainbowGeneratedPodLabelName] = fmt.Sprintf("%s.%s.generated", graph.Name, node.Name)
	}
	return labels
}
