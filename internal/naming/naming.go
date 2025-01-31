// Copyright 2024-2025 anza-labs contributors.
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

// It truncates the name to a maximum of 63 characters and formats it as "<base>-api-server".
func APIServer(base string) string {
	return DNSName(Truncate("%s-api-server", 63, base))
}

// It truncates the name to a maximum of 63 characters and formats it with the base string.
func APIServerCertificate(base string) string {
	return DNSName(Truncate("%s-api-server", 63, base))
}

// APIServerContainer returns the standard container name for the API server.
func APIServerContainer() string {
	return "api-server"
}

// The base string is used to create a unique identifier for the cluster's CA.
func ClusterCA(base string) string {
	return DNSName(Truncate("%s-ca", 63, base))
}

// The resulting name is a combination of the base string and hash, separated by a hyphen.
func ConfigMap(base, hash string) string {
	return DNSName(Truncate("%s-%s", 63, base, hash))
}

// The returned name is suitable for use in Kubernetes DNS naming conventions.
func ControllerManager(base string) string {
	return DNSName(Truncate("%s-controller-manager", 63, base))
}

// ControllerManagerContainer returns the standard container name for the Kubernetes controller manager.
func ControllerManagerContainer() string {
	return "controller-manager"
}

// It truncates the name to 63 characters and formats it with a "-proxy" suffix.
func FrontProxyCA(base string) string {
	return DNSName(Truncate("%s-proxy", 63, base))
}

// Kine generates a DNS name for Kine using the provided base string. It truncates the name to 63 characters and formats it with a "-kine" suffix.
func Kine(base string) string {
	return DNSName(Truncate("%s-kine", 63, base))
}

// It truncates the formatted name to a maximum of 63 characters and ensures a valid DNS name format.
func KineAPIServerClientCertificate(base string) string {
	return DNSName(Truncate("%s-etcd", 63, base))
}

// The base string is used to create a unique identifier for the Kine CA.
func KineCA(base string) string {
	return DNSName(Truncate("%s-etcd", 63, base))
}

// Returns a DNS-compatible name for the Kine server certificate.
func KineServerCertificate(base string) string {
	return DNSName(Truncate("%s-etcd", 63, base))
}

// It provides a consistent string identifier for the Kine container.
func KineContainer() string {
	return "kine"
}

// It truncates the name to a maximum of 63 characters and formats it with a "-root-ca" suffix.
func RootCA(base string) string {
	return DNSName(Truncate("%s-root-ca", 63, base))
}

// It truncates the name to a maximum of 63 characters and formats it as "<base>-scheduler".
func Scheduler(base string) string {
	return DNSName(Truncate("%s-scheduler", 63, base))
}

// SchedulerContainer returns the standard container name for the Kubernetes scheduler component.
func SchedulerContainer() string {
	return "scheduler"
}

// It truncates the formatted name to a maximum of 63 characters and ensures DNS name compatibility.
func ServiceAccountCertificate(base string) string {
	return DNSName(Truncate("%s-sa", 63, base))
}

// The resulting name is converted to a valid DNS name format.
func Node(base string) string {
	return DNSName(Truncate("%s-node", 63, base))
}

// This function provides a consistent identifier for the base node container across the system.
func NodeBaseContainer() string {
	return "base"
}
