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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImage(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		image, version, defaultImage string
		expected, expectedError      string
	}{
		"unversioned default": {
			image:         "",
			version:       "",
			defaultImage:  "registry.example.com/default/image",
			expected:      "registry.example.com/default/image",
			expectedError: "",
		},
		"versioned default": {
			image:         "",
			version:       "",
			defaultImage:  "registry.example.com/default/image:v1.0.0",
			expected:      "registry.example.com/default/image:v1.0.0",
			expectedError: "",
		},
		"empty image, version provided": {
			image:         "",
			version:       "v1.0.0",
			defaultImage:  "registry.example.com/default/image",
			expected:      "registry.example.com/default/image:v1.0.0",
			expectedError: "",
		},
		"empty image, version provided, versioned default": {
			image:         "",
			version:       "v1.2.3",
			defaultImage:  "registry.example.com/default/image:v1.0.0",
			expected:      "registry.example.com/default/image:v1.2.3",
			expectedError: "",
		},
		"image provided, no version": {
			image:         "registry.example.com/custom/image",
			version:       "",
			defaultImage:  "registry.example.com/default/image",
			expected:      "registry.example.com/custom/image",
			expectedError: "",
		},
		"versioned image provided, no version": {
			image:         "registry.example.com/custom/image:v1.0.0",
			version:       "",
			defaultImage:  "registry.example.com/default/image",
			expected:      "registry.example.com/custom/image:v1.0.0",
			expectedError: "",
		},
		"image and version provided": {
			image:         "registry.example.com/custom/image",
			version:       "v2.3.4",
			defaultImage:  "registry.example.com/default/image",
			expected:      "registry.example.com/custom/image:v2.3.4",
			expectedError: "",
		},
		"versioned image and version provided": {
			image:         "registry.example.com/custom/image:v1.0.0",
			version:       "v2.3.4",
			defaultImage:  "registry.example.com/default/image",
			expected:      "registry.example.com/custom/image:v1.0.0",
			expectedError: "",
		},
		"invalid image format": {
			image:         "invalid@image",
			version:       "v1.0.0",
			defaultImage:  "registry.example.com/default/image",
			expected:      "",
			expectedError: "invalid reference format",
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// test
			actual, err := Image(tc.image, tc.version, tc.defaultImage)

			// validate
			if tc.expectedError != "" {
				assert.ErrorContains(t, err, tc.expectedError)
			} else {
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
