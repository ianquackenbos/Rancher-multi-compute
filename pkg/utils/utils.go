package utils

import "fmt"

// GetVersion returns the current version
func GetVersion() string {
	return "v0.1.0"
}

// ValidateConfig validates configuration
func ValidateConfig(config map[string]interface{}) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}
	return nil
}
