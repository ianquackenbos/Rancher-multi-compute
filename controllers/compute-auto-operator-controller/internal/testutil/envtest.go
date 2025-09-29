package testutil

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	multisuseiov1alpha1 "github.com/suse/rancher-multi-compute/api/multi.suse.io/v1alpha1"
)

// TestEnv represents a test environment with envtest
type TestEnv struct {
	cfg       *rest.Config
	k8sClient client.Client
	testEnv   *envtest.Environment
	ctx       context.Context
	cancel    context.CancelFunc
}

// StartEnvTest starts an envtest environment
func StartEnvTest(t *testing.T) *TestEnv {
	ginkgo.By("bootstrapping test environment")
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "..", "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing: false,
		BinaryAssetsDirectory: "/tmp/k8sbin/k8s/1.30.0-darwin-arm64",
	}

	cfg, err := testEnv.Start()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(cfg).NotTo(gomega.BeNil())

	err = multisuseiov1alpha1.AddToScheme(scheme.Scheme)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(k8sClient).NotTo(gomega.BeNil())

	ctx, cancel := context.WithCancel(context.TODO())

	return &TestEnv{
		cfg:       cfg,
		k8sClient: k8sClient,
		testEnv:   testEnv,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Stop stops the test environment
func (te *TestEnv) Stop() {
	te.cancel()
	ginkgo.By("tearing down the test environment")
	err := te.testEnv.Stop()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

// GetConfig returns the REST config
func (te *TestEnv) GetConfig() *rest.Config {
	return te.cfg
}

// GetClient returns the Kubernetes client
func (te *TestEnv) GetClient() client.Client {
	return te.k8sClient
}

// GetContext returns the context
func (te *TestEnv) GetContext() context.Context {
	return te.ctx
}

// GetScheme returns the scheme
func (te *TestEnv) GetScheme() *runtime.Scheme {
	return scheme.Scheme
}

// GetManager returns a controller manager for testing
func (te *TestEnv) GetManager() ctrl.Manager {
	mgr, err := ctrl.NewManager(te.cfg, ctrl.Options{
		Scheme: scheme.Scheme,
		// Disable metrics server to avoid port conflicts in parallel tests
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
		// Disable health probe server to avoid port conflicts
		HealthProbeBindAddress: "0",
	})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	return mgr
}

// SetupTestLogger sets up the test logger
func SetupTestLogger() {
	logf.SetLogger(zap.New(zap.WriteTo(ginkgo.GinkgoWriter), zap.UseDevMode(true)))
}
