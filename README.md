# Rancher Multi-Compute

Rancher Multi-Compute provides unified management of GPU/accelerator stacks across multi-cluster environments using Fleet.

## Overview

This project orchestrates vendor GPU/accelerator stacks (NVIDIA, AMD, Intel) across multi-cluster environments, providing:

- **Unified GPU Management**: Single interface for all GPU vendors
- **Fleet Integration**: GitOps-based deployment using Rancher Fleet
- **Policy Enforcement**: Security and resource policies via Gatekeeper/Kyverno
- **Rancher UI Extension**: Dashboard integration for monitoring and management
- **Multi-Cluster Support**: Deploy across multiple Kubernetes clusters

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Rancher UI    │    │   Controllers   │    │   Fleet Bundles │
│   Extension     │    │                 │    │                 │
│                 │    │ • Profiler      │    │ • NVIDIA Stack  │
│ • Channel Mgmt  │    │ • Auto-Operator │    │ • AMD Stack     │
│ • Status View   │    │ • Drift Detector│    │ • Intel Stack   │
│ • Policy Config │    │ • Policy Ctrl   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Multi-Cluster │
                    │   Environment   │
                    │                 │
                    │ • Cluster A     │
                    │ • Cluster B     │
                    │ • Cluster C     │
                    └─────────────────┘
```

## Components

### Controllers

- **compute-profiler-controller**: Discovers GPU capabilities using Node Feature Discovery
- **compute-auto-operator-controller**: Manages Fleet bundles for vendor GPU operators
- **compute-drift-detector**: Monitors configuration drift across clusters
- **policy-controller**: Enforces security and resource policies

### Custom Resources

- **Channel**: Defines vendor and release channel for GPU operator deployment
- **MultiComputeConfig**: Global configuration for policies and vendor sources

### Fleet Bundles

- **NVIDIA**: GPU Operator with CUDA support
- **AMD**: ROCm device plugin and runtime
- **Intel**: Intel GPU plugin and oneAPI support

## Quick Start

### Prerequisites

- Kubernetes 1.24+
- Rancher 2.7+
- Fleet installed
- Gatekeeper or Kyverno for policy enforcement

### Installation

1. **Install CRDs**:
   ```bash
   kubectl apply -f config/crd/bases/
   ```

2. **Deploy Controllers**:
   ```bash
   kubectl apply -f config/default/
   ```

3. **Configure Version Channels**:
   ```bash
   kubectl apply -f fleet/overlays/stable/VERSION.yaml
   ```

### Create a Channel

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

### Deploy Example Workload

```bash
# Deploy vLLM with NVIDIA GPUs
helm install vllm-nvidia ./examples/vllm-nvidia

# Deploy PyTorch with AMD ROCm
helm install pytorch-rocm ./examples/pytorch-rocm
```

## Configuration

### Version Channels

- **stable**: Production-ready releases
- **lts**: Long-term support versions
- **canary**: Latest development builds

### Policy Configuration

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

## Examples

### NVIDIA vLLM Inference

```yaml
# examples/vllm-nvidia/values.yaml
resources:
  limits:
    nvidia.com/gpu: 1
nodeSelector:
  compute.multi.suse.io/vendor: nvidia
mig:
  enabled: true
  profile: "1g.5gb"
```

### AMD PyTorch Training

```yaml
# examples/pytorch-rocm/values.yaml
resources:
  limits:
    amd.com/gpu: 1
nodeSelector:
  compute.multi.suse.io/vendor: amd
runtimeClassName: rocm
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Linting

```bash
make lint
```

### E2E Testing

```bash
./hack/e2e-kind.sh
```

## Documentation

- [Operations Guide](docs/operations.md) - Deployment and troubleshooting
- [Support Matrix](docs/support-matrix.md) - Compatibility information
- [Examples](examples/) - Workload examples and Helm charts

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/suse/rancher-multi-compute/issues)
- **Discussions**: [GitHub Discussions](https://github.com/suse/rancher-multi-compute/discussions)
- **Documentation**: [Project Wiki](https://github.com/suse/rancher-multi-compute/wiki)
