# 3-Phase Scheduling Workflow

![Vela 3-Phase Scheduling Workflow](./assets/scheduling-workflow.png)

This document describes the 3-Phase Scheduling Workflow of Vela Scheduler in more detail to provide an overview of the involved interfaces and classes (we use the term "class" to refer to Go structs).
Note that the code uses the name `PolarisScheduler` instead of Vela.
The packages in the class diagrams below do not resemble actual package names, but are used to indicate whether a type resides within the scheduler or the cluster agent.

# Sampling Phase

![Main Interfaces and Classes in Sampling Phase](./assets/sampling-classes.svg)

The class diagram above shows the main interfaces and classes involved in the sampling phase.
The `DefaultPolarisScheduler` is the main class implementing the scheduler process.
It is initialized with a `PodSource` object;
The `PodSubmissionApi`, which provides the REST interface for receiving new pods implements this interface.
Upon receiving a new pod, from its `PodSource`, the `DefaultPolarisScheduler` adds it to the `PrioritySchedulingQueue`, which represents the sampling queue.
It uses a `SortPlugin` to establish an order among the incoming pods - only one `SortPlugin` can be configured.

The scheduler maintains a list of `SampleNodesPlugin` objects, which resemble the sampler pool.
All these objects are of the same type, because only one such plugin can be configured.
A new pods is dequeued from the sampling queue, whenever there is an idle sampler available.
By default, 2-Smart Sampling is used - its high-level algorithm is shown in the figure below.
To this end, the `RemoteNodesSamplerPlugin` is used at this stage, which relies on the `DefaultRemoteSamplerClientsManager` to contact the `Cp` percent of the remote cluster agents to ask them for node samples.

![2-Smart Sampling Algorithm](./assets/2-smart-sampling-algorithm.png)

In the cluster agent, the `DefaultPolarisNodesSampler` exposes a REST interface to accept sampling requests.
Each request contains the pod and its requirements, the sampling strategy to use (e.g., `random` or `round-robin`), as well as, the percent of nodes that should be sampled (`Np`).
The nodes sampler maintains a set of `SamplingPipeline` instances, each of them containing the configured set of plugins.
Each sampling strategy is implemented by a `SamplingStrategyPlugin` - there is only a single instance of each plugin that is shared across all sampling pipelines.
This decision was made intentionally to ensure that these plugins can maintain a global state within the cluster agent, e.g., the `RoundRobinSamplingPlugin` has to advance its current index with every request irrespective of the sampling pipeline that issues it, otherwise two subsequent requests from different pipelines would get the same nodes.
The sampling strategy is used to obtain node samples, which are then evaluated by the filter and score plugins.
This procedure is carried out in a loop until `Np` eligible nodes have been found or a timeout is reached.
The final list of sampled nodes is returned to the scheduler, where it is added to the decision pipeline queue.

