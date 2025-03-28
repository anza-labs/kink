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

package manifestutils

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Annotations return the annotations for Registry resources.
func Annotations(instance client.Object, filterAnnotations []string) map[string]string {
	// new map every time, so that we don't touch the instance's annotations
	annotations := map[string]string{}

	if nil != instance.GetAnnotations() {
		for k, v := range instance.GetAnnotations() {
			if !IsFilteredSet(k, filterAnnotations) {
				annotations[k] = v
			}
		}
	}

	return annotations
}

// PodAnnotations return the spec annotations for Registry pod.
func PodAnnotations(instance client.Object, filterAnnotations []string) map[string]string {
	// new map every time, so that we don't touch the instance's annotations
	podAnnotations := map[string]string{}

	annotations := Annotations(instance, filterAnnotations)
	// propagating annotations from metadata.annotations
	for kMeta, vMeta := range annotations {
		if _, found := podAnnotations[kMeta]; !found {
			podAnnotations[kMeta] = vMeta
		}
	}

	return podAnnotations
}
