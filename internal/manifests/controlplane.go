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

package manifests

import (
	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/naming"
	"github.com/anza-labs/kink/version"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ControlPlaneBuilder struct {
	spec *controlplanev1alpha1.KinkControlPlaneSpec
}

func (b *ControlPlaneBuilder) Kine() []runtime.Object {
	objects := []runtime.Object{}
	return objects
}

func (b *ControlPlaneBuilder) kineContainer() corev1.Container {
	containerImage := b.spec.Kine.Image
	if containerImage == "" {
		containerImage = version.Kine()
	}

	resources := b.spec.Kine.Resources

	return corev1.Container{
		Name:            naming.KineContainer(),
		Image:           containerImage, // TODO: set value here
		Command:         []string{"/kine"},
		Ports:           []corev1.ContainerPort{}, // TODO: set metrics port here
		Resources:       resources,                // TODO: set it from spec
		VolumeMounts:    []corev1.VolumeMount{},   // TODO: mount volume for db
		ImagePullPolicy: b.spec.Kine.ImagePullPolicy,
		SecurityContext: nil,
		// TODO: any of these shoul be set?
		LivenessProbe:  nil,
		ReadinessProbe: nil,
		StartupProbe:   nil,
	}
}

func (b *ControlPlaneBuilder) APIServer() []runtime.Object {
	objects := []runtime.Object{}
	return objects
}

func (b *ControlPlaneBuilder) ControllerManager() []runtime.Object {
	objects := []runtime.Object{}
	return objects
}

func (b *ControlPlaneBuilder) Scheduler() []runtime.Object {
	objects := []runtime.Object{}
	return objects
}
