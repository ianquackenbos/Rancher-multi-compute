package vendors

// Vendor represents a GPU vendor
type Vendor string

const (
	VendorNVIDIA Vendor = "nvidia"
	VendorAMD    Vendor = "amd"
	VendorIntel  Vendor = "intel"
)

// Source represents a vendor's Helm chart source
type Source struct {
	Repo      string `json:"repo"`
	Chart     string `json:"chart"`
	Namespace string `json:"namespace"`
}

// DefaultSources returns the default vendor sources
func DefaultSources() map[Vendor]Source {
	return map[Vendor]Source{
		VendorNVIDIA: {
			Repo:      "https://nvidia.github.io/helm-charts",
			Chart:     "gpu-operator",
			Namespace: "gpu-operator",
		},
		VendorAMD: {
			Repo:      "https://rocm.github.io/helm-charts",
			Chart:     "rocm-device-plugin",
			Namespace: "rocm-system",
		},
		VendorIntel: {
			Repo:      "https://intel.github.io/helm-charts",
			Chart:     "intel-gpu-plugin",
			Namespace: "intel-gpu",
		},
	}
}
