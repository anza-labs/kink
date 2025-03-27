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

package manifestutils

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapilatest "k8s.io/client-go/tools/clientcmd/api/latest"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	ErrRequiredFieldMissing = errors.New("secret is missing required field")
)

func NewKubeconfigFor(
	ctx context.Context,
	client client.Client,
	name, endpoint string,
	secretRef types.NamespacedName,
) (*clientcmdapiv1.Config, error) {
	secret := &corev1.Secret{}
	if err := client.Get(ctx, secretRef, secret); err != nil {
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

func ClientFromKubeconfig(config *clientcmdapiv1.Config) (client.Client, error) {
	buf := new(bytes.Buffer)
	if err := clientcmdapilatest.Codec.Encode(config, buf); err != nil {
		return nil, fmt.Errorf("failed to serialize kubeconfig: %w", err)
	}

	restConfig, err := clientcmd.RESTConfigFromKubeConfig(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to create rest.Config: %w", err)
	}

	cli, err := client.New(restConfig, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return cli, nil
}
