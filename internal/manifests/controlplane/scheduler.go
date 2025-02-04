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

type Scheduler struct {
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

func (b *Scheduler) Build() ([]client.Object, error) {
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

func (b *Scheduler) Deployment() (*appsv1.Deployment, error) {
	name := naming.Scheduler(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.Scheduler.Image,
		b.KinkControlPlane.Spec.Version,
		version.Scheduler(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentScheduler, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(b.KinkControlPlane.ObjectMeta, ComponentScheduler, ConceptControlPlane)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)
	podAnnotations := manifestutils.PodAnnotations(b.KinkControlPlane, nil)

	ha := false
	replicas := b.KinkControlPlane.Spec.Scheduler.Replicas
	if replicas > 1 {
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

func (b *Scheduler) Service() (*corev1.Service, error) {
	name := naming.Scheduler(b.KinkControlPlane.Name)

	image, err := manifestutils.Image(
		b.KinkControlPlane.Spec.Scheduler.Image,
		b.KinkControlPlane.Spec.Version,
		version.Scheduler(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to assess image: %w", err)
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentScheduler, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(b.KinkControlPlane.ObjectMeta, ComponentScheduler, ConceptControlPlane)
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
					Port:       10259,
					TargetPort: intstr.FromString("self"),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}, nil
}

func (b *Scheduler) volumes() []corev1.Volume {
	name := naming.Scheduler(b.KinkControlPlane.Name)

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
	}
}

func (b *Scheduler) volumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "kubeconfig",
			ReadOnly:  true,
			MountPath: kubeconfigPath,
		},
	}
}

func (b *Scheduler) container(image string, ha bool) corev1.Container {
	cfg := b.KinkControlPlane.Spec.Scheduler
	resources := cfg.Resources
	verbosity := cfg.Verbosity

	args := map[string]string{
		"v":                         fmt.Sprint(verbosity),
		"leader-elect":              fmt.Sprint(ha),
		"kubeconfig":                path.Join(kubeconfigPath, kubeconfigName),
		"authorization-kubeconfig":  path.Join(kubeconfigPath, kubeconfigName),
		"authentication-kubeconfig": path.Join(kubeconfigPath, kubeconfigName),
	}
	for arg, value := range cfg.ExtraArgs {
		if _, ok := args[arg]; !ok {
			args[arg] = value
		}
	}

	return corev1.Container{
		Name:      naming.SchedulerContainer(),
		Image:     image,
		Command:   []string{"kube-scheduler"},
		Args:      buildArgs(args),
		Resources: resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "self",
				ContainerPort: 10259,
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
