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

echo "[4] Updating Helm repositories..."
helm repo update
echo "✓ Helm repositories updated"

echo
echo "[5] Checking and installing KEDA..."
echo

if helm list -n keda | grep -q "keda"; then
    echo "✓ KEDA is already installed"
else
    echo "Installing KEDA..."
    helm install keda kedacore/keda --namespace keda --create-namespace
    echo "✓ KEDA installed successfully"
fi

echo