package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	multisuseiov1alpha1 "github.com/suse/rancher-multi-compute/api/multi.suse.io/v1alpha1"
	"github.com/suse/rancher-multi-compute/internal/fleetutil"
	"github.com/suse/rancher-multi-compute/internal/vendors"
	"github.com/suse/rancher-multi-compute/internal/versions"
)

// GVKs for Fleet
var (
	bundleGVK = schema.GroupVersionKind{
		Group:   "fleet.cattle.io",
		Version: "v1alpha1",
		Kind:    "Bundle",
	}
	bdGVK = schema.GroupVersionKind{
		Group:   "fleet.cattle.io",
		Version: "v1alpha1",
		Kind:    "BundleDeployment",
	}
)

const (
	finalizerName        = "channel.multi.suse.io/finalizer"
	fleetSystemNamespace = "cattle-fleet-system"
	bundleNamePrefix     = "rmc-"
	ownerLabelKey        = "multi.suse.io/owner"
	vendorLabelKey       = "multi.suse.io/vendor"
	channelLabelKey      = "multi.suse.io/channel"
	partOfLabelKey       = "app.kubernetes.io/part-of"
	partOfLabelValue     = "rancher-multi-compute"
)

// ChannelReconciler reconciles a Channel object
type ChannelReconciler struct {
	client.Client
	Scheme          *runtime.Scheme
	VersionResolver versions.Resolver
	VendorSources   map[vendors.Vendor]vendors.Source
}

//+kubebuilder:rbac:groups=multi.suse.io,resources=channels,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=multi.suse.io,resources=channels/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=multi.suse.io,resources=channels/finalizers,verbs=update
//+kubebuilder:rbac:groups=fleet.cattle.io,resources=bundles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fleet.cattle.io,resources=bundledeployments,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

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

	// Handle deletion
	if !channel.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, channel)
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(channel, finalizerName) {
		controllerutil.AddFinalizer(channel, finalizerName)
		if err := r.Update(ctx, channel); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Resolve versions
	vendorPins, err := r.VersionResolver.Resolve(ctx, channel.Spec.Channel)
	if err != nil {
		logger.Error(err, "Failed to resolve versions for channel", "channel", channel.Spec.Channel)
		return r.updateChannelStatus(ctx, channel, "Failed", "VersionResolutionError", err.Error())
	}

	var currentVendorPins versions.Pins
	var vendorName vendors.Vendor
	switch strings.ToLower(channel.Spec.Vendor) {
	case string(vendors.VendorNVIDIA):
		currentVendorPins = vendorPins.NVIDIA
		vendorName = vendors.VendorNVIDIA
	case string(vendors.VendorAMD):
		currentVendorPins = vendorPins.AMD
		vendorName = vendors.VendorAMD
	case string(vendors.VendorIntel):
		currentVendorPins = vendorPins.Intel
		vendorName = vendors.VendorIntel
	default:
		err := fmt.Errorf("unsupported vendor: %s", channel.Spec.Vendor)
		logger.Error(err, "Invalid vendor specified in Channel")
		return r.updateChannelStatus(ctx, channel, "Failed", "InvalidVendor", err.Error())
	}

	desiredVersion := fmt.Sprintf("%s/%s", currentVendorPins.OperatorTag, currentVendorPins.RuntimeTag)

	// Get vendor source
	vendorSource, exists := r.VendorSources[vendorName]
	if !exists {
		err := fmt.Errorf("no source configuration found for vendor: %s", vendorName)
		logger.Error(err, "Missing vendor source configuration")
		return r.updateChannelStatus(ctx, channel, "Failed", "MissingVendorSource", err.Error())
	}

	// Create Fleet targets
	targets := fleetutil.ConvertLabelSelectorToTargets(channel.Spec.ClusterSelector, &fleetutil.BundleDeploymentOptions{
		DefaultNamespace: vendorSource.Namespace,
		Helm: &fleetutil.HelmOptions{
			ReleaseName: fmt.Sprintf("%s-%s", strings.ToLower(channel.Spec.Vendor), channel.Spec.Channel),
			Repo:        vendorSource.Repo,
			Chart:       vendorSource.Chart,
			Values:      r.buildHelmValues(currentVendorPins),
		},
	})

	// Create or update Fleet Bundle
	err = r.upsertBundle(ctx, channel, strings.ToLower(channel.Spec.Vendor), targets)
	if err != nil {
		logger.Error(err, "Failed to create or update Bundle")
		return r.updateChannelStatus(ctx, channel, "Failed", "BundleCreationError", fmt.Sprintf("Failed to create/update Bundle: %v", err))
	}

	// Compute Channel status based on BundleDeployments
	newPhase := r.computeChannelPhase(ctx, channel)

	// Update Channel status
	if channel.Status.Phase != newPhase || channel.Status.ObservedVersion != desiredVersion {
		return r.updateChannelStatus(ctx, channel, newPhase, "Reconciled", fmt.Sprintf("Channel phase changed to %s, observed version %s", newPhase, desiredVersion))
	}

	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func (r *ChannelReconciler) handleDeletion(ctx context.Context, channel *multisuseiov1alpha1.Channel) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Handling Channel deletion", "channel", channel.Name)

	// List and delete owned Bundles using unstructured
	bundleList := &unstructured.UnstructuredList{}
	bundleList.SetGroupVersionKind(bundleGVK)
	listOpts := []client.ListOption{
		client.InNamespace(fleetSystemNamespace),
		client.MatchingLabels{ownerLabelKey: channel.Name},
	}
	if err := r.List(ctx, bundleList, listOpts...); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list owned Bundles for deletion: %w", err)
	}

	for _, bundle := range bundleList.Items {
		logger.Info("Deleting owned Bundle", "bundle", bundle.GetName())
		if err := r.Delete(ctx, &bundle); err != nil && !errors.IsNotFound(err) {
			return ctrl.Result{}, fmt.Errorf("failed to delete Bundle %s: %w", bundle.GetName(), err)
		}
	}

	// Remove finalizer if all owned resources are gone
	if controllerutil.ContainsFinalizer(channel, finalizerName) {
		controllerutil.RemoveFinalizer(channel, finalizerName)
		if err := r.Update(ctx, channel); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to remove finalizer from Channel %s: %w", channel.Name, err)
		}
	}

	return ctrl.Result{}, nil
}

func (r *ChannelReconciler) updateChannelStatus(ctx context.Context, channel *multisuseiov1alpha1.Channel, phase, reason, message string) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	channel.Status.Phase = phase

	// Update conditions
	now := metav1.Now()
	newCondition := metav1.Condition{
		Type:               "Ready",
		Status:             metav1.ConditionTrue,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: now,
	}
	switch phase {
	case "Failed":
		newCondition.Status = metav1.ConditionFalse
	case "Progressing", "Pending", "RollingOut":
		newCondition.Status = metav1.ConditionUnknown
	}

	// Find existing condition
	existingCondition := getChannelCondition(channel.Status.Conditions, newCondition.Type)
	if existingCondition == nil {
		channel.Status.Conditions = append(channel.Status.Conditions, newCondition)
	} else if existingCondition.Status != newCondition.Status ||
		existingCondition.Reason != newCondition.Reason ||
		existingCondition.Message != newCondition.Message {
		*existingCondition = newCondition
	}

	if err := r.Status().Update(ctx, channel); err != nil {
		logger.Error(err, "Failed to update Channel status", "channel", channel.Name)
		return ctrl.Result{}, fmt.Errorf("failed to update Channel status: %w", err)
	}
	return ctrl.Result{}, nil
}

func (r *ChannelReconciler) buildHelmValues(pins versions.Pins) map[string]interface{} {
	return map[string]interface{}{
		"image": map[string]string{
			"operatorTag": pins.OperatorTag,
			"runtimeTag":  pins.RuntimeTag,
		},
	}
}

func (r *ChannelReconciler) computeChannelPhase(ctx context.Context, channel *multisuseiov1alpha1.Channel) string {
	phase, err := r.summarizePhase(ctx, channel, strings.ToLower(channel.Spec.Vendor))
	if err != nil {
		return "Failed"
	}
	return phase
}

// upsertBundle creates or updates a Fleet Bundle using unstructured objects
func (r *ChannelReconciler) upsertBundle(ctx context.Context, ch *multisuseiov1alpha1.Channel, vendor string, targets []fleetutil.Target) error {
	name := fmt.Sprintf("rmc-%s-stack", vendor)

	b := &unstructured.Unstructured{}
	b.SetGroupVersionKind(bundleGVK)
	b.SetNamespace(fleetSystemNamespace)
	b.SetName(name)
	b.SetLabels(map[string]string{
		partOfLabelKey:  partOfLabelValue,
		vendorLabelKey:  vendor,
		channelLabelKey: ch.Spec.Channel,
		ownerLabelKey:   ch.Name,
	})
	// OwnerReference to Channel (cluster-scoped → leave Namespace empty)
	b.SetOwnerReferences([]metav1.OwnerReference{{
		APIVersion: "multi.suse.io/v1alpha1",
		Kind:       "Channel",
		Name:       ch.Name,
		UID:        ch.UID,
	}})

	// Spec: .spec.helm, .spec.namespace, .spec.targets
	spec := map[string]any{
		"targets": targets,
	}
	if err := unstructured.SetNestedField(b.Object, spec, "spec"); err != nil {
		return err
	}

	// Create or Patch
	current := &unstructured.Unstructured{}
	current.SetGroupVersionKind(bundleGVK)
	key := client.ObjectKey{Namespace: b.GetNamespace(), Name: b.GetName()}
	if err := r.Get(ctx, key, current); err != nil {
		// Not found → create
		return r.Create(ctx, b)
	}
	// Patch spec/labels if changed
	current.Object["spec"] = b.Object["spec"]
	current.SetLabels(b.GetLabels())
	return r.Update(ctx, current)
}

// summarizePhase determines the overall phase from BundleDeployments
func (r *ChannelReconciler) summarizePhase(ctx context.Context, ch *multisuseiov1alpha1.Channel, vendor string) (string, error) {
	bds := &unstructured.UnstructuredList{}
	bds.SetGroupVersionKind(bdGVK)
	// List by labels
	ls := client.MatchingLabels{
		vendorLabelKey: vendor,
		ownerLabelKey:  ch.Name,
	}
	if err := r.List(ctx, bds, ls); err != nil {
		return "Pending", err
	}
	if len(bds.Items) == 0 {
		return "RollingOut", nil
	}
	anyFailed, anyReady := false, false
	for i := range bds.Items {
		// Fleet BD readiness usually shows under status fields; treat lack of ready as not ready.
		ready, _, _ := unstructured.NestedBool(bds.Items[i].Object, "status", "ready")
		state, _, _ := unstructured.NestedString(bds.Items[i].Object, "status", "display", "state")
		if ready {
			anyReady = true
		}
		if state == "ErrApplied" || state == "Modified" || state == "NotReady" {
			anyFailed = true
		}
	}
	switch {
	case anyFailed:
		return "Failed", nil
	case anyReady:
		return "Completed", nil
	default:
		return "RollingOut", nil
	}
}

func getChannelCondition(conditions []metav1.Condition, condType string) *metav1.Condition {
	for i := range conditions {
		if conditions[i].Type == condType {
			return &conditions[i]
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ChannelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("channel-auto-operator").
		For(&multisuseiov1alpha1.Channel{}).
		Complete(r)
}
