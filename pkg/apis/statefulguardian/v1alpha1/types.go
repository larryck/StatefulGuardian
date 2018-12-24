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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	apps "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LifecycleOperation struct {
	Init OperationSet `json:"init,omitempty"`
	Quit OperationSet `json:"quit,omitempty"`
	Funcs OperationSet `json:"funcs,omitempty"`
}

type OperationSet struct {
	Cluster []Operation `json:"cluster,omitempty"`
	Pod []Operation `json:"pod,omitempty"`
}

type Operation struct {
	Name string `json:"name,omitempty"`
	Desc string `json:"desc,omitempty"`
	Cmd []string `json:"cmd,omitempty"`
	Container string `json:"container,omitempty"`
}

type AppMonitorOperation struct {
	Interval int `json:"interval,omitempty"`
	Metrics MetricsSrc `json:"metrics,omitempty"`
}

type MetricsSrc struct {
	Name string `json:"name,omitempty"`
	Desc string `json:"desc,omitempty"`
	Cmd []string `json:"cmd,omitempty"`
	Unit string `json:"unit,omitempty"`
	Container string `json:"container,omitempty"`
}


// StatefulguardianSpec defines the attributes a user can specify when creating a cluster
type StatefulguardianSpec struct {
        Lifecycle LifecycleOperation  `json:"lifecycle,omitempty"`
        AppMonitor AppMonitorOperation  `json:"appmonitor,omitempty"`
	App apps.StatefulSet  `json:"app,omitempty"`
	Resources *Resources `json:"resources,omitempty"`
}

// StatefulguardianConditionType represents a valid condition of a Statefulguardian.
type StatefulguardianConditionType string

const (
	// StatefulguardianReady means the Statefulguardian is able to service requests.
	StatefulguardianReady StatefulguardianConditionType = "Ready"
)

// StatefulguardianCondition describes the observed state of a Statefulguardian at a certain point.
type StatefulguardianCondition struct {
	Type   StatefulguardianConditionType
	Status corev1.ConditionStatus
	// +optional
	LastTransitionTime metav1.Time
	// +optional
	Reason string
	// +optional
	Message string
}

type StatefulguardianStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	// +optional
	Conditions []StatefulguardianCondition
}

// +genclient
// +genclient:noStatus
// +resourceName=statefulguardians
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Statefulguardian represents a cluster spec and associated metadata
type Statefulguardian struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   StatefulguardianSpec   `json:"spec"`
	Status StatefulguardianStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type StatefulguardianList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Statefulguardian `json:"items"`
}

type Resources struct {
	Agent  *corev1.ResourceRequirements `json:"agent,omitempty"`
	Server *corev1.ResourceRequirements `json:"server,omitempty"`
}


