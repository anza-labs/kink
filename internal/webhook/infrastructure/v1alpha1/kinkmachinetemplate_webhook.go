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

// kinkmachinetemplatelog is for logging in this file.
var kinkmachinetemplatelog = logf.Log.WithName("kinkmachinetemplate-resource")

// SetupKinkMachineTemplateWebhookWithManager registers the webhook for KinkMachineTemplate in the manager.
func SetupKinkMachineTemplateWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&infrastructurev1alpha1.KinkMachineTemplate{}).
		WithValidator(&KinkMachineTemplateCustomValidator{}).
		WithDefaulter(&KinkMachineTemplateCustomDefaulter{}).
		Complete()
}

//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha1-kinkmachinetemplate,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kinkmachinetemplates,verbs=create;update,versions=v1alpha1,name=mkinkmachinetemplate-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkMachineTemplateCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind KinkMachineTemplate when those are created or updated.
type KinkMachineTemplateCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &KinkMachineTemplateCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind KinkMachineTemplate.
func (d *KinkMachineTemplateCustomDefaulter) Default(
	ctx context.Context,
	obj runtime.Object,
) error {
	kinkmachinetemplate, ok := obj.(*infrastructurev1alpha1.KinkMachineTemplate)

	if !ok {
		return fmt.Errorf("expected an KinkMachineTemplate object but got %T", obj)
	}
	kinkmachinetemplatelog.Info("Defaulting for KinkMachineTemplate", "name", kinkmachinetemplate.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// NOTE: change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-kinkmachinetemplate,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kinkmachinetemplates,verbs=create;update,versions=v1alpha1,name=vkinkmachinetemplate-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkMachineTemplateCustomValidator struct is responsible for validating the KinkMachineTemplate resource
// when it is created, updated, or deleted.
type KinkMachineTemplateCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &KinkMachineTemplateCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type KinkMachineTemplate.
func (v *KinkMachineTemplateCustomValidator) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkmachinetemplate, ok := obj.(*infrastructurev1alpha1.KinkMachineTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a KinkMachineTemplate object but got %T", obj)
	}
	kinkmachinetemplatelog.Info("Validation for KinkMachineTemplate upon creation", "name", kinkmachinetemplate.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type KinkMachineTemplate.
func (v *KinkMachineTemplateCustomValidator) ValidateUpdate(
	ctx context.Context,
	oldObj, newObj runtime.Object,
) (admission.Warnings, error) {
	kinkmachinetemplate, ok := newObj.(*infrastructurev1alpha1.KinkMachineTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a KinkMachineTemplate object for the newObj but got %T", newObj)
	}
	kinkmachinetemplatelog.Info("Validation for KinkMachineTemplate upon update", "name", kinkmachinetemplate.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type KinkMachineTemplate.
func (v *KinkMachineTemplateCustomValidator) ValidateDelete(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkmachinetemplate, ok := obj.(*infrastructurev1alpha1.KinkMachineTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a KinkMachineTemplate object but got %T", obj)
	}
	kinkmachinetemplatelog.Info("Validation for KinkMachineTemplate upon deletion", "name", kinkmachinetemplate.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
