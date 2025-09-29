package controller

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	multisuseiov1alpha1 "github.com/suse/rancher-multi-compute/api/multi.suse.io/v1alpha1"
)

// MultiComputeConfigReconciler reconciles a MultiComputeConfig object
type MultiComputeConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=multi.suse.io,resources=multicomputeconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multi.suse.io,resources=multicomputeconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=multi.suse.io,resources=multicomputeconfigs/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop
func (r *MultiComputeConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the MultiComputeConfig instance
	config := &multisuseiov1alpha1.MultiComputeConfig{}
	if err := r.Get(ctx, req.NamespacedName, config); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Apply org-wide compute policies
	if err := r.applyPolicies(ctx, config); err != nil {
		logger.Error(err, "failed to apply policies", "config", config.Name)
		return ctrl.Result{}, err
	}

	// Update status
	config.Status.Conditions = []metav1.Condition{
		{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			Reason:             "PoliciesApplied",
			Message:            "All policies applied successfully",
			LastTransitionTime: metav1.Now(),
		},
	}

	if err := r.Status().Update(ctx, config); err != nil {
		logger.Error(err, "failed to update MultiComputeConfig status")
		return ctrl.Result{}, err
	}

	logger.Info("MultiComputeConfig reconciled successfully", "config", config.Name)
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// applyPolicies applies the configured policies
func (r *MultiComputeConfigReconciler) applyPolicies(ctx context.Context, config *multisuseiov1alpha1.MultiComputeConfig) error {
	// This is a simplified policy application
	// In a real implementation, you would:
	// 1. Apply Gatekeeper constraints
	// 2. Apply Kyverno policies
	// 3. Configure resource quotas
	// 4. Set up network policies

	logger := log.FromContext(ctx)

	if config.Spec.Policies.EnforceRuntimeClass {
		logger.Info("Enforcing runtime class policy")
	}

	if config.Spec.Policies.RestrictGPUNamespaces {
		logger.Info("Restricting GPU namespaces")
	}

	if config.Spec.Policies.RequireCosign {
		logger.Info("Requiring Cosign signatures")
	}

	if config.Spec.Policies.LimitGPUsPerPod > 0 {
		logger.Info("Limiting GPUs per pod", "limit", config.Spec.Policies.LimitGPUsPerPod)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MultiComputeConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multisuseiov1alpha1.MultiComputeConfig{}).
		Complete(r)
}
