package helpers

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"sigs.k8s.io/controller-runtime/pkg/client"

	autolabellerv1alpha1 "github.com/Joe-Bresee/Autolabeller/api/v1alpha1"
)

func SetCondition(rule *autolabellerv1alpha1.ClassificationRule, condType string, status metav1.ConditionStatus, reason, msg string) {
	cond := metav1.Condition{
		Type:               condType,
		Status:             status,
		ObservedGeneration: rule.GetGeneration(),
		Reason:             reason,
		Message:            msg,
		LastTransitionTime: metav1.Now(),
	}
	replaced := false
	for i := range rule.Status.Conditions {
		if rule.Status.Conditions[i].Type == condType {
			rule.Status.Conditions[i] = cond
			replaced = true
			break
		}
	}
	if !replaced {
		rule.Status.Conditions = append(rule.Status.Conditions, cond)
	}
}

// SetConditionWithLog sets a condition on the rule and logs the reason/message.
// Prefer this helper in controllers to ensure user-visible status and operator logs stay in sync.
func SetConditionWithLog(logger logr.Logger, rule *autolabellerv1alpha1.ClassificationRule, condType string, status metav1.ConditionStatus, reason, msg string) {
	SetCondition(rule, condType, status, reason, msg)
	logger.Info("condition updated", "type", condType, "status", string(status), "reason", reason, "message", msg)
}

func ApplyLabelsToObject(obj client.Object, labels map[string]string) bool {
	if len(labels) == 0 {
		return false
	}
	m := obj.GetLabels()
	if m == nil {
		m = map[string]string{}
	}
	changed := false
	for k, v := range labels {
		if m[k] != v {
			m[k] = v
			changed = true
		}
	}
	if changed {
		obj.SetLabels(m)
	}
	return changed
}

func FilterPodList(listOpts *[]client.ListOption, match *autolabellerv1alpha1.MatchCriteria) {
	if match == nil {
		return
	}

	// CommonMatch filters applied at API level (reduces objects fetched)
	// Name patterns are checked later in MatchesPodDetailed (requires in-memory inspection)
	// Annotations are checked later in MatchesPodDetailed (requires in-memory inspection)
	if cm := match.CommonMatch; cm != nil {
		if cm.Namespace != "" {
			*listOpts = append(*listOpts, client.InNamespace(cm.Namespace))
		}
		if len(cm.Labels) > 0 {
			*listOpts = append(*listOpts, client.MatchingLabels(cm.Labels))
		}
	}
}

func FilterNodeList(listOpts *[]client.ListOption, match *autolabellerv1alpha1.MatchCriteria) {
	if match == nil {
		return
	}

	// CommonMatch filters applied at API level (cluster-scoped, so no namespace)
	// Name patterns are checked later in MatchesNodeDetailed (requires in-memory inspection)
	// Annotations are checked later in MatchesNodeDetailed (requires in-memory inspection)
	if cm := match.CommonMatch; cm != nil {
		if len(cm.Labels) > 0 {
			*listOpts = append(*listOpts, client.MatchingLabels(cm.Labels))
		}
	}

	// NodeMatch filters - single or multi-value arch/os label filtering applied at API level
	// Taints, KernelVersion, ContainerRuntime are checked later in MatchesNodeDetailed (requires in-memory inspection)
	if nm := match.NodeMatch; nm != nil {
		// Single arch value → exact label selector
		if len(nm.ArchLabels) == 1 {
			*listOpts = append(*listOpts, client.MatchingLabels(map[string]string{
				"kubernetes.io/arch": nm.ArchLabels[0],
			}))
		}
		// Single OS value → exact label selector
		if len(nm.OSLabels) == 1 {
			*listOpts = append(*listOpts, client.MatchingLabels(map[string]string{
				"kubernetes.io/os": nm.OSLabels[0],
			}))
		}

		// Multi-value arch/os → set-based selectors with OR semantics
		selector := labels.NewSelector()
		added := false
		if len(nm.ArchLabels) > 1 {
			if req, err := labels.NewRequirement("kubernetes.io/arch", selection.In, nm.ArchLabels); err == nil {
				selector = selector.Add(*req)
				added = true
			}
		}
		if len(nm.OSLabels) > 1 {
			if req, err := labels.NewRequirement("kubernetes.io/os", selection.In, nm.OSLabels); err == nil {
				selector = selector.Add(*req)
				added = true
			}
		}
		if added {
			*listOpts = append(*listOpts, client.MatchingLabelsSelector{Selector: selector})
		}
	}
}

func FilterDeploymentList(listOpts *[]client.ListOption, match *autolabellerv1alpha1.MatchCriteria) {
	if match == nil {
		return
	}

	// CommonMatch filters applied at API level (reduces objects fetched, in turn reducing memory load and API calls)
	// Name patterns are checked later in MatchesDeploymentDetailed (requires in-memory inspection)
	// Annotations are checked later in MatchesDeploymentDetailed (requires in-memory inspection)
	if cm := match.CommonMatch; cm != nil {
		if cm.Namespace != "" {
			*listOpts = append(*listOpts, client.InNamespace(cm.Namespace))
		}
		if len(cm.Labels) > 0 {
			*listOpts = append(*listOpts, client.MatchingLabels(cm.Labels))
		}
	}
}
