#!/usr/bin/env bash
set -euo pipefail

echo
echo "[1] Deleting kubernetes manifests..."
echo

if ! kubectl delete -f deploy/ --ignore-not-found=true --wait=true --timeout=180s; then
    echo "failed to delete kubernetes manifests"
    exit 1
else
    echo "✓ Deleting kubernetes manifests successfully"
fi

echo
echo "[2] Deleting minikube..."
echo

if ! minikube delete -p minikube; then
    echo "failed to delete Minikube"
    exit 1
else
    echo "✓ Deleting Minikube successfully"
fi