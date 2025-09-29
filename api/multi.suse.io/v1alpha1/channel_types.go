package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ChannelSpec defines the desired state of Channel
type ChannelSpec struct {
	// Vendor specifies the GPU vendor (nvidia, amd, intel)
	// +kubebuilder:validation:Enum=nvidia;amd;intel
	Vendor string `json:"vendor"`

	// Channel specifies the release channel (stable, lts, canary)
	// +kubebuilder:validation:Enum=stable;lts;canary
	Channel string `json:"channel"`

	// ClusterSelector defines which clusters this channel applies to
	ClusterSelector metav1.LabelSelector `json:"clusterSelector"`
}

// ChannelStatus defines the observed state of Channel
type ChannelStatus struct {
	// ObservedVersion is the version currently deployed
	ObservedVersion string `json:"observedVersion,omitempty"`

	// Phase represents the current phase of the rollout
	// +kubebuilder:validation:Enum=Pending;RollingOut;Paused;Completed;Failed
	Phase string `json:"phase,omitempty"`

	// Conditions represent the latest available observations of the channel's state
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// Channel is the Schema for the channels API
type Channel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ChannelSpec   `json:"spec,omitempty"`
	Status ChannelStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ChannelList contains a list of Channel
type ChannelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Channel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Channel{}, &ChannelList{})
}
