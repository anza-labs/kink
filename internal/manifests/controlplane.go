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

package manifests

import "k8s.io/apimachinery/pkg/runtime"

type ControlPlaneBuilder struct{}

func (b *ControlPlaneBuilder) Kine() []runtime.Object {
	objects := []runtime.Object{}
	return objects
}

func (b *ControlPlaneBuilder) APIServer() []runtime.Object {
	objects := []runtime.Object{}
	return objects
}

func (b *ControlPlaneBuilder) ControllerManager() []runtime.Object {
	objects := []runtime.Object{}
	return objects
}

func (b *ControlPlaneBuilder) Scheduler() []runtime.Object {
	objects := []runtime.Object{}
	return objects
}
