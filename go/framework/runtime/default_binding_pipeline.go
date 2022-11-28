package runtime

import (
	"sync"

	"github.com/go-logr/logr"
	core "k8s.io/api/core/v1"

	"polaris-slo-cloud.github.io/polaris-scheduler/v2/framework/client"
	"polaris-slo-cloud.github.io/polaris-scheduler/v2/framework/collections"
	"polaris-slo-cloud.github.io/polaris-scheduler/v2/framework/pipeline"
	"polaris-slo-cloud.github.io/polaris-scheduler/v2/framework/util"
)

var (
	_ pipeline.BindingPipeline = (*DefaultBindingPipeline)(nil)
)

const (
	BindingPipelineStopwatchesStateKey = "polaris-internal.stopwatches.binding"
)

// Default implementation of the BindingPipeline.
type DefaultBindingPipeline struct {
	id                   int
	plugins              *pipeline.BindingPipelinePlugins
	clusterAgentServices pipeline.ClusterAgentServices
	nodesLocker          collections.EntityLocker
	logger               *logr.Logger
}

// A collection of stopwatches for timing various parts of the binding pipeline.
//
// ToDo: find a better way to handle the responsibilities for timing other than this shared struct.
type BindingPipelineStopwatches struct {
	QueueTime       *util.Stopwatch
	NodeLockTime    *util.Stopwatch
	FetchNodeInfo   *util.Stopwatch
	BindingPipeline *util.Stopwatch
	CommitDecision  *util.Stopwatch
}

// Creates a new instance of the default implementation of the BindingPipeline.
func NewDefaultBindingPipeline(
	id int,
	plugins *pipeline.BindingPipelinePlugins,
	clusterAgentServices pipeline.ClusterAgentServices,
	nodesLocker collections.EntityLocker,
) *DefaultBindingPipeline {
	bp := &DefaultBindingPipeline{
		id:                   id,
		plugins:              plugins,
		clusterAgentServices: clusterAgentServices,
		nodesLocker:          nodesLocker,
		logger:               clusterAgentServices.Logger(),
	}
	return bp
}

func NewBindingPipelineStopwatches() *BindingPipelineStopwatches {
	stopwatches := &BindingPipelineStopwatches{
		QueueTime:       util.NewStopwatch(),
		FetchNodeInfo:   util.NewStopwatch(),
		NodeLockTime:    util.NewStopwatch(),
		BindingPipeline: util.NewStopwatch(),
		CommitDecision:  util.NewStopwatch(),
	}
	return stopwatches
}

func (bp *DefaultBindingPipeline) CommitSchedulingDecision(schedCtx pipeline.SchedulingContext, schedDecision *client.ClusterSchedulingDecision) pipeline.Status {
	stopwatches := bp.getStopwatches(schedCtx)

	stopwatches.NodeLockTime.Start()
	nodeLock := bp.nodesLocker.Lock(schedDecision.NodeName)
	defer nodeLock.Unlock()
	stopwatches.NodeLockTime.Stop()

	// Fetch the NodeInfo.
	stopwatches.FetchNodeInfo.Start()
	updatedNodeInfo, err := bp.fetchNodeInfo(schedCtx, schedDecision.NodeName)
	stopwatches.FetchNodeInfo.Stop()
	if err != nil {
		return pipeline.NewStatus(pipeline.Unschedulable, "error fetching node information", err.Error())
	}
	decision := &pipeline.SchedulingDecision{
		Pod:        &pipeline.PodInfo{Pod: schedDecision.Pod},
		TargetNode: updatedNodeInfo,
	}

	// Run the pipeline.
	stopwatches.BindingPipeline.Start()
	status := bp.runCheckConflictsPlugins(schedCtx, decision)
	stopwatches.BindingPipeline.Stop()
	if !pipeline.IsSuccessStatus(status) {
		return status
	}

	stopwatches.CommitDecision.Start()
	status = bp.commitSchedulingDecision(schedCtx, decision)
	stopwatches.CommitDecision.Stop()
	return status
}

func (bp *DefaultBindingPipeline) getStopwatches(schedCtx pipeline.SchedulingContext) *BindingPipelineStopwatches {
	stopwatches, ok, err := pipeline.ReadTypedStateData[*BindingPipelineStopwatches](schedCtx, BindingPipelineStopwatchesStateKey)
	if !ok || err != nil {
		panic("could not read BindingPipelineStopwatches from SchedulingContext")
	}
	return stopwatches
}

func (bp *DefaultBindingPipeline) fetchNodeInfo(schedCtx pipeline.SchedulingContext, nodeName string) (*pipeline.NodeInfo, error) {
	clusterClient := bp.clusterAgentServices.ClusterClient()
	var node *core.Node
	var podsOnNode []core.Pod
	var nodeFetchErr, podsFetchErr error

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		node, nodeFetchErr = clusterClient.FetchNode(schedCtx.Context(), nodeName)
		wg.Done()
	}()

	go func() {
		podsOnNode, podsFetchErr = clusterClient.FetchPodsScheduledOnNode(schedCtx.Context(), nodeName)
		wg.Done()
	}()

	wg.Wait()
	if nodeFetchErr != nil {
		return nil, nodeFetchErr
	}
	if podsFetchErr != nil {
		return nil, podsFetchErr
	}

	clusterNode := &client.ClusterNode{
		Node:               node,
		AvailableResources: util.CalculateNodeAvailableResources(node, podsOnNode),
		TotalResources:     util.NewResourcesFromList(node.Status.Capacity),
	}

	return pipeline.NewNodeInfo(clusterClient.ClusterName(), clusterNode), nil
}

func (bp *DefaultBindingPipeline) runCheckConflictsPlugins(schedCtx pipeline.SchedulingContext, decision *pipeline.SchedulingDecision) pipeline.Status {
	var status pipeline.Status

	for _, plugin := range bp.plugins.CheckConflicts {
		status = plugin.CheckForConflicts(schedCtx, decision)
		if !pipeline.IsSuccessStatus(status) {
			status.SetFailedPlugin(plugin, pipeline.CheckConflictsStage)
			return status
		}
	}

	return status
}

func (bp *DefaultBindingPipeline) commitSchedulingDecision(schedCtx pipeline.SchedulingContext, decision *pipeline.SchedulingDecision) pipeline.Status {
	clusterSchedDecision := &client.ClusterSchedulingDecision{
		Pod:      decision.Pod.Pod,
		NodeName: decision.TargetNode.Node.Name,
	}

	if err := bp.clusterAgentServices.ClusterClient().CommitSchedulingDecision(schedCtx.Context(), clusterSchedDecision); err != nil {
		return pipeline.NewStatus(pipeline.Unschedulable, "error committing scheduling decision", err.Error())
	}
	return pipeline.NewSuccessStatus()
}
