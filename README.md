# Vela Scheduler

Vela Scheduler is an orchestrator-independent distributed scheduler for the Edge-to-Cloud continuum.
Detailed design and architecture documentation is available [here](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/docs).

The Edge-to-Cloud continuum may span hundreds of thousands of nodes and can be formed of multiple Edge and Cloud clusters, which may be operated by different orchestrators, e.g., one cluster could be managed by [Kubernetes](https://kubernetes.io) and another one by [Nomad](https://www.nomadproject.io).
Vela Scheduler can be run as an arbitrary number of instances.
To be compatible with the Vela Scheduler, each cluster must run the Vela Cluster Agent, which is the only interface needed by the scheduler to interact with the cluster.

![Vela Scheduler Workflow Overview](./assets/vela-scheduler-overview.svg)

The high-level scheduling workflow of Vela Scheduler is depicted in the diagram above. 
A user or a system component submits a job to an arbitrary scheduler instance through its Submit Job API, which adds this job to the scheduler’s queue.
Once dequeued, the job enters the scheduling pipeline, which begins by obtaining a sample of candidate nodes from every cluster.
To this end, the scheduler contacts the Vela Cluster Agent in each cluster and sends it the description of the job to be scheduled.
Using the Cluster Agent prevents the scheduler from having to communicate with each node directly. For each node picked by the sampling algorithm, the Cluster Agent performs initial filtering and scoring, i.e., it determines if the node meets the requirements set forth by the received job description and assigns a score to each node that passes filtering to indicates how well it is suited for the job.
The list suitable nodes along with their scores is returned to the scheduler.
The scheduler may run additional filtering and scoring steps and subsequently chooses the node with the highest score.
Finally, the scheduler submits the job to the Cluster Agent that is responsible for the chosen node to deploy the job.

## Design & Documentation

The documentation for the Vela Scheduler is available in the [docs](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/docs) folder.

Godoc documentation for the source code is available [here](./godoc/pkg).


## Experiments

Details on experiment results obtained with Vela Scheduler can be found [here](./experiments).


### Repository Organization

This project is realized in Go 1.19 and relies on Go workspaces to divide code into multiple modules.

| Directory                | Contents |
|--------------------------|----------|
| [`deployment`](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/deployment) | Configuration files for deploying Vela Scheduler and Vela Cluster Agent |
| [`docs`](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/docs) | Documentation files (Work in progress) |
| [`go`](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/go) | Go workspace containing all Vela Scheduler modules |
| [`go/cluster-agent`](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/go/cluster-agent) | The polaris-cluster-agent executable module |
| [`go/context-awareness`](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/go/context-awareness) | Context-aware scheduling plugins  |
| [`go/framework`](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/go/framework) | The orchestrator-independent Vela Scheduler framework library that defines the scheduling pipeline and plugin structures. |
| [`go/k8s-connector`](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/go/k8s-connector) | Kubernetes orchestrator connector |
| [`go/scheduler`](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/go/scheduler) | Vela Scheduler executable module |
| [`testbed`](https://github.com/polaris-slo-cloud/vela-scheduler/tree/master/testbed) | Scripts and configurations for setting testbeds for experimenting with Vela Scheduler |

