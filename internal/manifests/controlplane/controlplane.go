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

	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ConceptControlPlane = "kink-control-plane"

	ComponentCertificates      = "certificates"
	ComponentAPIServer         = "api-server"
	ComponentControllerManager = "controller-manager"
	ComponentKine              = "kine"
	ComponentScheduler         = "scheduler"
)

func buildArgs(args map[string]string) []string {
	cmd := []string{}
	for arg, val := range args {
		cmd = append(cmd, fmt.Sprintf("--%s=%s", arg, val))
	}
	return cmd
}

type ControlPlaneBuilder struct{}

func (b *ControlPlaneBuilder) Build(kcp *controlplanev1alpha1.KinkControlPlane) ([]runtime.Object, error) {
	objects := []runtime.Object{}

	objects = append(objects, (&Certificates{KinkControlPlane: kcp}).Build()...)
	objects = append(objects, (&APIServer{KinkControlPlane: kcp}).Build()...)
	objects = append(objects, (&ControllerManager{KinkControlPlane: kcp}).Build()...)
	objects = append(objects, (&Kine{KinkControlPlane: kcp}).Build()...)
	objects = append(objects, (&Scheduler{KinkControlPlane: kcp}).Build()...)

	return objects, nil
}
