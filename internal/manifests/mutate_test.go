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

// Additional copyrights:
// Copyright The OpenTelemetry Authors

//nolint:dupl // just don't
package manifests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/anza-labs/kink/internal/manifests/manifestutils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMutateDeploymentAdditionalContainers(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.Deployment
		desired  appsv1.Deployment
	}{
		{
			name: "add container to deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
								{
									Name:  "alpine",
									Image: "alpine:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "remove container from deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
								{
									Name:  "alpine",
									Image: "alpine:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "modify container in deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
								{
									Name:  "alpine",
									Image: "alpine:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
								{
									Name:  "alpine",
									Image: "alpine:1.0",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			require.NoError(t, err)
			assert.Equal(t, tt.existing, tt.desired)
		})
	}
}

func TestMutateStatefulSetAdditionalContainers(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.StatefulSet
		desired  appsv1.StatefulSet
	}{
		{
			name: "add container to statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
								{
									Name:  "alpine",
									Image: "alpine:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "remove container from statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
								{
									Name:  "alpine",
									Image: "alpine:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "modify container in statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
								{
									Name:  "alpine",
									Image: "alpine:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
								{
									Name:  "alpine",
									Image: "alpine:1.0",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			require.NoError(t, err)
			assert.Equal(t, tt.existing, tt.desired)
		})
	}
}

func TestMutateDeploymentAffinity(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.Deployment
		desired  appsv1.Deployment
	}{
		{
			name: "add affinity to deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
							Affinity: &corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchFields: []corev1.NodeSelectorRequirement{
													{
														Key:      "kubernetes.io/os",
														Operator: corev1.NodeSelectorOpIn,
														Values:   []string{"linux"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "remove affinity from deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
							Affinity: &corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchFields: []corev1.NodeSelectorRequirement{
													{
														Key:      "kubernetes.io/os",
														Operator: corev1.NodeSelectorOpIn,
														Values:   []string{"linux"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "modify affinity in deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
							Affinity: &corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchFields: []corev1.NodeSelectorRequirement{
													{
														Key:      "kubernetes.io/os",
														Operator: corev1.NodeSelectorOpIn,
														Values:   []string{"linux"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
							Affinity: &corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchFields: []corev1.NodeSelectorRequirement{
													{
														Key:      "kubernetes.io/os",
														Operator: corev1.NodeSelectorOpIn,
														Values:   []string{"windows"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			require.NoError(t, err)
			assert.Equal(t, tt.existing, tt.desired)
		})
	}
}

func TestMutateStatefulSetAffinity(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.StatefulSet
		desired  appsv1.StatefulSet
	}{
		{
			name: "add affinity to statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
							Affinity: &corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchFields: []corev1.NodeSelectorRequirement{
													{
														Key:      "kubernetes.io/os",
														Operator: corev1.NodeSelectorOpIn,
														Values:   []string{"linux"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "remove affinity from statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
							Affinity: &corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchFields: []corev1.NodeSelectorRequirement{
													{
														Key:      "kubernetes.io/os",
														Operator: corev1.NodeSelectorOpIn,
														Values:   []string{"linux"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "modify affinity in statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
							Affinity: &corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchFields: []corev1.NodeSelectorRequirement{
													{
														Key:      "kubernetes.io/os",
														Operator: corev1.NodeSelectorOpIn,
														Values:   []string{"linux"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
							Affinity: &corev1.Affinity{
								NodeAffinity: &corev1.NodeAffinity{
									RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
										NodeSelectorTerms: []corev1.NodeSelectorTerm{
											{
												MatchFields: []corev1.NodeSelectorRequirement{
													{
														Key:      "kubernetes.io/os",
														Operator: corev1.NodeSelectorOpIn,
														Values:   []string{"windows"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			require.NoError(t, err)
			assert.Equal(t, tt.existing, tt.desired)
		})
	}
}

func TestMutateDeploymentArgs(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.Deployment
		desired  appsv1.Deployment
	}{
		{
			name: "add argument to container in deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true"},
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true", "extra-arg=yes"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "remove extra arg from container in deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true", "extra-arg=yes"},
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "modify extra arg in container in deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true", "extra-arg=yes"},
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true", "extra-arg=no"},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			require.NoError(t, err)
			assert.Equal(t, tt.existing, tt.desired)
		})
	}
}

func TestMutateStatefulSetArgs(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.StatefulSet
		desired  appsv1.StatefulSet
	}{
		{
			name: "add argument to container in statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true"},
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true", "extra-arg=yes"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "remove extra arg from container in statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true", "extra-arg=yes"},
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "modify extra arg in container in statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true", "extra-arg=yes"},
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
									Args:  []string{"--default-arg=true", "extra-arg=no"},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			require.NoError(t, err)
			assert.Equal(t, tt.existing, tt.desired)
		})
	}
}

func TestNoImmutableLabelChange(t *testing.T) {
	existingSelectorLabels := map[string]string{
		manifestutils.LabelComponent: "test-component",
		manifestutils.LabelInstance:  "default.deployment",
		manifestutils.LabelManagedBy: "test-operator",
		manifestutils.LabelPartOf:    "test",
	}
	desiredLabels := map[string]string{
		manifestutils.LabelComponent: "test-component",
		manifestutils.LabelInstance:  "default.deployment",
		manifestutils.LabelManagedBy: "test-operator",
		manifestutils.LabelPartOf:    "test",
		"extra-label":                "true",
	}
	err := hasImmutableLabelChange(existingSelectorLabels, desiredLabels)
	require.NoError(t, err)
	assert.NoError(t, err)
}

func TestHasImmutableLabelChange(t *testing.T) {
	existingSelectorLabels := map[string]string{
		manifestutils.LabelComponent: "test-component",
		manifestutils.LabelInstance:  "default.deployment",
		manifestutils.LabelManagedBy: "test-operator",
		manifestutils.LabelPartOf:    "test",
	}
	desiredLabels := map[string]string{
		manifestutils.LabelComponent: "test-component",
		manifestutils.LabelInstance:  "default.deployment",
		manifestutils.LabelManagedBy: "test-operator",
		manifestutils.LabelPartOf:    "not-test",
	}
	err := hasImmutableLabelChange(existingSelectorLabels, desiredLabels)
	assert.Error(t, err)
}

func TestMissingImmutableLabelChange(t *testing.T) {
	existingSelectorLabels := map[string]string{
		manifestutils.LabelComponent: "test-component",
		manifestutils.LabelInstance:  "default.deployment",
		manifestutils.LabelManagedBy: "test-operator",
		manifestutils.LabelPartOf:    "test",
	}
	desiredLabels := map[string]string{
		manifestutils.LabelComponent: "test-component",
		manifestutils.LabelInstance:  "default.deployment",
		manifestutils.LabelManagedBy: "test-operator",
	}
	err := hasImmutableLabelChange(existingSelectorLabels, desiredLabels)
	assert.Error(t, err)
}

func TestMutateDeploymentError(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.Deployment
		desired  appsv1.Deployment
	}{
		{
			name: "modified immutable label in deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.Now(),
					Name:              "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.deployment",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.deployment",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.deployment",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.deployment",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "not-test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "modified immutable selector in deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.Now(),
					Name:              "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.deployment",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.deployment",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.deployment",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "not-test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.deployment",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			assert.Error(t, err)
		})
	}
}

func TestMutateStatefulSetError(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.StatefulSet
		desired  appsv1.StatefulSet
	}{
		{
			name: "modified immutable label in statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.Now(),
					Name:              "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.statefulset",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.statefulset",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.statefulset",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.statefulset",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "not-test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "modified immutable selector in statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.Now(),
					Name:              "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.statefulset",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.statefulset",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.statefulset",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "not-test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.statefulset",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			assert.Error(t, err)
		})
	}
}

func TestMutateDeploymentLabelChange(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.Deployment
		desired  appsv1.Deployment
	}{
		{
			name: "modified label in deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.Now(),
					Name:              "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.deployment",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.deployment",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
								"user-label":                 "existing",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.deployment",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.deployment",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
								"user-label":                 "desired",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "new label in deployment",
			existing: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.Now(),
					Name:              "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.deployment",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.deployment",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
								"user-label":                 "existing",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name: "deployment",
				},
				Spec: appsv1.DeploymentSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.deployment",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.deployment",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
								"user-label":                 "existing",
								"new-user-label":             "desired",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			require.NoError(t, err)
			assert.Equal(t, tt.existing.Spec, tt.desired.Spec)
		})
	}
}

func TestMutateStatefulSetLabelChange(t *testing.T) {
	tests := []struct {
		name     string
		existing appsv1.StatefulSet
		desired  appsv1.StatefulSet
	}{
		{
			name: "modified label in statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.Now(),
					Name:              "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.statefulset",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.statefulset",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
								"user-label":                 "existing",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.statefulset",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.statefulset",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
								"user-label":                 "desired",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "new label in statefulset",
			existing: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					CreationTimestamp: metav1.Now(),
					Name:              "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.statefulset",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.statefulset",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
								"user-label":                 "existing",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
			desired: appsv1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name: "statefulset",
				},
				Spec: appsv1.StatefulSetSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							manifestutils.LabelComponent: "test-component",
							manifestutils.LabelInstance:  "default.statefulset",
							manifestutils.LabelManagedBy: "test-operator",
							manifestutils.LabelPartOf:    "test",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								manifestutils.LabelComponent: "test-component",
								manifestutils.LabelInstance:  "default.statefulset",
								manifestutils.LabelManagedBy: "test-operator",
								manifestutils.LabelPartOf:    "test",
								"user-label":                 "existing",
								"new-user-label":             "desired",
							},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "test",
									Image: "test-image:latest",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mutateFn := MutateFuncFor(&tt.existing, &tt.desired)
			err := mutateFn()
			require.NoError(t, err)
			assert.Equal(t, tt.existing.Spec, tt.desired.Spec)
		})
	}
}
