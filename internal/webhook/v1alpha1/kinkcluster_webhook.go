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

package v1alpha1

import (
	"context"
	"fmt"

	infrastructurev1alpha1 "github.com/anza-labs/kink/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// kinkclusterlog is for logging in this file.
var kinkclusterlog = logf.Log.WithName("kinkcluster-resource")

// SetupKinkClusterWebhookWithManager registers the webhook for KinkCluster in the manager.
func SetupKinkClusterWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&infrastructurev1alpha1.KinkCluster{}).
		WithValidator(&KinkClusterCustomValidator{}).
		WithDefaulter(&KinkClusterCustomDefaulter{}).
		Complete()
}

//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha1-kinkcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kinkclusters,verbs=create;update,versions=v1alpha1,name=mkinkcluster-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkClusterCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind KinkCluster when those are created or updated.
type KinkClusterCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &KinkClusterCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind KinkCluster.
func (d *KinkClusterCustomDefaulter) Default(
	ctx context.Context,
	obj runtime.Object,
) error {
	kinkcluster, ok := obj.(*infrastructurev1alpha1.KinkCluster)

	if !ok {
		return fmt.Errorf("expected an KinkCluster object but got %T", obj)
	}
	kinkclusterlog.Info("Defaulting for KinkCluster", "name", kinkcluster.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// NOTE: change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-kinkcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kinkclusters,verbs=create;update,versions=v1alpha1,name=vkinkcluster-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkClusterCustomValidator struct is responsible for validating the KinkCluster resource
// when it is created, updated, or deleted.
type KinkClusterCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &KinkClusterCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type KinkCluster.
func (v *KinkClusterCustomValidator) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkcluster, ok := obj.(*infrastructurev1alpha1.KinkCluster)
	if !ok {
		return nil, fmt.Errorf("expected a KinkCluster object but got %T", obj)
	}
	kinkclusterlog.Info("Validation for KinkCluster upon creation", "name", kinkcluster.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type KinkCluster.
func (v *KinkClusterCustomValidator) ValidateUpdate(
	ctx context.Context,
	oldObj, newObj runtime.Object,
) (admission.Warnings, error) {
	kinkcluster, ok := newObj.(*infrastructurev1alpha1.KinkCluster)
	if !ok {
		return nil, fmt.Errorf("expected a KinkCluster object for the newObj but got %T", newObj)
	}
	kinkclusterlog.Info("Validation for KinkCluster upon update", "name", kinkcluster.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type KinkCluster.
func (v *KinkClusterCustomValidator) ValidateDelete(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkcluster, ok := obj.(*infrastructurev1alpha1.KinkCluster)
	if !ok {
		return nil, fmt.Errorf("expected a KinkCluster object but got %T", obj)
	}
	kinkclusterlog.Info("Validation for KinkCluster upon deletion", "name", kinkcluster.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
