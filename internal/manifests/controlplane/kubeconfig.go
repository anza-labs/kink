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
	"bytes"
	"context"
	"errors"
	"fmt"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/manifests/manifestutils"
	"github.com/anza-labs/kink/internal/naming"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientcmdapilatest "k8s.io/client-go/tools/clientcmd/api/latest"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	ErrRequiredFieldMissing = errors.New("secret is missing required field")
)

type Kubeconfig struct {
	client.Client
	KinkControlPlane *controlplanev1alpha1.KinkControlPlane
}

func (b *Kubeconfig) Build(ctx context.Context) ([]client.Object, error) {
	log := log.FromContext(ctx)
	obj := []client.Object{}

	var errs error

	log.V(4).Info("Building ClusterAPI kubeconfig")
	if kcCAPI, err := b.ClusterAPI(ctx); err != nil {
		errs = errors.Join(errs, fmt.Errorf("error bulding secret for ClusterAPI: %w", err))
	} else {
		log.V(4).Info("Created kubeconfig for ClusterAPI")
		obj = append(obj, kcCAPI)
	}

	log.V(4).Info("Building ControllerManager kubeconfig")
	if kcCM, err := b.ControllerManager(ctx); err != nil {
		errs = errors.Join(errs, fmt.Errorf("error bulding secret for ControllerManager: %w", err))
	} else {
		log.V(4).Info("Created kubeconfig for ControllerManager")
		obj = append(obj, kcCM)
	}

	log.V(4).Info("Building Scheduler kubeconfig")
	if kcS, err := b.Scheduler(ctx); err != nil {
		errs = errors.Join(errs, fmt.Errorf("error bulding secret for Scheduler: %w", err))
	} else {
		log.V(4).Info("Created kubeconfig for Scheduler")
		obj = append(obj, kcS)
	}

	return obj, errs
}

func (b *Kubeconfig) ClusterAPI(ctx context.Context) (*corev1.Secret, error) {
	endpoint := naming.PublicAPIServerEndpoint(
		b.KinkControlPlane.Spec.ControlPlaneEndpoint.Host,
		b.KinkControlPlane.Spec.ControlPlaneEndpoint.Port,
	)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	key := types.NamespacedName{
		Name:      naming.AdminCertificate(b.KinkControlPlane.Name),
		Namespace: b.KinkControlPlane.Namespace,
	}
	config, err := b.newFor(ctx, b.KinkControlPlane.Name, endpoint, key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate kubeconfig: %w", err)
	}

	buf := new(bytes.Buffer)
	if err := clientcmdapilatest.Codec.Encode(config, buf); err != nil {
		return nil, fmt.Errorf("failed to serialize kubeconfig: %w", err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        naming.Kubeconfig(b.KinkControlPlane.Name),
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Type: capiv1beta1.ClusterSecretType,
		Data: map[string][]byte{
			kubeconfigName: buf.Bytes(),
		},
	}, nil
}

func (b *Kubeconfig) Scheduler(ctx context.Context) (*corev1.Secret, error) {
	endpoint := naming.LocalAPIServerEndpoint(b.KinkControlPlane.Name, b.KinkControlPlane.Namespace)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	key := types.NamespacedName{
		Name:      naming.SchedulerCertificate(b.KinkControlPlane.Name),
		Namespace: b.KinkControlPlane.Namespace,
	}
	config, err := b.newFor(ctx, b.KinkControlPlane.Name, endpoint, key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate kubeconfig: %w", err)
	}

	buf := new(bytes.Buffer)
	if err := clientcmdapilatest.Codec.Encode(config, buf); err != nil {
		return nil, fmt.Errorf("failed to serialize kubeconfig: %w", err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        naming.Kubeconfig(naming.Scheduler(b.KinkControlPlane.Name)),
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Type: capiv1beta1.ClusterSecretType,
		Data: map[string][]byte{
			kubeconfigName: buf.Bytes(),
		},
	}, nil
}

func (b *Kubeconfig) ControllerManager(ctx context.Context) (*corev1.Secret, error) {
	endpoint := naming.LocalAPIServerEndpoint(b.KinkControlPlane.Name, b.KinkControlPlane.Namespace)

	selectorLabels := manifestutils.SelectorLabels(
		b.KinkControlPlane.ObjectMeta,
		ComponentCertificates, ConceptControlPlane,
	)
	annotations := manifestutils.Annotations(b.KinkControlPlane, nil)

	key := types.NamespacedName{
		Name:      naming.ControllerManagerCertificate(b.KinkControlPlane.Name),
		Namespace: b.KinkControlPlane.Namespace,
	}
	config, err := b.newFor(ctx, b.KinkControlPlane.Name, endpoint, key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate kubeconfig: %w", err)
	}

	buf := new(bytes.Buffer)
	if err := clientcmdapilatest.Codec.Encode(config, buf); err != nil {
		return nil, fmt.Errorf("failed to serialize kubeconfig: %w", err)
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        naming.Kubeconfig(naming.ControllerManager(b.KinkControlPlane.Name)),
			Namespace:   b.KinkControlPlane.Namespace,
			Labels:      selectorLabels,
			Annotations: annotations,
		},
		Type: capiv1beta1.ClusterSecretType,
		Data: map[string][]byte{
			kubeconfigName: buf.Bytes(),
		},
	}, nil
}

func (b *Kubeconfig) newFor(
	ctx context.Context,
	name, endpoint string,
	secretRef types.NamespacedName,
) (*clientcmdapiv1.Config, error) {
	secret := &corev1.Secret{}
	if err := b.Get(ctx, secretRef, secret); err != nil {
		return nil, fmt.Errorf("failed to fetch secret %s: %w", secretRef, err)
	}

	caCert, caExists := secret.Data["ca.crt"]
	clientCert, certExists := secret.Data["tls.crt"]
	clientKey, keyExists := secret.Data["tls.key"]

	if !caExists || !certExists || !keyExists {
		return nil, fmt.Errorf("%s: %w", secretRef, ErrRequiredFieldMissing)
	}

	return &clientcmdapiv1.Config{
		APIVersion: clientcmdapilatest.Version,
		Kind:       "Config",
		Clusters: []clientcmdapiv1.NamedCluster{
			{
				Name: name,
				Cluster: clientcmdapiv1.Cluster{
					Server:                   endpoint,
					CertificateAuthorityData: caCert,
				},
			},
		},
		AuthInfos: []clientcmdapiv1.NamedAuthInfo{
			{
				Name: name,
				AuthInfo: clientcmdapiv1.AuthInfo{
					ClientCertificateData: clientCert,
					ClientKeyData:         clientKey,
				},
			},
		},
		Contexts: []clientcmdapiv1.NamedContext{
			{
				Name: name,
				Context: clientcmdapiv1.Context{
					Cluster:  name,
					AuthInfo: name,
				},
			},
		},
		CurrentContext: name,
	}, nil
}
