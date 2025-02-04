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
	"slices"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ConceptControlPlane = "kink-control-plane"

	ComponentCertificates      = "certificates"
	ComponentAPIServer         = "api-server"
	ComponentControllerManager = "controller-manager"
	ComponentKine              = "kine"
	ComponentScheduler         = "scheduler"

	rootPKIPath  = "/etc/pki/kubernetes"
	rootCAFile   = "ca.crt"
	rootCertFile = "tls.crt"
	rootKeyFile  = "tls.key"

	kubeconfigPath = "/etc/kubernetes"
	kubeconfigName = "value"

	serviceAccountsPKIPath         = "/etc/pki/service-accounts"
	serviceAccountsCertificateFile = "tls.crt"
	serviceAccountsKeyFile         = "tls.key"
)

func buildArgs(args map[string]string) []string {
	cmd := []string{}
	for arg, val := range args {
		cmd = append(cmd, fmt.Sprintf("--%s=%s", arg, val))
	}
	slices.Sort(cmd)
	return cmd
}

type Builder struct{}

func (b *Builder) Build(kcp *controlplanev1alpha1.KinkControlPlane) ([]client.Object, error) {
	objects := []client.Object{}

	objects = append(objects, (&Certificates{KinkControlPlane: kcp}).Build()...)
	objects = append(objects, (&Kine{KinkControlPlane: kcp}).Build()...)

	kas, err := (&APIServer{KinkControlPlane: kcp}).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build API Server components: %w", err)
	}
	objects = append(objects, kas...)

	kcm, err := (&ControllerManager{KinkControlPlane: kcp}).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build Controller Manager components: %w", err)
	}
	objects = append(objects, kcm...)

	ks, err := (&Scheduler{KinkControlPlane: kcp}).Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build Scheduler components: %w", err)
	}
	objects = append(objects, ks...)

	return objects, nil
}
