package helpers

import (
	"fmt"

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

// MatchesPodDetailed returns whether the pod matches and a list of fields that matched.
func MatchesPodDetailed(mc *autolabellerv1alpha1.MatchCriteria, pod *corev1.Pod) (bool, []string) {
	matchedFields := []string{}
	if mc == nil {
		return true, matchedFields
	}

	if cm := mc.CommonMatch; cm != nil {
		if ns := cm.Namespace; ns != "" {
			if pod.Namespace != ns {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "commonMatch.namespace")
		}
		if name := cm.Name; name != "" {
			if pod.Name != name {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "commonMatch.name")
		}
		for k, v := range cm.Labels {
			if pod.Labels[k] != v {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, fmt.Sprintf("commonMatch.labels[%s]", k))
		}
		for k, v := range cm.Annotations {
			if pod.Annotations[k] != v {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, fmt.Sprintf("commonMatch.annotations[%s]", k))
		}
	}

	if pm := mc.PodMatch; pm != nil {
		if pm.HostNetwork != nil {
			if pod.Spec.HostNetwork != *pm.HostNetwork {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "podMatch.hostNetwork")
		}
		if pm.ServiceAccount != "" {
			if pod.Spec.ServiceAccountName != pm.ServiceAccount {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "podMatch.serviceAccount")
		}
		if len(pm.NodeSelector) > 0 {
			for k, v := range pm.NodeSelector {
				if pod.Spec.NodeSelector[k] != v {
					return false, matchedFields
				}
				matchedFields = append(matchedFields, fmt.Sprintf("podMatch.nodeSelector[%s]", k))
			}
		}
		if pm.RestartPolicy != "" {
			if string(pod.Spec.RestartPolicy) != pm.RestartPolicy {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "podMatch.restartPolicy")
		}
		if len(pm.Images) > 0 {
			images := map[string]struct{}{}
			for _, c := range pod.Spec.Containers {
				images[c.Image] = struct{}{}
			}
			matchedAny := ""
			for _, want := range pm.Images {
				if _, ok := images[want]; ok {
					matchedAny = want
					break
				}
			}
			if matchedAny == "" {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, fmt.Sprintf("podMatch.images:%s", matchedAny))
		}
	}

	return true, matchedFields
}

// MatchesPod keeps backward compatibility while providing details.
func MatchesPod(mc *autolabellerv1alpha1.MatchCriteria, pod *corev1.Pod) bool {
	ok, _ := MatchesPodDetailed(mc, pod)
	return ok
}
