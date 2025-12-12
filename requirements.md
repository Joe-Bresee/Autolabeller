# Autolabeller Requirements

## Project Overview
Autolabeller is a Kubernetes operator that provides dynamic, rule-driven classification and labeling of Kubernetes resources. It watches Pods, Nodes, and other workload objects, evaluates them against user-defined policies via `ClassificationRule` CRDs, and applies labels that represent workloads' characteristics, behaviors, or metadata.

---

## Functional Requirements

### FR1: CRD-Based Classification Rules
- **Description**: Support `ClassificationRule` custom resources that define labeling policies
- **Specification**: 
  - Define target resource kinds (Pods, Nodes, Deployments, etc.)
  - Specify match criteria (images, resource requests/limits, annotations, namespaces, labels)
  - Define labels to apply when criteria match
  - Configure conflict-resolution strategies (overwrite, merge, ignore)
- **Status**: Scaffolded, needs implementation

### FR2: Pod Classification
- **Description**: Automatically label Pods based on active `ClassificationRule` resources
- **Criteria Support**:
  - Container images
  - Resource requests and limits (CPU, memory)
  - Pod annotations
  - Namespace labels and properties
  - Pod labels
- **Status**: Not implemented

### FR3: Node Classification
- **Description**: Automatically label Nodes based on active `ClassificationRule` resources
- **Criteria Support**:
  - Node labels
  - Node annotations
  - Capacity and resource information
  - Taints
  - Kernel version and OS information
- **Status**: Not implemented

### FR4: Intelligent Label Patching
- **Description**: Apply labels in a non-destructive manner
- **Requirements**:
  - Only patch labels owned by the Autolabeller operator
  - Preserve user-defined and third-party labels
  - Support merging labels from multiple classification rules
  - Track label ownership via annotations
- **Status**: Not implemented

### FR5: Conflict Handling
- **Description**: Handle scenarios where multiple rules create conflicting labels
- **Strategies**:
  - **Overwrite**: Last rule wins
  - **Merge**: Combine labels from all matching rules
  - **Ignore**: Skip conflicting labels
  - **Report**: Log conflicts in resource status
- **Status**: Not implemented

### FR6: Resource Status Tracking
- **Description**: Maintain detailed status information for each `ClassificationRule`
- **Status Fields**:
  - `matchedResourceCount`: Number of resources matching this rule
  - `lastReconciled`: Timestamp of last reconciliation
  - `conflicts`: List of detected label conflicts
  - `lastError`: Error message from last failed reconciliation
  - `conditions`: Standard Kubernetes conditions (Available, Progressing, Degraded)
- **Status**: Partially scaffolded (conditions field exists)

### FR7: Observability & Metrics
- **Description**: Expose metrics for monitoring operator health and performance
- **Metrics**:
  - Number of rules processed
  - Number of resources labeled
  - Number of conflicts detected
  - Reconciliation duration
  - Rule processing errors
- **Status**: Framework exists, needs implementation

### FR8: Multi-Resource Support
- **Description**: Support labeling of multiple Kubernetes resource types
- **Initial Support**: Pods, Nodes
- **Future Support**: Deployments, Jobs, StatefulSets, PVCs, Namespaces
- **Status**: Architecture ready, Pod/Node support not implemented

---

## Non-Functional Requirements

### NFR1: Reliability
- Operator must gracefully handle malformed rules
- All operations must be idempotent
- Reconciliation loops must not get stuck or infinite loop

### NFR2: Performance
- Controller should process events with minimal latency
- Label patching should batch operations where possible
- Should efficiently handle thousands of resources

### NFR3: Security
- RBAC roles must follow least-privilege principle
- Pod Security Policy: Restricted
- Network policies should restrict ingress/egress
- Secure webhook and metrics endpoints with TLS

### NFR4: High Availability
- Support leader election for multiple operator replicas
- Graceful degradation if operator becomes unavailable
- Persistent state should not be required (stateless design)

### NFR5: Maintainability
- Code must be well-tested (unit and e2e tests)
- Must follow Go best practices and kubebuilder patterns
- Clear error messages and logging
- Comprehensive documentation

### NFR6: Deployment
- Container image size should be minimal (distroless base)
- Helm chart or Kustomize manifests for deployment
- Health checks (liveness and readiness probes)
- Metrics endpoint for Prometheus monitoring

---

## API Specification

### ClassificationRule Resource

```yaml
apiVersion: autolabeller.autolabeller.github.com/v1alpha1
kind: ClassificationRule
metadata:
  name: rule-name
  namespace: autolabeller-system
spec:
  targetKind: Pod | Node | Deployment | Job | StatefulSet | PVC | Namespace
  
  # Match criteria (all conditions must be true - AND logic)
  match:
    # Pod/Container specific
    images:
      - "*/compute-*"
      - "*/ml-*"
    
    # Pod resource requests/limits
    resourceRequests:
      cpu: "> 1"
      memory: "> 1Gi"
    
    # Annotations to match
    annotations:
      key1: "value1"
      key2: "regex:pattern.*"
    
    # Labels to match
    labels:
      environment: production
      team: platform
    
    # Namespace selector
    namespaceSelector:
      matchLabels:
        tier: core
    
    # Node-specific criteria
    nodeSelector:
      kubernetes.io/os: linux
    
    # Taint matches (for nodes)
    taints:
      - effect: NoSchedule
  
  # Labels to apply when match criteria are met
  labels:
    workloadProfile: "cpu-heavy"
    schedulingHint: "isolation-preferred"
    cost-center: "platform"
  
  # Conflict resolution strategy: Overwrite | Merge | Ignore | Report
  conflictPolicy: Merge
  
  # Suspend processing of this rule
  suspend: false

status:
  # Number of resources matched by this rule
  matchedResourceCount: 42
  
  # Timestamp of last reconciliation
  lastReconciled: "2025-12-11T10:30:00Z"
  
  # List of detected label conflicts
  conflicts:
    - resourceName: pod-1
      resourceNamespace: default
      labels:
        - key: tier
          values: ["frontend", "backend"]
  
  # Error from last failed reconciliation
  lastError: ""
  
  # Standard Kubernetes conditions
  conditions:
    - type: Available
      status: "True"
      reason: RulesProcessed
      message: All classification rules processed successfully
```

---

## Testing Requirements

### Unit Tests
- Test individual reconciliation logic for each resource type
- Test match criteria evaluation
- Test label patching logic
- Test conflict resolution strategies
- Test error handling and recovery
- Target: >80% code coverage

### Integration Tests
- Test controller startup and initialization
- Test RBAC enforcement
- Test webhook validation (if implemented)
- Test leader election

### E2E Tests
- Deploy operator to Kind cluster
- Create various ClassificationRule resources
- Verify Pods and Nodes are labeled correctly
- Verify conflict handling
- Verify metric collection
- Verify graceful operator upgrade

---

## Deployment Requirements

### Container Image
- Base: `gcr.io/distroless/static:nonroot`
- Binary: `/manager` (static Go binary)
- User: 65532:65532 (nonroot)

### RBAC
- Service Account: `autolabeller-controller-manager`
- ClusterRole: `manager-role` with permissions to:
  - Read/list/watch ClassificationRule resources
  - Read/list/watch Pods
  - Read/list/watch Nodes
  - Patch Pods and Nodes
  - Create/read ServiceMonitor (for Prometheus)
  - Manage finalizers

### Network
- Metrics endpoint: Default 0 (disabled) or :8443 (HTTPS) or :8080 (HTTP)
- Health probe: :8081
- Webhook: HTTPS with certificate management

### Monitoring
- Prometheus metrics at `/metrics`
- Metrics secured with mTLS or RBAC proxy
- ServiceMonitor resource for Prometheus Operator integration
- Health checks exposed at `/healthz` and `/readyz`

---

## Configuration

### Command-line Flags
- `--metrics-bind-address`: Metrics endpoint address (default: "0")
- `--health-probe-bind-address`: Health probe endpoint (default: ":8081")
- `--leader-elect`: Enable leader election (default: false)
- `--metrics-secure`: Secure metrics with HTTPS (default: true)
- `--webhook-cert-path`: Path to webhook certificates
- `--metrics-cert-path`: Path to metrics server certificates
- `--enable-http2`: Enable HTTP/2 (default: false, disabled for security)

### Environment Variables
- `KUBECONFIG`: Path to kubeconfig file
- Standard Kubernetes client-go auth plugins for OIDC, Azure, GCP, etc.

---

## Implementation Milestones

1. **Core API Implementation**: Complete ClassificationRule types with validation
2. **Pod Classifier**: Implement Pod labeling logic with match criteria
3. **Node Classifier**: Implement Node labeling logic
4. **Conflict Resolution**: Implement all conflict handling strategies
5. **Status Tracking**: Populate status fields during reconciliation
6. **Testing**: Unit tests, integration tests, e2e tests
7. **Documentation**: API docs, deployment guide, examples
8. **Production Hardening**: Security review, performance optimization, HA testing