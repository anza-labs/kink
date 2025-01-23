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

// KinkMachineTemplateSpec defines the desired state of KinkMachineTemplate.
type KinkMachineTemplateSpec struct {
	Template KinkMachineTemplateResource `json:"template"`
}

type KinkMachineTemplateResource struct {
	// +optional
	ObjectMeta metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec       KinkMachineSpec   `json:"spec"`
}

// KinkMachineTemplateStatus defines the observed state of KinkMachineTemplate.
type KinkMachineTemplateStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KinkMachineTemplate is the Schema for the kinkmachinetemplates API.
type KinkMachineTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KinkMachineTemplateSpec   `json:"spec,omitempty"`
	Status KinkMachineTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KinkMachineTemplateList contains a list of KinkMachineTemplate.
type KinkMachineTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KinkMachineTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KinkMachineTemplate{}, &KinkMachineTemplateList{})
}
