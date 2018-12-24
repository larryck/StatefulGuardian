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
	statefulguardian_v1alpha1 "statefulguardian/pkg/apis/statefulguardian/v1alpha1"
	versioned "statefulguardian/pkg/generated/clientset/versioned"
	internalinterfaces "statefulguardian/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "statefulguardian/pkg/generated/listers/statefulguardian/v1alpha1"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// StatefulguardianInformer provides access to a shared informer and lister for
// Statefulguardians.
type StatefulguardianInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.StatefulguardianLister
}

type statefulguardianInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewStatefulguardianInformer constructs a new informer for Statefulguardian type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewStatefulguardianInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredStatefulguardianInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredStatefulguardianInformer constructs a new informer for Statefulguardian type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredStatefulguardianInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.StatefulguardianV1alpha1().Statefulguardians(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.StatefulguardianV1alpha1().Statefulguardians(namespace).Watch(options)
			},
		},
		&statefulguardian_v1alpha1.Statefulguardian{},
		resyncPeriod,
		indexers,
	)
}

func (f *statefulguardianInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredStatefulguardianInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *statefulguardianInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&statefulguardian_v1alpha1.Statefulguardian{}, f.defaultInformer)
}

func (f *statefulguardianInformer) Lister() v1alpha1.StatefulguardianLister {
	return v1alpha1.NewStatefulguardianLister(f.Informer().GetIndexer())
}
