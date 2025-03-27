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
	"bytes"
	"fmt"
	"io"

	"github.com/anza-labs/kink/version/internal/values"

	yaml "sigs.k8s.io/yaml/goyaml.v3"

	_ "embed"
)

var (
	//go:embed values.yaml
	rawValues []byte

	apiServer          string
	controllerManager  string
	scheduler          string
	kine               string
	konnectivityServer string
	konnectivityAgent  string
)

const (
	dockerRegistry = "docker.io"
	ghcrRegistry   = "ghcr.io"
	k8sRegistry    = "registry.k8s.io"
	quayRegistry   = "quay.io"
)

func init() {
	vals := loadValues(bytes.NewReader(rawValues))

	apiServer = initVersion(vals["apiServer"].Image, k8sRegistry)
	controllerManager = initVersion(vals["controllerManager"].Image, k8sRegistry)
	scheduler = initVersion(vals["scheduler"].Image, k8sRegistry)
	kine = initVersion(vals["kine"].Image, ghcrRegistry)
	konnectivityServer = initVersion(vals["konnectivityServer"].Image, k8sRegistry)
	konnectivityAgent = initVersion(vals["konnectivityAgent"].Image, k8sRegistry)
}

func loadValues(r io.Reader) values.Values {
	v := values.Values{}
	if err := yaml.NewDecoder(r).Decode(&v); err != nil {
		panic(err)
	}
	return v
}

func initVersion(image values.Image, defaultRegistry string) string {
	registry := image.Registry
	if registry == "" {
		registry = defaultRegistry
	}
	repository := image.Repository
	tag := image.Tag
	if tag == "" {
		tag = "latest"
	}
	return fmt.Sprintf("%s/%s:%s", registry, repository, tag)
}

func APIServer() string {
	return apiServer
}

func ControllerManager() string {
	return controllerManager
}

func Scheduler() string {
	return scheduler
}

func Kine() string {
	return kine
}

func KonnectivityServer() string {
	return konnectivityServer
}

func KonnectivityAgent() string {
	return konnectivityAgent
}
