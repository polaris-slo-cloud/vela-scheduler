package pipeline

import (
	core "k8s.io/api/core/v1"
	"polaris-slo-cloud.github.io/polaris-scheduler/v2/framework/util"
)

// NodeInfo stores a Node and additional pre-computed scheduling-relevant information about it.
type NodeInfo struct {

	// The Node described by this NodeInfo.
	Node *core.Node `json:"node" yaml:"node"`

	// The accumulated score computed by the Score plugins of the sampling pipeline.
	// This is nil if no sampling score plugins are configured.
	SamplingScore *SamplingScore

	// The name of the cluster that the node is part of.
	ClusterName string `json:"clusterName" yaml:"clusterName"`

	// The resources that are currently available for allocation on the node.
	AllocatableResources *util.Resources `json:"allocatableResources" yaml:"allocatableResources"`

	// The total amount of resources that are available on the node.
	TotalResources *util.Resources `json:"totalResources" yaml:"totalResources"`
}

// NodeScore describes the score of a particular node.
type NodeScore struct {
	Node  *NodeInfo
	Score int64
}

// SamplingScore describes the accumulated score from all sampling score plugins.
type SamplingScore struct {
	// The accumulated score of all sampling score plugins.
	AccumulatedScore int64 `json:"accumulatedScore" yaml:"accumulatedScore"`

	// The number score plugins that contributed to the accumulated score.
	ScorePluginsCount int `json:"scorePluginsCount" yaml:"scorePluginsCount"`
}

// Creates a new NodeInfo object and computes its resources.
func NewNodeInfo(clusterName string, node *core.Node) *NodeInfo {
	return &NodeInfo{
		Node:                 node,
		ClusterName:          clusterName,
		AllocatableResources: util.NewResourcesFromList(node.Status.Allocatable),
		TotalResources:       util.NewResourcesFromList(node.Status.Capacity),
	}
}