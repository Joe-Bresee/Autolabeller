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

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	autolabellerv1alpha1 "github.com/Joe-Bresee/Autolabeller/api/v1alpha1"
	. "github.com/Joe-Bresee/Autolabeller/internal/controller/helpers"
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
		SetCondition(&rule, "Suspended", metav1.ConditionTrue, "RuleSuspended", "Rule is suspended")
		_ = r.Status().Update(ctx, &rule)
		return ctrl.Result{}, nil
	}

	matched := int32(0)
	switch rule.Spec.TargetKind {
	case "Pod":
		var pods corev1.PodList
		listOpts := []client.ListOption{}
		if rule.Spec.Match != nil && rule.Spec.Match.CommonMatch != nil {
			ns := rule.Spec.Match.CommonMatch.Namespace
			if ns != "" {
				listOpts = append(listOpts, client.InNamespace(ns))
			}
		}
		if err := r.List(ctx, &pods, listOpts...); err != nil {
			SetCondition(&rule, "Ready", metav1.ConditionFalse, "ListFailed", fmt.Sprintf("Failed to list pods: %v", err))
			_ = r.Status().Update(ctx, &rule)
			return ctrl.Result{}, err
		}
		for i := range pods.Items {
			pod := &pods.Items[i]
			if ok, fields := MatchesPodDetailed(rule.Spec.Match, pod); ok {
				log.Info("pod matched criteria", "pod", client.ObjectKeyFromObject(pod), "matchedFields", fields)
				if ApplyLabelsToObject(pod, rule.Spec.Labels) {
					if err := r.Update(ctx, pod); err != nil {
						log.Error(err, "failed to update pod labels", "pod", client.ObjectKeyFromObject(pod))
						continue
					}
					matched++
				}
			}
		}
	default:
		SetCondition(&rule, "Ready", metav1.ConditionFalse, "UnsupportedTarget", fmt.Sprintf("TargetKind %s not yet implemented", rule.Spec.TargetKind))
		_ = r.Status().Update(ctx, &rule)
		return ctrl.Result{}, nil
	}

	rule.Status.MatchedResourcesCount = matched
	rule.Status.ObservedGeneration = rule.GetGeneration()
	SetCondition(&rule, "Ready", metav1.ConditionTrue, "Applied", fmt.Sprintf("Applied labels to %d resources", matched))
	if err := r.Status().Update(ctx, &rule); err != nil {
		return ctrl.Result{}, err
	}

	// Compute reqeue Interval
	requeueAfter := 30 * time.Second // default
	if rule.Spec.RefreshInterval != "" {
		d, err := time.ParseDuration(rule.Spec.RefreshInterval)
		if err != nil {
			SetCondition(&rule, "Degraded", metav1.ConditionTrue, "InvalidRefreshInterval", "RefreshInterval must be a valid duration (e.g. 30s, 5m)")
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
