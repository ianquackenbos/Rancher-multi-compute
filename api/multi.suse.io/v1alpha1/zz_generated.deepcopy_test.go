package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestChannelDeepCopy(t *testing.T) {
	original := &Channel{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-channel",
		},
		Spec: ChannelSpec{
			Vendor:  "nvidia",
			Channel: "stable",
		},
		Status: ChannelStatus{
			Phase: "Progressing",
		},
	}

	// Test DeepCopy
	copied := original.DeepCopy()
	assert.NotNil(t, copied)
	assert.Equal(t, original.Name, copied.Name)
	assert.Equal(t, original.Spec.Vendor, copied.Spec.Vendor)
	assert.Equal(t, original.Spec.Channel, copied.Spec.Channel)
	assert.Equal(t, original.Status.Phase, copied.Status.Phase)

	// Test DeepCopyInto
	var target Channel
	original.DeepCopyInto(&target)
	assert.Equal(t, original.Name, target.Name)
	assert.Equal(t, original.Spec.Vendor, target.Spec.Vendor)
	assert.Equal(t, original.Spec.Channel, target.Spec.Channel)
	assert.Equal(t, original.Status.Phase, target.Status.Phase)

	// Test DeepCopyObject
	obj := original.DeepCopyObject()
	assert.NotNil(t, obj)
	channelObj, ok := obj.(*Channel)
	assert.True(t, ok)
	assert.Equal(t, original.Name, channelObj.Name)
}

func TestChannelListDeepCopy(t *testing.T) {
	original := &ChannelList{
		Items: []Channel{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "channel1"},
				Spec:       ChannelSpec{Vendor: "nvidia", Channel: "stable"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "channel2"},
				Spec:       ChannelSpec{Vendor: "amd", Channel: "stable"},
			},
		},
	}

	// Test DeepCopy
	copied := original.DeepCopy()
	assert.NotNil(t, copied)
	assert.Len(t, copied.Items, 2)
	assert.Equal(t, original.Items[0].Name, copied.Items[0].Name)
	assert.Equal(t, original.Items[1].Name, copied.Items[1].Name)

	// Test DeepCopyInto
	var target ChannelList
	original.DeepCopyInto(&target)
	assert.Len(t, target.Items, 2)
	assert.Equal(t, original.Items[0].Name, target.Items[0].Name)
	assert.Equal(t, original.Items[1].Name, target.Items[1].Name)

	// Test DeepCopyObject
	obj := original.DeepCopyObject()
	assert.NotNil(t, obj)
	listObj, ok := obj.(*ChannelList)
	assert.True(t, ok)
	assert.Len(t, listObj.Items, 2)
}

func TestMultiComputeConfigDeepCopy(t *testing.T) {
	original := &MultiComputeConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-config",
		},
		Spec: MultiComputeConfigSpec{
			Policies: PolicyConfig{
				EnforceRuntimeClass: true,
				LimitGPUsPerPod:     4,
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

	// Test DeepCopy
	copied := original.DeepCopy()
	assert.NotNil(t, copied)
	assert.Equal(t, original.Name, copied.Name)
	assert.Equal(t, original.Spec.Policies.EnforceRuntimeClass, copied.Spec.Policies.EnforceRuntimeClass)
	assert.Equal(t, original.Spec.Policies.LimitGPUsPerPod, copied.Spec.Policies.LimitGPUsPerPod)
	assert.Len(t, copied.Status.Conditions, 1)

	// Test DeepCopyInto
	var target MultiComputeConfig
	original.DeepCopyInto(&target)
	assert.Equal(t, original.Name, target.Name)
	assert.Equal(t, original.Spec.Policies.EnforceRuntimeClass, target.Spec.Policies.EnforceRuntimeClass)
	assert.Equal(t, original.Spec.Policies.LimitGPUsPerPod, target.Spec.Policies.LimitGPUsPerPod)
	assert.Len(t, target.Status.Conditions, 1)

	// Test DeepCopyObject
	obj := original.DeepCopyObject()
	assert.NotNil(t, obj)
	configObj, ok := obj.(*MultiComputeConfig)
	assert.True(t, ok)
	assert.Equal(t, original.Name, configObj.Name)
}

func TestMultiComputeConfigListDeepCopy(t *testing.T) {
	original := &MultiComputeConfigList{
		Items: []MultiComputeConfig{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "config1"},
				Spec:       MultiComputeConfigSpec{},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "config2"},
				Spec:       MultiComputeConfigSpec{},
			},
		},
	}

	// Test DeepCopy
	copied := original.DeepCopy()
	assert.NotNil(t, copied)
	assert.Len(t, copied.Items, 2)
	assert.Equal(t, original.Items[0].Name, copied.Items[0].Name)
	assert.Equal(t, original.Items[1].Name, copied.Items[1].Name)

	// Test DeepCopyInto
	var target MultiComputeConfigList
	original.DeepCopyInto(&target)
	assert.Len(t, target.Items, 2)
	assert.Equal(t, original.Items[0].Name, target.Items[0].Name)
	assert.Equal(t, original.Items[1].Name, target.Items[1].Name)

	// Test DeepCopyObject
	obj := original.DeepCopyObject()
	assert.NotNil(t, obj)
	listObj, ok := obj.(*MultiComputeConfigList)
	assert.True(t, ok)
	assert.Len(t, listObj.Items, 2)
}
