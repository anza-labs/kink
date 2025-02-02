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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestFilterer(t *testing.T) {
	t.Parallel()

	t.Run("Only", func(t *testing.T) {
		t.Parallel()

		for name, tc := range map[string]struct {
			filter   Only[*corev1.Pod]
			input    []client.Object
			expected []client.Object
		}{
			"only_pods": {
				filter: Only[*corev1.Pod]{},
				input: []client.Object{
					&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-1"}},
					&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc-1"}},
					&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-2"}},
				},
				expected: []client.Object{
					&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-1"}},
					&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-2"}},
				},
			},
			"no_matches": {
				filter: Only[*corev1.Pod]{},
				input: []client.Object{
					&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc-1"}},
				},
				expected: []client.Object{},
			},
			"empty_list": {
				filter:   Only[*corev1.Pod]{},
				input:    []client.Object{},
				expected: []client.Object{},
			},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				// test
				actual := tc.filter.Filter(tc.input)

				// validate
				assert.Equal(t, tc.expected, actual)
			})
		}
	})

	t.Run("Exclude", func(t *testing.T) {
		t.Parallel()

		for name, tc := range map[string]struct {
			filter   Exclude[*corev1.Pod]
			input    []client.Object
			expected []client.Object
		}{
			"exclude_pods": {
				filter: Exclude[*corev1.Pod]{},
				input: []client.Object{
					&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-1"}},
					&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc-1"}},
					&corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod-2"}},
				},
				expected: []client.Object{
					&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc-1"}},
				},
			},
			"no_exclusion": {
				filter: Exclude[*corev1.Pod]{},
				input: []client.Object{
					&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc-1"}},
				},
				expected: []client.Object{
					&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "svc-1"}},
				},
			},
			"empty_list": {
				filter:   Exclude[*corev1.Pod]{},
				input:    []client.Object{},
				expected: []client.Object{},
			},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				// test
				actual := tc.filter.Filter(tc.input)

				// validate
				assert.Equal(t, tc.expected, actual)
			})
		}
	})
}
