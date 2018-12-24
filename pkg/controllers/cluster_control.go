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
	"fmt"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/util/retry"

	v1alpha1 "statefulguardian/pkg/apis/statefulguardian/v1alpha1"
	clientset "statefulguardian/pkg/generated/clientset/versioned"
	listersv1alpha1 "statefulguardian/pkg/generated/listers/statefulguardian/v1alpha1"
)

type sgUpdaterInterface interface {
	UpdateSgStatus(sg *v1alpha1.Statefulguardian, status *v1alpha1.StatefulguardianStatus) error
	UpdateSgLabels(sg *v1alpha1.Statefulguardian, lbls labels.Set) error
}

type sgUpdater struct {
	client clientset.Interface
	lister listersv1alpha1.StatefulguardianLister
}

func newSgUpdater(client clientset.Interface, lister listersv1alpha1.StatefulguardianLister) sgUpdaterInterface {
	return &sgUpdater{client: client, lister: lister}
}

func (csu *sgUpdater) UpdateSgStatus(sg *v1alpha1.Statefulguardian, status *v1alpha1.StatefulguardianStatus) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		sg.Status = *status
		_, updateErr := csu.client.StatefulguardianV1alpha1().Statefulguardians(sg.Namespace).Update(sg)
		if updateErr == nil {
			return nil
		}

		updated, err := csu.lister.Statefulguardians(sg.Namespace).Get(sg.Name)
		if err != nil {
			glog.Errorf("Error getting updated Statefulguardian %s/%s: %v", sg.Namespace, sg.Name, err)
			return err
		}

		// Copy the Statefulguardian so we don't mutate the cache.
		sg = updated.DeepCopy()
		return updateErr
	})
}

func (csu *sgUpdater) UpdateSgLabels(sg *v1alpha1.Statefulguardian, lbls labels.Set) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		sg.Labels = labels.Merge(labels.Set(sg.Labels), lbls)
		_, updateErr := csu.client.StatefulguardianV1alpha1().Statefulguardians(sg.Namespace).Update(sg)
		if updateErr == nil {
			return nil
		}

		key := fmt.Sprintf("%s/%s", sg.GetNamespace(), sg.GetName())
		glog.V(4).Infof("Conflict updating Statefulguardian labels. Getting updated Statefulguardian %s from cache...", key)

		updated, err := csu.lister.Statefulguardians(sg.Namespace).Get(sg.Name)
		if err != nil {
			glog.Errorf("Error getting updated Statefulguardian %q: %v", key, err)
			return err
		}

		// Copy the Statefulguardian so we don't mutate the cache.
		sg = updated.DeepCopy()
		return updateErr
	})
}
