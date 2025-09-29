package fleetutil

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Target represents a Fleet target
type Target struct {
	ClusterSelector         *metav1.LabelSelector    `json:"clusterSelector,omitempty"`
	BundleDeploymentOptions *BundleDeploymentOptions `json:"bundleDeploymentOptions,omitempty"`
}

// BundleDeploymentOptions represents Fleet bundle deployment options
type BundleDeploymentOptions struct {
	DefaultNamespace string       `json:"defaultNamespace,omitempty"`
	Helm             *HelmOptions `json:"helm,omitempty"`
}

// HelmOptions represents Helm deployment options
type HelmOptions struct {
	ReleaseName string                 `json:"releaseName,omitempty"`
	Repo        string                 `json:"repo,omitempty"`
	Chart       string                 `json:"chart,omitempty"`
	Values      map[string]interface{} `json:"values,omitempty"`
}

// ConvertLabelSelectorToTargets converts a LabelSelector to Fleet targets
func ConvertLabelSelectorToTargets(selector metav1.LabelSelector, options *BundleDeploymentOptions) []Target {
	return []Target{
		{
			ClusterSelector:         &selector,
			BundleDeploymentOptions: options,
		},
	}
}
