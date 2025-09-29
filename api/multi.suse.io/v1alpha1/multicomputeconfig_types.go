package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MultiComputeConfigSpec defines the desired state of MultiComputeConfig
type MultiComputeConfigSpec struct {
	// Policies defines which policies to enable
	Policies PolicyConfig `json:"policies,omitempty"`

	// VendorSources allows overriding default vendor configurations
	VendorSources map[string]VendorSource `json:"vendorSources,omitempty"`
}

// PolicyConfig defines policy settings
type PolicyConfig struct {
	// EnforceRuntimeClass enables runtime class enforcement
	EnforceRuntimeClass bool `json:"enforceRuntimeClass,omitempty"`

	// RestrictGPUNamespaces restricts GPU workloads to specific namespaces
	RestrictGPUNamespaces bool `json:"restrictGPUNamespaces,omitempty"`

	// RequireCosign enables image signature verification
	RequireCosign bool `json:"requireCosign,omitempty"`

	// LimitGPUsPerPod sets maximum GPUs per pod
	LimitGPUsPerPod int32 `json:"limitGPUsPerPod,omitempty"`
}

// VendorSource defines vendor-specific configuration
type VendorSource struct {
	// Repo is the Helm repository URL
	Repo string `json:"repo"`

	// Chart is the Helm chart name
	Chart string `json:"chart"`

	// Namespace is the target namespace for deployment
	Namespace string `json:"namespace"`
}

// MultiComputeConfigStatus defines the observed state of MultiComputeConfig
type MultiComputeConfigStatus struct {
	// Conditions represent the latest available observations
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// MultiComputeConfig is the Schema for the multicomputeconfigs API
type MultiComputeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MultiComputeConfigSpec   `json:"spec,omitempty"`
	Status MultiComputeConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MultiComputeConfigList contains a list of MultiComputeConfig
type MultiComputeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MultiComputeConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MultiComputeConfig{}, &MultiComputeConfigList{})
}
