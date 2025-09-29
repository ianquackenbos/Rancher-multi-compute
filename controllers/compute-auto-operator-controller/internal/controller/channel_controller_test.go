package controller

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	multisuseiov1alpha1 "github.com/suse/rancher-multi-compute/api/multi.suse.io/v1alpha1"
	"github.com/suse/rancher-multi-compute/controllers/compute-auto-operator-controller/internal/testutil"
	"github.com/suse/rancher-multi-compute/internal/vendors"
	"github.com/suse/rancher-multi-compute/internal/versions"
)

// GVKs for Fleet (test)
var (
	testBundleGVK = schema.GroupVersionKind{
		Group:   "fleet.cattle.io",
		Version: "v1alpha1",
		Kind:    "Bundle",
	}
	testBdGVK = schema.GroupVersionKind{
		Group:   "fleet.cattle.io",
		Version: "v1alpha1",
		Kind:    "BundleDeployment",
	}
)

func TestChannelController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Channel Controller Suite")
}

// fakeVersionResolver for testing
type fakeVersionResolver struct {
	pins map[string]versions.VendorPins
}

func (f *fakeVersionResolver) Resolve(ctx context.Context, channel string) (versions.VendorPins, error) {
	if pins, ok := f.pins[channel]; ok {
		return pins, nil
	}
	return versions.VendorPins{}, nil
}

var _ = Describe("Channel Controller (EnvTest)", func() {
	var (
		testEnv      *testutil.TestEnv
		reconciler   *ChannelReconciler
		mockResolver *fakeVersionResolver
		mgr          ctrl.Manager
		ctx          context.Context
		cancel       context.CancelFunc
	)

	BeforeEach(func() {
		testutil.SetupTestLogger()
		testEnv = testutil.StartEnvTest(&testing.T{})
		ctx, cancel = context.WithCancel(context.Background())

		mockResolver = &fakeVersionResolver{
			pins: map[string]versions.VendorPins{
				"stable": {
					NVIDIA: versions.Pins{OperatorTag: "v24.9.0", RuntimeTag: "12.4.1"},
					AMD:    versions.Pins{OperatorTag: "v24.9.0", RuntimeTag: "5.7.1"},
					Intel:  versions.Pins{OperatorTag: "v24.9.0", RuntimeTag: "24.16.0"},
				},
			},
		}

		mgr = testEnv.GetManager()
		reconciler = &ChannelReconciler{
			Client:          testEnv.GetClient(),
			Scheme:          testEnv.GetScheme(),
			VersionResolver: mockResolver,
			VendorSources: map[vendors.Vendor]vendors.Source{
				vendors.VendorNVIDIA: {
					Repo:      "https://nvidia.github.io/helm-charts",
					Chart:     "gpu-operator",
					Namespace: "gpu-operator",
				},
			},
		}
		// Use a unique controller name to avoid conflicts
		Expect(reconciler.SetupWithManager(mgr)).To(Succeed())

		go func() {
			defer GinkgoRecover()
			err := mgr.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
		}()
	})

	AfterEach(func() {
		cancel()
		testEnv.Stop()
	})

	Context("Channel reconciliation", func() {
		It("should create a Fleet Bundle and update Channel status to RollingOut initially", func() {
			channel := &multisuseiov1alpha1.Channel{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nvidia-stable",
				},
				Spec: multisuseiov1alpha1.ChannelSpec{
					Vendor:  "nvidia",
					Channel: "stable",
					ClusterSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{"kubernetes.io/os": "linux"},
					},
				},
			}
			Expect(testEnv.GetClient().Create(ctx, channel)).To(Succeed())

			// Wait for reconciliation
			Eventually(func() bool {
				bundle := &unstructured.Unstructured{}
				bundle.SetGroupVersionKind(testBundleGVK)
				err := testEnv.GetClient().Get(ctx,
					types.NamespacedName{Name: "rmc-nvidia-stack", Namespace: "cattle-fleet-system"}, bundle)
				return err == nil
			}, 10*time.Second, 1*time.Second).Should(BeTrue())

			// Verify Bundle was created with correct labels
			bundle := &unstructured.Unstructured{}
			bundle.SetGroupVersionKind(testBundleGVK)
			Expect(testEnv.GetClient().Get(ctx,
				types.NamespacedName{Name: "rmc-nvidia-stack", Namespace: "cattle-fleet-system"}, bundle)).To(Succeed())

			Expect(bundle.GetLabels()).To(HaveKeyWithValue("multi.suse.io/vendor", "nvidia"))
			Expect(bundle.GetLabels()).To(HaveKeyWithValue("multi.suse.io/channel", "stable"))
			Expect(bundle.GetLabels()).To(HaveKeyWithValue("multi.suse.io/owner", "nvidia-stable"))

			// Verify Channel status
			Eventually(func() string {
				fetchedChannel := &multisuseiov1alpha1.Channel{}
				_ = testEnv.GetClient().Get(ctx, types.NamespacedName{Name: "nvidia-stable"}, fetchedChannel)
				return fetchedChannel.Status.Phase
			}, 10*time.Second, 1*time.Second).Should(Equal("RollingOut"))
		})

		It("should update Channel status to Completed when BundleDeployment is Ready", func() {
			channel := &multisuseiov1alpha1.Channel{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nvidia-ready",
				},
				Spec: multisuseiov1alpha1.ChannelSpec{
					Vendor:  "nvidia",
					Channel: "stable",
					ClusterSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{"kubernetes.io/os": "linux"},
					},
				},
			}
			Expect(testEnv.GetClient().Create(ctx, channel)).To(Succeed())

			// Wait for Bundle creation
			Eventually(func() bool {
				bundle := &unstructured.Unstructured{}
				bundle.SetGroupVersionKind(testBundleGVK)
				err := testEnv.GetClient().Get(ctx,
					types.NamespacedName{Name: "rmc-nvidia-stack", Namespace: "cattle-fleet-system"}, bundle)
				return err == nil
			}, 10*time.Second, 1*time.Second).Should(BeTrue())

			// Create a BundleDeployment with Ready status
			bundleDeployment := &unstructured.Unstructured{}
			bundleDeployment.SetGroupVersionKind(testBdGVK)
			bundleDeployment.SetName("rmc-nvidia-stack-bd-test")
			bundleDeployment.SetNamespace("cattle-fleet-system")
			bundleDeployment.SetLabels(map[string]string{
				"multi.suse.io/owner":  "nvidia-ready",
				"multi.suse.io/vendor": "nvidia",
			})
			Expect(unstructured.SetNestedField(bundleDeployment.Object, true, "status", "ready")).To(Succeed())
			Expect(unstructured.SetNestedField(bundleDeployment.Object, "Ready", "status", "display", "state")).To(Succeed())
			Expect(testEnv.GetClient().Create(ctx, bundleDeployment)).To(Succeed())

			// Verify Channel status
			Eventually(func() string {
				fetchedChannel := &multisuseiov1alpha1.Channel{}
				_ = testEnv.GetClient().Get(ctx, types.NamespacedName{Name: "nvidia-ready"}, fetchedChannel)
				return fetchedChannel.Status.Phase
			}, 10*time.Second, 1*time.Second).Should(Equal("Completed"))
		})

		It("should update Channel status to Failed when BundleDeployment fails", func() {
			channel := &multisuseiov1alpha1.Channel{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nvidia-failed",
				},
				Spec: multisuseiov1alpha1.ChannelSpec{
					Vendor:  "nvidia",
					Channel: "stable",
					ClusterSelector: metav1.LabelSelector{
						MatchLabels: map[string]string{"kubernetes.io/os": "linux"},
					},
				},
			}
			Expect(testEnv.GetClient().Create(ctx, channel)).To(Succeed())

			// Wait for Bundle creation
			Eventually(func() bool {
				bundle := &unstructured.Unstructured{}
				bundle.SetGroupVersionKind(testBundleGVK)
				err := testEnv.GetClient().Get(ctx,
					types.NamespacedName{Name: "rmc-nvidia-stack", Namespace: "cattle-fleet-system"}, bundle)
				return err == nil
			}, 10*time.Second, 1*time.Second).Should(BeTrue())

			// Create a BundleDeployment with Failed status
			bundleDeployment := &unstructured.Unstructured{}
			bundleDeployment.SetGroupVersionKind(testBdGVK)
			bundleDeployment.SetName("rmc-nvidia-stack-bd-fail")
			bundleDeployment.SetNamespace("cattle-fleet-system")
			bundleDeployment.SetLabels(map[string]string{
				"multi.suse.io/owner":  "nvidia-failed",
				"multi.suse.io/vendor": "nvidia",
			})
			Expect(unstructured.SetNestedField(bundleDeployment.Object, false, "status", "ready")).To(Succeed())
			Expect(unstructured.SetNestedField(bundleDeployment.Object, "ErrApplied", "status", "display", "state")).To(Succeed())
			Expect(testEnv.GetClient().Create(ctx, bundleDeployment)).To(Succeed())

			// Verify Channel status
			Eventually(func() string {
				fetchedChannel := &multisuseiov1alpha1.Channel{}
				_ = testEnv.GetClient().Get(ctx, types.NamespacedName{Name: "nvidia-failed"}, fetchedChannel)
				return fetchedChannel.Status.Phase
			}, 10*time.Second, 1*time.Second).Should(Equal("Failed"))
		})

		It("should handle invalid vendor specification", func() {
			channel := &multisuseiov1alpha1.Channel{
				ObjectMeta: metav1.ObjectMeta{
					Name: "invalid-vendor-channel",
				},
				Spec: multisuseiov1alpha1.ChannelSpec{
					Vendor:  "unsupported", // Invalid vendor
					Channel: "stable",
				},
			}
			Expect(testEnv.GetClient().Create(ctx, channel)).To(Succeed())

			// Verify Channel status
			Eventually(func() string {
				fetchedChannel := &multisuseiov1alpha1.Channel{}
				_ = testEnv.GetClient().Get(ctx, types.NamespacedName{Name: "invalid-vendor-channel"}, fetchedChannel)
				return fetchedChannel.Status.Phase
			}, 10*time.Second, 1*time.Second).Should(Equal("Failed"))

			Eventually(func() bool {
				fetchedChannel := &multisuseiov1alpha1.Channel{}
				_ = testEnv.GetClient().Get(ctx, types.NamespacedName{Name: "invalid-vendor-channel"}, fetchedChannel)
				for _, cond := range fetchedChannel.Status.Conditions {
					if cond.Type == "Ready" && cond.Status == metav1.ConditionFalse && cond.Reason == "InvalidVendor" {
						return true
					}
				}
				return false
			}, 10*time.Second, 1*time.Second).Should(BeTrue())
		})
	})
})
