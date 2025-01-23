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

// KinkControlPlaneSpec defines the desired state of KinkControlPlane.
type KinkControlPlaneSpec struct {
	// Version defines the desired Kubernetes version for the control plane.
	// The value must be a valid semantic version; also if the value provided by the user
	// does not start with the v prefix, it must be added.
	Version string `json:"version"`
}

// KinkControlPlaneStatus defines the observed state of KinkControlPlane.
type KinkControlPlaneStatus struct {
	// Version represents the minimum Kubernetes version for the control plane machines
	// in the cluster.
	// +optional
	Version *string `json:"version,omitempty"`

	// Initialized denotes that the kink control plane API Server is initialized and thus
	// it can accept requests.
	Initialized bool `json:"initialized"`

	// Ready denotes that the kink control plane is ready to serve requests.
	// +optional
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KinkControlPlane is the Schema for the kinkcontrolplanes API.
type KinkControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KinkControlPlaneSpec   `json:"spec,omitempty"`
	Status KinkControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KinkControlPlaneList contains a list of KinkControlPlane.
type KinkControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KinkControlPlane `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KinkControlPlane{}, &KinkControlPlaneList{})
}
