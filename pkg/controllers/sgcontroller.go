// Copyright 2018 Oracle and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"context"
	"fmt"
	//"strings"
	"time"

	apps "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	wait "k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1beta1"
	coreinformers "k8s.io/client-go/informers/core/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	scheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1beta1"
	corelisters "k8s.io/client-go/listers/core/v1"
	cache "k8s.io/client-go/tools/cache"
	record "k8s.io/client-go/tools/record"
	workqueue "k8s.io/client-go/util/workqueue"

	//"github.com/coreos/go-semver/semver"
	"github.com/golang/glog"

	v1alpha1 "statefulguardian/pkg/apis/statefulguardian/v1alpha1"
	clientset "statefulguardian/pkg/generated/clientset/versioned"
	opscheme "statefulguardian/pkg/generated/clientset/versioned/scheme"
	informersv1alpha1 "statefulguardian/pkg/generated/informers/externalversions/statefulguardian/v1alpha1"
	listersv1alpha1 "statefulguardian/pkg/generated/listers/statefulguardian/v1alpha1"

	metrics "statefulguardian/pkg/metrics"
	buildversion "statefulguardian/pkg/version"
)

const controllerAgentName = "sg-operator"

const (
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a
	// Cluster fails to sync due to a resource of the same name already
	// existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a resource already existing.
	MessageResourceExists = "%s %s/%s already exists and is not managed by Cluster"
	// MessageResourceSynced is the message used for an Event fired when a
	// Cluster is synced successfully
	MessageResourceSynced = "Cluster synced successfully"
)


// The SGController watches the Kubernetes API for changes to resources
type SGController struct {

	kubeClient kubernetes.Interface
	opClient   clientset.Interface
	ctx context.Context
	sgDeploy *apps.Deployment

	shutdown bool
	queue    workqueue.RateLimitingInterface

	// sgLister is able to list/get Statefulguardians from a shared informer's
	// store.
	sgLister listersv1alpha1.StatefulguardianLister
	// sgListerSynced returns true if the Cluster shared informer has
	// synced at least once.
	sgListerSynced cache.InformerSynced
	// sgUpdater implements control logic for updating Cluster
	// statuses. Implemented as an interface to enable testing.
	sgUpdater sgUpdaterInterface

	// statefulSetLister is able to list/get StatefulSets from a shared
	// informer's store.
	statefulSetLister appslisters.StatefulSetLister
	// statefulSetListerSynced returns true if the StatefulSet shared informer
	// has synced at least once.
	statefulSetListerSynced cache.InformerSynced
	// statefulSetControl enables control of StatefulSets associated with
	// Statefulguardians.
	statefulSetControl StatefulSetControlInterface

	// podLister is able to list/get Pods from a shared
	// informer's store.
	podLister corelisters.PodLister
	// podListerSynced returns true if the Pod shared informer
	// has synced at least once.
	podListerSynced cache.InformerSynced

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController creates a new SGController.
func NewSgController(
	opClient clientset.Interface,
	kubeClient kubernetes.Interface,
	sgInformer informersv1alpha1.StatefulguardianInformer,
	statefulSetInformer appsinformers.StatefulSetInformer,
	podInformer coreinformers.PodInformer,
	serviceInformer coreinformers.ServiceInformer,
	resyncPeriod time.Duration,
	namespace string,
	ctx context.Context,
	sgDeploy *apps.Deployment) *SGController {
	opscheme.AddToScheme(scheme.Scheme) // TODO: This shouldn't be done here I don't think.

	// Create event broadcaster.
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	m := SGController{
		opClient:   opClient,
		kubeClient: kubeClient,

		sgLister:       sgInformer.Lister(),
		sgListerSynced: sgInformer.Informer().HasSynced,
		sgUpdater:      newSgUpdater(opClient, sgInformer.Lister()),

		statefulSetLister:       statefulSetInformer.Lister(),
		statefulSetListerSynced: statefulSetInformer.Informer().HasSynced,
		statefulSetControl:      NewRealStatefulSetControl(kubeClient, statefulSetInformer.Lister()),

		podLister:       podInformer.Lister(),
		podListerSynced: podInformer.Informer().HasSynced,

		queue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "sgcluster"),
		recorder: recorder,
		ctx: ctx,
		sgDeploy: sgDeploy,
	}

	sgInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: m.enqueueSg,
		UpdateFunc: func(old, new interface{}) {
			m.enqueueSg(new)
		},
		DeleteFunc: func(obj interface{}) {
			sg, ok := obj.(*v1alpha1.Statefulguardian)
			if ok {
				m.onSgDeleted(sg)
			}
		},
	})

	statefulSetInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: m.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newStatefulSet := new.(*apps.StatefulSet)
			oldStatefulSet := old.(*apps.StatefulSet)
			if newStatefulSet.ResourceVersion == oldStatefulSet.ResourceVersion {
				return
			}

			// If sg is ready ...
			if newStatefulSet.Status.ReadyReplicas == newStatefulSet.Status.Replicas {
				_, ok := newStatefulSet.Labels[SgLabel]
				if ok {
					m.onSgReady(newStatefulSet)
				}
			}
			m.handleObject(new)
		},
		DeleteFunc: m.handleObject,
	})

	return &m
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (m *SGController) Run(ctx context.Context, threadiness int) {
	defer utilruntime.HandleCrash()
	defer m.queue.ShutDown()

	glog.Info("Starting Cluster controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for Cluster controller informer caches to sync")
	if !WaitForCacheSync("sg", ctx.Done(),
		m.sgListerSynced,
		m.statefulSetListerSynced,
		m.podListerSynced) {
		return
	}

	glog.Info("Starting Cluster controller workers")
	// Launch two workers to process Foo resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(m.runWorker, time.Second, ctx.Done())
	}

	glog.Info("Started Cluster controller workers")
	defer glog.Info("Shutting down Cluster controller workers")
	<-ctx.Done()
}

// worker runs a worker goroutine that invokes processNextWorkItem until the
// controller's queue is closed.
func (m *SGController) runWorker() {
	for m.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (m *SGController) processNextWorkItem() bool {
	obj, shutdown := m.queue.Get()

	if shutdown {
		return false
	}

	err := func(obj interface{}) error {
		defer m.queue.Done(obj)
		key, ok := obj.(string)
		if !ok {
			m.queue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in queue but got %#v", obj))
			return nil
		}
		if err := m.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		m.queue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Cluster
// resource with the current status of the resource.
func (m *SGController) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name.
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	nsName := types.NamespacedName{Namespace: namespace, Name: name}

	// Get the Cluster resource with this namespace/name.
	sg, err := m.sgLister.Statefulguardians(namespace).Get(name)
	if err != nil {
		// The Cluster resource may no longer exist, in which case we stop processing.
		if apierrors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("sg '%s' in work queue no longer exists", key))
			return nil
		}
		return err
	}

	//sg.EnsureDefaults()
	//if err = sg.Validate(); err != nil {
	//	return errors.Wrap(err, "validating Cluster")
	//}

	operatorVersion := buildversion.GetBuildVersion()
	// Ensure that the required labels are set on the sg.
	sel := combineSelectors(SelectorForSg(sg), SelectorForSgVersion(operatorVersion))
	if !sel.Matches(labels.Set(sg.Labels)) {
		glog.V(2).Infof("Setting labels on sg %s", SelectorForSg(sg).String())
		if sg.Labels == nil {
			sg.Labels = make(map[string]string)
		}
		sg.Labels[SgLabel] = sg.Name
		sg.Labels[SgVersionLabel] = buildversion.GetBuildVersion()
		return m.sgUpdater.UpdateSgLabels(sg.DeepCopy(), labels.Set(sg.Labels))
	}

	//// If the Service is not controlled by this Cluster resource, we should
	//// log a warning to the event recorder and return.
	//if !metav1.IsControlledBy(svc, sg) {
	//	msg := fmt.Sprintf(MessageResourceExists, "Service", svc.Namespace, svc.Name)
	//	m.recorder.Event(sg, corev1.EventTypeWarning, ErrResourceExists, msg)
	//	return errors.New(msg)
	//}

	ss := sg.Spec.App
	ess, err := m.statefulSetLister.StatefulSets(ss.Namespace).Get(ss.Name)
	// If the resource doesn't exist, we'll create it
	if apierrors.IsNotFound(err) {
		glog.Infof("Creating a new StatefulSet for sg %q", nsName)
		ss.Labels[SgLabel] = sg.Name
		ss.Labels[SgNsLabel] = sg.Namespace
		err = m.statefulSetControl.CreateStatefulSet(&ss, m.sgDeploy)
	} else {
		ss = *ess
        }

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the StatefulSet is not controlled by this Cluster resource, we
	// should log a warning to the event recorder and return.
	if !metav1.IsControlledBy(&ss, sg) {
		msg := fmt.Sprintf(MessageResourceExists, "StatefulSet", ss.Namespace, ss.Name)
		m.recorder.Event(sg, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}


	//// If this number of the members on the Cluster does not equal the
	//// current desired replicas on the StatefulSet, we should update the
	//// StatefulSet resource.
	//if sg.Spec.Members != *ss.Spec.Replicas {
	//	glog.V(4).Infof("Updating %q: sgMembers=%d statefulSetReplicas=%d",
	//		nsName, sg.Spec.Members, ss.Spec.Replicas)
	//	old := ss.DeepCopy()
	//	ss = statefulsets.NewForCluster(sg, m.opConfig.Images, svc.Name)
	//	if err := m.statefulSetControl.Patch(old, ss); err != nil {
	//		// Requeue the item so we can attempt processing again later.
	//		// This could have been caused by a temporary network failure etc.
	//		return err
	//	}
	//}

	// Finally, we update the status block of the Cluster resource to
	// reflect the current state of the world.
	err = m.updateSgStatus(sg, &ss)
	if err != nil {
		return err
	}

	m.recorder.Event(sg, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}







// updateSgStatusForSS updates Cluster statuses based on changes to their associated StatefulSets.
func (m *SGController) updateSgStatus(sg *v1alpha1.Statefulguardian, ss *apps.StatefulSet) error {
	glog.V(4).Infof("%s/%s: ss.Spec.Replicas=%d, ss.Status.ReadyReplicas=%d, ss.Status.Replicas=%d",
		sg.Namespace, sg.Name, *ss.Spec.Replicas, ss.Status.ReadyReplicas, ss.Status.Replicas)

	status := sg.Status.DeepCopy()
	_, condition := GetSgCondition(&sg.Status, v1alpha1.StatefulguardianReady)
	if condition == nil {
		condition = &v1alpha1.StatefulguardianCondition{Type: v1alpha1.StatefulguardianReady}
	}
	if ss.Status.ReadyReplicas == ss.Status.Replicas {
		condition.Status = corev1.ConditionTrue
	} else {
		condition.Status = corev1.ConditionFalse
	}

	if updated := UpdateSgCondition(status, condition); updated {
		return m.sgUpdater.UpdateSgStatus(sg.DeepCopy(), status)
	}
	return nil
}

// enqueueSg takes a Cluster resource and converts it into a
// namespace/name string which is then put onto the work queue. This method
// should *not* be passed resources of any type other than Cluster.
func (m *SGController) enqueueSg(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	m.queue.AddRateLimited(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Foo resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (m *SGController) handleObject(obj interface{}) {
	object, ok := obj.(metav1.Object)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		glog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}

	glog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Cluster, we should not do
		// anything more with it.
		if ownerRef.Kind != v1alpha1.StatefulguardianCRDResourceKind {
			return
		}

		sg, err := m.sgLister.Statefulguardians(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			glog.V(4).Infof("ignoring orphaned object '%s' of Cluster '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		m.enqueueSg(sg)
		return
	}
}

func (m *SGController) injectEnv() {
}


func (m *SGController) getSgBySs(ss *apps.StatefulSet) *v1alpha1.Statefulguardian {
	sgName := ss.Labels[SgLabel]
	sgNs := ss.Labels[SgNsLabel]

	sg, err := m.sgLister.Statefulguardians(sgNs).Get(sgName)
	if err==nil {
		return sg
	}
	glog.V(4).Infof("Error get '%s' in namespace '%s'", sgName, sgNs)

	return nil
}

func (m *SGController) onSgReady(ss *apps.StatefulSet) {
	sg := m.getSgBySs(ss)
	if sg==nil {
		return
	}
	glog.Infof("Cluster %s ready", sg.Name)
        m.injectEnv()
	execCont:= NewExecController(sg)
        execCont.ClusterInit(m.ctx)
	metrics.IncEventCounter(sgsCreatedCount)
	metrics.IncEventGauge(sgsTotalCount)
}

// Run quit operators and delete the statefulset
func (m *SGController) onSgDeleted(sg *v1alpha1.Statefulguardian) {
	glog.Infof("Cluster %s deleted", sg.Name)
	execCont:= NewExecController(sg)
        execCont.ClusterQuit(m.ctx)
	metrics.IncEventCounter(sgsDeletedCount)
	metrics.DecEventGauge(sgsTotalCount)
	glog.Infof("Delete statefulset")
	m.statefulSetControl.DeleteStatefulSet(sg)
}

func combineSelectors(first labels.Selector, rest ...labels.Selector) labels.Selector {
	res := first.DeepCopySelector()
	for _, s := range rest {
	     reqs, _ := s.Requirements()
	     res = res.Add(reqs...)
	}
   	return res
}

// SelectorForSg creates a labels.Selector to match a given clusters
// associated resources.
func SelectorForSg(c *v1alpha1.Statefulguardian) labels.Selector {
	return labels.SelectorFromSet(labels.Set{SgLabel: c.Name})
}

func SelectorForSgVersion(operatorVersion string) labels.Selector {
	return labels.SelectorFromSet(labels.Set{SgVersionLabel: operatorVersion})
}
