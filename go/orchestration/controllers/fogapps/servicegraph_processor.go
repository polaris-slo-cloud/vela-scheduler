package fogapps

import (
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	fogapps "k8s.rainbow-h2020.eu/rainbow/orchestration/apis/fogapps/v1"
	svcGraphUtil "k8s.rainbow-h2020.eu/rainbow/orchestration/internal/servicegraphutil"
	"k8s.rainbow-h2020.eu/rainbow/orchestration/pkg/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type serviceGraphChildObjectMaps struct {
	Deployments  map[string]*apps.Deployment
	StatefulSets map[string]*apps.StatefulSet
	Services     map[string]*core.Service
	Ingresses    map[string]*networking.Ingress
}

type serviceGraphProcessor struct {
	svcGraph *fogapps.ServiceGraph

	// The child objects that already existed prior to this processing.
	existingChildObjects serviceGraphChildObjectMaps

	// The child objects that were created during this processing.
	newChildObjects serviceGraphChildObjectMaps

	log        logr.Logger
	setOwnerFn controllerutil.SetOwnerReferenceFn
	changes    *controllerutil.ResourceChangesList
}

// ProcessServiceGraph assembles a list of changes that need to be applied due to the specified ServiceGraph.
func ProcessServiceGraph(
	graph *fogapps.ServiceGraph,
	childObjects *serviceGraphChildObjects,
	log logr.Logger,
	setOwnerFn controllerutil.SetOwnerReferenceFn,
) (*controllerutil.ResourceChangesList, error) {
	graphProcessor := newServiceGraphProcessor(graph, childObjects, log, setOwnerFn)

	if err := graphProcessor.assembleGraphChanges(); err != nil {
		return nil, err
	}

	return graphProcessor.changes, nil
}

func newServiceGraphChildObjectMaps(lists *serviceGraphChildObjects) serviceGraphChildObjectMaps {
	maps := serviceGraphChildObjectMaps{
		Deployments:  make(map[string]*apps.Deployment),
		StatefulSets: make(map[string]*apps.StatefulSet),
		Services:     make(map[string]*core.Service),
		Ingresses:    make(map[string]*networking.Ingress),
	}

	if lists != nil {
		for i, item := range lists.Deployments {
			// item is apparently an object that is allocated for the loop and overwritten with
			// the values of lists.Deployments[i], which means that the address of item is
			// the same on every iteration.
			maps.Deployments[item.Name] = &lists.Deployments[i]
		}
		for i, item := range lists.StatefulSets {
			maps.StatefulSets[item.Name] = &lists.StatefulSets[i]
		}
		for i, item := range lists.Services {
			maps.Services[item.Name] = &lists.Services[i]
		}
		for i, item := range lists.Ingresses {
			maps.Ingresses[item.Name] = &lists.Ingresses[i]
		}
	}

	return maps
}

func newServiceGraphProcessor(
	graph *fogapps.ServiceGraph,
	childObjects *serviceGraphChildObjects,
	log logr.Logger,
	setOwnerFn controllerutil.SetOwnerReferenceFn,
) *serviceGraphProcessor {
	return &serviceGraphProcessor{
		svcGraph:             graph,
		existingChildObjects: newServiceGraphChildObjectMaps(childObjects),
		newChildObjects:      newServiceGraphChildObjectMaps(nil),
		log:                  log,
		setOwnerFn:           setOwnerFn,
		changes:              controllerutil.NewResourceChangesList(),
	}
}

// assembleGraphChanges assembles the list of changes that need to be made to align the cluster state with
// the state of the ServiceGraph.
//
// We use the following approach:
// 1. Iterate through the ServiceGraph and create or update the child objects (Deployments, StatefulSets, SLOs, etc.) and store them in newChildObjects.
// 2. Iterate through the lists of existing child objects and check if a corresponding new child object exists. We have the following options
// 		- If a corresponding new child object exists, check if there are any changes between the two specs.
//		  If there are changes, create a ResourceUpdate. In any case, remove the new child object from the newChildObjects map.
// 		- Otherwise, create a ResourceDeletion.
// 3. For all new child objects that are still in the newChildObjects map, create a ResourceAddition.
func (me *serviceGraphProcessor) assembleGraphChanges() error {
	if err := me.createChildObjectsForServiceGraph(); err != nil {
		return err
	}

	if err := me.assembleUpdatesForDeployments(); err != nil {
		return err
	}
	if err := me.assembleUpdatesForStatefulSets(); err != nil {
		return err
	}
	if err := me.assembleUpdatesForServices(); err != nil {
		return err
	}
	if err := me.assembleUpdatesForIngresses(); err != nil {
		return err
	}
	if err := me.assembleAdditions(); err != nil {
		return err
	}

	return nil
}

func (me *serviceGraphProcessor) createChildObjectsForServiceGraph() error {
	for _, node := range me.svcGraph.Spec.Nodes {
		switch node.NodeType {
		case fogapps.ServiceNode:
			if err := me.createChildObjectsForServiceNode(&node); err != nil {
				return err
			}
		case fogapps.UserNode:
			// Nothing to be done here.
		default:
			return fmt.Errorf("unknown ServiceGraphNode.NodeType: %s", node.NodeType)
		}
	}
	return nil
}

func (me *serviceGraphProcessor) createChildObjectsForServiceNode(node *fogapps.ServiceGraphNode) error {
	var err error

	switch node.Replicas.SetType {
	case fogapps.SimpleReplicaSet:
		err = me.createOrUpdateDeployment(node)
	case fogapps.StatefulReplicaSet:
		err = me.createOrUpdateStatefulSet(node)
	}

	if err != nil {
		return err
	}

	// If ExposedPorts are set, create or update the Service and Ingress.
	// If no ExposedPorts are set, not creating any Service or Ingress will cause any existing ones to be deleted later.
	if node.ExposedPorts != nil {
		if err = me.createOrUpdateServiceAndIngress(node); err != nil {
			return err
		}
	}

	return nil
}

func (me *serviceGraphProcessor) createOrUpdateDeployment(node *fogapps.ServiceGraphNode) error {
	var deployment *apps.Deployment
	var err error

	if existingDeployment, isUpdate := me.existingChildObjects.Deployments[node.Name]; isUpdate {
		deployment, err = svcGraphUtil.UpdateDeployment(existingDeployment.DeepCopy(), node, me.svcGraph)
	} else {
		if deployment, err = svcGraphUtil.CreateDeployment(node, me.svcGraph); err != nil {
			return err
		}
		err = me.setOwner(deployment)
	}

	if err != nil {
		return err
	}

	me.newChildObjects.Deployments[deployment.Name] = deployment
	return nil
}

func (me *serviceGraphProcessor) createOrUpdateStatefulSet(node *fogapps.ServiceGraphNode) error {
	var statefulSet *apps.StatefulSet
	var err error

	if existingStatefulSet, isUpdate := me.existingChildObjects.StatefulSets[node.Name]; isUpdate {
		statefulSet, err = svcGraphUtil.UpdateStatefulSet(existingStatefulSet.DeepCopy(), node, me.svcGraph)
	} else {
		if statefulSet, err = svcGraphUtil.CreateStatefulSet(node, me.svcGraph); err != nil {
			return err
		}
		err = me.setOwner(statefulSet)
	}

	if err != nil {
		return err
	}

	me.newChildObjects.StatefulSets[statefulSet.Name] = statefulSet
	return nil
}

func (me *serviceGraphProcessor) createOrUpdateServiceAndIngress(node *fogapps.ServiceGraphNode) error {
	existingServiceAndIngress := &svcGraphUtil.ServiceAndIngressPair{}
	var serviceAndIngress *svcGraphUtil.ServiceAndIngressPair
	var err error

	// Check if we have an existing Service and Ingress.
	if existingService, ok := me.existingChildObjects.Services[node.Name]; ok {
		existingServiceAndIngress.Service = existingService.DeepCopy()
	}
	if existingIngress, ok := me.existingChildObjects.Ingresses[node.Name]; ok {
		existingServiceAndIngress.Ingress = existingIngress.DeepCopy()
	}

	if existingServiceAndIngress.Service != nil || existingServiceAndIngress.Ingress != nil {
		serviceAndIngress, err = svcGraphUtil.UpdateServiceAndIngress(existingServiceAndIngress, node, me.svcGraph)
	} else {
		if serviceAndIngress, err = svcGraphUtil.CreateServiceAndIngress(node, me.svcGraph); err == nil {
			if serviceAndIngress.Service != nil {
				if err = me.setOwner(serviceAndIngress.Service); err != nil {
					return err
				}
			}
			if serviceAndIngress.Ingress != nil {
				if err := me.setOwner(serviceAndIngress.Ingress); err != nil {
					return err
				}
			}
		}
	}

	if err != nil {
		return err
	}

	if serviceAndIngress.Service != nil {
		me.newChildObjects.Services[serviceAndIngress.Service.Name] = serviceAndIngress.Service
	}
	if serviceAndIngress.Ingress != nil {
		me.newChildObjects.Ingresses[serviceAndIngress.Ingress.Name] = serviceAndIngress.Ingress
	}

	return nil
}

func (me *serviceGraphProcessor) setOwner(childObj client.Object) error {
	if err := me.setOwnerFn(childObj); err != nil {
		return fmt.Errorf("could not set owner reference. Cause: %w", err)
	}
	return nil
}

func (me *serviceGraphProcessor) assembleUpdatesForDeployments() error {
	for _, existingDeployment := range me.existingChildObjects.Deployments {
		if updatedDeployment, ok := me.newChildObjects.Deployments[existingDeployment.Name]; ok {

			// ToDo: Containers are currently never equal, because Kubernetes sets some values, which are unset in new containers.
			// containersEqual := reflect.DeepEqual(existingDeployment.Spec.Template.Spec.Containers, updatedDeployment.Spec.Template.Spec.Containers)
			// _ = containersEqual

			if !reflect.DeepEqual(existingDeployment.Spec, updatedDeployment.Spec) {
				// Deployment was changed, we need to update it
				me.changes.AddChanges(controllerutil.NewResourceUpdate(updatedDeployment))
			}

			delete(me.newChildObjects.Deployments, updatedDeployment.Name)
		} else {
			// The corresponding ServiceGraphNode was deleted, so we delete the Deployment
			me.changes.AddChanges(controllerutil.NewResourceDeletion(existingDeployment))
		}
	}
	return nil
}

func (me *serviceGraphProcessor) assembleUpdatesForStatefulSets() error {
	for _, existingStatefulSet := range me.existingChildObjects.StatefulSets {
		if updatedStatefulSet, ok := me.newChildObjects.StatefulSets[existingStatefulSet.Name]; ok {

			if !reflect.DeepEqual(existingStatefulSet.Spec, updatedStatefulSet.Spec) {
				// StatefulSet was changed, we need to update it
				me.changes.AddChanges(controllerutil.NewResourceUpdate(updatedStatefulSet))
			}

			delete(me.newChildObjects.StatefulSets, updatedStatefulSet.Name)
		} else {
			// The corresponding ServiceGraphNode was deleted, so we delete the StatefulSet
			me.changes.AddChanges(controllerutil.NewResourceDeletion(existingStatefulSet))
		}
	}
	return nil
}

func (me *serviceGraphProcessor) assembleUpdatesForServices() error {
	for _, existingService := range me.existingChildObjects.Services {
		if updatedService, ok := me.newChildObjects.Services[existingService.Name]; ok {

			if !reflect.DeepEqual(existingService.Spec, updatedService.Spec) {
				// Service was changed, we need to update it
				me.changes.AddChanges(controllerutil.NewResourceUpdate(updatedService))
			}

			delete(me.newChildObjects.Services, updatedService.Name)
		} else {
			// The corresponding ServiceGraphNode was deleted, so we delete the Service
			me.changes.AddChanges(controllerutil.NewResourceDeletion(existingService))
		}
	}
	return nil
}

func (me *serviceGraphProcessor) assembleUpdatesForIngresses() error {
	for _, existingIngress := range me.existingChildObjects.Ingresses {
		if updatedIngress, ok := me.newChildObjects.Ingresses[existingIngress.Name]; ok {

			if !reflect.DeepEqual(existingIngress.Spec, updatedIngress.Spec) {
				// Ingress was changed, we need to update it
				me.changes.AddChanges(controllerutil.NewResourceUpdate(updatedIngress))
			}

			delete(me.newChildObjects.Ingresses, updatedIngress.Name)
		} else {
			// The corresponding ServiceGraphNode was deleted, so we delete the Ingress
			me.changes.AddChanges(controllerutil.NewResourceDeletion(existingIngress))
		}
	}
	return nil
}

func (me *serviceGraphProcessor) assembleAdditions() error {
	for _, value := range me.newChildObjects.Deployments {
		me.changes.AddChanges(controllerutil.NewResourceAddition(value))
	}
	for _, value := range me.newChildObjects.StatefulSets {
		me.changes.AddChanges(controllerutil.NewResourceAddition(value))
	}
	for _, value := range me.newChildObjects.Services {
		me.changes.AddChanges(controllerutil.NewResourceAddition(value))
	}
	for _, value := range me.newChildObjects.Ingresses {
		me.changes.AddChanges(controllerutil.NewResourceAddition(value))
	}
	return nil
}
