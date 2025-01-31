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

package naming

import "fmt"

func KineEndpoint(name, namespace string) string {
	if namespace == "" {
		return fmt.Sprintf("https://%s:2379", Kine(name))
	}
	return fmt.Sprintf("https://%s.%s.svc.cluster.local:2379", Kine(name), namespace)
}

func KineDNSNames(name, namespace string) []string {
	dnsNames := []string{
		Kine(name),
		"localhost",
	}
	if namespace != "" {
		dnsNames = append(dnsNames, fmt.Sprintf("%s.%s.svc.cluster.local", Kine(name), namespace))
	}
	return dnsNames
}

func KubernetesDNSNames(name, namespace string) []string {
	serviceName := APIServer(name)
	dnsNames := []string{
		serviceName,
		"kubernetes",
		"kubernetes.default",
		"kubernetes.default.svc",
		"kubernetes.default.svc.cluster.local",
	}
	if namespace != "" {
		dnsNames = append(dnsNames,
			fmt.Sprintf("%s.%s", serviceName, namespace),
			fmt.Sprintf("%s.%s.svc", serviceName, namespace),
			fmt.Sprintf("%s.%s.svc.cluster.local", serviceName, namespace),
		)
	}
	return dnsNames
}
