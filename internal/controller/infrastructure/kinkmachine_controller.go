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

package infrastructure

import (
	"context"
	"fmt"

	infrastructurev1alpha1 "github.com/anza-labs/kink/api/infrastructure/v1alpha1"
	"github.com/anza-labs/kink/internal/manifests/infrastructure"

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
	// machineFinalizer is name of machine finalizer.
	machineFinalizer = "machine.kink.anza-labs.com/finalizer"
)

// KinkMachineReconciler reconciles a KinkMachine object.
type KinkMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kinkmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kinkmachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kinkmachines/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch
// +kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=statefulsets/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=services/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
//nolint:dupl // just don't
func (r *KinkMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.V(3).Info("Fetching KinkMachine object")
	kinkM := &infrastructurev1alpha1.KinkMachine{}
	if err := r.Get(ctx, req.NamespacedName, kinkM); err != nil {
		log.Error(err, "Failed to fetch KinkMachine object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle finalizer logic
	if kinkM.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(kinkM, machineFinalizer) {
			log.V(3).Info("Adding finalizer")
			controllerutil.AddFinalizer(kinkM, machineFinalizer)
			if err := r.Update(ctx, kinkM); err != nil {
				log.Error(err, "Failed to update KinkMachine object with finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(kinkM, machineFinalizer) {
			// Perform cleanup
			log.V(3).Info("Performing cleanup and removing finalizer")
			if err := r.cleanupResources(ctx, kinkM); err != nil {
				log.Error(err, "Failed to clean up resources")
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(kinkM, machineFinalizer)
			if err := r.Update(ctx, kinkM); err != nil {
				log.Error(err, "Failed to update KinkMachine object during finalizer removal")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	log.V(3).Info("Starting ControlPlane reconciliation")
	if err := r.reconcile(ctx, kinkM); err != nil {
		log.Error(err, "Failed to reconcile resources")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcile ensures that the necessary resources for the given KinkMachine
// are built and applied in the cluster.
func (r *KinkMachineReconciler) reconcile(
	ctx context.Context,
	kinkM *infrastructurev1alpha1.KinkMachine,
) error {
	log := log.FromContext(ctx)

	log.V(2).Info("Building components")
	obj, err := (&infrastructure.Builder{}).Build(kinkM)
	if err != nil {
		return fmt.Errorf("failed to build components: %w", err)
	}

	log.V(2).Info("Ensuring components")
	if err := r.ensureResources(ctx, kinkM, obj...); err != nil {
		return fmt.Errorf("failed to ensure resources: %w", err)
	}

	// TODO: get status of dependents and set it

	return nil
}

// ensureResource ensures that a resource is created or updated.
func (r *KinkMachineReconciler) ensureResources(
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

// cleanupResources removes resources owned by the KinkMachine.
func (r *KinkMachineReconciler) cleanupResources(
	ctx context.Context,
	kinkM *infrastructurev1alpha1.KinkMachine,
) error {
	log := log.FromContext(ctx, "kinkControllPlane", klog.KRef(kinkM.Namespace, kinkM.Name))
	log.V(3).Info("Cleaning up resources")

	// Define a list of owned resources to delete
	resourceTypes := []client.ObjectList{
		&appsv1.StatefulSetList{},
		&corev1.ServiceList{},
	}

	ownerUID := kinkM.GetUID()
	for _, resourceType := range resourceTypes {
		list := resourceType.DeepCopyObject().(client.ObjectList)
		err := r.List(ctx, list, client.InNamespace(kinkM.Namespace))
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
func (r *KinkMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.KinkMachine{}).
		Named("kinkmachine").
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
