# Rancher Multi-Compute Operations Guide

## Overview

Rancher Multi-Compute provides unified management of GPU/accelerator stacks across multi-cluster environments.

## Core Components

### Controllers

1. **compute-profiler-controller**: Discovers and profiles GPU capabilities using Node Feature Discovery
2. **compute-auto-operator-controller**: Manages Fleet bundles for vendor GPU operators
3. **compute-drift-detector**: Monitors configuration drift across clusters
4. **policy-controller**: Enforces security and resource policies

### Custom Resources

- **Channel**: Defines vendor and release channel for GPU operator deployment
- **MultiComputeConfig**: Global configuration for policies and vendor sources

## Deployment

### Prerequisites

- Kubernetes 1.24+
- Rancher 2.7+
- Fleet installed
- Gatekeeper or Kyverno for policy enforcement

### Installation

```bash
# Install CRDs
kubectl apply -f config/crd/bases/

# Install controllers
kubectl apply -f config/default/

# Configure vendor sources
kubectl apply -f fleet/overlays/stable/VERSION.yaml
```

## Configuration

### Channel Management

Create a Channel to deploy NVIDIA GPU operator:

```yaml
apiVersion: multi.suse.io/v1alpha1
kind: Channel
metadata:
  name: nvidia-stable
spec:
  vendor: nvidia
  channel: stable
  clusterSelector:
    matchLabels:
      multi.suse.io/cluster-group: gpu-clusters
```

### Policy Configuration

Enable policy enforcement:

```yaml
apiVersion: multi.suse.io/v1alpha1
kind: MultiComputeConfig
metadata:
  name: default
spec:
  policies:
    enforceRuntimeClass: true
    restrictGPUNamespaces: true
    requireCosign: true
    limitGPUsPerPod: 4
```

## Monitoring

### Status Checking

```bash
# Check Channel status
kubectl get channels

# Check controller logs
kubectl logs -n cattle-fleet-system deployment/compute-auto-operator-controller

# Check Fleet bundle status
kubectl get bundles -n cattle-fleet-system
```

### Troubleshooting

Common issues and solutions:

1. **Bundle deployment fails**: Check Fleet agent logs and network connectivity
2. **Policy violations**: Review Gatekeeper/Kyverno logs and policy configuration
3. **GPU discovery issues**: Verify Node Feature Discovery is running

## Security

- All images are signed with Cosign
- SBOMs are generated for supply chain security
- Policies enforce runtime class and resource limits
- Air-gapped deployments supported