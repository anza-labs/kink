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

// Additional copyrights:
// Copyright The OpenTelemetry Authors

package util

import (
	"context"
	"errors"
	"fmt"

	"github.com/anza-labs/kink/internal/manifests"

	rbacv1 "k8s.io/api/rbac/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func ShouldGVK(obj client.Object, scheme *runtime.Scheme) schema.GroupVersionKind {
	gvk, _ := apiutil.GVKForObject(obj, scheme)
	return gvk
}

// getList queries the Kubernetes API to list the requested resource, setting the list l of type T.
func getList[T client.Object](
	ctx context.Context,
	cl client.Client,
	scheme *runtime.Scheme,
	l T,
	options ...client.ListOption,
) (map[types.UID]client.Object, error) {
	ownedObjects := map[types.UID]client.Object{}
	gvk, err := apiutil.GVKForObject(l, scheme)
	if err != nil {
		return nil, err
	}
	gvk.Kind = fmt.Sprintf("%sList", gvk.Kind)
	list, err := scheme.New(gvk)
	if err != nil {
		return nil, fmt.Errorf("unable to list objects of type %s: %w", gvk.Kind, err)
	}

	objList := list.(client.ObjectList)

	err = cl.List(ctx, objList, options...)
	if err != nil {
		return ownedObjects, fmt.Errorf("error listing %T: %w", l, err)
	}
	objs, err := apimeta.ExtractList(objList)
	if err != nil {
		return ownedObjects, fmt.Errorf("error listing %T: %w", l, err)
	}
	for i := range objs {
		typedObj, ok := objs[i].(T)
		if !ok {
			return ownedObjects, fmt.Errorf("error listing %T: %w", l, err)
		}
		ownedObjects[typedObj.GetUID()] = typedObj
	}
	return ownedObjects, nil
}

func FindOwnedObjects(
	ctx context.Context,
	kubeClient client.Client,
	scheme *runtime.Scheme,
	owner metav1.Object,
	ownedObjectTypes []client.Object,
) (map[types.UID]client.Object, error) {
	ownedObjects := map[types.UID]client.Object{}

	listOpts := []client.ListOption{
		client.InNamespace(owner.GetNamespace()),
	}

	for _, objectType := range ownedObjectTypes {
		objs, err := getList(ctx, kubeClient, scheme, objectType, listOpts...)
		if err != nil {
			return nil, err
		}
		for uid, object := range objs {
			ownerUID := owner.GetUID()
			for _, ref := range object.GetOwnerReferences() {
				if ref.UID == ownerUID {
					ownedObjects[uid] = object
				}
			}
		}
	}

	return ownedObjects, nil
}

func isNamespaceScoped(obj client.Object) bool {
	switch obj.(type) {
	case *rbacv1.ClusterRole, *rbacv1.ClusterRoleBinding:
		return false
	default:
		return true
	}
}

// ReconcileDesiredObjects runs the reconcile process using the mutateFn over the given list of objects.
func ReconcileDesiredObjects(
	ctx context.Context,
	kubeClient client.Client,
	owner metav1.Object,
	scheme *runtime.Scheme,
	desiredObjects []client.Object,
	ownedObjects map[types.UID]client.Object,
) error {
	log := log.FromContext(ctx)

	var errs []error
	for _, desired := range desiredObjects {
		l := log.WithValues(
			"object_name", desired.GetName(),
			"object_kind", ShouldGVK(desired, scheme),
		)
		if isNamespaceScoped(desired) {
			if setErr := ctrl.SetControllerReference(owner, desired, scheme); setErr != nil {
				l.Error(setErr, "Failed to set controller owner reference to desired")
				errs = append(errs, setErr)
				continue
			}
		}
		// existing is an object the controller runtime will hydrate for us
		// we obtain the existing object by deep copying the desired object because it's the most convenient way
		existing := desired.DeepCopyObject().(client.Object)
		mutateFn := manifests.MutateFuncFor(existing, desired)

		var op controllerutil.OperationResult
		crudErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			result, createOrUpdateErr := ctrl.CreateOrUpdate(ctx, kubeClient, existing, mutateFn)
			op = result
			return createOrUpdateErr
		})
		if crudErr != nil && errors.As(crudErr, &manifests.ImmutableChangeErr) {
			l.Error(crudErr, "Detected immutable field change, trying to delete, new object will be created on next reconcile",
				"existing", existing.GetName())
			delErr := kubeClient.Delete(ctx, existing)
			if delErr != nil {
				return delErr
			}
			continue
		} else if crudErr != nil {
			l.Error(crudErr, "Failed to configure desired")
			errs = append(errs, crudErr)
			continue
		}

		l.V(3).Info("Desired object reconicled", "result", op)
		// This object is still managed by the operator, remove it from the list of objects to prune
		delete(ownedObjects, existing.GetUID())
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to create objects for %s: %w", owner.GetName(), errors.Join(errs...))
	}

	// Pruning owned objects in the cluster which are not should not be present after the reconciliation.
	err := DeleteObjects(ctx, kubeClient, scheme, ownedObjects)
	if err != nil {
		return fmt.Errorf("failed to prune objects for %s: %w", owner.GetName(), err)
	}

	return nil
}

func DeleteObjects(
	ctx context.Context,
	kubeClient client.Client,
	scheme *runtime.Scheme,
	objects map[types.UID]client.Object,
) error {
	log := log.FromContext(ctx)

	// Pruning owned objects in the cluster which are not should not be present after the reconciliation.
	pruneErrs := []error{}
	for _, obj := range objects {
		l := log.WithValues(
			"object_name", obj.GetName(),
			"object_kind", ShouldGVK(obj, scheme),
		)

		l.V(1).Info("Pruning unmanaged resource")
		err := kubeClient.Delete(ctx, obj)
		if err != nil {
			l.Error(err, "Failed to delete resource")
			pruneErrs = append(pruneErrs, err)
		}
	}

	return errors.Join(pruneErrs...)
}
