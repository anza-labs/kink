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

// KinkMachineSpec defines the desired state of KinkMachine.
type KinkMachineSpec struct {
	// providerID must match the provider ID as seen on the node object corresponding to this machine.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`
}

// KinkMachineStatus defines the observed state of KinkMachine.
type KinkMachineStatus struct {
	// Ready denotes that the kink machine infrastructure is fully provisioned.
	// +optional
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KinkMachine is the Schema for the kinkmachines API.
type KinkMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KinkMachineSpec   `json:"spec,omitempty"`
	Status KinkMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KinkMachineList contains a list of KinkMachine.
type KinkMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KinkMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KinkMachine{}, &KinkMachineList{})
}
