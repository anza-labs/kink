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

// kinkmachinelog is for logging in this file.
var kinkmachinelog = logf.Log.WithName("kinkmachine-resource")

// SetupKinkMachineWebhookWithManager registers the webhook for KinkMachine in the manager.
func SetupKinkMachineWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&infrastructurev1alpha1.KinkMachine{}).
		WithValidator(&KinkMachineCustomValidator{}).
		WithDefaulter(&KinkMachineCustomDefaulter{}).
		Complete()
}

//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha1-kinkmachine,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kinkmachines,verbs=create;update,versions=v1alpha1,name=mkinkmachine-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkMachineCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind KinkMachine when those are created or updated.
type KinkMachineCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &KinkMachineCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind KinkMachine.
func (d *KinkMachineCustomDefaulter) Default(
	ctx context.Context,
	obj runtime.Object,
) error {
	kinkmachine, ok := obj.(*infrastructurev1alpha1.KinkMachine)

	if !ok {
		return fmt.Errorf("expected an KinkMachine object but got %T", obj)
	}
	kinkmachinelog.Info("Defaulting for KinkMachine", "name", kinkmachine.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// NOTE: change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-kinkmachine,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=kinkmachines,verbs=create;update,versions=v1alpha1,name=vkinkmachine-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkMachineCustomValidator struct is responsible for validating the KinkMachine resource
// when it is created, updated, or deleted.
type KinkMachineCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &KinkMachineCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type KinkMachine.
func (v *KinkMachineCustomValidator) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkmachine, ok := obj.(*infrastructurev1alpha1.KinkMachine)
	if !ok {
		return nil, fmt.Errorf("expected a KinkMachine object but got %T", obj)
	}
	kinkmachinelog.Info("Validation for KinkMachine upon creation", "name", kinkmachine.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type KinkMachine.
func (v *KinkMachineCustomValidator) ValidateUpdate(
	ctx context.Context,
	oldObj, newObj runtime.Object,
) (admission.Warnings, error) {
	kinkmachine, ok := newObj.(*infrastructurev1alpha1.KinkMachine)
	if !ok {
		return nil, fmt.Errorf("expected a KinkMachine object for the newObj but got %T", newObj)
	}
	kinkmachinelog.Info("Validation for KinkMachine upon update", "name", kinkmachine.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type KinkMachine.
func (v *KinkMachineCustomValidator) ValidateDelete(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkmachine, ok := obj.(*infrastructurev1alpha1.KinkMachine)
	if !ok {
		return nil, fmt.Errorf("expected a KinkMachine object but got %T", obj)
	}
	kinkmachinelog.Info("Validation for KinkMachine upon deletion", "name", kinkmachine.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
