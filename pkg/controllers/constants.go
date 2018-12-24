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

// ClusterLabel is applied to all components of a cluster
const SgLabel = "v1alpha1.ck.com/statefulguardian"
const SgNsLabel = "v1alpha1.ck.com/statefulguardianNamespace"

const SgVersionLabel = "v1alpha1.ck.com/version"

// LabelClusterRole specifies the role of a Pod within a Cluster.
const LabelSgRole = "v1alpha1.ck.com/role"

// ClusterRolePrimary denotes a primary InnoDB cluster member.
const SgRolePrimary = "primary"

// ClusterRoleSecondary denotes a secondary InnoDB cluster member.
const SgRoleSecondary = "secondary"
