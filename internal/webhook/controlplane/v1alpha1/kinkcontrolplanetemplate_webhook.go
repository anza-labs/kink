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

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// kinkcontrolplanetemplatelog is for logging in this file.
var kinkcontrolplanetemplatelog = logf.Log.WithName("kinkcontrolplanetemplate-resource")

// SetupKinkControlPlaneTemplateWebhookWithManager registers the webhook for KinkControlPlaneTemplate in the manager.
func SetupKinkControlPlaneTemplateWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&controlplanev1alpha1.KinkControlPlaneTemplate{}).
		WithValidator(&KinkControlPlaneTemplateCustomValidator{}).
		WithDefaulter(&KinkControlPlaneTemplateCustomDefaulter{}).
		Complete()
}

//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/mutate-controlplane-cluster-x-k8s-io-v1alpha1-kinkcontrolplanetemplate,mutating=true,failurePolicy=fail,sideEffects=None,groups=controlplane.cluster.x-k8s.io,resources=kinkcontrolplanetemplates,verbs=create;update,versions=v1alpha1,name=mkinkcontrolplanetemplate-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkControlPlaneTemplateCustomDefaulter struct is responsible for setting default values
// on the custom resource of the
// Kind KinkControlPlaneTemplate when those are created or updated.
type KinkControlPlaneTemplateCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &KinkControlPlaneTemplateCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind KinkControlPlaneTemplate.
func (d *KinkControlPlaneTemplateCustomDefaulter) Default(
	ctx context.Context,
	obj runtime.Object,
) error {
	kinkcontrolplanetemplate, ok := obj.(*controlplanev1alpha1.KinkControlPlaneTemplate)

	if !ok {
		return fmt.Errorf("expected an KinkControlPlaneTemplate object but got %T", obj)
	}
	kinkcontrolplanetemplatelog.Info("Defaulting for KinkControlPlaneTemplate", "name", kinkcontrolplanetemplate.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// NOTE: change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/validate-controlplane-cluster-x-k8s-io-v1alpha1-kinkcontrolplanetemplate,mutating=false,failurePolicy=fail,sideEffects=None,groups=controlplane.cluster.x-k8s.io,resources=kinkcontrolplanetemplates,verbs=create;update,versions=v1alpha1,name=vkinkcontrolplanetemplate-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkControlPlaneTemplateCustomValidator struct is responsible for validating the KinkControlPlaneTemplate resource
// when it is created, updated, or deleted.
type KinkControlPlaneTemplateCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &KinkControlPlaneTemplateCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered
// for the type KinkControlPlaneTemplate.
func (v *KinkControlPlaneTemplateCustomValidator) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkcontrolplanetemplate, ok := obj.(*controlplanev1alpha1.KinkControlPlaneTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a KinkControlPlaneTemplate object but got %T", obj)
	}
	kinkcontrolplanetemplatelog.Info("Validation for KinkControlPlaneTemplate upon creation",
		"name", kinkcontrolplanetemplate.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered
// for the type KinkControlPlaneTemplate.
func (v *KinkControlPlaneTemplateCustomValidator) ValidateUpdate(
	ctx context.Context,
	oldObj, newObj runtime.Object,
) (admission.Warnings, error) {
	kinkcontrolplanetemplate, ok := newObj.(*controlplanev1alpha1.KinkControlPlaneTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a KinkControlPlaneTemplate object for the newObj but got %T", newObj)
	}
	kinkcontrolplanetemplatelog.Info("Validation for KinkControlPlaneTemplate upon update",
		"name", kinkcontrolplanetemplate.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered
// for the type KinkControlPlaneTemplate.
func (v *KinkControlPlaneTemplateCustomValidator) ValidateDelete(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkcontrolplanetemplate, ok := obj.(*controlplanev1alpha1.KinkControlPlaneTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a KinkControlPlaneTemplate object but got %T", obj)
	}
	kinkcontrolplanetemplatelog.Info("Validation for KinkControlPlaneTemplate upon deletion",
		"name", kinkcontrolplanetemplate.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
