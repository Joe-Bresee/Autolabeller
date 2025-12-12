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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClassificationRuleSpec defines the desired state of ClassificationRule
type ClassificationRuleSpec struct {

	// TargetKind specifies the Kubernetes resource type to apply the rule to
	// +kubebuilder:validation:Enum=Pod;Node;Namespace;Service;Deployment;StatefulSet;DaemonSet;ReplicaSet;Job;CronJob
	// +kubebuilder:default=Pod
	TargetKind string `json:"targetKind"`

	// Match defines the resource fields to match for labelling
	// Key = field name (e.g., "image", "name", "namespace")
	// Value = expected value to match (e.g., "nginx", "prod_proxy", "production")
	// + optional
	Match map[string]string `json:"match,omitempty"`

	// Labels defines the labels to apply when a resource matches the rule
	// Key = label name
	// Value = label value
	// + optional
	Labels map[string]string `json:"labels,omitempty"`

	// ConflictPolicy defines the policy to apply when there is a conflict in labeling
	// +kubebuilder:validation:Enum=Overwrite;Merge;Ignore;Error
	// +kubebuilder:default=Merge
	ConflictPolicy string `json:"conflictPolicy,omitempty"`
}

// ClassificationRuleStatus defines the observed state of ClassificationRule.
type ClassificationRuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the ClassificationRule resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClassificationRule is the Schema for the classificationrules API
type ClassificationRule struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of ClassificationRule
	// +required
	Spec ClassificationRuleSpec `json:"spec"`

	// status defines the observed state of ClassificationRule
	// +optional
	Status ClassificationRuleStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ClassificationRuleList contains a list of ClassificationRule
type ClassificationRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []ClassificationRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClassificationRule{}, &ClassificationRuleList{})
}
