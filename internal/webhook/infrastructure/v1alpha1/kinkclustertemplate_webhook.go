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

	infrastructurev1alpha1 "github.com/anza-labs/kink/api/infrastructure/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// kinkclustertemplatelog is for logging in this file.
var kinkclustertemplatelog = logf.Log.WithName("kinkclustertemplate-resource")

// SetupKinkClusterTemplateWebhookWithManager registers the webhook for KinkClusterTemplate in the manager.
func SetupKinkClusterTemplateWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&infrastructurev1alpha1.KinkClusterTemplate{}).
		WithValidator(&KinkClusterTemplateCustomValidator{}).
		WithDefaulter(&KinkClusterTemplateCustomDefaulter{}).
		Complete()
}

//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha1-kinkclustertemplate,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kinkclustertemplates,verbs=create;update,versions=v1alpha1,name=mkinkclustertemplate-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkClusterTemplateCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind KinkClusterTemplate when those are created or updated.
type KinkClusterTemplateCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &KinkClusterTemplateCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind KinkClusterTemplate.
func (d *KinkClusterTemplateCustomDefaulter) Default(
	ctx context.Context,
	obj runtime.Object,
) error {
	kinkclustertemplate, ok := obj.(*infrastructurev1alpha1.KinkClusterTemplate)

	if !ok {
		return fmt.Errorf("expected an KinkClusterTemplate object but got %T", obj)
	}
	kinkclustertemplatelog.Info("Defaulting for KinkClusterTemplate", "name", kinkclustertemplate.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// NOTE: change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-kinkclustertemplate,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kinkclustertemplates,verbs=create;update,versions=v1alpha1,name=vkinkclustertemplate-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkClusterTemplateCustomValidator struct is responsible for validating the KinkClusterTemplate resource
// when it is created, updated, or deleted.
type KinkClusterTemplateCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &KinkClusterTemplateCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type KinkClusterTemplate.
func (v *KinkClusterTemplateCustomValidator) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkclustertemplate, ok := obj.(*infrastructurev1alpha1.KinkClusterTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a KinkClusterTemplate object but got %T", obj)
	}
	kinkclustertemplatelog.Info("Validation for KinkClusterTemplate upon creation", "name", kinkclustertemplate.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type KinkClusterTemplate.
func (v *KinkClusterTemplateCustomValidator) ValidateUpdate(
	ctx context.Context,
	oldObj, newObj runtime.Object,
) (admission.Warnings, error) {
	kinkclustertemplate, ok := newObj.(*infrastructurev1alpha1.KinkClusterTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a KinkClusterTemplate object for the newObj but got %T", newObj)
	}
	kinkclustertemplatelog.Info("Validation for KinkClusterTemplate upon update", "name", kinkclustertemplate.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type KinkClusterTemplate.
func (v *KinkClusterTemplateCustomValidator) ValidateDelete(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkclustertemplate, ok := obj.(*infrastructurev1alpha1.KinkClusterTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a KinkClusterTemplate object but got %T", obj)
	}
	kinkclustertemplatelog.Info("Validation for KinkClusterTemplate upon deletion", "name", kinkclustertemplate.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
