package vendors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSources(t *testing.T) {
	sources := DefaultSources()

	assert.NotNil(t, sources)
	assert.Contains(t, sources, VendorNVIDIA)
	assert.Contains(t, sources, VendorAMD)
	assert.Contains(t, sources, VendorIntel)

	// Check NVIDIA source
	nvidiaSource := sources[VendorNVIDIA]
	assert.Equal(t, "https://nvidia.github.io/helm-charts", nvidiaSource.Repo)
	assert.Equal(t, "gpu-operator", nvidiaSource.Chart)
	assert.Equal(t, "gpu-operator", nvidiaSource.Namespace)

	// Check AMD source
	amdSource := sources[VendorAMD]
	assert.Equal(t, "https://rocm.github.io/helm-charts", amdSource.Repo)
	assert.Equal(t, "rocm-device-plugin", amdSource.Chart)
	assert.Equal(t, "rocm-system", amdSource.Namespace)

	// Check Intel source
	intelSource := sources[VendorIntel]
	assert.Equal(t, "https://intel.github.io/helm-charts", intelSource.Repo)
	assert.Equal(t, "intel-gpu-plugin", intelSource.Chart)
	assert.Equal(t, "intel-gpu", intelSource.Namespace)
}

func TestVendorString(t *testing.T) {
	assert.Equal(t, "nvidia", string(VendorNVIDIA))
	assert.Equal(t, "amd", string(VendorAMD))
	assert.Equal(t, "intel", string(VendorIntel))
}
