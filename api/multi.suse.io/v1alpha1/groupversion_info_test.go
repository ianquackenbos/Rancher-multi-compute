package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestGroupVersion(t *testing.T) {
	expected := schema.GroupVersion{Group: "multi.suse.io", Version: "v1alpha1"}
	assert.Equal(t, expected, GroupVersion)
}

func TestSchemeBuilder(t *testing.T) {
	assert.NotNil(t, SchemeBuilder)
	assert.Equal(t, GroupVersion, SchemeBuilder.GroupVersion)
}

func TestAddToScheme(t *testing.T) {
	assert.NotNil(t, AddToScheme)
}
