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

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/manifests/manifestutils"
	"github.com/anza-labs/kink/internal/naming"
	"github.com/anza-labs/kink/version"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type APIServer struct {
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

func (b *APIServer) Build() ([]runtime.Object, error) {
	objects := []runtime.Object{}

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

func (b *APIServer) Deployment() (*appsv1.Deployment, error) {
	name := naming.APIServer(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.APIServer.Image,
		b.KinkControlPlane.Spec.Version,
		version.APIServer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentAPIServer, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(b.KinkControlPlane.ObjectMeta, ComponentAPIServer, ConceptControlPlane)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)
	podAnnotations := manifestutils.PodAnnotations(b.KinkControlPlane, nil)

	ha := false
	replicas := b.KinkControlPlane.Spec.APIServer.Replicas
	if replicas > 1 {
		ha = true
	}

	podSpec := corev1.PodSpec{
		Affinity:   manifestutils.Affinity(b.KinkControlPlane),
		Containers: []corev1.Container{b.container(image, ha)},
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
			Replicas: &replicas,
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

func (b *APIServer) Service() (*corev1.Service, error) {
	name := naming.APIServer(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.APIServer.Image,
		b.KinkControlPlane.Spec.Version,
		version.APIServer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentAPIServer, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(b.KinkControlPlane.ObjectMeta, ComponentAPIServer, ConceptControlPlane)
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
			Ports:    []corev1.ServicePort{}, // TODO
		},
	}, nil
}

func (b *APIServer) container(image string, ha bool) corev1.Container {
	cfg := b.KinkControlPlane.Spec.APIServer
	resources := cfg.Resources
	verbosity := cfg.Verbosity

	args := map[string]string{ // TODO: add default flags
		"v":             fmt.Sprint(verbosity),
		"leader-elect":  fmt.Sprint(ha),
		"etcd-cafile":   "/etc/etcd/ca.crt",          // TODO: validate with cert-manager
		"etcd-certfile": "/etc/etcd/etcd-client.crt", // TODO: validate with cert-manager
		"etcd-keyfile":  "/etc/etcd/etcd-client.key", // TODO: validate with cert-manager
		"etcd-servers":  naming.KineEndpoint(b.KinkControlPlane.Name, b.KinkControlPlane.Namespace),
	}
	for arg, value := range cfg.ExtraArgs {
		if _, ok := args[arg]; !ok {
			args[arg] = value
		}
	}

	return corev1.Container{
		Name:      naming.APIServerContainer(),
		Image:     image,
		Command:   []string{"kube-apiserver"},
		Args:      buildArgs(args),
		Resources: resources,
		Ports: []corev1.ContainerPort{
			{Name: "metrics", ContainerPort: 6443},
		},
		ImagePullPolicy: b.KinkControlPlane.Spec.APIServer.ImagePullPolicy,
		SecurityContext: nil,
		LivenessProbe:   nil, // TODO
		ReadinessProbe:  nil, // TODO
		StartupProbe:    nil, // TODO
	}
}
