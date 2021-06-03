// +build !ignore_autogenerated

/*
Copyright 2021 Rainbow Project.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	corev1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	clusterv1 "k8s.rainbow-h2020.eu/rainbow/orchestration/apis/cluster/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiVersionKind) DeepCopyInto(out *ApiVersionKind) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiVersionKind.
func (in *ApiVersionKind) DeepCopy() *ApiVersionKind {
	if in == nil {
		return nil
	}
	out := new(ApiVersionKind)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArbitraryObject) DeepCopyInto(out *ArbitraryObject) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArbitraryObject.
func (in *ArbitraryObject) DeepCopy() *ArbitraryObject {
	if in == nil {
		return nil
	}
	out := new(ArbitraryObject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigMap) DeepCopyInto(out *ConfigMap) {
	*out = *in
	if in.Data != nil {
		in, out := &in.Data, &out.Data
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.BinaryData != nil {
		in, out := &in.BinaryData, &out.BinaryData
		*out = make(map[string][]byte, len(*in))
		for key, val := range *in {
			var outVal []byte
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]byte, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigMap.
func (in *ConfigMap) DeepCopy() *ConfigMap {
	if in == nil {
		return nil
	}
	out := new(ConfigMap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExposedPorts) DeepCopyInto(out *ExposedPorts) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]corev1.ServicePort, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExposedPorts.
func (in *ExposedPorts) DeepCopy() *ExposedPorts {
	if in == nil {
		return nil
	}
	out := new(ExposedPorts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeoLocation) DeepCopyInto(out *GeoLocation) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeoLocation.
func (in *GeoLocation) DeepCopy() *GeoLocation {
	if in == nil {
		return nil
	}
	out := new(GeoLocation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LinkQosRequirements) DeepCopyInto(out *LinkQosRequirements) {
	*out = *in
	in.LinkType.DeepCopyInto(&out.LinkType)
	in.Throughput.DeepCopyInto(&out.Throughput)
	in.Latency.DeepCopyInto(&out.Latency)
	out.PacketLoss = in.PacketLoss
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LinkQosRequirements.
func (in *LinkQosRequirements) DeepCopy() *LinkQosRequirements {
	if in == nil {
		return nil
	}
	out := new(LinkQosRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LinkTrustRequirements) DeepCopyInto(out *LinkTrustRequirements) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LinkTrustRequirements.
func (in *LinkTrustRequirements) DeepCopy() *LinkTrustRequirements {
	if in == nil {
		return nil
	}
	out := new(LinkTrustRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LinkType) DeepCopyInto(out *LinkType) {
	*out = *in
	if in.Protocol != nil {
		in, out := &in.Protocol, &out.Protocol
		*out = new(LinkProtocol)
		**out = **in
	}
	if in.MinQualityClass != nil {
		in, out := &in.MinQualityClass, &out.MinQualityClass
		*out = new(clusterv1.NetworkQualityClass)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LinkType.
func (in *LinkType) DeepCopy() *LinkType {
	if in == nil {
		return nil
	}
	out := new(LinkType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringConfig) DeepCopyInto(out *MonitoringConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringConfig.
func (in *MonitoringConfig) DeepCopy() *MonitoringConfig {
	if in == nil {
		return nil
	}
	out := new(MonitoringConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkLatencyRequirements) DeepCopyInto(out *NetworkLatencyRequirements) {
	*out = *in
	if in.MaxPacketDelayVariance != nil {
		in, out := &in.MaxPacketDelayVariance, &out.MaxPacketDelayVariance
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkLatencyRequirements.
func (in *NetworkLatencyRequirements) DeepCopy() *NetworkLatencyRequirements {
	if in == nil {
		return nil
	}
	out := new(NetworkLatencyRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkPacketLossRequirements) DeepCopyInto(out *NetworkPacketLossRequirements) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkPacketLossRequirements.
func (in *NetworkPacketLossRequirements) DeepCopy() *NetworkPacketLossRequirements {
	if in == nil {
		return nil
	}
	out := new(NetworkPacketLossRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NetworkThroughputRequirements) DeepCopyInto(out *NetworkThroughputRequirements) {
	*out = *in
	if in.MaxBandwidthVariance != nil {
		in, out := &in.MaxBandwidthVariance, &out.MaxBandwidthVariance
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NetworkThroughputRequirements.
func (in *NetworkThroughputRequirements) DeepCopy() *NetworkThroughputRequirements {
	if in == nil {
		return nil
	}
	out := new(NetworkThroughputRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeHardware) DeepCopyInto(out *NodeHardware) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeHardware.
func (in *NodeHardware) DeepCopy() *NodeHardware {
	if in == nil {
		return nil
	}
	out := new(NodeHardware)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeTrustRequirements) DeepCopyInto(out *NodeTrustRequirements) {
	*out = *in
	if in.MinTpmVersion != nil {
		in, out := &in.MinTpmVersion, &out.MinTpmVersion
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeTrustRequirements.
func (in *NodeTrustRequirements) DeepCopy() *NodeTrustRequirements {
	if in == nil {
		return nil
	}
	out := new(NodeTrustRequirements)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RainbowService) DeepCopyInto(out *RainbowService) {
	*out = *in
	out.Type = in.Type
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(ArbitraryObject)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RainbowService.
func (in *RainbowService) DeepCopy() *RainbowService {
	if in == nil {
		return nil
	}
	out := new(RainbowService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicasConfig) DeepCopyInto(out *ReplicasConfig) {
	*out = *in
	if in.InitialCount != nil {
		in, out := &in.InitialCount, &out.InitialCount
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicasConfig.
func (in *ReplicasConfig) DeepCopy() *ReplicasConfig {
	if in == nil {
		return nil
	}
	out := new(ReplicasConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceGraph) DeepCopyInto(out *ServiceGraph) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceGraph.
func (in *ServiceGraph) DeepCopy() *ServiceGraph {
	if in == nil {
		return nil
	}
	out := new(ServiceGraph)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceGraph) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceGraphList) DeepCopyInto(out *ServiceGraphList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ServiceGraph, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceGraphList.
func (in *ServiceGraphList) DeepCopy() *ServiceGraphList {
	if in == nil {
		return nil
	}
	out := new(ServiceGraphList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ServiceGraphList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceGraphNode) DeepCopyInto(out *ServiceGraphNode) {
	*out = *in
	if in.ServiceAccountName != nil {
		in, out := &in.ServiceAccountName, &out.ServiceAccountName
		*out = new(string)
		**out = **in
	}
	if in.PodLabels != nil {
		in, out := &in.PodLabels, &out.PodLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]corev1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]corev1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]corev1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Replicas.DeepCopyInto(&out.Replicas)
	if in.ExposedPorts != nil {
		in, out := &in.ExposedPorts, &out.ExposedPorts
		*out = new(ExposedPorts)
		(*in).DeepCopyInto(*out)
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(corev1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.SLOs != nil {
		in, out := &in.SLOs, &out.SLOs
		*out = make([]ServiceLevelObjective, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RainbowServices != nil {
		in, out := &in.RainbowServices, &out.RainbowServices
		*out = make([]RainbowService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.TrustRequirements != nil {
		in, out := &in.TrustRequirements, &out.TrustRequirements
		*out = new(NodeTrustRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.NodeHardware != nil {
		in, out := &in.NodeHardware, &out.NodeHardware
		*out = new(NodeHardware)
		**out = **in
	}
	if in.GeoLocation != nil {
		in, out := &in.GeoLocation, &out.GeoLocation
		*out = new(GeoLocation)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceGraphNode.
func (in *ServiceGraphNode) DeepCopy() *ServiceGraphNode {
	if in == nil {
		return nil
	}
	out := new(ServiceGraphNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceGraphSpec) DeepCopyInto(out *ServiceGraphSpec) {
	*out = *in
	if in.ServiceAccountName != nil {
		in, out := &in.ServiceAccountName, &out.ServiceAccountName
		*out = new(string)
		**out = **in
	}
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]ServiceGraphNode, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Links != nil {
		in, out := &in.Links, &out.Links
		*out = make([]ServiceLink, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SLOs != nil {
		in, out := &in.SLOs, &out.SLOs
		*out = make([]ServiceLevelObjective, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RainbowServices != nil {
		in, out := &in.RainbowServices, &out.RainbowServices
		*out = make([]RainbowService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ConfigMaps != nil {
		in, out := &in.ConfigMaps, &out.ConfigMaps
		*out = make([]ConfigMap, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceGraphSpec.
func (in *ServiceGraphSpec) DeepCopy() *ServiceGraphSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceGraphSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceGraphStatus) DeepCopyInto(out *ServiceGraphStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceGraphStatus.
func (in *ServiceGraphStatus) DeepCopy() *ServiceGraphStatus {
	if in == nil {
		return nil
	}
	out := new(ServiceGraphStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceLevelObjective) DeepCopyInto(out *ServiceLevelObjective) {
	*out = *in
	out.SloType = in.SloType
	out.ElasticityStrategy = in.ElasticityStrategy
	if in.SloConfig != nil {
		in, out := &in.SloConfig, &out.SloConfig
		*out = new(ArbitraryObject)
		**out = **in
	}
	if in.StaticElasticityStrategyConfig != nil {
		in, out := &in.StaticElasticityStrategyConfig, &out.StaticElasticityStrategyConfig
		*out = new(ArbitraryObject)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceLevelObjective.
func (in *ServiceLevelObjective) DeepCopy() *ServiceLevelObjective {
	if in == nil {
		return nil
	}
	out := new(ServiceLevelObjective)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceLink) DeepCopyInto(out *ServiceLink) {
	*out = *in
	if in.QosRequirements != nil {
		in, out := &in.QosRequirements, &out.QosRequirements
		*out = new(LinkQosRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.TrustRequirements != nil {
		in, out := &in.TrustRequirements, &out.TrustRequirements
		*out = new(LinkTrustRequirements)
		**out = **in
	}
	if in.SLOs != nil {
		in, out := &in.SLOs, &out.SLOs
		*out = make([]ServiceLevelObjective, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceLink.
func (in *ServiceLink) DeepCopy() *ServiceLink {
	if in == nil {
		return nil
	}
	out := new(ServiceLink)
	in.DeepCopyInto(out)
	return out
}
