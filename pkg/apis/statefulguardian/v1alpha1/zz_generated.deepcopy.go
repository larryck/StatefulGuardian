// +build !ignore_autogenerated

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppMonitorOperation) DeepCopyInto(out *AppMonitorOperation) {
	*out = *in
	in.Metrics.DeepCopyInto(&out.Metrics)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppMonitorOperation.
func (in *AppMonitorOperation) DeepCopy() *AppMonitorOperation {
	if in == nil {
		return nil
	}
	out := new(AppMonitorOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LifecycleOperation) DeepCopyInto(out *LifecycleOperation) {
	*out = *in
	in.Init.DeepCopyInto(&out.Init)
	in.Quit.DeepCopyInto(&out.Quit)
	in.Funcs.DeepCopyInto(&out.Funcs)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LifecycleOperation.
func (in *LifecycleOperation) DeepCopy() *LifecycleOperation {
	if in == nil {
		return nil
	}
	out := new(LifecycleOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MetricsSrc) DeepCopyInto(out *MetricsSrc) {
	*out = *in
	if in.Cmd != nil {
		in, out := &in.Cmd, &out.Cmd
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MetricsSrc.
func (in *MetricsSrc) DeepCopy() *MetricsSrc {
	if in == nil {
		return nil
	}
	out := new(MetricsSrc)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Operation) DeepCopyInto(out *Operation) {
	*out = *in
	if in.Cmd != nil {
		in, out := &in.Cmd, &out.Cmd
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Operation.
func (in *Operation) DeepCopy() *Operation {
	if in == nil {
		return nil
	}
	out := new(Operation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperationSet) DeepCopyInto(out *OperationSet) {
	*out = *in
	if in.Cluster != nil {
		in, out := &in.Cluster, &out.Cluster
		*out = make([]Operation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Pod != nil {
		in, out := &in.Pod, &out.Pod
		*out = make([]Operation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperationSet.
func (in *OperationSet) DeepCopy() *OperationSet {
	if in == nil {
		return nil
	}
	out := new(OperationSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resources) DeepCopyInto(out *Resources) {
	*out = *in
	if in.Agent != nil {
		in, out := &in.Agent, &out.Agent
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.ResourceRequirements)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Server != nil {
		in, out := &in.Server, &out.Server
		if *in == nil {
			*out = nil
		} else {
			*out = new(v1.ResourceRequirements)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resources.
func (in *Resources) DeepCopy() *Resources {
	if in == nil {
		return nil
	}
	out := new(Resources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Statefulguardian) DeepCopyInto(out *Statefulguardian) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Statefulguardian.
func (in *Statefulguardian) DeepCopy() *Statefulguardian {
	if in == nil {
		return nil
	}
	out := new(Statefulguardian)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Statefulguardian) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatefulguardianCondition) DeepCopyInto(out *StatefulguardianCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulguardianCondition.
func (in *StatefulguardianCondition) DeepCopy() *StatefulguardianCondition {
	if in == nil {
		return nil
	}
	out := new(StatefulguardianCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatefulguardianList) DeepCopyInto(out *StatefulguardianList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Statefulguardian, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulguardianList.
func (in *StatefulguardianList) DeepCopy() *StatefulguardianList {
	if in == nil {
		return nil
	}
	out := new(StatefulguardianList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StatefulguardianList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatefulguardianSpec) DeepCopyInto(out *StatefulguardianSpec) {
	*out = *in
	in.Lifecycle.DeepCopyInto(&out.Lifecycle)
	in.AppMonitor.DeepCopyInto(&out.AppMonitor)
	in.App.DeepCopyInto(&out.App)
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		if *in == nil {
			*out = nil
		} else {
			*out = new(Resources)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulguardianSpec.
func (in *StatefulguardianSpec) DeepCopy() *StatefulguardianSpec {
	if in == nil {
		return nil
	}
	out := new(StatefulguardianSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatefulguardianStatus) DeepCopyInto(out *StatefulguardianStatus) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]StatefulguardianCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatefulguardianStatus.
func (in *StatefulguardianStatus) DeepCopy() *StatefulguardianStatus {
	if in == nil {
		return nil
	}
	out := new(StatefulguardianStatus)
	in.DeepCopyInto(out)
	return out
}
