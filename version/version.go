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
	_ "embed"
	"fmt"
	"io"

	"github.com/anza-labs/kink/version/internal/values"

	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

var (
	//go:embed values.yaml
	rawValues []byte

	apiServer         string
	controllerManager string
	scheduler         string
	nodeBase          string
	kine              string
)

const (
	dockerRegistry = "docker.io"
	ghcrRegistry   = "ghcr.io"
	k8sRegistry    = "registry.k8s.io"
	quayRegistry   = "quay.io"
)

func init() {
	vals := loadValues(bytes.NewReader(rawValues))
	apiServer = initAPIServer(vals.APIServer.Image)
	controllerManager = initControllerManager(vals.ControllerManager.Image)
	scheduler = initScheduler(vals.Scheduler.Image)
	nodeBase = initNodeBase(vals.NodeBase.Image)
	kine = initKine(vals.Kine.Image)

}

func loadValues(r io.Reader) values.Values {
	v := values.Values{}
	if err := yaml.NewDecoder(r).Decode(&v); err != nil {
		panic(err)
	}
	return v
}

func initAPIServer(image values.Image) string {
	registry := image.Registry
	if registry == "" {
		registry = k8sRegistry
	}
	repository := image.Repository
	tag := image.Tag
	return fmt.Sprintf("%s/%s:%s", registry, repository, tag)
}

func initControllerManager(image values.Image) string {
	registry := image.Registry
	if registry == "" {
		registry = k8sRegistry
	}
	repository := image.Repository
	tag := image.Tag
	return fmt.Sprintf("%s/%s:%s", registry, repository, tag)
}

func initScheduler(image values.Image) string {
	registry := image.Registry
	if registry == "" {
		registry = k8sRegistry
	}
	repository := image.Repository
	tag := image.Tag
	return fmt.Sprintf("%s/%s:%s", registry, repository, tag)
}

func initNodeBase(image values.Image) string {
	registry := image.Registry
	if registry == "" {
		registry = dockerRegistry
	}
	repository := image.Repository
	tag := image.Tag
	return fmt.Sprintf("%s/%s:%s", registry, repository, tag)
}

func initKine(image values.Image) string {
	registry := image.Registry
	if registry == "" {
		registry = ghcrRegistry
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

func NodeBase() string {
	return nodeBase
}

func Kine() string {
	return kine
}
