package versions

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileResolver_Resolve(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()

	// Create stable channel directory
	stableDir := filepath.Join(tempDir, "stable")
	err := os.MkdirAll(stableDir, 0755)
	require.NoError(t, err)

	// Create VERSION.yaml file
	versionFile := filepath.Join(stableDir, "VERSION.yaml")
	versionContent := `
nvidia:
  operatorTag: "v24.9.0"
  runtimeTag: "12.4.1"
amd:
  operatorTag: "v24.9.0"
  runtimeTag: "5.7.1"
intel:
  operatorTag: "v24.9.0"
  runtimeTag: "24.16.0"
`
	err = os.WriteFile(versionFile, []byte(versionContent), 0644)
	require.NoError(t, err)

	// Create resolver
	resolver := NewFileResolver(tempDir)

	// Test resolve
	ctx := context.Background()
	pins, err := resolver.Resolve(ctx, "stable")
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "v24.9.0", pins.NVIDIA.OperatorTag)
	assert.Equal(t, "12.4.1", pins.NVIDIA.RuntimeTag)
	assert.Equal(t, "v24.9.0", pins.AMD.OperatorTag)
	assert.Equal(t, "5.7.1", pins.AMD.RuntimeTag)
	assert.Equal(t, "v24.9.0", pins.Intel.OperatorTag)
	assert.Equal(t, "24.16.0", pins.Intel.RuntimeTag)
}

func TestFileResolver_Resolve_NonExistentChannel(t *testing.T) {
	tempDir := t.TempDir()
	resolver := NewFileResolver(tempDir)

	ctx := context.Background()
	_, err := resolver.Resolve(ctx, "nonexistent")
	assert.Error(t, err)
}

func TestFileResolver_Resolve_InvalidYAML(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	
	// Create stable channel directory
	stableDir := filepath.Join(tempDir, "stable")
	err := os.MkdirAll(stableDir, 0755)
	require.NoError(t, err)
	
	// Create invalid VERSION.yaml file
	versionFile := filepath.Join(stableDir, "VERSION.yaml")
	versionContent := `invalid: yaml: content: [`
	err = os.WriteFile(versionFile, []byte(versionContent), 0644)
	require.NoError(t, err)
	
	// Create resolver
	resolver := NewFileResolver(tempDir)
	
	// Test resolve
	ctx := context.Background()
	_, err = resolver.Resolve(ctx, "stable")
	assert.Error(t, err)
}

func TestLoadSources(t *testing.T) {
	ctx := context.Background()
	
	// Test with valid YAML data
	configMapData := map[string]string{
		"nvidia": `
repo: "https://nvidia.github.io/helm-charts"
chart: "gpu-operator"
namespace: "gpu-operator"
`,
		"amd": `
repo: "https://rocm.github.io/helm-charts"
chart: "rocm-device-plugin"
namespace: "rocm-system"
`,
	}
	
	sources, err := LoadSources(ctx, configMapData)
	require.NoError(t, err)
	assert.NotNil(t, sources)
	assert.Contains(t, sources, "nvidia")
	assert.Contains(t, sources, "amd")
	
	// Verify nvidia source structure
	nvidiaSource := sources["nvidia"].(map[string]interface{})
	assert.Equal(t, "https://nvidia.github.io/helm-charts", nvidiaSource["repo"])
	assert.Equal(t, "gpu-operator", nvidiaSource["chart"])
	assert.Equal(t, "gpu-operator", nvidiaSource["namespace"])
}

func TestLoadSources_InvalidYAML(t *testing.T) {
	ctx := context.Background()
	
	// Test with invalid YAML data
	configMapData := map[string]string{
		"nvidia": `invalid: yaml: content: [`,
	}
	
	_, err := LoadSources(ctx, configMapData)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal source for vendor nvidia")
}

func TestLoadSources_EmptyData(t *testing.T) {
	ctx := context.Background()
	
	// Test with empty data
	configMapData := map[string]string{}
	
	sources, err := LoadSources(ctx, configMapData)
	require.NoError(t, err)
	assert.NotNil(t, sources)
	assert.Empty(t, sources)
}
