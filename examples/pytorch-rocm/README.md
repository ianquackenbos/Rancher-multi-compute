# PyTorch ROCm Example

This example demonstrates how to deploy PyTorch training workloads with AMD ROCm using Rancher Multi-Compute.

## Prerequisites

- AMD ROCm operator installed via Rancher Multi-Compute
- Nodes labeled with `compute.multi.suse.io/vendor=amd`
- ROCm-compatible AMD GPUs

## Features

- **ROCm Runtime**: Uses `rocm` runtime class for GPU access
- **AMD GPU Resources**: Requests `amd.com/gpu` resources
- **Node Selection**: Targets AMD-enabled nodes
- **Distributed Training**: Configurable for multi-GPU training

## Configuration

### Basic Deployment

```bash
helm install pytorch-rocm ./examples/pytorch-rocm
```

### Distributed Training

```yaml
pytorch:
  distributed: true
  worldSize: 4
  backend: "nccl"
```

### Custom ROCm Configuration

```yaml
rocm:
  visibleDevices: "0,1"
  debugLevel: 1
```

## Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | PyTorch ROCm container image | `rocm/pytorch` |
| `image.tag` | Image tag | `rocm5.7_ubuntu22.04_py3.11_pytorch_2.1.0` |
| `resources.limits.amd.com/gpu` | GPU limit | `1` |
| `pytorch.distributed` | Enable distributed training | `false` |
| `pytorch.backend` | Distributed backend | `nccl` |
| `rocm.visibleDevices` | Visible GPU devices | `all` |

## Usage

Once deployed, check the logs to verify GPU detection:

```bash
# Get pod logs
kubectl logs deployment/pytorch-rocm

# Expected output:
# PyTorch version: 2.1.0+rocm5.7
# ROCm available: True
# GPU count: 1
# GPU 0: AMD Radeon RX 7900 XTX
```

## Troubleshooting

1. **ROCm not available**: Ensure ROCm operator is installed and nodes have compatible GPUs
2. **Runtime class issues**: Verify `rocm` runtime class exists
3. **GPU not detected**: Check node labels and GPU driver installation
4. **Distributed training**: Ensure all nodes have the same GPU configuration
