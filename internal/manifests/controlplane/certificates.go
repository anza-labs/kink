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
	"time"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmetav1 "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/manifests/manifestutils"
	"github.com/anza-labs/kink/internal/naming"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IssuerKind defines the kind of issuer used for certificates.
const IssuerKind = "Issuer"

var (
	// defaultRenewBefore specifies the default time before expiration
	// when certificates should be renewed (30 days).
	defaultRenewBefore = metav1.Duration{Duration: time.Hour * 24 * 30}

	// defaultCertResidualTime defines the default validity period
	// for non-CA certificates (1 year).
	defaultCertResidualTime = metav1.Duration{Duration: time.Hour * 24 * 365}

	// defaultCAResidualTime defines the default validity period
	// for CA certificates (10 years).
	defaultCAResidualTime = metav1.Duration{Duration: time.Hour * 24 * 365 * 10}
)

// Certificates manages the generation of control plane certificates.
type Certificates struct {
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

// Build constructs and returns a list of certificate-related runtime objects.
func (b *Certificates) Build() []client.Object {
	objects := []client.Object{
		b.RootCA(),
		b.ClusterCA(), b.ClusterCAIssuer(),
		b.APIServer(), b.ServiceAccountCertificate(),
		b.AdminCertificate(), b.SchedulerCertificate(), b.ControllerManagerCertificate(),
		b.FrontProxyCA(),
		b.KineCA(),
		b.KineCAIssuer(), b.KineServer(), b.KineAPIServerClient(),
	}
	return objects
}

// RootCA generates a self-signed root CA issuer.
func (b *Certificates) RootCA() *cmv1.Issuer {
	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Issuer{
		ObjectMeta: metav1.ObjectMeta{
			Name:        naming.RootCA(b.KinkControlPlane.Name),
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.IssuerSpec{
			IssuerConfig: cmv1.IssuerConfig{
				SelfSigned: &cmv1.SelfSignedIssuer{},
			},
		},
	}
}

// ClusterCA creates a certificate to act as the cluster's CA.
func (b *Certificates) ClusterCA() *cmv1.Certificate {
	name := naming.ClusterCA(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			IsCA:        true,
			CommonName:  "Kubernetes API",
			Duration:    &defaultCAResidualTime,
			RenewBefore: &defaultRenewBefore,
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.RootCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			PrivateKey: &cmv1.CertificatePrivateKey{
				Algorithm: cmv1.RSAKeyAlgorithm,
				Size:      2048,
			},
			SecretName: name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
		},
	}
}

// ClusterCAIssuer defines an issuer that uses the Cluster CA.
func (b *Certificates) ClusterCAIssuer() *cmv1.Issuer {
	name := naming.ClusterCA(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Issuer{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.IssuerSpec{
			IssuerConfig: cmv1.IssuerConfig{
				CA: &cmv1.CAIssuer{
					SecretName: name,
				},
			},
		},
	}
}

// APIServer generates a server certificate for the Kubernetes API server.
func (b *Certificates) APIServer() *cmv1.Certificate {
	name := naming.APIServerCertificate(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			CommonName: "kubernetes",
			Subject: &cmv1.X509Subject{
				Organizations: []string{"system:masters"},
			},
			Duration:    &defaultCertResidualTime,
			RenewBefore: &defaultRenewBefore,
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.ClusterCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			SecretName: name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
			Usages: []cmv1.KeyUsage{
				cmv1.UsageDigitalSignature,
				cmv1.UsageKeyEncipherment,
				cmv1.UsageServerAuth,
			},
			DNSNames: naming.KubernetesDNSNames(
				b.KinkControlPlane.Name,
				b.KinkControlPlane.Namespace,
				b.KinkControlPlane.Spec.ControlPlaneEndpoint.Host,
			),
			IPAddresses: []string{"127.0.0.1"},
		},
	}
}

func (b *Certificates) AdminCertificate() *cmv1.Certificate {
	name := naming.AdminCertificate(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			CommonName: "cluster-admin",
			Subject: &cmv1.X509Subject{
				Organizations: []string{"system:masters"},
			},
			Duration:    &defaultCertResidualTime,
			RenewBefore: &defaultRenewBefore,
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.ClusterCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			SecretName: name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
			Usages: []cmv1.KeyUsage{
				cmv1.UsageDigitalSignature,
				cmv1.UsageKeyEncipherment,
				cmv1.UsageClientAuth,
			},
		},
	}
}

func (b *Certificates) SchedulerCertificate() *cmv1.Certificate {
	name := naming.SchedulerCertificate(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			CommonName: "system:kube-scheduler",
			Subject: &cmv1.X509Subject{
				Organizations: []string{"system:kube-scheduler"},
			},
			Duration:    &defaultCertResidualTime,
			RenewBefore: &defaultRenewBefore,
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.ClusterCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			SecretName: name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
			Usages: []cmv1.KeyUsage{
				cmv1.UsageDigitalSignature,
				cmv1.UsageKeyEncipherment,
				cmv1.UsageClientAuth,
			},
		},
	}
}

func (b *Certificates) ControllerManagerCertificate() *cmv1.Certificate {
	name := naming.ControllerManagerCertificate(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			CommonName: "system:kube-controller-manager",
			Subject: &cmv1.X509Subject{
				Organizations: []string{"system:kube-controller-manager"},
			},
			Duration:    &defaultCertResidualTime,
			RenewBefore: &defaultRenewBefore,
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.ClusterCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			SecretName: name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
			Usages: []cmv1.KeyUsage{
				cmv1.UsageDigitalSignature,
				cmv1.UsageKeyEncipherment,
				cmv1.UsageClientAuth,
			},
		},
	}
}

// KineCA generates a CA certificate for Kine (etcd alternative).
func (b *Certificates) KineCA() *cmv1.Certificate {
	name := naming.KineCA(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			IsCA:        true,
			CommonName:  "ETCD CA",
			Duration:    &defaultCAResidualTime,
			RenewBefore: &defaultRenewBefore,
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.RootCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			PrivateKey: &cmv1.CertificatePrivateKey{
				Algorithm: cmv1.RSAKeyAlgorithm,
				Size:      2048,
			},
			SecretName: name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
		},
	}
}

// KineCAIssuer creates an issuer that uses the Kine CA.
func (b *Certificates) KineCAIssuer() *cmv1.Issuer {
	name := naming.KineCA(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Issuer{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.IssuerSpec{
			IssuerConfig: cmv1.IssuerConfig{
				CA: &cmv1.CAIssuer{
					SecretName: name,
				},
			},
		},
	}
}

// KineServer generates a certificate for the Kine server.
func (b *Certificates) KineServer() *cmv1.Certificate {
	name := naming.KineServerCertificate(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			CommonName: "ETCD Server",
			Subject: &cmv1.X509Subject{
				Organizations: []string{"etcd"},
			},
			Duration:    &defaultCertResidualTime,
			RenewBefore: &defaultRenewBefore,
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.KineCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			SecretName: name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
			Usages: []cmv1.KeyUsage{
				cmv1.UsageDigitalSignature,
				cmv1.UsageKeyEncipherment,
				cmv1.UsageServerAuth,
				cmv1.UsageClientAuth,
			},
			DNSNames:    naming.KineDNSNames(b.KinkControlPlane.Name, b.KinkControlPlane.Namespace),
			IPAddresses: []string{"127.0.0.1"},
		},
	}
}

// KineAPIServerClient generates a client certificate for the API server to authenticate with Kine.
func (b *Certificates) KineAPIServerClient() *cmv1.Certificate {
	name := naming.KineAPIServerClientCertificate(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			CommonName: "Kubernetes",
			Subject: &cmv1.X509Subject{
				Organizations: []string{"apiserver"},
			},
			Duration:    &defaultCertResidualTime,
			RenewBefore: &defaultRenewBefore,
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.KineCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			SecretName: name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
			Usages: []cmv1.KeyUsage{
				cmv1.UsageDigitalSignature,
				cmv1.UsageKeyEncipherment,
				cmv1.UsageClientAuth,
			},
			DNSNames:    naming.KineDNSNames(b.KinkControlPlane.Name, b.KinkControlPlane.Namespace),
			IPAddresses: []string{"127.0.0.1"},
		},
	}
}

// FrontProxyCA generates a CA certificate for the front proxy.
func (b *Certificates) FrontProxyCA() *cmv1.Certificate {
	name := naming.FrontProxyCA(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			IsCA:        true,
			CommonName:  "Front-End Proxy",
			Duration:    &defaultCAResidualTime,
			RenewBefore: &defaultRenewBefore,
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.RootCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			PrivateKey: &cmv1.CertificatePrivateKey{
				Algorithm: cmv1.RSAKeyAlgorithm,
				Size:      2048,
			},
			SecretName: name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
		},
	}
}

// ServiceAccountCertificate generates a certificate used for signing service account tokens.
func (b *Certificates) ServiceAccountCertificate() *cmv1.Certificate {
	name := naming.ServiceAccountCertificate(b.KinkControlPlane.Name)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Spec: cmv1.CertificateSpec{
			CommonName: "service-accounts",
			Subject: &cmv1.X509Subject{
				Organizations: []string{"system:serviceaccounts"},
			},
			Duration:    &defaultCertResidualTime,
			RenewBefore: &defaultRenewBefore,
			SecretName:  name,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: selectorLabels,
			},
			IssuerRef: cmmetav1.ObjectReference{
				Name: naming.ClusterCA(b.KinkControlPlane.Name),
				Kind: IssuerKind,
			},
			DNSNames: naming.KubernetesDNSNames(
				b.KinkControlPlane.Name,
				b.KinkControlPlane.Namespace,
				b.KinkControlPlane.Spec.ControlPlaneEndpoint.Host,
			),
			IPAddresses: []string{"127.0.0.1"},
		},
	}
}
