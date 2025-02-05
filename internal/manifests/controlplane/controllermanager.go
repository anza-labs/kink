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

package controlplane

import (
	"fmt"
	"path"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/manifests/manifestutils"
	"github.com/anza-labs/kink/internal/naming"
	"github.com/anza-labs/kink/version"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	frontProxyPKIPath = "/etc/pki/front-proxy"
	frontProxyCAFile  = "ca.crt"
)

type ControllerManager struct {
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

func (b *ControllerManager) Build() ([]client.Object, error) {
	objects := []client.Object{}

	svc, err := b.Service()
	if err != nil {
		return nil, fmt.Errorf("failed to build Service: %w", err)
	}
	objects = append(objects, svc)

	depl, err := b.Deployment()
	if err != nil {
		return nil, fmt.Errorf("failed to build Deployment: %w", err)
	}
	objects = append(objects, depl)

	return objects, nil
}

func (b *ControllerManager) Deployment() (*appsv1.Deployment, error) {
	name := naming.ControllerManager(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.ControllerManager.Image,
		b.KinkControlPlane.Spec.Version,
		version.ControllerManager(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentControllerManager, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentControllerManager, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)
	podAnnotations := manifestutils.PodAnnotations(b.KinkControlPlane, nil)

	ha := false
	replicas := b.KinkControlPlane.Spec.Replicas
	if replicas == nil {
		replicas = ptr.To[int32](1)
	}
	if *replicas > 1 {
		ha = true
	}

	podSpec := corev1.PodSpec{
		Affinity:         manifestutils.Affinity(b.KinkControlPlane),
		Containers:       []corev1.Container{b.container(image, ha)},
		Volumes:          b.volumes(),
		ImagePullSecrets: b.KinkControlPlane.Spec.ImagePullSecrets,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Replicas: replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: podAnnotations,
				},
				Spec: podSpec,
			},
		},
	}, nil
}

func (b *ControllerManager) Service() (*corev1.Service, error) {
	name := naming.ControllerManager(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.ControllerManager.Image,
		b.KinkControlPlane.Spec.Version,
		version.ControllerManager(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentControllerManager, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentControllerManager, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Selector: selectorLabels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "self",
					Port:       10257,
					TargetPort: intstr.FromString("self"),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}, nil
}

func (b *ControllerManager) volumes() []corev1.Volume {
	name := naming.ControllerManager(b.KinkControlPlane.Name)

	return []corev1.Volume{
		{
			Name: "kubeconfig",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  naming.Kubeconfig(name),
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
		{
			Name: "root-ca",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  naming.ClusterCA(b.KinkControlPlane.Name),
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
		{
			Name: "front-proxy-ca",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  naming.FrontProxyCA(b.KinkControlPlane.Name),
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
		{
			Name: "service-accounts-cert",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  naming.ServiceAccountCertificate(b.KinkControlPlane.Name),
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
	}
}

func (b *ControllerManager) volumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "kubeconfig",
			ReadOnly:  true,
			MountPath: kubeconfigPath,
		},
		{
			Name:      "root-ca",
			ReadOnly:  true,
			MountPath: rootPKIPath,
		},
		{
			Name:      "front-proxy-ca",
			ReadOnly:  true,
			MountPath: frontProxyPKIPath,
		},
		{
			Name:      "service-accounts-cert",
			ReadOnly:  true,
			MountPath: serviceAccountsPKIPath,
		},
	}
}

func (b *ControllerManager) container(image string, ha bool) corev1.Container {
	cfg := b.KinkControlPlane.Spec.ControllerManager
	resources := cfg.Resources
	verbosity := cfg.Verbosity

	args := map[string]string{
		"v":                                fmt.Sprint(verbosity),
		"leader-elect":                     fmt.Sprint(ha),
		"kubeconfig":                       path.Join(kubeconfigPath, kubeconfigName),
		"authorization-kubeconfig":         path.Join(kubeconfigPath, kubeconfigName),
		"authentication-kubeconfig":        path.Join(kubeconfigPath, kubeconfigName),
		"cluster-signing-cert-file":        path.Join(rootPKIPath, rootCertFile),
		"cluster-signing-key-file":         path.Join(rootPKIPath, rootKeyFile),
		"service-account-private-key-file": path.Join(serviceAccountsPKIPath, serviceAccountsKeyFile),
		"requestheader-client-ca-file":     path.Join(frontProxyPKIPath, frontProxyCAFile),
		"controllers":                      "*,bootstrapsigner,tokencleaner",
		"use-service-account-credentials":  "true",
		"cluster-cidr":                     "10.200.0.0/16", // TODO
		"service-cluster-ip-range":         "10.32.0.0/24",  // TODO
	}
	for arg, value := range cfg.ExtraArgs {
		if _, ok := args[arg]; !ok {
			args[arg] = value
		}
	}

	return corev1.Container{
		Name:      naming.ControllerManagerContainer(),
		Image:     image,
		Command:   []string{"kube-controller-manager"},
		Args:      buildArgs(args),
		Resources: resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "self",
				ContainerPort: 10257,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts:    b.volumeMounts(),
		ImagePullPolicy: b.KinkControlPlane.Spec.APIServer.ImagePullPolicy,
		SecurityContext: nil,
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path:   "/healthz",
					Port:   intstr.FromString("self"),
					Scheme: corev1.URISchemeHTTPS,
				},
			},
		},
	}
}
