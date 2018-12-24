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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// StatefulguardianLister helps list Statefulguardians.
type StatefulguardianLister interface {
	// List lists all Statefulguardians in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Statefulguardian, err error)
	// Statefulguardians returns an object that can list and get Statefulguardians.
	Statefulguardians(namespace string) StatefulguardianNamespaceLister
	StatefulguardianListerExpansion
}

// statefulguardianLister implements the StatefulguardianLister interface.
type statefulguardianLister struct {
	indexer cache.Indexer
}

// NewStatefulguardianLister returns a new StatefulguardianLister.
func NewStatefulguardianLister(indexer cache.Indexer) StatefulguardianLister {
	return &statefulguardianLister{indexer: indexer}
}

// List lists all Statefulguardians in the indexer.
func (s *statefulguardianLister) List(selector labels.Selector) (ret []*v1alpha1.Statefulguardian, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Statefulguardian))
	})
	return ret, err
}

// Statefulguardians returns an object that can list and get Statefulguardians.
func (s *statefulguardianLister) Statefulguardians(namespace string) StatefulguardianNamespaceLister {
	return statefulguardianNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// StatefulguardianNamespaceLister helps list and get Statefulguardians.
type StatefulguardianNamespaceLister interface {
	// List lists all Statefulguardians in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Statefulguardian, err error)
	// Get retrieves the Statefulguardian from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Statefulguardian, error)
	StatefulguardianNamespaceListerExpansion
}

// statefulguardianNamespaceLister implements the StatefulguardianNamespaceLister
// interface.
type statefulguardianNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Statefulguardians in the indexer for a given namespace.
func (s statefulguardianNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Statefulguardian, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Statefulguardian))
	})
	return ret, err
}

// Get retrieves the Statefulguardian from the indexer for a given namespace and name.
func (s statefulguardianNamespaceLister) Get(name string) (*v1alpha1.Statefulguardian, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("statefulguardian"), name)
	}
	return obj.(*v1alpha1.Statefulguardian), nil
}
