#!/usr/bin/env bash
set -euo pipefail

# e2e-kind.sh - End-to-end testing with KIND

CLUSTER_NAME="rmc-e2e"
KIND_VERSION="v0.20.0"

red() { printf "\033[31m%s\033[0m\n" "$*"; }
grn() { printf "\033[32m%s\033[0m\n" "$*"; }
ylw() { printf "\033[33m%s\033[0m\n" "$*"; }

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*"
}

cleanup() {
    log "Cleaning up..."
    kind delete cluster --name "$CLUSTER_NAME" 2>/dev/null || true
}

trap cleanup EXIT

# Check prerequisites
check_prereqs() {
    log "Checking prerequisites..."
    
    if ! command -v kind &> /dev/null; then
        red "kind is not installed. Please install kind $KIND_VERSION"
        exit 1
    fi
    
    if ! command -v kubectl &> /dev/null; then
        red "kubectl is not installed"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        red "docker is not installed or not running"
        exit 1
    fi
    
    grn "Prerequisites check passed"
}

# Create KIND cluster
create_cluster() {
    log "Creating KIND cluster: $CLUSTER_NAME"
    
    cat <<EOF | kind create cluster --name "$CLUSTER_NAME" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "multi.suse.io/cluster-group=e2e-test"
- role: worker
  kubeadmConfigPatches:
  - |
    kind: JoinConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "multi.suse.io/cluster-group=e2e-test"
EOF
    
    grn "KIND cluster created successfully"
}

# Install Fleet
install_fleet() {
    log "Installing Fleet..."
    
    # Install Fleet CRDs
    kubectl apply -f https://github.com/rancher/fleet/releases/download/v0.9.4/fleet-crd-0.9.4.yaml
    
    # Install Fleet
    kubectl apply -f https://github.com/rancher/fleet/releases/download/v0.9.4/fleet-0.9.4.yaml
    
    # Wait for Fleet to be ready
    kubectl wait --for=condition=available --timeout=300s deployment/fleet-controller -n cattle-fleet-system
    
    grn "Fleet installed successfully"
}

# Deploy test resources
deploy_test_resources() {
    log "Deploying test resources..."
    
    # Create test namespace
    kubectl create namespace gpu-workloads --dry-run=client -o yaml | kubectl apply -f -
    
    # Apply CRDs first
    kubectl apply -f config/crd/bases/multi.suse.io_channels.yaml
    kubectl apply -f config/crd/bases/multi.suse.io_multicomputeconfigs.yaml
    
    # Wait for CRDs to be established
    kubectl wait --for condition=established --timeout=60s crd/channels.multi.suse.io
    kubectl wait --for condition=established --timeout=60s crd/multicomputeconfigs.multi.suse.io
    
    # Apply test Channel
    cat <<EOF | kubectl apply -f -
apiVersion: multi.suse.io/v1alpha1
kind: Channel
metadata:
  name: nvidia-stable
spec:
  vendor: nvidia
  channel: stable
  clusterSelector:
    matchLabels:
      multi.suse.io/cluster-group: e2e-test
EOF
    
    grn "Test resources deployed"
}

# Run tests
run_tests() {
    log "Running e2e tests..."
    
    # Wait for Channel to be processed (with longer timeout)
    kubectl wait --for=condition=Ready --timeout=300s channel/nvidia-stable || {
        red "Channel not ready after 5 minutes"
        kubectl describe channel nvidia-stable
        return 1
    }
    
    # Check if Channel status is updated
    if kubectl get channel nvidia-stable -o jsonpath='{.status.phase}' | grep -q "RollingOut\|Completed"; then
        grn "Channel status updated correctly"
    else
        red "Channel status not updated"
        kubectl describe channel nvidia-stable
        return 1
    fi
    
    grn "E2E tests passed"
}

# Deploy controller
deploy_controller() {
    log "Deploying auto-operator controller..."
    
    # Deploy the controller
    kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: auto-operator-controller
  namespace: cattle-fleet-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: auto-operator-controller
  template:
    metadata:
      labels:
        app: auto-operator-controller
    spec:
      containers:
      - name: controller
        image: ghcr.io/ianquackenbos/rancher-multi-compute/auto-operator-controller:latest
        command: ["./auto-operator-controller"]
        env:
        - name: VERSION_DIR
          value: "./fleet/overlays"
        - name: VENDOR_SOURCES_CM
          value: "vendor-sources"
EOF
    
    # Wait for controller to be ready
    kubectl wait --for=condition=available --timeout=120s deployment/auto-operator-controller -n cattle-fleet-system
    
    grn "Controller deployed successfully"
}

# Main execution
main() {
    log "Starting Rancher Multi-Compute E2E tests with KIND"
    
    check_prereqs
    create_cluster
    install_fleet
    deploy_test_resources
    deploy_controller
    run_tests
    
    grn "All E2E tests completed successfully!"
}

main "$@"