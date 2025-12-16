package matchinglogic

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"

	autolabellerv1alpha1 "github.com/Joe-Bresee/Autolabeller/api/v1alpha1"
)

// MatchesDeploymentDetailed returns whether the Deployment matches and a list of fields that matched.
// Note: Namespace and CommonMatch.Labels are pre-filtered by FilterDeploymentList, so we only check items
// that require in-memory inspection (name, annotations, and Deployment-specific criteria).
func MatchesDeploymentDetailed(mc *autolabellerv1alpha1.MatchCriteria, deployment *appsv1.Deployment) (bool, []string) {
	matchedFields := []string{}
	if mc == nil {
		return true, matchedFields
	}

	if cm := mc.CommonMatch; cm != nil {
		// Namespace and Labels are already pre-filtered by FilterDeploymentList, skip them here
		if name := cm.Name; name != "" {
			if deployment.Name != name {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "commonMatch.name")
		}
		for k, v := range cm.Annotations {
			if deployment.Annotations[k] != v {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, fmt.Sprintf("commonMatch.annotations[%s]", k))
		}
	}

	if pm := mc.DeploymentMatch; pm != nil {
		// Replicas: exact match on desired replicas when provided
		if pm.Replicas != "" {
			desired := int32(1)
			if deployment.Spec.Replicas != nil {
				desired = *deployment.Spec.Replicas
			}
			if fmt.Sprintf("%d", desired) != pm.Replicas {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "deploymentMatch.replicas")
		}

		// Strategy: RollingUpdate or Recreate
		if pm.Strategy != "" {
			if string(deployment.Spec.Strategy.Type) != pm.Strategy {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "deploymentMatch.strategy")
		}

		// ImagePullPolicy: check any container in the pod template
		if pm.ImagePullPolicy != "" {
			matched := false
			for _, c := range deployment.Spec.Template.Spec.Containers {
				if string(c.ImagePullPolicy) == pm.ImagePullPolicy {
					matched = true
					break
				}
			}
			if !matched {
				return false, matchedFields
			}
			matchedFields = append(matchedFields, "deploymentMatch.imagePullPolicy")
		}
	}

	return true, matchedFields
}
