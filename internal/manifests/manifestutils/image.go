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
	"path"
	"strings"

	"github.com/distribution/reference"
)

func Image(image, version, defaultImage string) (string, error) {
	if image == "" {
		if version == "" {
			return defaultImage, nil
		}
		return setVersion(defaultImage, version, true)
	}

	return setVersion(image, version, false)
}

func setVersion(image, version string, force bool) (string, error) {
	ref, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return "", err
	}

	registry := reference.Domain(ref)
	baseImage := reference.Path(ref)

	var result strings.Builder

	result.WriteString(path.Join(registry, baseImage))

	versionSet := false
	if !force {
		if tagged, ok := ref.(reference.Tagged); ok {
			result.WriteString(":" + tagged.Tag())
			versionSet = true
		}
		if digested, ok := ref.(reference.Digested); ok {
			result.WriteString("@" + digested.Digest().String())
			versionSet = true
		}
	}

	if force || !versionSet {
		if version != "" {
			result.WriteString(":" + version)
		}
	}

	return result.String(), nil
}
