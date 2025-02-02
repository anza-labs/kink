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

package controlplane

import (
	"context"
	"fmt"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/manifests/controlplane"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// controlplaneFinalizer is name of control plane finalizer.
	controlplaneFinalizer = "control-plane.kink.anza-labs.com/finalizer"
)

// KinkControlPlaneReconciler reconciles a KinkControlPlane object.
type KinkControlPlaneReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=kinkcontrolplanes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=kinkcontrolplanes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=controlplane.cluster.x-k8s.io,resources=kinkcontrolplanes/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=services/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets/finalizers,verbs=update
// +kubebuilder:rbac:groups=cert-manager.io,resources=issuers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=issuers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cert-manager.io,resources=issuers/finalizers,verbs=update
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *KinkControlPlaneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.V(2).Info("Fetching KinkControlPlane object")
	kinkCP := &controlplanev1alpha1.KinkControlPlane{}
	if err := r.Get(ctx, req.NamespacedName, kinkCP); err != nil {
		log.Error(err, "Failed to fetch KinkControlPlane object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle finalizer logic
	if kinkCP.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(kinkCP, controlplaneFinalizer) {
			log.V(3).Info("Adding finalizer")
			controllerutil.AddFinalizer(kinkCP, controlplaneFinalizer)
			if err := r.Update(ctx, kinkCP); err != nil {
				log.Error(err, "Failed to update KinkControlPlane object with finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(kinkCP, controlplaneFinalizer) {
			// Perform cleanup
			log.V(3).Info("Performing cleanup and removing finalizer")
			if err := r.cleanupResources(ctx, kinkCP); err != nil {
				log.Error(err, "Failed to clean up resources")
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(kinkCP, controlplaneFinalizer)
			if err := r.Update(ctx, kinkCP); err != nil {
				log.Error(err, "Failed to update KinkControlPlane object during finalizer removal")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	log.V(2).Info("Starting ControlPlane reconciliation")
	if err := r.reconcile(ctx, kinkCP); err != nil {
		log.Error(err, "Failed to reconcile resources")
		return ctrl.Result{}, err
	}

	log.V(2).Info("Reconciliation successful")
	return ctrl.Result{}, nil
}

// reconcile ensures that the necessary resources for the given KinkControlPlane
// are built and applied in the cluster.
func (r *KinkControlPlaneReconciler) reconcile(
	ctx context.Context,
	kinkCP *controlplanev1alpha1.KinkControlPlane,
) error {
	log := log.FromContext(ctx)

	log.V(2).Info("Building components")
	obj, err := (&controlplane.Builder{}).Build(kinkCP)
	if err != nil {
		return fmt.Errorf("failed to build components: %w", err)
	}

	log.V(2).Info("Ensuring components")
	if err := r.ensureResources(ctx, kinkCP, obj...); err != nil {
		return fmt.Errorf("failed to ensure resources: %w", err)
	}

	// TODO: get status of dependents and set it

	return nil
}

// ensureResource ensures that a resource is created or updated.
func (r *KinkControlPlaneReconciler) ensureResources(
	ctx context.Context,
	owner client.Object,
	objs ...client.Object,
) error {
	log := log.FromContext(ctx, "owner", klog.KObj(owner))

	for _, resource := range objs {
		log.V(3).Info("Ensuring object exists",
			"name", resource.GetName(),
			"kind", resource.GetObjectKind().GroupVersionKind().Kind)

		_, err := controllerutil.CreateOrUpdate(ctx, r.Client, resource, func() error {
			return ctrl.SetControllerReference(owner, resource, r.Scheme)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// cleanupResources removes resources owned by the KinkControlPlane.
func (r *KinkControlPlaneReconciler) cleanupResources(
	ctx context.Context,
	kinkCP *controlplanev1alpha1.KinkControlPlane,
) error {
	log := log.FromContext(ctx, "kinkControllPlane", klog.KRef(kinkCP.Namespace, kinkCP.Name))
	log.V(3).Info("Cleaning up resources")

	// Define a list of owned resources to delete
	resourceTypes := []client.ObjectList{
		&appsv1.DeploymentList{},
		&appsv1.StatefulSetList{},
		&corev1.ServiceList{},
		&corev1.SecretList{},
		&cmv1.IssuerList{},
		&cmv1.CertificateList{},
	}

	ownerUID := kinkCP.GetUID()
	for _, resourceType := range resourceTypes {
		list := resourceType.DeepCopyObject().(client.ObjectList)
		err := r.List(ctx, list, client.InNamespace(kinkCP.Namespace))
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}

		// Iterate over resources and delete them
		items, err := meta.ExtractList(list)
		if err != nil {
			return fmt.Errorf("failed to extract list: %w", err)
		}
		for _, item := range items {
			resource := item.(client.Object)
			for _, ref := range resource.GetOwnerReferences() {
				if ref.UID == ownerUID {
					log.V(3).Info("Deleting resource",
						"name", resource.GetName(),
						"kind", resource.GetObjectKind().GroupVersionKind().Kind)
					if err := r.Delete(ctx, resource); err != nil {
						return fmt.Errorf("failed to delete resource: %w", err)
					}
				}
			}
		}
	}

	log.V(3).Info("Cleanup complete")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KinkControlPlaneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&controlplanev1alpha1.KinkControlPlane{}).
		Named("kinkcontrolplane").
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Owns(&cmv1.Issuer{}).
		Owns(&cmv1.Certificate{}).
		Complete(r)
}
