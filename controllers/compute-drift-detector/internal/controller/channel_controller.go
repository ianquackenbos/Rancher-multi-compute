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

// ChannelReconciler reconciles a Channel object for drift detection
type ChannelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=multi.suse.io,resources=channels,verbs=get;list;watch
//+kubebuilder:rbac:groups=multi.suse.io,resources=channels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fleet.cattle.io,resources=bundledeployments,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop
func (r *ChannelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Channel instance
	channel := &multisuseiov1alpha1.Channel{}
	if err := r.Get(ctx, req.NamespacedName, channel); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Detect drift between declared Channel and actual deployments
	driftDetected, driftDetails := r.detectDrift(ctx, channel)

	if driftDetected {
		logger.Info("Drift detected", "channel", channel.Name, "details", driftDetails)

		// Update Channel status with drift information
		channel.Status.Phase = "DriftDetected"

		// Add drift condition
		condition := metav1.Condition{
			Type:               "DriftDetected",
			Status:             metav1.ConditionTrue,
			Reason:             "ConfigurationDrift",
			Message:            driftDetails,
			LastTransitionTime: metav1.Now(),
		}
		channel.Status.Conditions = append(channel.Status.Conditions, condition)

		if err := r.Status().Update(ctx, channel); err != nil {
			logger.Error(err, "failed to update Channel status with drift")
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("No drift detected", "channel", channel.Name)
	}

	return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
}

// detectDrift checks for configuration drift
func (r *ChannelReconciler) detectDrift(ctx context.Context, channel *multisuseiov1alpha1.Channel) (bool, string) {
	// This is a simplified drift detection implementation
	// In a real implementation, you would:
	// 1. Compare Channel spec with actual BundleDeployment status
	// 2. Check for version mismatches
	// 3. Verify resource configurations
	// 4. Check for policy violations

	// For now, return no drift
	return false, ""
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChannelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&multisuseiov1alpha1.Channel{}).
		Complete(r)
}
