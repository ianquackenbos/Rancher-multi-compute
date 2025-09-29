package controllers

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

var _ = Describe("BaseController", func() {
	var controller *BaseController
	var ctx context.Context

	BeforeEach(func() {
		controller = NewBaseController("test-controller")
		ctx = context.Background()
	})

	Describe("NewBaseController", func() {
		It("should create a controller with the correct name", func() {
			Expect(controller.Name).To(Equal("test-controller"))
		})
	})

	Describe("Start", func() {
		It("should start without error", func() {
			err := controller.Start(ctx)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Stop", func() {
		It("should stop without error", func() {
			err := controller.Stop()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
