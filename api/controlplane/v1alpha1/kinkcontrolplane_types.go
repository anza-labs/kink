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
	"net"

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

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// Based on it, an ingress will be provisioned.
	// +optional
	ControlPlaneEndpoint *APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	// Number of desired ControlPlane replicas. Defaults to 1.
	// +optional
	// +default=1
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=5
	Replicas *int32 `json:"replicas,omitempty"`

	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use
	// for pulling any of the images used by KinkControlPlane. If specified, these secrets will
	// be passed to individual puller implementations for them to use.
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

// APIEndpoint represents a reachable Kubernetes API endpoint.
type APIEndpoint struct {
	// host is the hostname on which the API server is serving.
	Host HostnameOrIP `json:"host"`

	// port is the port on which the API server is serving.
	Port int32 `json:"port"`

	// ServiceType
	// +optional
	// +default="ClusterIP"
	// +kubebuilder:default="ClusterIP"
	ServiceType corev1.ServiceType `json:"serviceType"`

	// Gateway.
	// +optional
	Gateway *Gateway `json:"gateway,omitempty"`

	// Ingress.
	// +optional
	Ingress *Ingress `json:"ingress,omitempty"`
}

// HostnameOrIP.
type HostnameOrIP string

// IsIP checks if the value is a valid IP address.
func (h HostnameOrIP) IsIP() bool {
	return net.ParseIP(string(h)) != nil
}

// IsHostname checks if the value is a valid hostname.
func (h HostnameOrIP) IsHostname() bool {
	return !h.IsIP()
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

// Gateway.
type Gateway struct {
	// GatewayClassName used for this Gateway. This is the name of a
	// GatewayClass resource.
	GatewayClassName string `json:"gatewayClassName"`
}

// Ingress.
type Ingress struct {
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// GatewayClassName used for this Gateway. This is the name of a
	// GatewayClass resource.
	IngressClassName string `json:"ingressClassName"`
}

// KinkControlPlaneStatus defines the observed state of KinkControlPlane.
type KinkControlPlaneStatus struct {
	// Version represents the minimum Kubernetes version for the control plane replicas
	// in the cluster.
	// +optional
	Version *string `json:"version,omitempty"`

	// Selector is the label selector in string format to avoid introspection
	// by clients, and is used to provide the CRD-based integration for the
	// scale subresource and additional integrations for things like kubectl
	// describe. The string will be in the same format as the query-param syntax.
	// More info about label selectors: http://kubernetes.io/docs/user-guide/labels#label-selectors
	// +optional
	Selector string `json:"selector,omitempty"`

	// Replicas is the total number of replicas targeted by this control plane
	// (their labels match the selector).
	// +optional
	Replicas int32 `json:"replicas"`

	// UpdatedReplicas is the total number of replicas targeted by this control plane
	// that have the desired template spec.
	// +optional
	UpdatedReplicas int32 `json:"updatedReplicas"`

	// ReadyReplicas is the total number of fully running and ready control plane replicas.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas"`

	// UnavailableReplicas is the total number of unavailable replicas targeted by this control plane.
	// This is the total number of replicas that are still required for the deployment to have 100% available capacity.
	// They may either be replicas that are running but not yet ready or replicas
	// that still have not been created.
	// +optional
	UnavailableReplicas int32 `json:"unavailableReplicas"`

	// Conditions defines current service state of the KinkControlPlane.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Initialized denotes that the kink control plane API Server is initialized and thus
	// it can accept requests.
	// +optional
	Initialized bool `json:"initialized"`

	// Ready denotes that the kink control plane is ready to serve requests.
	// +optional
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=kinkcontrolplanes,scope=Namespaced,categories=cluster-api,shortName=kink
// +kubebuilder:conversion:hub
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name"
// +kubebuilder:printcolumn:name="Replicas",type="string",JSONPath=".status.replicas"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.readyReplicas"
// +kubebuilder:printcolumn:name="Updated",type="string",JSONPath=".status.updatedReplicas"
// +kubebuilder:printcolumn:name="Unavailable",type="string",JSONPath=".status.unavailableReplicas"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".status.version"

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
