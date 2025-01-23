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

// kinkcontrolplanelog is for logging in this file.
var kinkcontrolplanelog = logf.Log.WithName("kinkcontrolplane-resource")

// SetupKinkControlPlaneWebhookWithManager registers the webhook for KinkControlPlane in the manager.
func SetupKinkControlPlaneWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&controlplanev1alpha1.KinkControlPlane{}).
		WithValidator(&KinkControlPlaneCustomValidator{}).
		WithDefaulter(&KinkControlPlaneCustomDefaulter{}).
		Complete()
}

//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/mutate-controlplane-cluster-x-k8s-io-v1alpha1-kinkcontrolplane,mutating=true,failurePolicy=fail,sideEffects=None,groups=controlplane.cluster.x-k8s.io,resources=kinkcontrolplanes,verbs=create;update,versions=v1alpha1,name=mkinkcontrolplane-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkControlPlaneCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind KinkControlPlane when those are created or updated.
type KinkControlPlaneCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &KinkControlPlaneCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind KinkControlPlane.
func (d *KinkControlPlaneCustomDefaulter) Default(
	ctx context.Context,
	obj runtime.Object,
) error {
	kinkcontrolplane, ok := obj.(*controlplanev1alpha1.KinkControlPlane)

	if !ok {
		return fmt.Errorf("expected an KinkControlPlane object but got %T", obj)
	}
	kinkcontrolplanelog.Info("Defaulting for KinkControlPlane", "name", kinkcontrolplane.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// NOTE: change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:webhook:path=/validate-controlplane-cluster-x-k8s-io-v1alpha1-kinkcontrolplane,mutating=false,failurePolicy=fail,sideEffects=None,groups=controlplane.cluster.x-k8s.io,resources=kinkcontrolplanes,verbs=create;update,versions=v1alpha1,name=vkinkcontrolplane-v1alpha1.kb.io,admissionReviewVersions=v1

// KinkControlPlaneCustomValidator struct is responsible for validating the KinkControlPlane resource
// when it is created, updated, or deleted.
type KinkControlPlaneCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &KinkControlPlaneCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type KinkControlPlane.
func (v *KinkControlPlaneCustomValidator) ValidateCreate(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkcontrolplane, ok := obj.(*controlplanev1alpha1.KinkControlPlane)
	if !ok {
		return nil, fmt.Errorf("expected a KinkControlPlane object but got %T", obj)
	}
	kinkcontrolplanelog.Info("Validation for KinkControlPlane upon creation", "name", kinkcontrolplane.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type KinkControlPlane.
func (v *KinkControlPlaneCustomValidator) ValidateUpdate(
	ctx context.Context,
	oldObj, newObj runtime.Object,
) (admission.Warnings, error) {
	kinkcontrolplane, ok := newObj.(*controlplanev1alpha1.KinkControlPlane)
	if !ok {
		return nil, fmt.Errorf("expected a KinkControlPlane object for the newObj but got %T", newObj)
	}
	kinkcontrolplanelog.Info("Validation for KinkControlPlane upon update", "name", kinkcontrolplane.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type KinkControlPlane.
func (v *KinkControlPlaneCustomValidator) ValidateDelete(
	ctx context.Context,
	obj runtime.Object,
) (admission.Warnings, error) {
	kinkcontrolplane, ok := obj.(*controlplanev1alpha1.KinkControlPlane)
	if !ok {
		return nil, fmt.Errorf("expected a KinkControlPlane object but got %T", obj)
	}
	kinkcontrolplanelog.Info("Validation for KinkControlPlane upon deletion", "name", kinkcontrolplane.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
