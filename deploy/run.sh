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
echo "[2] Loading docker image into minikube..."
echo

if ! minikube image load worker-pool:latest ; then
    echo "Failed to load image into minikube!"
    exit 1
else
    echo "✓ The docker image was successfully load into minikube"
fi

echo
echo "[3] Applying kubernetes manifests..."
echo

if ! kubectl apply -f ./deploy/; then
    echo "Failed to apply manifests!"
    exit 1
else
    echo "✓ Manifests successfully applied - application is now running"
fi