// Copyright 2025 anza-labs contributors.
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KinkClusterSpec defines the desired state of KinkCluster.
type KinkClusterSpec struct {
}

// KinkClusterStatus defines the observed state of KinkCluster.
type KinkClusterStatus struct {
	// Ready denotes that the kink cluster infrastructure is fully provisioned.
	// +optional
	Ready bool `json:"ready"`

	// Conditions defines current service state of the KinkCluster.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KinkCluster is the Schema for the kinkclusters API.
type KinkCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KinkClusterSpec   `json:"spec,omitempty"`
	Status KinkClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KinkClusterList contains a list of KinkCluster.
type KinkClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KinkCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KinkCluster{}, &KinkClusterList{})
}
