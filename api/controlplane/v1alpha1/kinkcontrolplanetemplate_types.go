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

// KinkControlPlaneTemplateSpec defines the desired state of KinkControlPlaneTemplate.
type KinkControlPlaneTemplateSpec struct {
	Template KinkControlPlaneTemplateResource `json:"template"`
}

type KinkControlPlaneTemplateResource struct {
	// +optional
	ObjectMeta metav1.ObjectMeta    `json:"metadata,omitempty"`
	Spec       KinkControlPlaneSpec `json:"spec"`
}

// KinkControlPlaneTemplateStatus defines the observed state of KinkControlPlaneTemplate.
type KinkControlPlaneTemplateStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KinkControlPlaneTemplate is the Schema for the kinkcontrolplanetemplates API.
type KinkControlPlaneTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KinkControlPlaneTemplateSpec   `json:"spec,omitempty"`
	Status KinkControlPlaneTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KinkControlPlaneTemplateList contains a list of KinkControlPlaneTemplate.
type KinkControlPlaneTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KinkControlPlaneTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KinkControlPlaneTemplate{}, &KinkControlPlaneTemplateList{})
}
