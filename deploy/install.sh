#!/usr/bin/env bash
set -euo pipefail

echo
echo "[1] Checking and installing Minikube..."
echo

if ! command -v minikube &> /dev/null; then
    curl -LO https://github.com/kubernetes/minikube/releases/latest/download/minikube-linux-amd64
    sudo install minikube-linux-amd64 /usr/local/bin/minikube
    rm minikube-linux-amd64
    echo "✓ Minikube installed successfully"
else
    echo "✓ Minikube is already installed"
fi

echo
echo "[2] Checking Minikube cluster status..."
echo

if minikube status | grep -q "Running"; then
    echo "✓ Minikube cluster is already running"
else
    echo "Starting Minikube cluster..."
    minikube start
    echo "✓ Minikube cluster started"
fi

echo
echo "[3] Checking and configuring Helm..."
echo

if ! command -v helm &> /dev/null; then
    echo "Helm not found. Please install Helm first."
    exit 1
fi

if helm repo list | grep -q "kedacore"; then
    echo "✓ KEDA repository already added"
else
    echo "Adding KEDA repository..."
    helm repo add kedacore https://kedacore.github.io/charts
    echo "✓ KEDA repository added"
fi

if helm repo list | grep -q "prometheus-community"; then
    echo "✓ prometheus-community repository already added"
else
    echo "Adding prometheus-community repository..."
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    echo "✓ prometheus-community repository added"
fi

echo "[4] Updating Helm repositories..."
helm repo update
echo "✓ Helm repositories updated"

echo
echo "[5] Checking and installing KEDA..."
echo

if helm status keda -n keda >/dev/null 2>&1; then
    echo "✓ KEDA is already installed"
else
    echo "Installing KEDA..."
    if ! helm install keda kedacore/keda \
            --namespace keda --create-namespace --wait --timeout 180s; then
        echo "failed to install KEDA"
        exit 1
    else
        echo "✓ KEDA installed successfully"
    fi
fi

echo
echo "[6] Checking and installing kube-prometheus-stack..."
echo

if helm status monitoring -n monitoring >/dev/null 2>&1; then
    echo "✓ Monitoring is already installed"
else
    echo "Installing Monitoring..."
    if ! helm install monitoring prometheus-community/kube-prometheus-stack \
            --namespace monitoring --create-namespace --wait --timeout 5m; then
        echo "failed to install Monitoring"
        exit 1
    else
        echo "✓ Monitoring installed successfully"
    fi
fi

echo