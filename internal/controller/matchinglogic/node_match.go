package matchinglogic

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"

	autolabellerv1alpha1 "github.com/Joe-Bresee/Autolabeller/api/v1alpha1"
)

// MatchesNodeDetailed returns whether the node matches and a list of fields that matched.
// Note: CommonMatch.Labels and NodeMatch.ArchLabels/OSLabels are pre-filtered by FilterNodeList,
// so we only check items that require in-memory inspection (name patterns, annotations, taints, kernel, runtime).
func MatchesNodeDetailed(mc *autolabellerv1alpha1.MatchCriteria, node *corev1.Node) (bool, []string) {
	matchedFields := []string{}
	if mc == nil {
		return true, matchedFields
	}

	if cm := mc.CommonMatch; cm != nil {
		// Labels are already pre-filtered by FilterNodeList, skip them here
		if name := cm.Name; name != "" {
			if node.Name != name {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "commonMatch.name")
		}
		for k, v := range cm.Annotations {
			if node.Annotations[k] != v {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, fmt.Sprintf("commonMatch.annotations[%s]", k))
		}
	}

	if nm := mc.NodeMatch; nm != nil {
		// Architecture and OS labels are already pre-filtered by FilterNodeList, skip them here

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
