# Auto-Labeler Operator

The Auto-Labeler Operator provides dynamic, rule-driven classification for Kubernetes resources. It watches Pods, Nodes, and other workload objects, evaluates them against user-defined policies, and applies labels that represent workloads’ characteristics, behaviors, or metadata. These labels can later be consumed by scheduling components or policy systems.

The operator enables platform teams to enforce consistent labeling, create intelligent scheduling hints, and build a foundation for higher-level automation.

---

## Goals

- Provide a CRD-driven labeling engine for Pods, Nodes, and other Kubernetes resources.
- Support both static metadata classification and dynamic behavior-driven classification.
- Maintain idempotent, conflict-aware label patching.
- Enable downstream systems (e.g., schedulers, cost allocators, governance engines) to use labels as first-class signals.
- Serve as a general-purpose platform component, not tied to a specific domain such as ML or GitOps.

---

## Key Features

### 1. CRD-Based Classification Rules
Users define `ClassificationRule` objects describing:
- the resource kinds to target  
- match criteria (images, resources, annotations, namespaces, behavior)  
- labels to apply  
- conflict-resolution strategy  

### 2. Policy-Driven Label Application
The operator applies labels based on:
- resource metadata
- resource specifications
- container images
- namespace properties
- resource requests/limits
- optional runtime metrics (future extension)

### 3. Multi-Resource Support
Initial support:
- Pods  
- Nodes  

Future extensions:
- Deployments  
- Jobs  
- StatefulSets  
- PVCs  
- Namespaces  

### 4. Intelligent Diffing & Patch Logic
- Only patches labels that belong to the Auto-Labeler.
- Avoids overwriting user-defined labels.
- Supports merging labels from multiple policies.

### 5. Conflict Handling
Policies can be configured to:
- overwrite  
- merge  
- ignore conflicts  
- report conflicts via status  

### 6. Observability & Status
Each `ClassificationRule` maintains:
- matchedResourceCount  
- lastReconciled  
- conflicts  
- lastError  

---

## CRD Example

```yaml
apiVersion: labeling.example.com/v1alpha1
kind: ClassificationRule
metadata:
  name: cpu-heavy-detector
spec:
  targetKind: Pod
  match:
    resourceRequests:
      cpu: "> 1"
    images:
      - "*/compute-*"
  labels:
    workloadProfile: "cpu-heavy"
    schedulingHint: "isolation-preferred"
  conflictPolicy: Merge
```

---

## Architecture

### Controllers
- **RuleController**: Watches ClassificationRule objects.
- **PodClassifier**: Applies rules to Pods.
- **NodeClassifier**: Applies rules to Nodes.

### Components
- Policy evaluator  
- Resource indexer  
- Label patch manager  
- Optional metrics ingestion (future)

---

## Roadmap (Initial)

### Phase 1
- Basic CRD definition  
- Pod classification  
- Label patch logic  
- Rule status reporting  

### Phase 2
- Node classification  
- Conflict resolution  
- Multi-policy merging  

### Phase 3
- Resource behavior classification (metrics-based)  
- Namespace-level defaults  
- Webhook for creation-time labeling  

### Phase 4
- Integration points for external schedulers  
- Advanced heuristics  
- Cost/capacity-aware classification  

---

## Project Structure

```
cmd/
  controller-manager/
pkg/
  apis/
    labeling.example.com/
      v1alpha1/
  controllers/
    rule_controller.go
    pod_classifier.go
    node_classifier.go
  internal/
    evaluator/
    patcher/
    matcher/
config/
  crd/
  manager/
  rbac/
```

---

## Future Directions

- Integration with custom scheduler
- Traffic-aware or load-aware classification
- Label-based compliance enforcement
- Multi-cluster rule propagation

---

## License

Apache 2.0 (recommended)
