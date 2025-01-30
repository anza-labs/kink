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

// +kubebuilder:object:generate=true
package api

import (
	corev1 "k8s.io/api/core/v1"
)

// Container defines the base container configuration for Kink control plane components.
type Container struct {
	// Image specifies the container image to use.
	// +optional
	Image string `json:"image,omitempty"`

	// Image pull policy. One of Always, Never, IfNotPresent.
	// +optional
	ImagePullPolicy corev1.PullPolicy `json:"imagePullPolicy,omitempty"`

	// Resources describes the compute resource requirements for the container.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// Container defines the minimum persistence configuration. It always defaults to EmptyDir.
type Persistence struct {
	// EmptyDir represents a temporary directory that shares a pod's lifetime.
	// +optional
	EmptyDir *corev1.EmptyDirVolumeSource `json:"emptyDir,omitempty"`

	// PersistentVolumeClaimVolumeSource represents a reference to a
	// PersistentVolumeClaim in the same namespace.
	// +optional
	PersistentVolumeClaim *corev1.PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`
}
