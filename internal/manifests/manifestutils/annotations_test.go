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

package manifestutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestAnnotationsPropagateDown(t *testing.T) {
	t.Parallel()

	// prepare
	for _, instance := range []client.Object{
		&controlplanev1alpha1.KinkControlPlane{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{"myapp": "mycomponent"},
			},
		},
	} {
		// test
		annotations := Annotations(instance, []string{})
		podAnnotations := PodAnnotations(instance, []string{})

		// verify
		assert.Len(t, annotations, 1)
		assert.Equal(t, "mycomponent", annotations["myapp"])
		assert.Equal(t, "mycomponent", podAnnotations["myapp"])
	}
}

func TestAnnotationsFilter(t *testing.T) {
	t.Parallel()

	// prepare
	for _, instance := range []client.Object{
		&controlplanev1alpha1.KinkControlPlane{
			ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"test.bar.io":  "foo",
					"test.io/port": "1234",
					"test.io/path": "/test",
				},
			},
		},
	} {

		// This requires the filter to be in regex match form and not the other simpler wildcard one.
		annotations := Annotations(instance, []string{".*\\.bar\\.io"})

		// verify
		assert.Len(t, annotations, 2)
		assert.NotContains(t, annotations, "test.bar.io")
		assert.Equal(t, "1234", annotations["test.io/port"])
	}
}
