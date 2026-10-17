#!/usr/bin/env bash
set -euo pipefail

echo
echo "[1] Building Docker image..."
echo

if ! docker build -t worker-pool:latest -f ./deploy/Dockerfile .; then
    echo "Docker build failed"
    exit 1
else
    echo "✓ Docker image successfully built"
fi

echo
echo "[2] Deleting kubernetes manifests (if it was launched previously)..."
echo

kubectl delete -f deploy/ --ignore-not-found=true --wait=true --timeout=180s

echo
echo "[3] Loading docker image into minikube..."
echo

if ! minikube image load worker-pool:latest ; then
    echo "Failed to load image into minikube!"
    exit 1
else
    echo "✓ The docker image was successfully load into minikube"
fi

echo
echo "[4] Applying kubernetes namespaces..."
echo

if ! (kubectl apply \
        -f ./deploy/namespace-apps.yaml \
        -f ./deploy/namespace-monitoring.yaml --wait=true --timeout=180s); then
    echo "Failed to apply namespaces!"
    exit 1
else
    echo "✓ Namespaces successfully applied"
fi

echo
echo "[5] Applying kubernetes manifests..."
echo

if ! kubectl apply -f ./deploy/ --wait=true --timeout=180s; then
    echo "Failed to apply manifests!"
    exit 1
else
    echo "✓ Manifests successfully applied - application is now running"
fi

MINIKUBE_IP=$(minikube ip)

echo
echo "Prometheus:  http://${MINIKUBE_IP}:30081"
echo "Grafana:     http://${MINIKUBE_IP}:30082"
echo