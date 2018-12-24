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
	v1alpha1 "statefulguardian/pkg/apis/statefulguardian/v1alpha1"
	scheme "statefulguardian/pkg/generated/clientset/versioned/scheme"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// StatefulguardiansGetter has a method to return a StatefulguardianInterface.
// A group's client should implement this interface.
type StatefulguardiansGetter interface {
	Statefulguardians(namespace string) StatefulguardianInterface
}

// StatefulguardianInterface has methods to work with Statefulguardian resources.
type StatefulguardianInterface interface {
	Create(*v1alpha1.Statefulguardian) (*v1alpha1.Statefulguardian, error)
	Update(*v1alpha1.Statefulguardian) (*v1alpha1.Statefulguardian, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Statefulguardian, error)
	List(opts v1.ListOptions) (*v1alpha1.StatefulguardianList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Statefulguardian, err error)
	StatefulguardianExpansion
}

// statefulguardians implements StatefulguardianInterface
type statefulguardians struct {
	client rest.Interface
	ns     string
}

// newStatefulguardians returns a Statefulguardians
func newStatefulguardians(c *StatefulguardianV1alpha1Client, namespace string) *statefulguardians {
	return &statefulguardians{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the statefulguardian, and returns the corresponding statefulguardian object, and an error if there is any.
func (c *statefulguardians) Get(name string, options v1.GetOptions) (result *v1alpha1.Statefulguardian, err error) {
	result = &v1alpha1.Statefulguardian{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("statefulguardians").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Statefulguardians that match those selectors.
func (c *statefulguardians) List(opts v1.ListOptions) (result *v1alpha1.StatefulguardianList, err error) {
	result = &v1alpha1.StatefulguardianList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("statefulguardians").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested statefulguardians.
func (c *statefulguardians) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("statefulguardians").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a statefulguardian and creates it.  Returns the server's representation of the statefulguardian, and an error, if there is any.
func (c *statefulguardians) Create(statefulguardian *v1alpha1.Statefulguardian) (result *v1alpha1.Statefulguardian, err error) {
	result = &v1alpha1.Statefulguardian{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("statefulguardians").
		Body(statefulguardian).
		Do().
		Into(result)
	return
}

// Update takes the representation of a statefulguardian and updates it. Returns the server's representation of the statefulguardian, and an error, if there is any.
func (c *statefulguardians) Update(statefulguardian *v1alpha1.Statefulguardian) (result *v1alpha1.Statefulguardian, err error) {
	result = &v1alpha1.Statefulguardian{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("statefulguardians").
		Name(statefulguardian.Name).
		Body(statefulguardian).
		Do().
		Into(result)
	return
}

// Delete takes name of the statefulguardian and deletes it. Returns an error if one occurs.
func (c *statefulguardians) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("statefulguardians").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *statefulguardians) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("statefulguardians").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched statefulguardian.
func (c *statefulguardians) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Statefulguardian, err error) {
	result = &v1alpha1.Statefulguardian{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("statefulguardians").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
