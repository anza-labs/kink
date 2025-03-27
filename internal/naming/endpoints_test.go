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

package naming

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinDomainPrefix(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		prefix, domain string
		expected       string
	}{
		{
			prefix:   "k8s",
			domain:   "example.com",
			expected: "k8s.example.com",
		},
		{
			prefix:   "k8s.",
			domain:   "example.com",
			expected: "k8s.example.com",
		},
		{
			prefix:   "k8s",
			domain:   ".example.com",
			expected: "k8s.example.com",
		},
		{
			prefix:   "k8s.",
			domain:   ".example.com",
			expected: "k8s.example.com",
		},
		{
			prefix:   ".k8s.",
			domain:   "..example.com",
			expected: "k8s.example.com",
		},
		{
			prefix:   "",
			domain:   "example.com",
			expected: "example.com",
		},
		{
			prefix:   "",
			domain:   ".example.com",
			expected: "example.com",
		},
		{
			prefix:   "",
			domain:   "..example.com",
			expected: "example.com",
		},
		{
			prefix:   "k8s",
			domain:   "",
			expected: "k8s",
		},
		{
			prefix:   "k8s.",
			domain:   "",
			expected: "k8s",
		},
		{
			prefix:   "k8s..",
			domain:   "",
			expected: "k8s",
		},
		{
			prefix:   ".k8s.",
			domain:   "",
			expected: "k8s",
		},
		{
			prefix:   "..k8s.",
			domain:   "",
			expected: "k8s",
		},
		{
			prefix:   "",
			domain:   "",
			expected: "",
		},
	} {
		name := fmt.Sprintf("%s & %s => %s", tc.prefix, tc.domain, tc.expected)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := JoinDomainPrefix(tc.prefix, tc.domain)
			assert.Equal(t, tc.expected, actual)
		})
	}

}
