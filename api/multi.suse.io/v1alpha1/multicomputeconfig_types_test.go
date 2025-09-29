package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPolicyConfig(t *testing.T) {
	policy := PolicyConfig{
		EnforceRuntimeClass:   true,
		RestrictGPUNamespaces: true,
		RequireCosign:         true,
		LimitGPUsPerPod:       4,
	}

	assert.True(t, policy.EnforceRuntimeClass)
	assert.True(t, policy.RestrictGPUNamespaces)
	assert.True(t, policy.RequireCosign)
	assert.Equal(t, int32(4), policy.LimitGPUsPerPod)
}

func TestVendorSource(t *testing.T) {
	source := VendorSource{
		Repo:      "https://nvidia.github.io/helm-charts",
		Chart:     "gpu-operator",
		Namespace: "gpu-operator",
	}

	assert.Equal(t, "https://nvidia.github.io/helm-charts", source.Repo)
	assert.Equal(t, "gpu-operator", source.Chart)
	assert.Equal(t, "gpu-operator", source.Namespace)
}

func TestMultiComputeConfigSpec(t *testing.T) {
	spec := MultiComputeConfigSpec{
		Policies: PolicyConfig{
			EnforceRuntimeClass: true,
			LimitGPUsPerPod:     2,
		},
		VendorSources: map[string]VendorSource{
			"nvidia": {
				Repo:      "https://nvidia.github.io/helm-charts",
				Chart:     "gpu-operator",
				Namespace: "gpu-operator",
			},
		},
	}

	assert.True(t, spec.Policies.EnforceRuntimeClass)
	assert.Equal(t, int32(2), spec.Policies.LimitGPUsPerPod)
	assert.Contains(t, spec.VendorSources, "nvidia")
	assert.Equal(t, "gpu-operator", spec.VendorSources["nvidia"].Chart)
}

func TestMultiComputeConfigStatus(t *testing.T) {
	status := MultiComputeConfigStatus{
		Conditions: []metav1.Condition{
			{
				Type:   "Ready",
				Status: metav1.ConditionTrue,
				Reason: "Reconciled",
			},
		},
	}

	assert.Len(t, status.Conditions, 1)
	assert.Equal(t, "Ready", status.Conditions[0].Type)
	assert.Equal(t, metav1.ConditionTrue, status.Conditions[0].Status)
}

func TestMultiComputeConfig(t *testing.T) {
	config := &MultiComputeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default",
		},
		Spec: MultiComputeConfigSpec{
			Policies: PolicyConfig{
				EnforceRuntimeClass: true,
			},
		},
		Status: MultiComputeConfigStatus{
			Conditions: []metav1.Condition{
				{
					Type:   "Ready",
					Status: metav1.ConditionTrue,
				},
			},
		},
	}

	assert.Equal(t, "default", config.Name)
	assert.True(t, config.Spec.Policies.EnforceRuntimeClass)
	assert.Len(t, config.Status.Conditions, 1)
}

func TestMultiComputeConfigList(t *testing.T) {
	config1 := MultiComputeConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "config1"},
		Spec:       MultiComputeConfigSpec{},
	}
	config2 := MultiComputeConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "config2"},
		Spec:       MultiComputeConfigSpec{},
	}

	list := &MultiComputeConfigList{
		Items: []MultiComputeConfig{config1, config2},
	}

	assert.Len(t, list.Items, 2)
	assert.Equal(t, "config1", list.Items[0].Name)
	assert.Equal(t, "config2", list.Items[1].Name)
}
