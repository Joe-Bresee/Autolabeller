package matchinglogic

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"

	autolabellerv1alpha1 "github.com/Joe-Bresee/Autolabeller/api/v1alpha1"
)

// MatchesNodeDetailed returns whether the node matches and a list of fields that matched.
// It mirrors the Pod matcher but uses NodeMatchCriteria fields.
func MatchesNodeDetailed(mc *autolabellerv1alpha1.MatchCriteria, node *corev1.Node) (bool, []string) {
	matchedFields := []string{}
	if mc == nil {
		return true, matchedFields
	}

	if cm := mc.CommonMatch; cm != nil {
		if name := cm.Name; name != "" {
			if node.Name != name {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "commonMatch.name")
		}
		for k, v := range cm.Labels {
			if node.Labels[k] != v {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, fmt.Sprintf("commonMatch.labels[%s]", k))
		}
		for k, v := range cm.Annotations {
			if node.Annotations[k] != v {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, fmt.Sprintf("commonMatch.annotations[%s]", k))
		}
	}

	if nm := mc.NodeMatch; nm != nil {
		// Architecture label match (kubernetes.io/arch)
		if len(nm.ArchLabels) > 0 {
			nodeArch := node.Labels["kubernetes.io/arch"]
			found := false
			for _, want := range nm.ArchLabels {
				if nodeArch == want {
					found = true
					matchedFields = append(matchedFields, fmt.Sprintf("nodeMatch.archLabels:%s", want))
					break
				}
			}
			if !found {
				return false, matchedFields
			}
		}

		// OS label match (kubernetes.io/os)
		if len(nm.OSLabels) > 0 {
			nodeOS := node.Labels["kubernetes.io/os"]
			found := false
			for _, want := range nm.OSLabels {
				if nodeOS == want {
					found = true
					matchedFields = append(matchedFields, fmt.Sprintf("nodeMatch.osLabels:%s", want))
					break
				}
			}
			if !found {
				return false, matchedFields
			}
		}

		// Taints match (expects key=value:effect)
		if len(nm.Taints) > 0 {
			existing := map[string]struct{}{}
			for _, t := range node.Spec.Taints {
				existing[fmt.Sprintf("%s=%s:%s", t.Key, t.Value, t.Effect)] = struct{}{}
			}
			for _, want := range nm.Taints {
				if _, ok := existing[want]; !ok {
					return false, matchedFields
				}
				matchedFields = append(matchedFields, fmt.Sprintf("nodeMatch.taints:%s", want))
			}
		}

		// Kernel version (simple contains/equality match)
		if nm.KernelVersion != "" {
			if !strings.Contains(node.Status.NodeInfo.KernelVersion, nm.KernelVersion) {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "nodeMatch.kernelVersion")
		}

		// Container runtime (contains match against runtime version string)
		if nm.ContainerRuntime != "" {
			if !strings.Contains(node.Status.NodeInfo.ContainerRuntimeVersion, nm.ContainerRuntime) {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "nodeMatch.containerRuntime")
		}
	}

	return true, matchedFields
}

// MatchesNode is a convenience wrapper returning only the boolean match result.
func MatchesNode(mc *autolabellerv1alpha1.MatchCriteria, node *corev1.Node) bool {
	ok, _ := MatchesNodeDetailed(mc, node)
	return ok
}
