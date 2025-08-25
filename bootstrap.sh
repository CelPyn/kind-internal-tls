#!/bin/bash
set -e

ACTIVE_CLUSTERS=$(kind get clusters -q | wc -l)
if [ "$ACTIVE_CLUSTERS" -eq 0 ]; then
  echo "No active kind clusters found. Creating a new kind cluster..."
  kind create cluster --config cluster.yaml
else
  echo "Active kind cluster(s) found. Skipping cluster creation."
fi

echo "Building Docker images and making them available to kind..."
docker build -t kind-tls-server:v1.0.0 --build-arg VARIANT=server .
docker build -t kind-tls-client:v1.0.0 --build-arg VARIANT=client .
kind load docker-image kind-tls-server:v1.0.0
kind load docker-image kind-tls-client:v1.0.0

echo "Setting kubectl context to kind-kind - you know, just in case..."
kubectl config use-context kind-kind

echo "Installing cert-manager..."
helm repo add jetstack https://charts.jetstack.io --force-update
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.18.2 \
  --set crds.enabled=true
helm upgrade trust-manager jetstack/trust-manager \
  --install \
  --namespace cert-manager \
  --wait

echo "Applying deployment manifests..."
kubectl apply -f deploy/namespace.yaml
kubectl apply -f deploy/
