package kubernetes

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"

	"polaris-slo-cloud.github.io/polaris-scheduler/v2/framework/client"
	"polaris-slo-cloud.github.io/polaris-scheduler/v2/framework/util"
)

const (
	// The name of the scheduler, to which pods are assigned.
	// This is needed to avoid interference from the default kube-scheduler.
	polarisClusterAgentSchedulerName = "polaris-cluster-agent"
)

var (
	_ KubernetesClusterClient = (*KubernetesClusterClientImpl)(nil)
)

// Default KubernetesClusterClient implementation.
type KubernetesClusterClientImpl struct {
	clusterName   string
	k8sClientSet  *clientset.Clientset
	eventRecorder record.EventRecorder
	logger        *logr.Logger
}

// Creates a new KubernetesClusterClientImpl using the specified kubeconfig.
//
// - clusterName is the name of the cluster to connect to
// - kubeconfig is the respective kubeconfig
// - parentComponentName is the name of the component that creates this client (this is used as the source name in the event recorder)
// - logger the Logger that should be used for logging
func NewKubernetesClusterClientImpl(
	clusterName string,
	kubeconfig *rest.Config,
	parentComponentName string,
	logger *logr.Logger,
) (*KubernetesClusterClientImpl, error) {
	k8sClientSet, err := clientset.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	eventSink := coreclient.EventSinkImpl{Interface: k8sClientSet.CoreV1().Events("")}
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(
		func(format string, args ...interface{}) {
			msg := fmt.Sprintf(format, args)
			logger.Info(msg)
		},
	)
	eventBroadcaster.StartRecordingToSink(&eventSink)
	eventRecorder := eventBroadcaster.NewRecorder(scheme.Scheme, core.EventSource{Component: parentComponentName})

	clusterClient := KubernetesClusterClientImpl{
		clusterName:   clusterName,
		k8sClientSet:  k8sClientSet,
		eventRecorder: eventRecorder,
		logger:        logger,
	}

	return &clusterClient, nil
}

func (c *KubernetesClusterClientImpl) ClusterName() string {
	return c.clusterName
}

func (c *KubernetesClusterClientImpl) ClientSet() clientset.Interface {
	return c.k8sClientSet
}

func (c *KubernetesClusterClientImpl) EventRecorder() record.EventRecorder {
	return c.eventRecorder
}

func (c *KubernetesClusterClientImpl) CommitSchedulingDecision(ctx context.Context, schedulingDecision *client.ClusterSchedulingDecision) (*client.CommitSchedulingDecisionSuccess, error) {
	// ToDo: check if pod already exists, before trying to create it.
	createPodStopwatch := util.NewStopwatch()
	createPodStopwatch.Start()
	pod, err := c.createPod(ctx, schedulingDecision)
	createPodStopwatch.Stop()
	if err != nil {
		return nil, err
	}

	binding := &core.Binding{
		ObjectMeta: meta.ObjectMeta{
			Namespace: pod.Namespace,
			Name:      pod.Name,
			UID:       pod.UID,
		},
		Target: core.ObjectReference{
			Kind: "Node",
			Name: schedulingDecision.NodeName,
		},
	}

	createBindingStopwatch := util.NewStopwatch()
	createBindingStopwatch.Start()
	err = c.k8sClientSet.CoreV1().Pods(pod.Namespace).Bind(ctx, binding, meta.CreateOptions{})
	createBindingStopwatch.Stop()
	if err != nil {
		c.logger.Error(err, "could not bind Pod", "pod", pod, "binding", binding)
		c.eventRecorder.Eventf(pod, core.EventTypeWarning, "FailedScheduling", "Could not bind pod to node %s", &binding.Target.Name)
		return nil, err
	}

	result := &client.CommitSchedulingDecisionSuccess{
		Namespace: pod.Namespace,
		PodName:   pod.Name,
		NodeName:  schedulingDecision.NodeName,
		Timings: &client.CommitSchedulingDecisionTimings{
			CreatePod:     createPodStopwatch.Duration().Milliseconds(),
			CreateBinding: createBindingStopwatch.Duration().Milliseconds(),
		},
	}

	fullyQualifiedPodName := pod.Namespace + "." + pod.Name
	c.logger.Info(
		"PodBindingSuccess",
		"pod", fullyQualifiedPodName,
		"cluster", c.clusterName,
		"node", binding.Target.Name,
		"createPodDurationMs", result.Timings.CreatePod,
		"createBindingDurationMs", result.Timings.CreateBinding,
	)
	return result, nil
}

func (c *KubernetesClusterClientImpl) FetchNode(ctx context.Context, name string) (*core.Node, error) {
	return c.k8sClientSet.CoreV1().Nodes().Get(ctx, name, meta.GetOptions{})
}

func (c *KubernetesClusterClientImpl) FetchPodsScheduledOnNode(ctx context.Context, nodeName string) ([]core.Pod, error) {
	// Based on how kubectl does this: https://github.com/kubernetes/kubernetes/blob/85f1205efa89008e6e94a92d66a9dd691cb59825/pkg/kubectl/describe/versioned/describe.go#L2819
	fieldSelector, err := fields.ParseSelector("spec.nodeName=" + nodeName + ",status.phase!=" + string(core.PodSucceeded) + ",status.phase!=" + string(core.PodFailed))
	if err != nil {
		return nil, err
	}

	podsList, err := c.k8sClientSet.CoreV1().Pods("").List(ctx, meta.ListOptions{FieldSelector: fieldSelector.String()})
	if err != nil {
		return nil, err
	}
	return podsList.Items, nil
}

func (c *KubernetesClusterClientImpl) createPod(ctx context.Context, schedulingDecision *client.ClusterSchedulingDecision) (*core.Pod, error) {
	pod := schedulingDecision.Pod.DeepCopy()
	pod.Spec.SchedulerName = polarisClusterAgentSchedulerName
	return c.k8sClientSet.CoreV1().Pods(pod.Namespace).Create(ctx, pod, meta.CreateOptions{})
}
