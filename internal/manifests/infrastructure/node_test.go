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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	infrastructurev1alpha1 "github.com/anza-labs/kink/api/infrastructure/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestNode(t *testing.T) {
	t.Parallel()

	t.Run("Build", func(t *testing.T) {
		// this is a smoke test, it's not meant to validate _everything_
		t.Parallel()

		// prepare
		b := (&Builder{})

		// test
		actual, err := b.Build(&infrastructurev1alpha1.KinkMachine{})

		// validate
		assert.NoError(t, err)
		assert.Len(t, actual, 2)
	})

	t.Run("StatefulSet", func(t *testing.T) {
		t.Parallel()

		for name, tc := range map[string]struct {
			km       *infrastructurev1alpha1.KinkMachine
			expected *appsv1.StatefulSet
		}{} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				// prepare
				node := (&Node{KinkMachine: tc.km})

				// test
				actual := node.StatefulSet()

				// validate
				assert.Equal(t, tc.expected, actual)
			})
		}
	})

	t.Run("Service", func(t *testing.T) {
		t.Parallel()

		// prepare
		expected := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-node",
				Namespace: "test-namespace",
			},
			Spec: corev1.ServiceSpec{
				Type: corev1.ServiceTypeClusterIP,
				Selector: map[string]string{
					"app.kubernetes.io/component":   "node",
					"app.kubernetes.io/instance":    "test-namespace.test",
					"app.kubernetes.io/managed-by":  "kink",
					"app.kubernetes.io/part-of":     "kink-infrastructure",
					"cluster.x-k8s.io/cluster-name": "test",
				},
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
		node := (&Node{KinkMachine: &infrastructurev1alpha1.KinkMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "test-namespace",
			},
		}})

		// test
		actual := node.Service()

		// validate
		assert.Equal(t, expected.ObjectMeta.Name, actual.ObjectMeta.Name)
		assert.Equal(t, expected.ObjectMeta.Namespace, actual.ObjectMeta.Namespace)
		require.Contains(t, actual.Labels, "kink.anza-labs.dev/node-port-range")
		assert.Equal(t, actual.Labels["kink.anza-labs.dev/node-port-range"], "30000-32767")
		assert.Equal(t, expected.Spec, actual.Spec)
	})
}
