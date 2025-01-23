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

// KinkClusterTemplateSpec defines the desired state of KinkClusterTemplate.
type KinkClusterTemplateSpec struct {
	Template KinkClusterTemplateResource `json:"template"`
}

type KinkClusterTemplateResource struct {
	// +optional
	ObjectMeta metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec       KinkClusterSpec   `json:"spec"`
}

// KinkClusterTemplateStatus defines the observed state of KinkClusterTemplate.
type KinkClusterTemplateStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KinkClusterTemplate is the Schema for the kinkclustertemplates API.
type KinkClusterTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KinkClusterTemplateSpec   `json:"spec,omitempty"`
	Status KinkClusterTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KinkClusterTemplateList contains a list of KinkClusterTemplate.
type KinkClusterTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KinkClusterTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KinkClusterTemplate{}, &KinkClusterTemplateList{})
}
