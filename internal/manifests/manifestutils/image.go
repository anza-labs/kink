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

// Returns the final image reference and any error encountered during processing.
func Image(image, version, defaultImage string) (string, error) {
	if image == "" {
		if version == "" {
			return defaultImage, nil
		}
		return setVersion(defaultImage, version, true)
	}

	return setVersion(image, version, false)
}

// setVersion constructs an image reference with an optional version tag.
// 
// It parses the provided image reference and builds a normalized image reference string.
// The function handles different scenarios for version specification:
// - If force is false, it preserves existing tags or digests from the original image
// - If force is true or no existing version is found, it applies the provided version
//
// Parameters:
//   - image: The original image reference to parse
//   - version: The version tag to potentially apply
//   - force: Flag to force applying the version, overriding existing tags
//
// Returns:
//   - A fully qualified image reference string
//   - An error if the image reference cannot be parsed
//
// Examples:
//   - setVersion("nginx", "1.2.3", false) returns "nginx:1.2.3"
//   - setVersion("nginx:latest", "1.2.3", false) returns "nginx:latest"
//   - setVersion("nginx:latest", "1.2.3", true) returns "nginx:1.2.3"
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
