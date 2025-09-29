# vLLM NVIDIA Example

This example demonstrates how to deploy vLLM inference workloads with NVIDIA GPUs using Rancher Multi-Compute.

## Prerequisites

- NVIDIA GPU operator installed via Rancher Multi-Compute
- Nodes labeled with `compute.multi.suse.io/vendor=nvidia`
- Sufficient GPU memory for the model

## Features

- **GPU Resource Management**: Automatically requests NVIDIA GPUs
- **MIG Support**: Optional Multi-Instance GPU configuration
- **Node Selection**: Targets NVIDIA-enabled nodes
- **Health Checks**: Built-in liveness and readiness probes

## Configuration

### Basic Deployment

```bash
helm install vllm-nvidia ./examples/vllm-nvidia
```

### With MIG Profile

```yaml
mig:
  enabled: true
  profile: "1g.5gb"
```

### Custom Model

```yaml
vllm:
  model: "microsoft/DialoGPT-large"
  maxModelLen: 4096
  gpuMemoryUtilization: 0.8
```

## Values

| Parameter | Description | Default |
|-----------|-------------|---------|
| `image.repository` | vLLM container image | `vllm/vllm-openai` |
| `image.tag` | Image tag | `v0.2.7` |
| `resources.limits.nvidia.com/gpu` | GPU limit | `1` |
| `mig.enabled` | Enable MIG support | `false` |
| `mig.profile` | MIG profile | `1g.5gb` |
| `vllm.model` | Model to load | `microsoft/DialoGPT-medium` |
| `vllm.gpuMemoryUtilization` | GPU memory usage | `0.9` |

## Usage

Once deployed, the vLLM service will be available at:

```bash
# Get the service URL
kubectl get svc vllm-nvidia

# Test the API
curl http://<service-ip>:8000/v1/models
```

## Troubleshooting

1. **GPU not found**: Ensure nodes have NVIDIA GPUs and are labeled correctly
2. **Out of memory**: Reduce `gpuMemoryUtilization` or use a smaller model
3. **MIG issues**: Verify MIG is configured on the GPU nodes
