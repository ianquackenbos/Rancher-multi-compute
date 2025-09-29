package fleetutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConvertLabelSelectorToTargets(t *testing.T) {
	// Test with matchLabels
	selector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"kubernetes.io/os":               "linux",
			"node-role.kubernetes.io/worker": "",
		},
	}

	options := &BundleDeploymentOptions{
		DefaultNamespace: "gpu-operator",
		Helm: &HelmOptions{
			ReleaseName: "nvidia-stable",
			Repo:        "https://nvidia.github.io/helm-charts",
			Chart:       "gpu-operator",
			Values: map[string]interface{}{
				"image": map[string]string{
					"operatorTag": "v24.9.0",
					"runtimeTag":  "12.4.1",
				},
			},
		},
	}

	targets := ConvertLabelSelectorToTargets(selector, options)

	assert.Len(t, targets, 1)
	target := targets[0]

	// Check cluster selector
	assert.Equal(t, "linux", target.ClusterSelector.MatchLabels["kubernetes.io/os"])
	assert.Equal(t, "", target.ClusterSelector.MatchLabels["node-role.kubernetes.io/worker"])

	// Check bundle deployment options
	assert.Equal(t, "gpu-operator", target.BundleDeploymentOptions.DefaultNamespace)
	assert.NotNil(t, target.BundleDeploymentOptions.Helm)
	assert.Equal(t, "nvidia-stable", target.BundleDeploymentOptions.Helm.ReleaseName)
	assert.Equal(t, "https://nvidia.github.io/helm-charts", target.BundleDeploymentOptions.Helm.Repo)
	assert.Equal(t, "gpu-operator", target.BundleDeploymentOptions.Helm.Chart)
}

func TestConvertLabelSelectorToTargets_WithMatchExpressions(t *testing.T) {
	// Test with matchExpressions
	selector := metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "kubernetes.io/os",
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{"linux", "windows"},
			},
		},
	}

	options := &BundleDeploymentOptions{
		DefaultNamespace: "test-namespace",
	}

	targets := ConvertLabelSelectorToTargets(selector, options)

	assert.Len(t, targets, 1)
	target := targets[0]

	// Check cluster selector
	assert.Len(t, target.ClusterSelector.MatchExpressions, 1)
	req := target.ClusterSelector.MatchExpressions[0]
	assert.Equal(t, "kubernetes.io/os", req.Key)
	assert.Equal(t, metav1.LabelSelectorOpIn, req.Operator)
	assert.Equal(t, []string{"linux", "windows"}, req.Values)
}

func TestConvertLabelSelectorToTargets_EmptySelector(t *testing.T) {
	// Test with empty selector
	selector := metav1.LabelSelector{}

	options := &BundleDeploymentOptions{
		DefaultNamespace: "default-namespace",
	}

	targets := ConvertLabelSelectorToTargets(selector, options)

	assert.Len(t, targets, 1)
	target := targets[0]

	// Check that cluster selector is empty
	assert.Empty(t, target.ClusterSelector.MatchLabels)
	assert.Empty(t, target.ClusterSelector.MatchExpressions)
}
