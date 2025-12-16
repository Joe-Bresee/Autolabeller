package matchinglogic

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	autolabellerv1alpha1 "github.com/Joe-Bresee/Autolabeller/api/v1alpha1"
)

// MatchesPodDetailed returns whether the pod matches and a list of fields that matched.
// Note: Namespace and CommonMatch.Labels are pre-filtered by FilterPodList, so we only check items
// that require in-memory inspection (name patterns, annotations, pod-specific criteria).
func MatchesPodDetailed(mc *autolabellerv1alpha1.MatchCriteria, pod *corev1.Pod) (bool, []string) {
	matchedFields := []string{}
	if mc == nil {
		return true, matchedFields
	}

	if cm := mc.CommonMatch; cm != nil {
		// Namespace and Labels are already pre-filtered by FilterPodList, skip them here
		if name := cm.Name; name != "" {
			if pod.Name != name {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "commonMatch.name")
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
