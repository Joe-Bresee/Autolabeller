# Autolabeller TODO List

## Priority 1: Core API Definition ✅ COMPLETE
- [X] **T1.1**: Define ClassificationRuleSpec fields (targetKind, match, labels, conflictPolicy, suspend, refreshInterval)
- [X] **T1.2**: Define match criteria types split by resource (CommonMatchCriteria, PodMatchCriteria, NodeMatchCriteria, DeploymentMatchCriteria) in separate file
- [ ] **T1.3**: Add advanced validation markers and CEL rules for ClassificationRule
- [X] **T1.4**: Define ClassificationRuleStatus fields (matchedResourceCount, conditions, lastError, observedGeneration)
- [X] **T1.5**: Create/refactor sample ClassificationRule manifests
- [X] **T1.6**: Generate CRD manifests successfully
- [X] **T1.7**: Create sample test Pod (funnypod.yaml) for manual testing

## Priority 2: Match Criteria Engine
- [X] **T2.1**: Implement basic Reconcile loop (fetch rule, guard suspend, list resources, match, label, update status)
- [X] **T2.2**: Implement MatchesPodDetailed for comprehensive Pod matching (common & pod-specific criteria)
- [X] **T2.3**: Implement ApplyLabelsToObject for idempotent label patching
- [X] **T2.4**: Implement SetCondition for status condition management
- [X] **T2.5**: Add Pod RBAC markers (get, list, watch, update, patch)
- [X] **T2.6**: Implement RefreshInterval parsing and RequeueAfter logic
- [X] **T2.7**: Add detailed logging for label application (conditional on changes only)
- [ ] **T2.8**: Write unit tests for matcher functions (>85% coverage target)
- [ ] **T2.9**: Write integration tests for full Pod reconciliation scenarios
- [ ] **T2.10**: Test with multiple classification rules and conflict handling

## Priority 3: Pod Classification Logic
- [ ] **T3.1**: Implement Node classification (MatchesNodeDetailed)
- [ ] **T3.2**: Add Node case to reconciler switch statement
- [X] **T3.3**: Add Node RBAC markers in controller (get, list, watch, update, patch)
- [ ] **T3.4**: Implement node-specific label patching
- [ ] **T3.5**: Test Node labeling end-to-end
- [ ] **T3.6**: Handle node taint and label matching edge cases
- [ ] **T3.7**: Write unit and integration tests for Node classification

## Priority 4: Node Classification Logic
- [ ] **T4.1**: Create NodeClassifier struct with Node evaluation logic
- [ ] **T4.2**: Implement Node listing and filtering logic
- [ ] **T4.3**: Implement Node matching against ClassificationRule criteria
- [ ] **T4.4**: Create NodeLabelPatcher for non-destructive label application
- [ ] **T4.5**: Implement label ownership tracking for Nodes
- [ ] **T4.6**: Implement conflict detection for Node labels
- [ ] **T4.7**: Handle conflict resolution strategies for Nodes
- [ ] **T4.8**: Implement Node reconciliation in ClassificationRuleReconciler
- [ ] **T4.9**: Add Node event handlers to trigger reconciliation
- [ ] **T4.10**: Write unit tests for NodeClassifier
- [ ] **T4.11**: Write integration tests for Node labeling

## Priority 5: Status & Observability ✅ PARTIALLY COMPLETE
- [X] **T5.1**: Implement status update logic in reconciler
- [X] **T5.2**: Populate matchedResourceCount during reconciliation
- [X] **T5.3**: Implement condition tracking (Ready, Suspended, Degraded)
- [X] **T5.4**: Update observedGeneration in status
- [ ] **T5.5**: Add lastReconciled timestamp population
- [ ] **T5.6**: Implement error tracking in status.lastError
- [ ] **T5.7**: Create Prometheus metrics for rules processed
- [ ] **T5.8**: Create Prometheus metrics for resources labeled
- [ ] **T5.9**: Create Prometheus metrics for conflicts detected
- [ ] **T5.10**: Create Prometheus metrics for reconciliation duration
- [ ] **T5.11**: Create Prometheus metrics for errors
- [ ] **T5.12**: Implement proper logging throughout reconciliation

## Priority 6: RBAC & Security
- [X] **T6.1**: Add Pod watching/patching RBAC markers
- [X] **T6.2**: Add Node watching/patching RBAC markers
- [ ] **T6.3**: Verify ClusterRole is properly generated from kubebuilder markers
- [ ] **T6.4**: Add permissions for ServiceMonitor (Prometheus)
- [ ] **T6.5**: Test RBAC enforcement in actual cluster
- [ ] **T6.6**: Review network policies
- [ ] **T6.7**: Security audit

## Priority 7: Webhook Validation (Optional for v1)
- [ ] **T7.1**: Create ValidatingWebhookConfiguration for ClassificationRule
- [ ] **T7.2**: Implement webhook validation logic
- [ ] **T7.3**: Add CEL validation rules for match criteria
- [ ] **T7.4**: Test webhook validation
- [ ] **T7.5**: Configure mutual TLS for webhooks

## Priority 8: Testing
- [ ] **T8.1**: Complete unit tests for all matchers
- [ ] **T8.2**: Complete unit tests for PodClassifier
- [ ] **T8.3**: Complete unit tests for NodeClassifier
- [ ] **T8.4**: Complete unit tests for reconciler logic
- [ ] **T8.5**: Add integration test for full reconciliation flow
- [ ] **T8.6**: Implement e2e tests with Kind cluster
- [ ] **T8.7**: Test conflict resolution scenarios in e2e
- [ ] **T8.8**: Test operator upgrade scenarios
- [ ] **T8.9**: Test multiple operators (leader election)
- [ ] **T8.10**: Achieve >80% code coverage
- [ ] **T8.11**: Fix any failing test suite

## Priority 9: Deployment & Documentation
- [ ] **T9.1**: Create Helm chart (or update Kustomize manifests)
- [ ] **T9.2**: Add deployment instructions to README
- [ ] **T9.3**: Create example ClassificationRule manifests
- [ ] **T9.4**: Document API specification
- [ ] **T9.5**: Create troubleshooting guide
- [ ] **T9.6**: Add Prometheus ServiceMonitor example
- [ ] **T9.7**: Document all command-line flags
- [ ] **T9.8**: Create operator architecture documentation

## Priority 10: Performance & Production Hardening
- [ ] **T10.1**: Benchmark label patching performance
- [ ] **T10.2**: Optimize list/watch operations
- [ ] **T10.3**: Implement proper reconciliation backoff
- [ ] **T10.4**: Test with large numbers of rules and resources
- [ ] **T10.5**: Implement resource limits in manifests
- [ ] **T10.6**: Verify metrics don't cause memory bloat
- [ ] **T10.7**: Performance load testing
- [ ] **T10.8**: Security audit and penetration testing
- [ ] **T10.9**: Ensure all error paths are tested
- [ ] **T10.10**: Verify graceful shutdown behavior

## Priority 11: Future Enhancements
- [ ] **T11.1**: Support additional resource types (Deployments, Jobs, StatefulSets, PVCs)
- [ ] **T11.2**: Implement webhooks for validation and mutation
- [ ] **T11.3**: Add support for runtime metrics (future extension)
- [ ] **T11.4**: Implement rule templating/parameterization
- [ ] **T11.5**: Add UI/dashboard for rule management
- [ ] **T11.6**: Support multi-cluster scenarios
- [ ] **T11.7**: Add policy conflict analysis tools
- [ ] **T11.8**: Implement cost allocation based on labels

---

## Legend
- [ ] Task not started
- [x] Task completed

## Notes
- All code changes should include appropriate error handling and logging
- Every feature should have corresponding unit and integration tests
- Security should be reviewed at each priority level
- Documentation should be updated alongside code changes

---

## Current Status
**Last Updated**: 2025-12-15

### ✅ Completed (Core v1alpha1 MVP)
- CRD types with split-by-resource-type MatchCriteria in separate file
- Basic reconciler with Pod matching and idempotent labeling
- Condition-based status tracking (Ready, Suspended, Degraded)
- RefreshInterval with configurable requeue intervals
- Conditional logging (only logs when labels actually change)
- RBAC markers for Pod/Node operations
- Sample CR and test Pod manifest
- Manual testing shows successful label application in Kind cluster

### 🔄 In Progress
- Testing and validation of current implementation
- Manual e2e testing in local Kind cluster

### ⏭️ Next Steps (High Priority)
1. Write unit tests for matcher functions (MatchesPodDetailed, ApplyLabelsToObject)
2. Integration tests for reconciliation scenarios
3. Implement Node classification (MatchesNodeDetailed)
4. Add lastReconciled timestamp population
5. Performance testing with larger deployments

**MILESTONE 1: ✅ CORE POD MATCHING (COMPLETE)**
- Scaffold project and CRD ✅
- Implement Pod labeling controller ✅
- Apply labels based on simple metadata rules ✅
- Status reporting ✅
- Test in local cluster ✅


MILESTONE 1:
Scaffold project and CRD
Implement Pod labeling controller
Apply labels based on simple metadata rules
Status reporting
Test in local cluster


Immediate tasks (current feature completion):

Wire up filter helpers in controller - Already done (Pod & Node cases call FilterPodList/FilterNodeList)
Test the filters - Verify namespace/label selectors actually reduce list results
Handle edge cases - Empty match criteria, nil checks (already covered in helpers)
Next priority resources to implement (from TargetKind enum & todo):

Deployment - Has DeploymentMatchCriteria defined but no matcher/filter/reconciler case
Namespace - Listed in TargetKind enum, no match criteria yet
Service - Listed in TargetKind enum, no match criteria yet
StatefulSet, DaemonSet, ReplicaSet - Listed in TargetKind enum, no match criteria yet
Job, CronJob - Listed in TargetKind enum, no match criteria yet
Recommended next steps:

Test current Pod/Node filtering works correctly
Implement Deployment matching (already has DeploymentMatchCriteria spec):
Create matchinglogic/deployment_match.go with MatchesDeploymentDetailed
Create FilterDeploymentList helper (namespace + labels only)
Add Deployment case to reconciler switch
Add RBAC markers for Deployments
Add Service/Namespace matchers (simpler—mostly common criteria)
Tackle workload resources (StatefulSet, DaemonSet, etc.) using similar patterns
Current state: Pod & Node are functionally complete for filtering. Deployment is next logical target since match criteria already exist.

