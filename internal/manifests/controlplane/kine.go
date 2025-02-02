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
	kineDataMountPoint = "/etc/kine/db"
	kineTLSMountPoint  = "/etc/kine/tls"
	kineEndpoint       = "sqlite://%s/db.sqlite"
)

type Kine struct {
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

func (b *Kine) Build() []client.Object {
	objects := []client.Object{}

	svc := b.Service()
	objects = append(objects, svc)

	ss := b.Deployment()
	objects = append(objects, ss)

	return objects
}

func (b *Kine) Deployment() *appsv1.Deployment {
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
	podAnnotations := manifestutils.PodAnnotations(b.KinkControlPlane, nil)

	podSpec := corev1.PodSpec{
		Affinity:         manifestutils.Affinity(b.KinkControlPlane),
		Containers:       []corev1.Container{b.container(image)},
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
				{
					Name:       "metrics",
					Port:       8080,
					TargetPort: intstr.FromString("metrics"),
					Protocol:   corev1.ProtocolTCP},
				{
					Name:       "kine",
					Port:       2379,
					TargetPort: intstr.FromString("kine"),
					Protocol:   corev1.ProtocolTCP},
			},
		},
	}
}

func (b *Kine) volumes() []corev1.Volume {
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

	return []corev1.Volume{
		dataVolume,
		{
			Name: "tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName:  naming.KineServerCertificate(b.KinkControlPlane.Name),
					DefaultMode: ptr.To[int32](420),
				},
			},
		},
	}
}

func (b *Kine) volumeMounts() []corev1.VolumeMount {
	return []corev1.VolumeMount{
		{
			Name:      "data",
			ReadOnly:  false,
			MountPath: kineDataMountPoint,
		},
		{
			Name:      "tls",
			ReadOnly:  true,
			MountPath: kineTLSMountPoint,
		},
	}
}

func (b *Kine) container(image string) corev1.Container {
	resources := b.KinkControlPlane.Spec.Kine.Resources

	probe := []string{
		"grpc_health_probe",
		"-addr", "127.0.0.1:2379",
		"-tls",
		"-tls-ca-cert", path.Join(kineTLSMountPoint, "ca.crt"),
		"-tls-client-cert", path.Join(kineTLSMountPoint, "tls.crt"),
		"-tls-client-key", path.Join(kineTLSMountPoint, "tls.key"),
	}

	return corev1.Container{
		Name:    naming.KineContainer(),
		Image:   image,
		Command: []string{"kine"},
		Args: []string{
			"--listen-address", "0.0.0.0:2379",
			"--endpoint", fmt.Sprintf(kineEndpoint, kineDataMountPoint),
			"--server-cert-file", path.Join(kineTLSMountPoint, "tls.crt"),
			"--server-key-file", path.Join(kineTLSMountPoint, "tls.key"),
		},
		Resources: resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "metrics",
				ContainerPort: 8080,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          "kine",
				ContainerPort: 2379,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts:    b.volumeMounts(),
		ImagePullPolicy: b.KinkControlPlane.Spec.Kine.ImagePullPolicy,
		SecurityContext: nil,
		LivenessProbe: &corev1.Probe{
			InitialDelaySeconds: 5,
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: probe,
				},
			},
		},
		ReadinessProbe: &corev1.Probe{
			InitialDelaySeconds: 5,
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: probe,
				},
			},
		},
		StartupProbe: nil,
	}
}
