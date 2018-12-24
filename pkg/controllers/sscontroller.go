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
	//"fmt"
	"encoding/json"
	"github.com/golang/glog"

	apps "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appslistersv1beta1 "k8s.io/client-go/listers/apps/v1beta1"
	//corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	v1alpha1 "statefulguardian/pkg/apis/statefulguardian/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"

)

// StatefulSetControlInterface defines the interface that the
// ClusterController uses to create and update StatefulSets. It
// is implemented as an interface to enable testing.
type StatefulSetControlInterface interface {
	CreateStatefulSet(ss *apps.StatefulSet, depl *apps.Deployment) error
	DeleteStatefulSet(sg *v1alpha1.Statefulguardian) error 
	Patch(old *apps.StatefulSet, new *apps.StatefulSet) error
}

type realStatefulSetControl struct {
	client            kubernetes.Interface
	statefulSetLister appslistersv1beta1.StatefulSetLister
}

// NewRealStatefulSetControl creates a concrete implementation of the
// StatefulSetControlInterface.
func NewRealStatefulSetControl(client kubernetes.Interface, statefulSetLister appslistersv1beta1.StatefulSetLister) StatefulSetControlInterface {
	return &realStatefulSetControl{client: client, statefulSetLister: statefulSetLister}
}

func (rssc *realStatefulSetControl) CreateStatefulSet(ss *apps.StatefulSet, depl *apps.Deployment) error {
	// config statefulset
	ss = rssc.SetPriority(ss)
	ss = rssc.SetOwnership(ss, depl)

	_, err := rssc.client.AppsV1beta1().StatefulSets(ss.Namespace).Create(ss)
	return err
}

func (rssc *realStatefulSetControl) DeleteStatefulSet(sg *v1alpha1.Statefulguardian) error {
	ss := sg.Spec.App
	glog.Infof("Delete statefulset %s namespaces %s" , ss.Name, ss.Namespace)
	backgroundPolicy := metav1.DeletePropagationBackground
	err := rssc.client.AppsV1beta1().StatefulSets(ss.Namespace).Delete(ss.Name, &metav1.DeleteOptions{PropagationPolicy: &backgroundPolicy})
	return err
}

// PatchStatefulSet performs a direct patch update for the specified StatefulSet.
func PatchStatefulSet(kubeClient kubernetes.Interface, oldData *apps.StatefulSet, newData *apps.StatefulSet) (*apps.StatefulSet, error) {
	originalJSON, err := json.Marshal(oldData)
	if err != nil {
		return nil, err
	}

	updatedJSON, err := json.Marshal(newData)
	if err != nil {
		return nil, err
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(
		originalJSON, updatedJSON, apps.StatefulSet{})
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Patching StatefulSet %q: %s", types.NamespacedName{Namespace: oldData.Namespace, Name: oldData.Name}, string(patchBytes))

	result, err := kubeClient.AppsV1beta1().StatefulSets(oldData.Namespace).Patch(oldData.Name, types.StrategicMergePatchType, patchBytes)
	if err != nil {
		glog.Errorf("Failed to patch StatefulSet: %v", err)
		return nil, err
	}

	return result, nil
}

func (rssc *realStatefulSetControl) Patch(old *apps.StatefulSet, new *apps.StatefulSet) error {
	_, err := PatchStatefulSet(rssc.client, old, new)
	return err
}


func (rssc *realStatefulSetControl) SetPriority(old *apps.StatefulSet) *apps.StatefulSet {
	return old
}

// set the owner of each ss to the sg operator so that we can delete it gracefully
func (rssc *realStatefulSetControl) SetOwnership(old *apps.StatefulSet, depl *apps.Deployment) *apps.StatefulSet {
	old.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
				*metav1.NewControllerRef(depl, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    v1alpha1.StatefulguardianCRDResourceKind,
				}),
	}
	glog.Infof("Setting StatefulSet owner: uid %s name %s" , string(old.ObjectMeta.OwnerReferences[0].UID),string(old.ObjectMeta.OwnerReferences[0].Name))

	return old
}
