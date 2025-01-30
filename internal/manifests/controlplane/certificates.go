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
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/manifests/manifestutils"
	"github.com/anza-labs/kink/internal/naming"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Certificates struct {
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

func (b *Certificates) Build() []runtime.Object {
	objects := []runtime.Object{
		b.RootCA(),
	}
	return objects
}

// TODO:
// - Issuer (Self-signed?)
// - ETCD Server
// - ETCD Client
// https://github.com/kelseyhightower/kubernetes-the-hard-way/blob/master/docs/04-certificate-authority.md

func (b *Certificates) RootCA() *certmanagerv1.Issuer {
	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &certmanagerv1.Issuer{
		ObjectMeta: metav1.ObjectMeta{
			Name:        naming.RootCA(),
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: certmanagerv1.IssuerSpec{
			IssuerConfig: certmanagerv1.IssuerConfig{
				SelfSigned: &certmanagerv1.SelfSignedIssuer{},
			},
		},
	}
}
