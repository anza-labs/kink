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
	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"

	corev1 "k8s.io/api/core/v1"
)

const (
	LabelOS   = "kubernetes.io/os"
	LabelArch = "kubernetes.io/arch"
)

// Affinity return the affinty rules for Registry pod.
func Affinity(instance *controlplanev1alpha1.KinkControlPlane) *corev1.Affinity {
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      LabelArch,
								Operator: corev1.NodeSelectorOpIn,
								Values: []string{
									"amd64",
									"arm64",
								},
							},
							{
								Key:      LabelOS,
								Operator: corev1.NodeSelectorOpIn,
								Values: []string{
									"linux",
								},
							},
						},
					},
				},
			},
		},
	}
}
