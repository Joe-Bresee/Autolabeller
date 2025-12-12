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
	Match *MatchCriteria `json:"match,omitempty"`

	// Labels defines the labels to apply when a resource matches the rule
	// Key = label name
	// Value = label value
	// + optional
	Labels map[string]string `json:"labels,omitempty"`

	// ConflictPolicy defines the policy to apply when there is a conflict in labeling
	// +kubebuilder:validation:Enum=Overwrite;Merge;Ignore;Error
	// +kubebuilder:default=Merge
	ConflictPolicy string `json:"conflictPolicy,omitempty"`

	// Suspend temporarily disables the application of this classification rule
	// +optional
	// +kubebuilder:default=false
	Suspend bool `json:"suspend,omitempty"`
}

// MatchCriteria defines the criteria that maps resource fields to expected values
type MatchCriteria struct {

	// General information
	Images       []string          `json:"images,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	Name         string            `json:"name,omitempty"`
	Namespace    string            `json:"namespace,omitempty"`
	Status       string            `json:"status,omitempty"`
	ReplicaCount int32             `json:"replicaCount,omitempty"`

	// Deployment related information
	Strategy       string `json:"strategy,omitempty"`
	UpdateStrategy string `json:"updateStrategy,omitempty"`

	// RBAC related information
	Roles          []string `json:"roles,omitempty"`
	ClusterRoles   []string `json:"clusterRoles,omitempty"`
	ServiceAccount string   `json:"serviceAccount,omitempty"`

	// Networking related information
	Ports       []int32 `json:"ports,omitempty"`
	HostNetwork *bool   `json:"hostNetwork,omitempty"`
	IPFamily    string  `json:"ipFamily,omitempty"`

	// Storage related information
	Volumes       []string `json:"volumes,omitempty"`
	StorageClass  string   `json:"storageClass,omitempty"`
	AccessModes   []string `json:"accessModes,omitempty"`
	VolumeBinding string   `json:"volumeBinding,omitempty"`

	// Scheduling related information
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	Tolerations  []string          `json:"tolerations,omitempty"`
	Affinity     string            `json:"affinity,omitempty"`
	Topology     string            `json:"topology,omitempty"`
	Taints       []string          `json:"taints,omitempty"`

	// Resource requirements
	CPURequests    string `json:"cpuRequests,omitempty"`
	CPULimits      string `json:"cpuLimits,omitempty"`
	MemoryRequests string `json:"memoryRequests,omitempty"`
	MemoryLimits   string `json:"memoryLimits,omitempty"`
}

// ClassificationRuleStatus defines the observed state of ClassificationRule.
type ClassificationRuleStatus struct {
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	//+ patchStrategy=merge
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// lastAppliedTime indicates the last time the rule was successfully applied.
	// +optional
	LastReconciled *metav1.Time `json:"lastReconciled,omitempty"`

	// matchedResourcesCount indicates the number of resources that matched this rule.
	// +optional
	MatchedResourcesCount int32 `json:"MatchedResourceCount,omitempty"`

	// lastError provides details of the last error encountered while applying the rule.
	// +optional
	LastError string `json:"lastError,omitempty"`

	// observedGeneration is the most recent generation observed for this ClassificationRule.
	// It corresponds to the ClassificationRule's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
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
