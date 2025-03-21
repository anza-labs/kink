---
title: Introduction
weight: 0
---

## Installing the Operator

To install the `kink`, run the following commands. This will ensure you're always pulling the latest stable release from the operator’s GitHub repository.

!!! note
    This software requires Cert Manager to manage SSL/TLS certificates. Cert Manager automates the creation, renewal, and management of certificates within Kubernetes clusters.

    For installation instructions, please visit the official documentation: [Cert Manager Installation](https://cert-manager.io/docs/installation/).

### Using kustomization

This method applies the latest configuration by fetching the latest release tag from GitHub.

```sh
LATEST="$(curl -s 'https://api.github.com/repos/anza-labs/kink/releases/latest' | jq -r '.tag_name')"
kubectl apply -k "https://github.com/anza-labs/kink//config/default?ref=${LATEST}"
```

<!-- ### Using release manifests

Alternatively, you can deploy the operator using the release manifest directly from GitHub.

```sh
LATEST="$(curl -s 'https://api.github.com/repos/anza-labs/kink/releases/latest' | jq -r '.tag_name')"
kubectl apply -f "https://github.com/anza-labs/kink/releases/download/${LATEST}/registry-operator.yaml"
``` -->

## Updating the Operator

To update to the latest version, rerun the installation command for your chosen method.
