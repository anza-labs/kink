// Copyright The OpenTelemetry Authors
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

// Additional copyrights:
// Copyright The OpenTelemetry Authors

//nolint:lll // no need to split lines here
package manifestutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	instanceName      = "my-instance"
	instanceNamespace = "my-ns"
	taname            = "my-instance"
	tanamespace       = "my-ns"
)

func TestLabels(t *testing.T) {
	t.Parallel()

	t.Run("CommonSet", func(t *testing.T) {
		t.Parallel()

		// prepare
		objectMeta := metav1.ObjectMeta{
			Name:      instanceName,
			Namespace: instanceNamespace,
		}
		image := "some/image:v1.31.5"

		// test
		labels := Labels(objectMeta, instanceName, image, "instance", "concept", []string{})

		// validate
		assert.Equal(t, "kink", labels["app.kubernetes.io/managed-by"])
		assert.Equal(t, "my-ns.my-instance", labels["app.kubernetes.io/instance"])
		assert.Equal(t, "v1.31.5", labels["app.kubernetes.io/version"])
		assert.Equal(t, "concept", labels["app.kubernetes.io/part-of"])
		assert.Equal(t, "instance", labels["app.kubernetes.io/component"])
	})

	t.Run("Sha256Set", func(t *testing.T) {
		t.Parallel()

		// prepare
		objectMeta := metav1.ObjectMeta{
			Name:      instanceName,
			Namespace: instanceNamespace,
		}
		image := "some/image@sha256:ac0192b549007e22998eb74e8d8488dcfe70f1489520c3b144a6047ac5efbe90"

		// test
		labels := Labels(objectMeta, instanceName, image, "instance", "concept", []string{})

		// validate
		assert.Equal(t, "kink", labels["app.kubernetes.io/managed-by"])
		assert.Equal(t, "my-ns.my-instance", labels["app.kubernetes.io/instance"])
		assert.Equal(t, "ac0192b549007e22998eb74e8d8488dcfe70f1489520c3b144a6047ac5efbe9", labels["app.kubernetes.io/version"])
		assert.Equal(t, "concept", labels["app.kubernetes.io/part-of"])
		assert.Equal(t, "instance", labels["app.kubernetes.io/component"])
	})

	t.Run("TagSha256Set", func(t *testing.T) {
		t.Parallel()

		// prepare
		objectMeta := metav1.ObjectMeta{
			Name:      instanceName,
			Namespace: instanceNamespace,
		}
		image := "some/image:v1.31.5@sha256:ac0192b549007e22998eb74e8d8488dcfe70f1489520c3b144a6047ac5efbe90"

		// test
		labels := Labels(objectMeta, instanceName, image, "instance", "concept", []string{})

		// validate
		assert.Equal(t, "kink", labels["app.kubernetes.io/managed-by"])
		assert.Equal(t, "my-ns.my-instance", labels["app.kubernetes.io/instance"])
		assert.Equal(t, "v1.31.5", labels["app.kubernetes.io/version"])
		assert.Equal(t, "concept", labels["app.kubernetes.io/part-of"])
		assert.Equal(t, "instance", labels["app.kubernetes.io/component"])
	})

	t.Run("TagUnset", func(t *testing.T) {
		t.Parallel()

		// prepare
		objectMeta := metav1.ObjectMeta{
			Name:      instanceName,
			Namespace: instanceNamespace,
		}
		image := "some/image"

		// test
		labels := Labels(objectMeta, instanceName, image, "instance", "concept", []string{})

		// validate
		assert.Equal(t, "kink", labels["app.kubernetes.io/managed-by"])
		assert.Equal(t, "my-ns.my-instance", labels["app.kubernetes.io/instance"])
		assert.Equal(t, "latest", labels["app.kubernetes.io/version"])
		assert.Equal(t, "concept", labels["app.kubernetes.io/part-of"])
		assert.Equal(t, "instance", labels["app.kubernetes.io/component"])
	})

	t.Run("PropagateDown", func(t *testing.T) {
		t.Parallel()

		// prepare
		objectMeta := metav1.ObjectMeta{
			Labels: map[string]string{
				"myapp":                  "mycomponent",
				"app.kubernetes.io/name": "test",
			},
		}
		image := "some/image"

		// test
		labels := Labels(objectMeta, instanceName, image, "instance", "concept", []string{})

		// verify
		assert.Len(t, labels, 7)
		assert.Equal(t, "mycomponent", labels["myapp"])
		assert.Equal(t, "test", labels["app.kubernetes.io/name"])
	})

	t.Run("Filter", func(t *testing.T) {
		t.Parallel()

		// prepare
		objectMeta := metav1.ObjectMeta{
			Labels: map[string]string{"test.bar.io": "foo", "test.foo.io": "bar"},
		}

		// test
		// This requires the filter to be in regex match form and not the other simpler wildcard one.
		labels := Labels(objectMeta, instanceName, "latest", "instance", "concept", []string{".*.bar.io"})

		// verify
		assert.Len(t, labels, 7)
		assert.NotContains(t, labels, "test.bar.io")
		assert.Equal(t, "bar", labels["test.foo.io"])
	})
}
