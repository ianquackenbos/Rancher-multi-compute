package utils

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = Describe("Utils", func() {
	Describe("GetVersion", func() {
		It("should return the correct version", func() {
			version := GetVersion()
			Expect(version).To(Equal("v0.1.0"))
		})
	})

	Describe("ValidateConfig", func() {
		It("should return error for nil config", func() {
			err := ValidateConfig(nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("config cannot be nil"))
		})

		It("should pass for valid config", func() {
			config := map[string]interface{}{
				"key": "value",
			}
			err := ValidateConfig(config)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
