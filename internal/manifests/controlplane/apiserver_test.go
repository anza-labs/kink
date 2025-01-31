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
	"testing"

	"github.com/stretchr/testify/assert"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// The tests ensure that method calls do not produce errors and return expected results.
func TestAPIServer(t *testing.T) {
	t.Parallel()

	t.Run("Build", func(t *testing.T) {
		// this is a smoke test, it's not meant to validate _everything_
		t.Parallel()

		// prepare
		apiServer := (&APIServer{KinkControlPlane: &controlplanev1alpha1.KinkControlPlane{}})

		// test
		actual, err := apiServer.Build()

		// validate
		assert.NoError(t, err)
		assert.Len(t, actual, 2)
	})

	t.Run("Deployment", func(t *testing.T) {
		t.Parallel()

		for name, tc := range map[string]struct {
			kcp      *controlplanev1alpha1.KinkControlPlane
			expected *appsv1.Deployment
		}{} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				// prepare
				apiServer := (&APIServer{KinkControlPlane: tc.kcp})

				// test
				actual, err := apiServer.Deployment()

				// validate
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			})
		}
	})

	t.Run("Service", func(t *testing.T) {
		t.Parallel()

		for name, tc := range map[string]struct {
			kcp      *controlplanev1alpha1.KinkControlPlane
			expected *corev1.Service
		}{} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				// prepare
				apiServer := (&APIServer{KinkControlPlane: tc.kcp})

				// test
				actual, err := apiServer.Service()

				// validate
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			})
		}
	})
}
