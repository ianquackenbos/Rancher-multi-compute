package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestChannelSpec(t *testing.T) {
	spec := ChannelSpec{
		Vendor:  "nvidia",
		Channel: "stable",
		ClusterSelector: metav1.LabelSelector{
			MatchLabels: map[string]string{
				"kubernetes.io/os": "linux",
			},
		},
	}

	assert.Equal(t, "nvidia", spec.Vendor)
	assert.Equal(t, "stable", spec.Channel)
	assert.Equal(t, "linux", spec.ClusterSelector.MatchLabels["kubernetes.io/os"])
}

func TestChannelStatus(t *testing.T) {
	status := ChannelStatus{
		Phase:           "Progressing",
		ObservedVersion: "v24.9.0/12.4.1",
		Conditions: []metav1.Condition{
			{
				Type:   "Ready",
				Status: metav1.ConditionUnknown,
				Reason: "Progressing",
			},
		},
	}

	assert.Equal(t, "Progressing", status.Phase)
	assert.Equal(t, "v24.9.0/12.4.1", status.ObservedVersion)
	assert.Len(t, status.Conditions, 1)
	assert.Equal(t, "Ready", status.Conditions[0].Type)
}

func TestChannel(t *testing.T) {
	channel := &Channel{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nvidia-stable",
		},
		Spec: ChannelSpec{
			Vendor:  "nvidia",
			Channel: "stable",
		},
		Status: ChannelStatus{
			Phase: "Progressing",
		},
	}

	assert.Equal(t, "nvidia-stable", channel.Name)
	assert.Equal(t, "nvidia", channel.Spec.Vendor)
	assert.Equal(t, "stable", channel.Spec.Channel)
	assert.Equal(t, "Progressing", channel.Status.Phase)
}

func TestChannelList(t *testing.T) {
	channel1 := Channel{
		ObjectMeta: metav1.ObjectMeta{Name: "nvidia-stable"},
		Spec:       ChannelSpec{Vendor: "nvidia", Channel: "stable"},
	}
	channel2 := Channel{
		ObjectMeta: metav1.ObjectMeta{Name: "amd-stable"},
		Spec:       ChannelSpec{Vendor: "amd", Channel: "stable"},
	}

	list := &ChannelList{
		Items: []Channel{channel1, channel2},
	}

	assert.Len(t, list.Items, 2)
	assert.Equal(t, "nvidia-stable", list.Items[0].Name)
	assert.Equal(t, "amd-stable", list.Items[1].Name)
}
