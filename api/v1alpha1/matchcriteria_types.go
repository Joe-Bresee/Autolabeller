package v1alpha1

// CommonMatchCriteria contains criteria common to most resource types
type CommonMatchCriteria struct {
	// Labels is a map of label keys and values to match
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations is a map of annotation keys and values to match
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Namespace is the namespace name to match. Exact string match.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Name is the resource name to match. Supports wildcard patterns (* and ?).
	// +optional
	Name string `json:"name,omitempty"`
}

// PodMatchCriteria contains Pod-specific match criteria
type PodMatchCriteria struct {
	// Images is a list of container image patterns to match.
	// Each pattern supports wildcard matching (* and ?).
	// +optional
	Images []string `json:"images,omitempty"`

	// CPURequests matches Pods with CPU requests matching the specified value.
	// Supports comparison operators (e.g., ">1", "<=500m").
	// +optional
	CPURequests string `json:"cpuRequests,omitempty"`

	// MemoryRequests matches Pods with memory requests matching the specified value.
	// Supports comparison operators (e.g., ">1Gi", "<=512Mi").
	// +optional
	MemoryRequests string `json:"memoryRequests,omitempty"`

	// CPULimits matches Pods with CPU limits matching the specified value.
	// Supports comparison operators.
	// +optional
	CPULimits string `json:"cpuLimits,omitempty"`

	// MemoryLimits matches Pods with memory limits matching the specified value.
	// Supports comparison operators.
	// +optional
	MemoryLimits string `json:"memoryLimits,omitempty"`

	// NodeSelector is a map of node labels to match for Pod scheduling.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// ServiceAccount is the name of the ServiceAccount to match. Exact match.
	// +optional
	ServiceAccount string `json:"serviceAccount,omitempty"`

	// HostNetwork matches Pods with hostNetwork setting.
	// +optional
	HostNetwork *bool `json:"hostNetwork,omitempty"`

	// RestartPolicy matches Pods with specific restart policy.
	// Valid values: Always, OnFailure, Never
	// +optional
	RestartPolicy string `json:"restartPolicy,omitempty"`
}

// NodeMatchCriteria contains Node-specific match criteria
type NodeMatchCriteria struct {
	// ArchLabels matches nodes with specific architecture labels.
	// Valid values: amd64, arm64, arm, ppc64le, s390x
	// +optional
	ArchLabels []string `json:"archLabels,omitempty"`

	// OSLabels matches nodes with specific OS labels.
	// Valid values: linux, windows
	// +optional
	OSLabels []string `json:"osLabels,omitempty"`

	// Taints matches nodes with specific taints.
	// Each taint is matched as key=value:effect format.
	// +optional
	Taints []string `json:"taints,omitempty"`

	// KernelVersion matches nodes with kernel versions matching the pattern.
	// Supports comparison operators.
	// +optional
	KernelVersion string `json:"kernelVersion,omitempty"`

	// ContainerRuntime matches nodes with specific container runtime.
	// Examples: docker, containerd, cri-o
	// +optional
	ContainerRuntime string `json:"containerRuntime,omitempty"`
}

// DeploymentMatchCriteria contains Deployment-specific match criteria
type DeploymentMatchCriteria struct {
	// Replicas matches Deployments with specific replica count.
	// Supports comparison operators (e.g., ">3", "==5").
	// +optional
	Replicas string `json:"replicas,omitempty"`

	// Strategy matches Deployments with specific update strategy.
	// Valid values: RollingUpdate, Recreate
	// +optional
	Strategy string `json:"strategy,omitempty"`

	// ImagePullPolicy matches Deployments with specific image pull policy.
	// Valid values: Always, Never, IfNotPresent
	// +optional
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}
