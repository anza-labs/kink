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

// KinkControlPlaneSpec defines the desired state of KinkControlPlane.
type KinkControlPlaneSpec struct {
	// Version defines the desired Kubernetes version for the control plane.
	// The value must be a valid semantic version; also if the value provided by the user
	// does not start with the v prefix, it must be added.
	Version string `json:"version"`

	// DNSName specifies the cluster endpoint, most likely backed by Ingress.
	// +default="localhost"
	// +kubebuilder:default="localhost"
	DNSName string `json:"dnsName"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Affinity specifies the scheduling constraints for Pods.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// APIServer defines the configuration for the Kubernetes API server.
	APIServer APIServer `json:"apiServer"`

	// Kine defines the configuration for the Kine component.
	Kine Kine `json:"kine"`

	// Scheduler defines the configuration for the Kubernetes scheduler.
	Scheduler Scheduler `json:"scheduler"`

	// ControllerManager defines the configuration for the Kubernetes controller manager.
	ControllerManager ControllerManager `json:"controllerManager"`
}

// Kine represents ETCD-shim container.
type Kine struct {
	kinkcorev1alpha1.Container `json:",inline"`

	// Persistence specifies volume configuration for Kine data persistence.
	// Defaults to EmptyDir.
	// +optional
	Persistence *kinkcorev1alpha1.Persistence `json:"persistence,omitempty"`
}

// Scheduler represents a Kubernetes scheduler.
//
// Image:
//   - If specified image contains tag or sha, those are ignored.
//   - Defaults to registry.k8s.io/kube-scheduler
type Scheduler struct {
	KubeComponent `json:",inline"`
}

// ControllerManager represents a Kubernetes controller manager.
//
// Image:
//   - If specified image contains tag or sha, those are ignored.
//   - Defaults to registry.k8s.io/kube-controller-manager
type ControllerManager struct {
	KubeComponent `json:",inline"`
}

// APIServer represents a Kubernetes API server.
//
// Image:
//   - If specified image contains tag or sha, those are ignored.
//   - Defaults to registry.k8s.io/kube-apiserver
type APIServer struct {
	KubeComponent `json:",inline"`
}

// KubeComponent defines the base configuration for Kink control plane components.
type KubeComponent struct {
	kinkcorev1alpha1.Container `json:",inline"`

	// Number of desired pods. Defaults to 1.
	// +optional
	// +default=1
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=5
	Replicas int32 `json:"replicas,omitempty"`

	// Verbosity specifies the log verbosity level for the container. Valid values range from 0 (silent) to 10 (most verbose).
	// +optional
	// +default=4
	// +kubebuilder:default=4
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=10
	Verbosity uint8 `json:"verbosity"`

	// ExtraArgs defines additional arguments to be passed to the container executable.
	// +optional
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`
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
