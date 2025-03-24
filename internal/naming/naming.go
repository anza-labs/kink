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

func AdminCertificate(base string) string {
	return DNSName(Truncate("%s-admin-cert", 63, base))
}

func APIServer(base string) string {
	return DNSName(Truncate("%s-api-server", 63, base))
}

func APIServerCertificate(base string) string {
	return DNSName(Truncate("%s-api-server", 63, base))
}

func APIServerContainer() string {
	return "api-server"
}

func ClusterCA(base string) string {
	return DNSName(Truncate("%s-ca", 63, base))
}

func ConfigMap(base, hash string) string {
	return DNSName(Truncate("%s-%s", 63, base, hash))
}

func ControllerManager(base string) string {
	return DNSName(Truncate("%s-controller-manager", 63, base))
}

func ControllerManagerCertificate(base string) string {
	return DNSName(Truncate("%s-controller-manager-cert", 63, base))
}

func ControllerManagerContainer() string {
	return "controller-manager"
}

func FrontProxyCA(base string) string {
	return DNSName(Truncate("%s-proxy", 63, base))
}

func Kine(base string) string {
	return DNSName(Truncate("%s-kine", 63, base))
}

func KineAPIServerClientCertificate(base string) string {
	return DNSName(Truncate("%s-etcd-client", 63, base))
}

func KineCA(base string) string {
	return DNSName(Truncate("%s-etcd", 63, base))
}

func KineServerCertificate(base string) string {
	return DNSName(Truncate("%s-etcd-server", 63, base))
}

func KinePersistentVolumeClaim(base string) string {
	return DNSName(Truncate("%s-kine", 63, base))
}

func KineContainer() string {
	return "kine"
}

func Kubeconfig(base string) string {
	return DNSName(Truncate("%s-kubeconfig", 63, base))
}

func RootCA(base string) string {
	return DNSName(Truncate("%s-root-ca", 63, base))
}

func Scheduler(base string) string {
	return DNSName(Truncate("%s-scheduler", 63, base))
}

func SchedulerCertificate(base string) string {
	return DNSName(Truncate("%s-scheduler-cert", 63, base))
}

func SchedulerContainer() string {
	return "scheduler"
}

func ServiceAccountCertificate(base string) string {
	return DNSName(Truncate("%s-sa", 63, base))
}

func Node(base string) string {
	return DNSName(Truncate("%s-node", 63, base))
}

func NodeBaseContainer() string {
	return "base"
}
