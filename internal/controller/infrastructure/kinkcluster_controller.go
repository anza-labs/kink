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

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// clusterFinalizer is name of cluster finalizer.
	clusterFinalizer = "cluster.kink.anza-labs.com/finalizer"
)

// KinkClusterReconciler reconciles a KinkCluster object.
type KinkClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//nolint:lll // kubebuilder directives cannot be split into lines
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kinkclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kinkclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=kinkclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
//nolint:dupl // just don't
func (r *KinkClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.V(3).Info("Fetching KinkCluster object")
	kinkC := &infrastructurev1alpha1.KinkCluster{}
	if err := r.Get(ctx, req.NamespacedName, kinkC); err != nil {
		log.Error(err, "Failed to fetch KinkCluster object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle finalizer logic
	if kinkC.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(kinkC, clusterFinalizer) {
			log.V(3).Info("Adding finalizer")
			controllerutil.AddFinalizer(kinkC, clusterFinalizer)
			if err := r.Update(ctx, kinkC); err != nil {
				log.Error(err, "Failed to update KinkCluster object with finalizer")
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(kinkC, clusterFinalizer) {
			// Perform cleanup
			log.V(3).Info("Performing cleanup and removing finalizer")
			if err := r.cleanupResources(ctx, kinkC); err != nil {
				log.Error(err, "Failed to clean up resources")
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(kinkC, clusterFinalizer)
			if err := r.Update(ctx, kinkC); err != nil {
				log.Error(err, "Failed to update KinkCluster object during finalizer removal")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	log.V(3).Info("Starting ControlPlane reconciliation")
	if err := r.reconcile(ctx, kinkC); err != nil {
		log.Error(err, "Failed to reconcile resources")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// reconcile ensures that the necessary resources for the given KinkCluster
// are built and applied in the cluster.
func (r *KinkClusterReconciler) reconcile(
	ctx context.Context,
	kinkC *infrastructurev1alpha1.KinkCluster,
) error {
	log := log.FromContext(ctx)

	// log.V(2).Info("Building components")
	// obj, err := (&controlplane.Builder{}).Build(kinkC)
	// if err != nil {
	// 	return fmt.Errorf("failed to build components: %w", err)
	// }

	log.V(2).Info("Ensuring components")
	if err := r.ensureResources(ctx, kinkC); err != nil {
		return fmt.Errorf("failed to ensure resources: %w", err)
	}

	// TODO: get status of dependents and set it

	return nil
}

// ensureResource ensures that a resource is created or updated.
func (r *KinkClusterReconciler) ensureResources(
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

// cleanupResources removes resources owned by the KinkCluster.
func (r *KinkClusterReconciler) cleanupResources(
	ctx context.Context,
	kinkC *infrastructurev1alpha1.KinkCluster,
) error {
	log := log.FromContext(ctx, "kinkControllPlane", klog.KRef(kinkC.Namespace, kinkC.Name))
	log.V(3).Info("Cleaning up resources")

	// Define a list of owned resources to delete
	resourceTypes := []client.ObjectList{}

	ownerUID := kinkC.GetUID()
	for _, resourceType := range resourceTypes {
		list := resourceType.DeepCopyObject().(client.ObjectList)
		err := r.List(ctx, list, client.InNamespace(kinkC.Namespace))
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
func (r *KinkClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.KinkCluster{}).
		Named("kinkcluster").
		Complete(r)
}
