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

package util

import (
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Filterer defines an interface for filtering a list of Kubernetes objects.
type Filterer interface {
	// Filter takes a list of client.Object and returns a filtered subset of objects.
	Filter(list []client.Object) []client.Object
}

// Only is a generic filter that retains only objects of type T.
type Only[T client.Object] struct{}

// interface guard.
var _ Filterer = (*Only[*corev1.Pod])(nil)

// Filter returns a new list containing only objects of type T.
// Any objects in the input list that do not match type T are excluded.
func (Only[T]) Filter(list []client.Object) []client.Object {
	objs := []client.Object{}
	for _, item := range list {
		switch item.(type) {
		case T:
			objs = append(objs, item)
		default:
			// no-op
		}
	}
	return objs
}

// Exclude is a generic filter that removes objects of type T from a list.
type Exclude[T client.Object] struct{}

// interface guard.
var _ Filterer = (*Exclude[*corev1.Pod])(nil)

// Filter returns a new list excluding all objects of type T.
// Any objects in the input list that match type T are removed.
func (Exclude[T]) Filter(list []client.Object) []client.Object {
	objs := []client.Object{}
	for _, item := range list {
		switch item.(type) {
		case T:
			// no-op
		default:
			objs = append(objs, item)
		}
	}
	return objs
}

// Include is a generic filter that adds objects of type T to a list.
type Include[T client.Object] struct{}

// interface guard.
var _ Filterer = (*Include[*corev1.Pod])(nil)

// Filter returns a new list including new zero objects of type T.
// If there are objects in the input list that match type T, this is no-op.
func (Include[T]) Filter(list []client.Object) []client.Object {
	found := false
	for _, item := range list {
		switch item.(type) {
		case T:
			found = true
		default:
			// no-op
		}
	}
	if !found {
		var zero T
		list = append(list, zero)
	}
	return list
}
