# Rancher Multi-Compute Support Matrix

## Supported Kubernetes Versions

| Kubernetes Version | Status | Notes |
|-------------------|--------|-------|
| 1.24.x | ✅ Supported | Recommended minimum |
| 1.25.x | ✅ Supported | |
| 1.26.x | ✅ Supported | |
| 1.27.x | ✅ Supported | |
| 1.28.x | ✅ Supported | Latest stable |

## Supported GPU Vendors

### NVIDIA

| Component | Version | Status | Notes |
|-----------|---------|--------|-------|
| GPU Operator | v24.9.0 | ✅ Stable | Latest stable release |
| GPU Operator | v24.8.0 | ✅ LTS | Long-term support |
| GPU Operator | v25.0.0-rc1 | 🧪 Canary | Release candidate |
| Driver | 12.4.1 | ✅ Supported | |
| CUDA | 12.4 | ✅ Supported | |

### AMD

| Component | Version | Status | Notes |
|-----------|---------|--------|-------|
| ROCm Operator | v1.0.0 | ✅ Stable | |
| ROCm Operator | v0.9.0 | ✅ LTS | |
| ROCm Operator | v1.1.0-rc1 | 🧪 Canary | |
| ROCm | 5.7.1 | ✅ Supported | |
| ROCm | 5.6.0 | ✅ Supported | |

### Intel

| Component | Version | Status | Notes |
|-----------|---------|--------|-------|
| GPU Plugin | v0.4.0 | ✅ Stable | |
| GPU Plugin | v0.3.0 | ✅ LTS | |
| GPU Plugin | v0.5.0-rc1 | 🧪 Canary | |
| Driver | 23.3.0 | ✅ Supported | |
| Driver | 23.2.0 | ✅ Supported | |

## Supported Rancher Versions

| Rancher Version | Status | Notes |
|----------------|--------|-------|
| 2.7.x | ✅ Supported | Minimum required |
| 2.8.x | ✅ Supported | Recommended |

## Supported Fleet Versions

| Fleet Version | Status | Notes |
|---------------|--------|-------|
| 0.9.x | ✅ Supported | |
| 0.10.x | ✅ Supported | |
| 0.11.x | ✅ Supported | Recommended |

## Policy Engine Support

| Policy Engine | Version | Status | Notes |
|---------------|---------|--------|-------|
| Gatekeeper | 3.14+ | ✅ Supported | |
| Kyverno | 1.10+ | ✅ Supported | |

## Operating Systems

| OS | Version | Status | Notes |
|----|---------|--------|-------|
| Ubuntu | 20.04 LTS | ✅ Supported | |
| Ubuntu | 22.04 LTS | ✅ Supported | Recommended |
| RHEL | 8.x | ✅ Supported | |
| RHEL | 9.x | ✅ Supported | |
| SLES | 15 SP4+ | ✅ Supported | |

## Cloud Providers

| Provider | Status | Notes |
|----------|--------|-------|
| AWS | ✅ Supported | EKS, EC2 with GPU instances |
| Azure | ✅ Supported | AKS, Azure VM with GPU |
| GCP | ✅ Supported | GKE, Compute Engine with GPU |
| On-Premises | ✅ Supported | Bare metal, VMware, OpenStack |

## Architecture Support

| Architecture | Status | Notes |
|--------------|--------|-------|
| x86_64 | ✅ Supported | Primary architecture |
| ARM64 | 🧪 Experimental | Limited GPU support |
| PowerPC | ❌ Not Supported | |

## Container Runtimes

| Runtime | Version | Status | Notes |
|---------|---------|--------|-------|
| containerd | 1.6+ | ✅ Supported | Recommended |
| Docker | 20.10+ | ✅ Supported | |
| CRI-O | 1.24+ | ✅ Supported | |

## Network Plugins

| Plugin | Status | Notes |
|--------|--------|-------|
| Calico | ✅ Supported | |
| Flannel | ✅ Supported | |
| Weave | ✅ Supported | |
| Cilium | ✅ Supported | |
| Antrea | ✅ Supported | |
