package versions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Pins represents version pins for a vendor
type Pins struct {
	OperatorTag string `yaml:"operatorTag"`
	RuntimeTag  string `yaml:"runtimeTag"`
}

// VendorPins represents version pins for all vendors
type VendorPins struct {
	NVIDIA Pins `yaml:"nvidia"`
	AMD    Pins `yaml:"amd"`
	Intel  Pins `yaml:"intel"`
}

// Resolver interface for resolving versions
type Resolver interface {
	Resolve(ctx context.Context, channel string) (VendorPins, error)
}

// FileResolver resolves versions from VERSION.yaml files
type FileResolver struct {
	VersionDir string
}

// NewFileResolver creates a new FileResolver
func NewFileResolver(versionDir string) *FileResolver {
	return &FileResolver{
		VersionDir: versionDir,
	}
}

// Resolve resolves versions for a given channel
func (r *FileResolver) Resolve(ctx context.Context, channel string) (VendorPins, error) {
	versionFile := filepath.Join(r.VersionDir, channel, "VERSION.yaml")

	data, err := os.ReadFile(versionFile)
	if err != nil {
		return VendorPins{}, fmt.Errorf("failed to read version file %s: %w", versionFile, err)
	}

	var pins VendorPins
	if err := yaml.Unmarshal(data, &pins); err != nil {
		return VendorPins{}, fmt.Errorf("failed to unmarshal version file: %w", err)
	}

	return pins, nil
}

// LoadSources loads vendor sources from ConfigMap
func LoadSources(ctx context.Context, configMapData map[string]string) (map[string]interface{}, error) {
	sources := make(map[string]interface{})

	for vendor, data := range configMapData {
		var source interface{}
		if err := yaml.Unmarshal([]byte(data), &source); err != nil {
			return nil, fmt.Errorf("failed to unmarshal source for vendor %s: %w", vendor, err)
		}
		sources[vendor] = source
	}

	return sources, nil
}
