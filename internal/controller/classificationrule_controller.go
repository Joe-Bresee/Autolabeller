/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	autolabellerv1alpha1 "github.com/Joe-Bresee/Autolabeller/api/v1alpha1"
	"github.com/Joe-Bresee/Autolabeller/internal/controller/helpers"
	"github.com/Joe-Bresee/Autolabeller/internal/controller/matchinglogic"
)

// ClassificationRuleReconciler reconciles a ClassificationRule object
type ClassificationRuleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=autolabeller.autolabeller.github.com,resources=classificationrules,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autolabeller.autolabeller.github.com,resources=classificationrules/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=autolabeller.autolabeller.github.com,resources=classificationrules/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;update;patch
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *ClassificationRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// Fetch rule
	var rule autolabellerv1alpha1.ClassificationRule
	if err := r.Get(ctx, req.NamespacedName, &rule); err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Guard suspend
	if rule.Spec.Suspend {
		helpers.SetConditionWithLog(log, &rule, "Suspended", metav1.ConditionTrue, "RuleSuspended", "Rule is suspended")
		_ = r.Status().Update(ctx, &rule)
		return ctrl.Result{}, nil
	}

	matched := int32(0)
	switch rule.Spec.TargetKind {
	case "Pod":
		var pods corev1.PodList
		listOpts := []client.ListOption{}
		if rule.Spec.Match != nil && rule.Spec.Match.CommonMatch != nil {
			// filter listing options
			helpers.FilterPodList(&listOpts, rule.Spec.Match)
		}
		if err := r.List(ctx, &pods, listOpts...); err != nil {
			helpers.SetConditionWithLog(log, &rule, "Ready", metav1.ConditionFalse, "ListFailed", fmt.Sprintf("Failed to list pods: %v", err))
			_ = r.Status().Update(ctx, &rule)
			return ctrl.Result{}, err
		}
		for i := range pods.Items {
			pod := &pods.Items[i]
			if ok, fields := matchinglogic.MatchesPodDetailed(rule.Spec.Match, pod); ok {
				if helpers.ApplyLabelsToObject(pod, rule.Spec.Labels) {
					log.Info("pod matched criteria, applying labels", "pod", client.ObjectKeyFromObject(pod), "matchedFields", fields)
					if err := r.Update(ctx, pod); err != nil {
						helpers.SetConditionWithLog(log, &rule, "Ready", metav1.ConditionFalse, "UpdateFailed", fmt.Sprintf("Failed to update pod labels: %v", err))
						_ = r.Status().Update(ctx, &rule)
						continue
					}
					matched++
				}
			}
		}
	case "Node":
		var nodes corev1.NodeList
		listOpts := []client.ListOption{}
		if rule.Spec.Match != nil && rule.Spec.Match.CommonMatch != nil && rule.Spec.Match.CommonMatch.Namespace != "" {
			helpers.SetConditionWithLog(log, &rule, "Degraded", metav1.ConditionTrue, "NamespaceIgnoredForNode", "commonMatch.namespace is ignored for Node targetKind")
		}
		if rule.Spec.Match != nil && rule.Spec.Match.NodeMatch != nil {
			helpers.FilterNodeList(&listOpts, rule.Spec.Match)
		}
		if err := r.List(ctx, &nodes, listOpts...); err != nil {
			helpers.SetConditionWithLog(log, &rule, "Ready", metav1.ConditionFalse, "ListFailed", fmt.Sprintf("Failed to list nodes: %v", err))
			_ = r.Status().Update(ctx, &rule)
			return ctrl.Result{}, err
		}
		for i := range nodes.Items {
			node := &nodes.Items[i]
			if ok, fields := matchinglogic.MatchesNodeDetailed(rule.Spec.Match, node); ok {
				if helpers.ApplyLabelsToObject(node, rule.Spec.Labels) {
					log.Info("node matched criteria, applying labels", "node", client.ObjectKeyFromObject(node), "matchedFields", fields)
					if err := r.Update(ctx, node); err != nil {
						helpers.SetConditionWithLog(log, &rule, "Ready", metav1.ConditionFalse, "UpdateFailed", fmt.Sprintf("Failed to update node labels: %v", err))
						_ = r.Status().Update(ctx, &rule)
						continue
					}
					matched++
				}
			}
		}
	case "Deployment":
		var deployments appsv1.DeploymentList
		listOpts := []client.ListOption{}
		if rule.Spec.Match != nil && rule.Spec.Match.CommonMatch != nil {
			helpers.FilterPodList(&listOpts, rule.Spec.Match)
		}
		if err := r.List(ctx, &deployments, listOpts...); err != nil {
			helpers.SetConditionWithLog(log, &rule, "Ready", metav1.ConditionFalse, "ListFailed", fmt.Sprintf("Failed to list deployments: %v", err))
			_ = r.Status().Update(ctx, &rule)
			return ctrl.Result{}, err
		}
		for i := range deployments.Items {
			deployment := &deployments.Items[i]
			if ok, fields := matchinglogic.MatchesDeploymentDetailed(rule.Spec.Match, deployment); ok {
				if helpers.ApplyLabelsToObject(deployment, rule.Spec.Labels) {
					log.Info("deployment matched criteria, applying labels", "deployment", client.ObjectKeyFromObject(deployment), "matchedFields", fields)
					if err := r.Update(ctx, deployment); err != nil {
						helpers.SetConditionWithLog(log, &rule, "Ready", metav1.ConditionFalse, "UpdateFailed", fmt.Sprintf("Failed to update deployment labels: %v", err))
						_ = r.Status().Update(ctx, &rule)
						continue
					}
					matched++
				}
			}
		}
	default:
		helpers.SetConditionWithLog(log, &rule, "Ready", metav1.ConditionFalse, "UnsupportedTarget", fmt.Sprintf("TargetKind %s not yet implemented", rule.Spec.TargetKind))
		_ = r.Status().Update(ctx, &rule)
		return ctrl.Result{}, nil
	}

	rule.Status.MatchedResourcesCount = matched
	rule.Status.ObservedGeneration = rule.GetGeneration()
	helpers.SetConditionWithLog(log, &rule, "Ready", metav1.ConditionTrue, "Applied", fmt.Sprintf("Applied labels to %d resources", matched))
	if err := r.Status().Update(ctx, &rule); err != nil {
		return ctrl.Result{}, err
	}

	// Compute reqeue Interval
	requeueAfter := 30 * time.Second // default
	if rule.Spec.RefreshInterval != "" {
		d, err := time.ParseDuration(rule.Spec.RefreshInterval)
		if err != nil {
			helpers.SetConditionWithLog(log, &rule, "Degraded", metav1.ConditionTrue, "InvalidRefreshInterval", "RefreshInterval must be a valid duration (e.g. 30s, 5m)")
			_ = r.Status().Update(ctx, &rule)
			return ctrl.Result{}, nil
		}
		requeueAfter = d
	}

	log.Info("Reconcile completed", "requeueAfter", requeueAfter.String())
	return ctrl.Result{RequeueAfter: requeueAfter}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClassificationRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&autolabellerv1alpha1.ClassificationRule{}).
		Named("classificationrule").
		Complete(r)
}
