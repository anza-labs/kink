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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientcmdapilatest "k8s.io/client-go/tools/clientcmd/api/latest"
	clientcmdapiv1 "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestKubeconfig(t *testing.T) {
	t.Parallel()

	t.Run("newFor", func(t *testing.T) {
		t.Parallel()

		testKey := types.NamespacedName{Namespace: "test", Name: "test"}
		testData := []byte("test")

		for name, tc := range map[string]struct {
			name, endpoint string
			key            types.NamespacedName
			secret         *corev1.Secret
			expected       *clientcmdapiv1.Config
			expectedError  string
		}{
			"success": {
				key:      testKey,
				name:     "test",
				endpoint: "https://test.example.com:6443",
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testKey.Namespace,
						Name:      testKey.Name,
					},
					Data: map[string][]byte{
						"ca.crt":  testData,
						"tls.crt": testData,
						"tls.key": testData,
					},
				},
				expected: &clientcmdapiv1.Config{
					APIVersion: clientcmdapilatest.Version,
					Kind:       "Config",
					Clusters: []clientcmdapiv1.NamedCluster{
						{
							Name: "test",
							Cluster: clientcmdapiv1.Cluster{
								Server:                   "https://test.example.com:6443",
								CertificateAuthorityData: testData,
							},
						},
					},
					AuthInfos: []clientcmdapiv1.NamedAuthInfo{
						{
							Name: "test",
							AuthInfo: clientcmdapiv1.AuthInfo{
								ClientCertificateData: testData,
								ClientKeyData:         testData,
							},
						},
					},
					Contexts: []clientcmdapiv1.NamedContext{
						{
							Name: "test",
							Context: clientcmdapiv1.Context{
								AuthInfo: "test",
								Cluster:  "test",
							},
						},
					},
					CurrentContext: "test",
				},
			},
			"missing_secret": {
				key:           testKey,
				expectedError: "not found",
			},
			"missing_ca": {
				key: testKey,
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testKey.Namespace,
						Name:      testKey.Name,
					},
					Data: map[string][]byte{
						"tls.crt": testData,
						"tls.key": testData,
					},
				},
				expectedError: ErrRequiredFieldMissing.Error(),
			},
			"missing_cert": {
				key: testKey,
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testKey.Namespace,
						Name:      testKey.Name,
					},
					Data: map[string][]byte{
						"ca.crt":  testData,
						"tls.key": testData,
					},
				},
				expectedError: ErrRequiredFieldMissing.Error(),
			},
			"missing_key": {
				key: testKey,
				secret: &corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: testKey.Namespace,
						Name:      testKey.Name,
					},
					Data: map[string][]byte{
						"ca.crt":  testData,
						"tls.crt": testData,
					},
				},
				expectedError: ErrRequiredFieldMissing.Error(),
			},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				// prepare
				objs := []runtime.Object{}
				if tc.secret != nil {
					objs = append(objs, tc.secret)
				}

				cli := fake.NewFakeClient(objs...)
				ctx := context.Background()

				// test
				actual, err := NewKubeconfigFor(ctx, cli, tc.name, tc.endpoint, tc.key)

				// validate
				if tc.expectedError == "" {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, actual)
				} else {
					assert.ErrorContains(t, err, tc.expectedError)

				}
			})
		}
	})
}
