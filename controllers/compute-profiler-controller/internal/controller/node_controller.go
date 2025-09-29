package controller

import (
	"context"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NodeReconciler reconciles a Node object
type NodeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;update;patch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile is part of the main kubernetes reconciliation loop
func (r *NodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the Node instance
	node := &corev1.Node{}
	if err := r.Get(ctx, req.NamespacedName, node); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Detect GPU capabilities using Node Feature Discovery labels
	gpuVendor := r.detectGPUVendor(node)
	if gpuVendor != "" {
		// Add normalized labels
		labels := node.Labels
		if labels == nil {
			labels = make(map[string]string)
		}

		labels["compute.multi.suse.io/vendor"] = gpuVendor
		labels["compute.multi.suse.io/gpu-available"] = "true"

		// Add vendor-specific labels
		switch gpuVendor {
		case "nvidia":
			labels["compute.multi.suse.io/nvidia-gpu"] = "true"
			if migProfile := r.detectMIGProfile(node); migProfile != "" {
				labels["compute.multi.suse.io/mig-profile"] = migProfile
			}
		case "amd":
			labels["compute.multi.suse.io/amd-gpu"] = "true"
		case "intel":
			labels["compute.multi.suse.io/intel-gpu"] = "true"
		}

		node.Labels = labels
		if err := r.Update(ctx, node); err != nil {
			logger.Error(err, "failed to update Node labels", "node", node.Name)
			return ctrl.Result{}, err
		}

		logger.Info("Updated Node with GPU labels", "node", node.Name, "vendor", gpuVendor)
	}

	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// detectGPUVendor detects GPU vendor from NFD labels
func (r *NodeReconciler) detectGPUVendor(node *corev1.Node) string {
	labels := node.Labels

	// Check for NVIDIA GPU
	if labels["feature.node.kubernetes.io/pci-10de.present"] == "true" {
		return "nvidia"
	}

	// Check for AMD GPU
	if labels["feature.node.kubernetes.io/pci-1002.present"] == "true" {
		return "amd"
	}

	// Check for Intel GPU
	if labels["feature.node.kubernetes.io/pci-8086.present"] == "true" {
		return "intel"
	}

	return ""
}

// detectMIGProfile detects NVIDIA MIG profile
func (r *NodeReconciler) detectMIGProfile(node *corev1.Node) string {
	labels := node.Labels

	// Check for MIG profile labels
	for key, value := range labels {
		if strings.HasPrefix(key, "nvidia.com/mig-") && value == "true" {
			return strings.TrimPrefix(key, "nvidia.com/mig-")
		}
	}

	return ""
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Node{}).
		Complete(r)
}
