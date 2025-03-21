# KinK

[![GitHub License](https://img.shields.io/github/license/anza-labs/kink)][license]
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](code_of_conduct.md)
[![GitHub issues](https://img.shields.io/github/issues/anza-labs/kink)](https://github.com/anza-labs/kink/issues)
[![GitHub release](https://img.shields.io/github/release/anza-labs/kink)](https://GitHub.com/anza-labs/kink/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/anza-labs/kink)](https://goreportcard.com/report/github.com/anza-labs/kink)

<p align="center">
  <img src="assets/kink.png" width="256p"/>
</p>

| Component               | Image                                     | Mode       |
| ----------------------- | ----------------------------------------- | ---------- |
| Kube Scheduler          | `registry.k8s.io/kube-scheduler`          | Deployment |
| Kube API-Server         | `registry.k8s.io/kube-apiserver`          | Deployment |
| Kube Controller Manager | `registry.k8s.io/kube-controller-manager` | Deployment |
| ETCD                    | `ghcr.io/anza-labs/library/kine`          | Deployment |

```yaml
kubectl apply -f - <<EOF
---
apiVersion: v1
kind: Namespace
metadata:
  name: example
---
apiVersion: controlplane.cluster.x-k8s.io/v1alpha1
kind: KinkControlPlane
metadata:
  name: example
  namespace: example
spec:
  version: v1.32.0
EOF
```

## License

`kink` is licensed under the [Apache-2.0][license].

<!-- Resources -->

[license]: https://github.com/anza-labs/kink/blob/main/LICENSE
