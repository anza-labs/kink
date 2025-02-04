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

package infrastructure

import (
	"maps"

	infrastructurev1alpha1 "github.com/anza-labs/kink/api/infrastructure/v1alpha1"
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
	ConceptInfra  = "kink-infrastructure"
	ComponentNode = "node"

	dataMountPoint = "/data" // TODO: this should persist kube configuration for node
)

var extraServiceLabels = map[string]string{
	"kink.anza-labs.com/node-port-range": "30000-32767",
}

type Node struct {
	KinkMachine *infrastructurev1alpha1.KinkMachine
}

func (b *Node) Build() []client.Object {
	objects := []client.Object{
		b.Service(),
		b.StatefulSet(),
	}
	return objects
}

func (b *Node) StatefulSet() *appsv1.StatefulSet {
	name := naming.Node(b.KinkMachine.Name)

	image := b.KinkMachine.Spec.Image
	if image == "" {
		image = version.NodeBase()
	}

	dataVolume := corev1.Volume{Name: "data"}
	if b.KinkMachine.Spec.Persistence != nil {
		persistence := b.KinkMachine.Spec.Persistence
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
		b.KinkMachine.ObjectMeta,
		name, image, ComponentNode, ConceptInfra,
		nil,
	)
	selectorLabels := manifestutils.SelectorLabels(b.KinkMachine.ObjectMeta, ComponentNode, ConceptInfra)
	annotations := manifestutils.Annotations(b.KinkMachine, nil)
	podAnnotations := manifestutils.PodAnnotations(b.KinkMachine, nil)

	podSpec := corev1.PodSpec{
		Affinity:         manifestutils.Affinity(b.KinkMachine),
		Containers:       []corev1.Container{b.container(image)},
		Volumes:          []corev1.Volume{dataVolume},
		ImagePullSecrets: b.KinkMachine.Spec.ImagePullSecrets,
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkMachine.Namespace,
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
		}}
}

func (b *Node) Service() *corev1.Service {
	name := naming.Node(b.KinkMachine.Name)

	image := b.KinkMachine.Spec.Image
	if image == "" {
		image = version.NodeBase()
	}

	labels := manifestutils.Labels(
		b.KinkMachine.ObjectMeta,
		name, image, ComponentNode,
		ConceptInfra,
		nil,
	)
	maps.Insert(labels, maps.All(extraServiceLabels))
	selectorLabels := manifestutils.SelectorLabels(b.KinkMachine.ObjectMeta, ComponentNode, ConceptInfra)
	annotations := manifestutils.Annotations(b.KinkMachine, nil)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkMachine.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Selector:  selectorLabels,
			Type:      corev1.ServiceTypeClusterIP,
			ClusterIP: corev1.ClusterIPNone,
			Ports: []corev1.ServicePort{
				{
					Name:       "kubelet",
					Port:       10250,
					TargetPort: intstr.FromString("kubelet"),
					Protocol:   corev1.ProtocolTCP,
				},
				{
					Name:       "kube-proxy",
					Port:       10256,
					TargetPort: intstr.FromString("kube-proxy"),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}
}

func (b *Node) container(image string) corev1.Container {
	resources := b.KinkMachine.Spec.Resources

	return corev1.Container{
		Name:      naming.NodeBaseContainer(),
		Image:     image,
		Command:   []string{}, // TODO
		Args:      []string{},
		Resources: resources,
		Ports: []corev1.ContainerPort{
			{
				Name:          "kubelet",
				ContainerPort: 10250,
				Protocol:      corev1.ProtocolTCP,
			},
			{
				Name:          "kube-proxy",
				ContainerPort: 10256,
				Protocol:      corev1.ProtocolTCP,
			},
		},
		VolumeMounts: []corev1.VolumeMount{
			b.volumeMount(),
		},
		ImagePullPolicy: b.KinkMachine.Spec.ImagePullPolicy,
		SecurityContext: nil,
		LivenessProbe:   nil,
		ReadinessProbe:  nil,
		StartupProbe:    nil,
	}
}

func (b *Node) volumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      "data",
		ReadOnly:  false,
		MountPath: dataMountPoint,
	}
}
