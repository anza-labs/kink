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
	"errors"
	"fmt"
	"math"

	"github.com/Masterminds/semver/v3"
	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"

	controlplanev1alpha1 "github.com/anza-labs/kink/api/controlplane/v1alpha1"
	"github.com/anza-labs/kink/internal/controller/util"
	"github.com/anza-labs/kink/internal/manifests/controlplane"
	"github.com/anza-labs/kink/internal/manifests/manifestutils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	gatewayapiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	// controlplaneFinalizer is name of control plane finalizer.
	controlplaneFinalizer = "control-plane.kink.anza-labs.dev/finalizer"
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
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=configmaps/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=secrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=secrets/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=services/finalizers,verbs=update
// +kubebuilder:rbac:groups=cert-manager.io,resources=issuers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=issuers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cert-manager.io,resources=issuers/finalizers,verbs=update
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates/finalizers,verbs=update
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/finalizers,verbs=update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/finalizers,verbs=update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes/finalizers,verbs=update

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
	if !kinkCP.DeletionTimestamp.IsZero() {
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

	log.V(2).Info("Starting owner reconciliation")
	if err := r.checkOwnership(kinkCP); err != nil {
		log.Error(err, "Failed to reconcile owner")
		return ctrl.Result{}, err
	}

	if !controllerutil.ContainsFinalizer(kinkCP, controlplaneFinalizer) {
		log.V(3).Info("Adding finalizer")
		controllerutil.AddFinalizer(kinkCP, controlplaneFinalizer)
		if err := r.Update(ctx, kinkCP); err != nil {
			log.Error(err, "Failed to update KinkControlPlane object with finalizer")
			return ctrl.Result{}, err
		}
	}

	log.V(2).Info("Starting ControlPlane resource reconciliation")
	if err := r.reconcileResources(ctx, kinkCP); err != nil {
		log.Error(err, "Failed to reconcile resources")
		return ctrl.Result{}, err
	}

	if err := r.reconcileStatus(ctx, kinkCP); err != nil {
		log.Error(err, "Failed to reconcile status")
		return ctrl.Result{}, err
	}

	log.V(2).Info("Reconciliation successful")
	return ctrl.Result{}, nil
}

func (r *KinkControlPlaneReconciler) checkOwnership(
	kinkCP *controlplanev1alpha1.KinkControlPlane,
) error {
	owners := kinkCP.OwnerReferences

	for _, owner := range owners {
		if owner.Kind == "Cluster" {
			return nil
		}
	}

	return errors.New("missing OwnerReference from the Cluster controller, waiting for it")
}

// reconcile ensures that the necessary resources for the given KinkControlPlane
// are built and applied in the cluster.
func (r *KinkControlPlaneReconciler) reconcileResources(
	ctx context.Context,
	kinkCP *controlplanev1alpha1.KinkControlPlane,
) error {
	log := log.FromContext(ctx)

	log.V(2).Info("Building components")
	obj, err := (&controlplane.Builder{}).Build(kinkCP)
	if err != nil {
		return fmt.Errorf("failed to build components: %w", err)
	}

	ownedObjects, err := util.FindOwnedObjects(
		ctx,
		r.Client,
		r.Scheme,
		kinkCP,
		r.GetOwnedResourceTypes(util.Exclude[*corev1.Secret]{}),
	)
	if err != nil {
		return fmt.Errorf("failed to find owned objects: %w", err)
	}
	log.V(8).Info("Found objects", "objects", len(ownedObjects))

	log.V(2).Info("Reconciling components", "object_count", len(ownedObjects), "expected_count", len(obj))
	if err := util.ReconcileDesiredObjects(
		ctx,
		r.Client,
		kinkCP,
		r.Scheme,
		obj,
		ownedObjects,
	); err != nil {
		return fmt.Errorf("failed to ensure resources: %w", err)
	}

	log.V(2).Info("Building kubeconfigs")
	kc, err := (&controlplane.Kubeconfig{
		Client:           r.Client,
		KinkControlPlane: kinkCP,
	}).Build(ctx)
	if err != nil {
		return fmt.Errorf("failed to build kubeconfigs: %w", err)
	}

	ownedSecrets, err := util.FindOwnedObjects(
		ctx,
		r.Client,
		r.Scheme,
		kinkCP,
		r.GetOwnedResourceTypes(util.Only[*corev1.Secret]{}),
	)
	if err != nil {
		return fmt.Errorf("failed to find owned secrets: %w", err)
	}
	log.V(8).Info("Found objects", "objects", ownedSecrets)

	log.V(2).Info("Reconciling kubeconfigs", "object_count", len(ownedSecrets), "expected_count", len(kc))
	if err := util.ReconcileDesiredObjects(
		ctx,
		r.Client,
		kinkCP,
		r.Scheme,
		kc,
		ownedSecrets,
	); err != nil {
		return fmt.Errorf("failed to ensure secrets: %w", err)
	}

	return nil
}

func (r *KinkControlPlaneReconciler) reconcileStatus(
	ctx context.Context,
	kinkCP *controlplanev1alpha1.KinkControlPlane,
) error {
	logger := log.FromContext(ctx)

	ownedObjects, err := util.FindOwnedObjects(
		ctx,
		r.Client,
		r.Scheme,
		kinkCP,
		r.GetOwnedResourceTypes(util.Exclude[*corev1.Secret]{}),
	)
	if err != nil {
		return fmt.Errorf("failed to find owned objects: %w", err)
	}
	logger.V(8).Info("Found objects", "objects", len(ownedObjects))

	// Aggregated status variables.
	hasReadyAPIServer := false
	allReady := true
	// Use a high initial value to find the minimums.
	minReplicas := int32(math.MaxInt32)
	minAvailable := int32(math.MaxInt32)
	minReady := int32(math.MaxInt32)
	minUpdated := int32(math.MaxInt32)
	minUnavailable := int32(math.MaxInt32)
	lowestVersion := semver.MustParse("v999.999.999")

	var errs error
	// Loop through all owned objects (e.g. Deployments).
	for _, obj := range ownedObjects {
		// Only consider Deployments in this example.
		deployment, ok := obj.(*appsv1.Deployment)
		if !ok {
			continue
		}

		// Determine desired replicas (default to 1 if not set).
		var desired int32 = 1
		if deployment.Spec.Replicas != nil {
			desired = *deployment.Spec.Replicas
		}

		// Extract status counts.
		var (
			replicas    = deployment.Status.Replicas
			ready       = deployment.Status.ReadyReplicas
			available   = deployment.Status.AvailableReplicas
			updated     = deployment.Status.UpdatedReplicas
			unavailable = deployment.Status.UnavailableReplicas
		)

		// If any deployment is not fully ready, then not all are ready.
		if ready < desired {
			allReady = false
		}

		// Track the smallest available, ready, updated replica count across components.
		if replicas < minReplicas {
			minReplicas = replicas
		}
		if available < minAvailable {
			minAvailable = available
		}
		if ready < minReady {
			minReady = ready
		}
		if updated < minUpdated {
			minUpdated = updated
		}
		if unavailable < minUnavailable {
			minUnavailable = unavailable
		}

		// Check if this Deployment is for the APIServer component.
		// Here we assume that API server deployments are labeled with "component=apiserver".
		if val, ok := deployment.Labels[manifestutils.LabelComponent]; ok && val == controlplane.ComponentAPIServer {
			if ready > 0 {
				hasReadyAPIServer = true
			}
		}

		// Kine version is not considered for the ControlPlane version.
		if val, ok := deployment.Labels[manifestutils.LabelComponent]; ok && val == controlplane.ComponentKine {
			continue
		}

		// Retrieve the version from an label.
		// We choose the lowest version (lexicographically) among all components.
		if version, ok := deployment.Labels[manifestutils.LabelVersion]; ok && version != "" {
			v, err := semver.NewVersion(version)
			if err != nil {
				errs = errors.Join(errs, err)
				continue
			}
			if v.Compare(lowestVersion) < 0 {
				lowestVersion = v
			}
		}
	}

	// Set status fields.
	kinkCP.Status.Initialized = hasReadyAPIServer
	kinkCP.Status.Ready = allReady
	kinkCP.Status.Replicas = minReplicas
	kinkCP.Status.ReadyReplicas = minReady
	kinkCP.Status.UpdatedReplicas = minUpdated
	kinkCP.Status.UnavailableReplicas = minUnavailable

	// Set the Version to the lowest found, or fallback to the Spec version.
	if lowestVersion.Major() != 999 {
		kinkCP.Status.Version = ptr.To(lowestVersion.Original())
	}

	if !allReady {
		errs = errors.Join(errs, errors.New("not all components are ready"))
	}

	if err := r.Status().Update(ctx, kinkCP); err != nil {
		errs = errors.Join(errs, fmt.Errorf("failed to apply status changes: %w", err))
	}

	return errs
}

// GetOwnedResourceTypes returns all the resource types the controller can own.
// Even though this method returns an array of client.Object, these are (empty)
// example structs rather than actual resources.
func (r *KinkControlPlaneReconciler) GetOwnedResourceTypes(filters ...util.Filterer) []client.Object {
	objs := []client.Object{
		&appsv1.Deployment{},
		&corev1.ConfigMap{},
		&corev1.Secret{},
		&corev1.Service{},
		&cmv1.Issuer{},
		&cmv1.Certificate{},
		&netv1.Ingress{},
		&gatewayapiv1.Gateway{},
		&gatewayapiv1.HTTPRoute{},
	}
	for _, filter := range filters {
		objs = filter.Filter(objs)
	}
	return objs
}

// cleanupResources removes resources owned by the KinkControlPlane.
func (r *KinkControlPlaneReconciler) cleanupResources(
	ctx context.Context,
	kinkCP *controlplanev1alpha1.KinkControlPlane,
) error {
	log := log.FromContext(ctx)
	log.V(3).Info("Cleaning up resources")

	ownedObjects, err := util.FindOwnedObjects(
		ctx,
		r.Client,
		r.Scheme,
		kinkCP,
		r.GetOwnedResourceTypes(),
	)
	if err != nil {
		return fmt.Errorf("failed to find owned objects: %w", err)
	}

	if err := util.DeleteObjects(ctx, r.Client, r.Scheme, ownedObjects); err != nil {
		return fmt.Errorf("failed to delete owned objects: %w", err)
	}

	log.V(3).Info("Cleanup complete")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KinkControlPlaneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c := ctrl.NewControllerManagedBy(mgr).
		For(&controlplanev1alpha1.KinkControlPlane{}).
		Named("kinkcontrolplane")

	for _, obj := range r.GetOwnedResourceTypes() {
		c = c.Owns(obj, builder.WithPredicates(predicate.GenerationChangedPredicate{}))
	}

	return c.Complete(r)
}
