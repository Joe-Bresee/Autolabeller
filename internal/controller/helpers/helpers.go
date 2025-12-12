package helpers

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	autolabellerv1alpha1 "github.com/Joe-Bresee/Autolabeller/api/v1alpha1"
)

// Helpers
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

func MatchesPod(mc *autolabellerv1alpha1.MatchCriteria, pod *corev1.Pod) bool {
	if mc == nil {
		return true
	}
	if cm := mc.CommonMatch; cm != nil {
		if ns := cm.Namespace; ns != "" && pod.Namespace != ns {
			return false
		}
		if name := cm.Name; name != "" && pod.Name != name {
			return false
		}
		for k, v := range cm.Labels {
			if pod.Labels[k] != v {
				return false
			}
		}
		for k, v := range cm.Annotations {
			if pod.Annotations[k] != v {
				return false
			}
		}
	}
	if pm := mc.PodMatch; pm != nil {
		if pm.HostNetwork != nil && pod.Spec.HostNetwork != *pm.HostNetwork {
			return false
		}
		if pm.ServiceAccount != "" && pod.Spec.ServiceAccountName != pm.ServiceAccount {
			return false
		}
		if len(pm.NodeSelector) > 0 {
			for k, v := range pm.NodeSelector {
				if pod.Spec.NodeSelector[k] != v {
					return false
				}
			}
		}
		if pm.RestartPolicy != "" && string(pod.Spec.RestartPolicy) != pm.RestartPolicy {
			return false
		}
		if len(pm.Images) > 0 {
			images := map[string]struct{}{}
			for _, c := range pod.Spec.Containers {
				images[c.Image] = struct{}{}
			}
			matchedAny := false
			for _, want := range pm.Images {
				if _, ok := images[want]; ok {
					matchedAny = true
					break
				}
			}
			if !matchedAny {
				return false
			}
		}
	}
	return true
}
