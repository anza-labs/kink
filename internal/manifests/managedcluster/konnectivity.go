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

package managedcluster

import (
	"fmt"
	"path"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/manifests/manifestutils"
	"github.com/anza-labs/kink/internal/naming"
	"github.com/anza-labs/kink/version"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	konnectivityTokenPath = "/var/run/secrets/tokens"
	konnectivityTokenFile = "konnectivity-agent-token"
)

type Konnectivity struct {
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

func (b *Konnectivity) Build() []client.Object {
	objects := []client.Object{}

	ds := b.DaemonSet()
	objects = append(objects, ds)

	crb := b.ClusterRoleBinding()
	objects = append(objects, crb)

	sa := b.ServiceAccount()
	objects = append(objects, sa)

	return objects
}

func (b *Konnectivity) DaemonSet() *appsv1.DaemonSet {
	name := naming.KonnectivityAgent()

	image := b.KinkControlPlane.Spec.KonnectivityAgent.Image
	if image == "" {
		image = version.KonnectivityAgent()
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentKonnectivity, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentKonnectivity,
		ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)
	podAnnotations := manifestutils.PodAnnotations(b.KinkControlPlane, nil)

	podSpec := corev1.PodSpec{
		PriorityClassName: "system-cluster-critical",
		Tolerations: []corev1.Toleration{
			{
				Key:      "CriticalAddonsOnly",
				Operator: corev1.TolerationOpExists,
			},
		},
		HostNetwork:        true,
		Containers:         []corev1.Container{b.container(image)},
		ServiceAccountName: "konnectivity-agent",
		Volumes:            b.volumes(),
	}

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   "kube-system",
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: podAnnotations,
				},
				Spec: podSpec,
			},
		},
	}
}

func (b *Konnectivity) volumes() []corev1.Volume {
	return []corev1.Volume{
		{
			Name: "konnectivity-token",
			VolumeSource: corev1.VolumeSource{
				Projected: &corev1.ProjectedVolumeSource{
					Sources: []corev1.VolumeProjection{
						{
							ServiceAccountToken: &corev1.ServiceAccountTokenProjection{
								Path:     konnectivityTokenFile,
								Audience: "system:konnectivity-server",
							},
						},
					},
				},
			},
		},
	}
}

func (b *Konnectivity) volumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "konnectivity-token",
			ReadOnly:  true,
			MountPath: "/var/run/secrets/tokens",
		},
	}
}

func (b *Konnectivity) container(image string) corev1.Container {
	cfg := b.KinkControlPlane.Spec.KonnectivityAgent
	resources := cfg.Resources
	verbosity := cfg.Verbosity

	args := map[string]string{
		"v":                          fmt.Sprint(verbosity),
		"logtostderr":                "true",
		"ca-cert":                    "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		"proxy-server-host":          b.KinkControlPlane.Status.IP,
		"proxy-server-port":          "8132",
		"service-account-token-path": path.Join(konnectivityTokenPath, konnectivityTokenFile),
	}

	return corev1.Container{
		Name:      naming.KonnectivityContainer(),
		Image:     image,
		Command:   []string{"/proxy-agent"},
		Args:      manifestutils.BuildKubernetesArgs(args),
		Resources: resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "health",
				ContainerPort: 8093,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts:    b.volumeMounts(),
		ImagePullPolicy: cfg.ImagePullPolicy,
		SecurityContext: nil,
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/healthz",
					Port:   intstr.FromString("health"),
					Scheme: corev1.URISchemeHTTP,
				},
			},
			InitialDelaySeconds: 10,
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/readyz",
					Port:   intstr.FromString("health"),
					Scheme: corev1.URISchemeHTTP,
				},
			},
			InitialDelaySeconds: 30,
		},
	}
}

func (b *Konnectivity) ClusterRoleBinding() *rbacv1.ClusterRoleBinding {
	name := "system:konnectivity-server"

	image := b.KinkControlPlane.Spec.KonnectivityAgent.Image
	if image == "" {
		image = version.KonnectivityAgent()
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentKonnectivity, ConceptControlPlane,
		nil,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	labels["kubernetes.io/cluster-service"] = "true"
	labels["addonmanager.kubernetes.io/mode"] = "Reconcile"

	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "ClusterRole",
			Name:     "system:auth-delegator",
		},
		Subjects: []rbacv1.Subject{
			{
				APIGroup: rbacv1.GroupName,
				Kind:     "User",
				Name:     name,
			},
		},
	}
}

func (b *Konnectivity) ServiceAccount() *corev1.ServiceAccount {
	name := naming.KonnectivityAgent()

	image := b.KinkControlPlane.Spec.KonnectivityAgent.Image
	if image == "" {
		image = version.KonnectivityAgent()
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentKonnectivity, ConceptControlPlane,
		nil,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	labels["kubernetes.io/cluster-service"] = "true"
	labels["addonmanager.kubernetes.io/mode"] = "Reconcile"

	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   "kube-system",
			Labels:      labels,
			Annotations: annotations,
		},
	}
}
