#!/bin/bash

CLUSTER_CONFIG_DIR=$(dirname "${BASH_SOURCE}")
source "${CLUSTER_CONFIG_DIR}/common.sh"

# The name of the kind cluster.
kindClusterName="kind-07-edge-frankfurt"

# The port on localhost, where the polaris-cluster-agent of this cluster should be exposed.
clusterAgentPortLocalhost=30007

# (optional) Additional kind node config.
# For config options see https://kind.sigs.k8s.io/docs/user/configuration/
kindExtraConfig=$(cat <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  # port forward ${clusterAgentPortLocalhost} on the host to 30033 on the control-plane node
  extraPortMappings:
  - containerPort: 30033
    hostPort: ${clusterAgentPortLocalhost}
    # optional: set the bind address on the host
    # 0.0.0.0 is the current default
    # listenAddress: "127.0.0.1"
    # optional: set the protocol to one of TCP, UDP, SCTP.
    # TCP is the default
    # protocol: TCP
EOF
)

# (optional) Create a namespace with this name.
testNamespace="test"

locationA="50.113485878139905_8.664028308860287"
locationB="50.078026635390955_8.486530443727274"
cloudletLocation="50.113485878139905_8.664028308860287"

# Declares the types of fake nodes and how many nodes of each type to create.
# For each fake node type, the amount of CPUs and memory must be added to fakeNodeTypeCpus and fakeNodeTypeMemory respectively.
declare -A fakeNodeTypes=(
    ["raspi-3b-plus"]="30"
    ["raspi-4b-2gi"]="20"
    ["raspi-4b-4gi-loc-a"]="10"
    ["raspi-4b-4gi-loc-a-low-batt"]="10"
    ["raspi-4b-4gi-loc-b"]="10"
    ["raspi-4b-4gi-loc-b-low-batt"]="10"
    ["cloudlet"]="10"
)

# Each node's CPUs are configured as `cpu` and `polaris-slo-cloud.github.io/fake-milli-cpu`.
declare -A fakeNodeTypeCpus=(
    ["raspi-3b-plus"]="4000m"
    ["raspi-4b-2gi"]="4000m"
    ["raspi-4b-4gi-loc-a"]="4000m"
    ["raspi-4b-4gi-loc-a-low-batt"]="4000m"
    ["raspi-4b-4gi-loc-b"]="4000m"
    ["raspi-4b-4gi-loc-b-low-batt"]="4000m"
    ["cloudlet"]="4000m"
)

# Each node's memory is configured as `memory` and `polaris-slo-cloud.github.io/fake-memory`.
declare -A fakeNodeTypeMemory=(
    ["raspi-3b-plus"]="1Gi"
    ["raspi-4b-2gi"]="2Gi"
    ["raspi-4b-4gi-loc-a"]="4Gi"
    ["raspi-4b-4gi-loc-a-low-batt"]="4Gi"
    ["raspi-4b-4gi-loc-b"]="4Gi"
    ["raspi-4b-4gi-loc-b-low-batt"]="4Gi"
    ["cloudlet"]="8Gi"
)

# Optional extra node labels for each node type.
# The value for each node type has to be a string of the following format (slashes and quotes must be escaped):
# "<domain1.io>\/<label1>: <value1>;<domain2.io>\/<label2>: <value2>;<...>"
declare -A extraNodeLabels=(
    ["raspi-3b-plus"]="polaris-slo-cloud.github.io\/battery.capacity-mah: \"2000\";polaris-slo-cloud.github.io\/battery.level: \"90\";polaris-slo-cloud.github.io\/geo-location: \"${locationA}\""
    ["raspi-4b-2gi"]="polaris-slo-cloud.github.io\/battery.capacity-mah: \"4000\";polaris-slo-cloud.github.io\/battery.level: \"50\";polaris-slo-cloud.github.io\/geo-location: \"${locationA}\""
    ["raspi-4b-4gi-loc-a"]="polaris-slo-cloud.github.io\/battery.capacity-mah: \"4000\";polaris-slo-cloud.github.io\/battery.level: \"70\";polaris-slo-cloud.github.io\/geo-location: \"${locationA}\""
    ["raspi-4b-4gi-loc-a-low-batt"]="polaris-slo-cloud.github.io\/battery.capacity-mah: \"4000\";polaris-slo-cloud.github.io\/battery.level: \"20\";polaris-slo-cloud.github.io\/geo-location: \"${locationB}\""
    ["raspi-4b-4gi-loc-b"]="polaris-slo-cloud.github.io\/battery.capacity-mah: \"4000\";polaris-slo-cloud.github.io\/battery.level: \"70\";polaris-slo-cloud.github.io\/geo-location: \"${locationB}\""
    ["raspi-4b-4gi-loc-b-low-batt"]="polaris-slo-cloud.github.io\/battery.capacity-mah: \"4000\";polaris-slo-cloud.github.io\/battery.level: \"20\";polaris-slo-cloud.github.io\/geo-location: \"${locationA}\""
    ["cloudlet"]="polaris-slo-cloud.github.io\/geo-location: \"${cloudletLocation}\""
)

# Extended resources.
# The value for each node type has to be a string of the following format (slashes must be escaped):
# "<domain1.io>\/<resource1>: <count1>;<domain2.io>\/<resource2>: <count2>;<...>"
declare -A extendedResources=(
    # ["cell-5g-base-station"]="polaris-slo-cloud.github.io\/base-station-5g: 1;polaris-slo-cloud.github.io\/test-resource: 1"
)
