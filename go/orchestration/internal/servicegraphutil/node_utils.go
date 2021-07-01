package servicegraphutil

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	fogappsCRDs "k8s.rainbow-h2020.eu/rainbow/orchestration/apis/fogapps/v1"
	"k8s.rainbow-h2020.eu/rainbow/orchestration/pkg/kubeutil"
)

const (
	kubernetesCpuArchLabel = "kubernetes.io/arch"
)

// CreatePodSpec creates a PodSpec from the specified node.
func CreatePodTemplate(node *fogappsCRDs.ServiceGraphNode, graph *fogappsCRDs.ServiceGraph) (*core.PodTemplateSpec, error) {
	podTemplate := core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{},
		Spec:       core.PodSpec{},
	}

	return &podTemplate, nil
}

// CreateDeployment creates a new Deployment from the specified node.
func CreateDeployment(node *fogappsCRDs.ServiceGraphNode, graph *fogappsCRDs.ServiceGraph) (*apps.Deployment, error) {
	deployment := apps.Deployment{
		ObjectMeta: *createNodeObjectMeta(node, graph),
		Spec:       apps.DeploymentSpec{},
	}

	return UpdateDeployment(&deployment, node, graph)
}

// UpdateDeployment updates an existing Deployment, based on the specified node.
func UpdateDeployment(deployment *apps.Deployment, node *fogappsCRDs.ServiceGraphNode, graph *fogappsCRDs.ServiceGraph) (*apps.Deployment, error) {
	replicas := getInitialReplicas(node)

	updateNodeObjectMeta(&deployment.ObjectMeta, node, graph)
	updatePodTemplate(&deployment.Spec.Template, node, graph)
	deployment.Spec.Selector = createLabelSelector(node, graph)
	deployment.Spec.Replicas = &replicas

	return deployment, nil
}

// CreateStatefulSet creates a StatefulSet from the specified node.
func CreateStatefulSet(node *fogappsCRDs.ServiceGraphNode, graph *fogappsCRDs.ServiceGraph) (*apps.StatefulSet, error) {
	statefulSet := apps.StatefulSet{
		ObjectMeta: *createNodeObjectMeta(node, graph),
		Spec:       apps.StatefulSetSpec{},
	}

	return UpdateStatefulSet(&statefulSet, node, graph)
}

// UpdateStatefulSet updates an existing StatefulSet, based on the specified node.
func UpdateStatefulSet(statefulSet *apps.StatefulSet, node *fogappsCRDs.ServiceGraphNode, graph *fogappsCRDs.ServiceGraph) (*apps.StatefulSet, error) {
	replicas := getInitialReplicas(node)

	updateNodeObjectMeta(&statefulSet.ObjectMeta, node, graph)
	updatePodTemplate(&statefulSet.Spec.Template, node, graph)
	statefulSet.Spec.Selector = createLabelSelector(node, graph)
	statefulSet.Spec.Replicas = &replicas

	return statefulSet, nil
}

func updatePodTemplate(podTemplate *core.PodTemplateSpec, node *fogappsCRDs.ServiceGraphNode, graph *fogappsCRDs.ServiceGraph) {
	podTemplate.Spec.SchedulerName = kubeutil.RainbowSchedulerName
	podTemplate.ObjectMeta.Labels = getPodLabels(node, graph)
	podTemplate.Spec.InitContainers = node.InitContainers
	podTemplate.Spec.Containers = node.Containers
	podTemplate.Spec.Volumes = node.Volumes
	podTemplate.Spec.Affinity = node.Affinity

	if node.ImagePullSecrets != nil && len(node.ImagePullSecrets) > 0 {
		podTemplate.Spec.ImagePullSecrets = node.ImagePullSecrets
	}

	if node.NodeHardware != nil {
		addNodeHardwareRequirements(podTemplate, node.NodeHardware)
	}

	if serviceAccountName := getServiceAccountName(node, graph); serviceAccountName != nil {
		podTemplate.Spec.ServiceAccountName = *serviceAccountName
	} else {
		podTemplate.Spec.ServiceAccountName = ""
	}
}

func createLabelSelector(node *fogappsCRDs.ServiceGraphNode, graph *fogappsCRDs.ServiceGraph) *meta.LabelSelector {
	return &meta.LabelSelector{
		MatchLabels: getPodLabels(node, graph),
	}
}

func getInitialReplicas(node *fogappsCRDs.ServiceGraphNode) int32 {
	if node.Replicas.InitialCount != nil {
		return *node.Replicas.InitialCount
	}
	return node.Replicas.Min
}

func getServiceAccountName(node *fogappsCRDs.ServiceGraphNode, graph *fogappsCRDs.ServiceGraph) *string {
	if node.ServiceAccountName != nil {
		return node.ServiceAccountName
	}
	return graph.Spec.ServiceAccountName
}

// addNodeHardwareRequirements adds/configures the affinity of the podTemplate to ensure that only nodes
// that match the nodeHardwareReq are eligible.
func addNodeHardwareRequirements(podTemplate *core.PodTemplateSpec, nodeHardwareReq *fogappsCRDs.NodeHardware) {
	if nodeHardwareReq.NodeType == nil && nodeHardwareReq.CpuInfo == nil && nodeHardwareReq.GpuInfo == nil {
		return
	}

	nodeSelector := ensureAffinityNodeSelectorExists(podTemplate)

	if nodeHardwareReq.CpuInfo != nil {
		addCpuSelectionTerms(nodeSelector, nodeHardwareReq.CpuInfo)
	}

	if len(nodeSelector.NodeSelectorTerms) == 0 {
		nodeSelector.NodeSelectorTerms = nil
	}

	// ToDo: Add support for other NodeHardware fields!
}

// ensureAffinityNodeSelectorExists creates the required during scheduling NodeSelector for the pod's node affinity,
// if it does not exist, and returns the NodeSelector.
func ensureAffinityNodeSelectorExists(podTemplate *core.PodTemplateSpec) *core.NodeSelector {
	if podTemplate.Spec.Affinity == nil {
		podTemplate.Spec.Affinity = &core.Affinity{}
	}
	if podTemplate.Spec.Affinity.NodeAffinity == nil {
		podTemplate.Spec.Affinity.NodeAffinity = &core.NodeAffinity{}
	}
	if podTemplate.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		podTemplate.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = &core.NodeSelector{}
	}
	nodeSelector := podTemplate.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution
	if nodeSelector.NodeSelectorTerms == nil {
		nodeSelector.NodeSelectorTerms = make([]core.NodeSelectorTerm, 0)
	}
	return nodeSelector
}

func addCpuSelectionTerms(nodeSelector *core.NodeSelector, cpuInfo *fogappsCRDs.CpuInfo) {
	architecturesCount := len(cpuInfo.Architectures)
	if architecturesCount > 0 {
		architectures := make([]string, architecturesCount)
		for i, cpuArch := range cpuInfo.Architectures {
			architectures[i] = string(cpuArch)
		}

		cpuArchTerm := core.NodeSelectorTerm{
			MatchExpressions: []core.NodeSelectorRequirement{
				{
					Key:      kubernetesCpuArchLabel,
					Operator: core.NodeSelectorOpIn,
					Values:   architectures,
				},
			},
		}
		nodeSelector.NodeSelectorTerms = append(nodeSelector.NodeSelectorTerms, cpuArchTerm)
	}
}
