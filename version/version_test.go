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

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmoke(t *testing.T) {
	assert.Regexp(t, "^registry.k8s.io/kube-apiserver:v.+$", APIServer())
	assert.Regexp(t, "^registry.k8s.io/kube-controller-manager:v.+$", ControllerManager())
	assert.Regexp(t, "^registry.k8s.io/kube-scheduler:v.+$", Scheduler())
	assert.Regexp(t, "^docker.io/rancher/kine:.+$", Kine())
	assert.Regexp(t, "^registry.k8s.io/kas-network-proxy/proxy-server:v.+$", KonnectivityServer())
	assert.Regexp(t, "^registry.k8s.io/kas-network-proxy/proxy-agent:v.+$", KonnectivityAgent())
	assert.Regexp(t, "^ghcr.io/grpc-ecosystem/grpc-health-probe:v.+$", GRPCHealthProbe())
}
