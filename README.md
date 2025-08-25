# Kind internal TLS

A proof-of-concept for using one-way TLS on cluster-local traffic.

Using:
- [cert-manager](https://cert-manager.io/docs/) to issue certificates and manage a private CA
- [kind](https://kind.sigs.k8s.io/) to run a local Kubernetes cluster
- [trust-manager](https://cert-manager.io/docs/trust/) to distribute the CA to all namespaces

## Prerequisites
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [helm](https://helm.sh/docs/intro/install/)
- [docker](https://docs.docker.com/get-docker/) (or another container runtime)

## Getting started
A quick setup is available via [bootstrap.sh](./bootstrap.sh):

```bash
./bootstrap.sh
```

This script does the following:
- Creates a kind cluster if none already exists
- Builds the `client` and `server` docker images (source code available in the `cmd/` folder)
- Loads the images into the kind cluster
- Installs cert-manager and trust-manager via Helm
- Deploys the `client` and `server` applications, alongside the necessary TLS resources
