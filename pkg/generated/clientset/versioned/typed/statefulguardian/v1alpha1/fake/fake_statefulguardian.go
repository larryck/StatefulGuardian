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

package fake

import (
	v1alpha1 "statefulguardian/pkg/apis/statefulguardian/v1alpha1"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeStatefulguardians implements StatefulguardianInterface
type FakeStatefulguardians struct {
	Fake *FakeStatefulguardianV1alpha1
	ns   string
}

var statefulguardiansResource = schema.GroupVersionResource{Group: "ck.com", Version: "v1alpha1", Resource: "statefulguardians"}

var statefulguardiansKind = schema.GroupVersionKind{Group: "ck.com", Version: "v1alpha1", Kind: "Statefulguardian"}

// Get takes name of the statefulguardian, and returns the corresponding statefulguardian object, and an error if there is any.
func (c *FakeStatefulguardians) Get(name string, options v1.GetOptions) (result *v1alpha1.Statefulguardian, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(statefulguardiansResource, c.ns, name), &v1alpha1.Statefulguardian{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Statefulguardian), err
}

// List takes label and field selectors, and returns the list of Statefulguardians that match those selectors.
func (c *FakeStatefulguardians) List(opts v1.ListOptions) (result *v1alpha1.StatefulguardianList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(statefulguardiansResource, statefulguardiansKind, c.ns, opts), &v1alpha1.StatefulguardianList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.StatefulguardianList{}
	for _, item := range obj.(*v1alpha1.StatefulguardianList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested statefulguardians.
func (c *FakeStatefulguardians) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(statefulguardiansResource, c.ns, opts))

}

// Create takes the representation of a statefulguardian and creates it.  Returns the server's representation of the statefulguardian, and an error, if there is any.
func (c *FakeStatefulguardians) Create(statefulguardian *v1alpha1.Statefulguardian) (result *v1alpha1.Statefulguardian, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(statefulguardiansResource, c.ns, statefulguardian), &v1alpha1.Statefulguardian{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Statefulguardian), err
}

// Update takes the representation of a statefulguardian and updates it. Returns the server's representation of the statefulguardian, and an error, if there is any.
func (c *FakeStatefulguardians) Update(statefulguardian *v1alpha1.Statefulguardian) (result *v1alpha1.Statefulguardian, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(statefulguardiansResource, c.ns, statefulguardian), &v1alpha1.Statefulguardian{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Statefulguardian), err
}

// Delete takes name of the statefulguardian and deletes it. Returns an error if one occurs.
func (c *FakeStatefulguardians) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(statefulguardiansResource, c.ns, name), &v1alpha1.Statefulguardian{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeStatefulguardians) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(statefulguardiansResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.StatefulguardianList{})
	return err
}

// Patch applies the patch and returns the patched statefulguardian.
func (c *FakeStatefulguardians) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Statefulguardian, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(statefulguardiansResource, c.ns, name, data, subresources...), &v1alpha1.Statefulguardian{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Statefulguardian), err
}
