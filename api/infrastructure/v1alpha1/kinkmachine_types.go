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
	kinkcorev1alpha1 "github.com/anza-labs/kink/api/core/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KinkMachineSpec defines the desired state of KinkMachine.
type KinkMachineSpec struct {
	kinkcorev1alpha1.Container `json:",inline"`

	// ProviderID must match the provider ID as seen on the node object corresponding to this machine.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Affinity specifies the scheduling constraints for Pods.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Persistence specifies volume configuration for Kine data persistence.
	// Defaults to EmptyDir.
	// +optional
	Persistence *kinkcorev1alpha1.Persistence `json:"persistence,omitempty"`
}

// KinkMachineStatus defines the observed state of KinkMachine.
type KinkMachineStatus struct {
	// Ready denotes that the kink machine infrastructure is fully provisioned.
	// +optional
	Ready bool `json:"ready"`

	// Conditions defines current service state of the KinkMachine.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
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
