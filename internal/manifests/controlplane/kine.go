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
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

const (
	kineDataMountPoint = "/etc/kine/db"
	kineEndpoint       = "sqlite://%s/db.sqlite"
)

type Kine struct {
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

func (b *Kine) Build() []runtime.Object {
	objects := []runtime.Object{}

	svc := b.Service()
	objects = append(objects, svc)

	ss := b.StatefulSet()
	objects = append(objects, ss)

	return objects
}

func (b *Kine) StatefulSet() *appsv1.StatefulSet {
	name := naming.Kine(b.KinkControlPlane.Name)

	image := b.KinkControlPlane.Spec.Kine.Image
	if image == "" {
		image = version.Kine()
	}

	dataVolume := corev1.Volume{Name: "data"}
	if b.KinkControlPlane.Spec.Kine.Persistence != nil {
		persistence := b.KinkControlPlane.Spec.Kine.Persistence
		if persistence.PersistentVolumeClaim != nil {
			dataVolume.PersistentVolumeClaim = persistence.PersistentVolumeClaim
		} else {
			dataVolume.EmptyDir = &corev1.EmptyDirVolumeSource{}
			if persistence.EmptyDir != nil {
				dataVolume.EmptyDir.Medium = persistence.EmptyDir.Medium
				dataVolume.EmptyDir.SizeLimit = persistence.EmptyDir.SizeLimit
			}
		}
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentKine, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(b.KinkControlPlane.ObjectMeta, ComponentKine, ConceptControlPlane)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)
	podAnnotations := manifestutils.PodAnnotations(b.KinkControlPlane, nil)

	podSpec := corev1.PodSpec{
		Affinity:   manifestutils.Affinity(b.KinkControlPlane),
		Containers: []corev1.Container{b.container(image)},
		Volumes:    []corev1.Volume{dataVolume},
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			Replicas: ptr.To[int32](1),
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

func (b *Kine) Service() *corev1.Service {
	name := naming.Kine(b.KinkControlPlane.Name)

	image := b.KinkControlPlane.Spec.Kine.Image
	if image == "" {
		image = version.Kine()
	}

	labels := manifestutils.Labels(
		b.KinkControlPlane.ObjectMeta,
		name, image, ComponentKine, ConceptControlPlane,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(b.KinkControlPlane.ObjectMeta, ComponentKine, ConceptControlPlane)
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
				{Name: "metrics", Port: 8080, TargetPort: intstr.FromString("metrics")},
				{Name: "kine", Port: 2379, TargetPort: intstr.FromString("kine")},
			},
		},
	}
}

func (b *Kine) volumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      "data",
		ReadOnly:  false,
		MountPath: kineDataMountPoint,
	}
}

func (b *Kine) container(image string) corev1.Container {
	resources := b.KinkControlPlane.Spec.Kine.Resources

	return corev1.Container{
		Name:    naming.KineContainer(),
		Image:   image,
		Command: []string{"/kine"},
		Args: []string{
			"--listen-address", "0.0.0.0:2379",
			"--endpoint", fmt.Sprintf(kineEndpoint, kineDataMountPoint),
		},
		Resources: resources,
		Ports: []corev1.ContainerPort{
			{Name: "metrics", ContainerPort: 8080},
			{Name: "kine", ContainerPort: 2379},
		},
		VolumeMounts: []corev1.VolumeMount{
			b.volumeMount(),
		},
		ImagePullPolicy: b.KinkControlPlane.Spec.Kine.ImagePullPolicy,
		SecurityContext: nil,
		LivenessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Port: intstr.FromString("metrics"),
				},
			},
		},
		ReadinessProbe: &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Port: intstr.FromString("metrics"),
				},
			},
		},
		StartupProbe: nil,
	}
}
