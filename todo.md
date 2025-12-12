# Autolabeller TODO List

## Priority 1: Core API Definition
- [ ] **T1.1**: Define ClassificationRuleSpec fields (targetKind, match, labels, conflictPolicy)
- [ ] **T1.2**: Define match criteria types (ImageMatch, ResourceMatch, AnnotationMatch, LabelMatch, NamespaceSelectorMatch, NodeSelectorMatch, TaintMatch)
- [ ] **T1.3**: Add validation markers and CEL rules for ClassificationRule
- [ ] **T1.4**: Define ClassificationRuleStatus fields (matchedResourceCount, lastReconciled, conflicts, lastError)
- [ ] **T1.5**: Run `make generate` and `make manifests` to update CRD
- [ ] **T1.6**: Create sample ClassificationRule manifests in config/samples/

## Priority 2: Match Criteria Engine
- [ ] **T2.1**: Implement match criteria interface and registry
- [ ] **T2.2**: Implement ImageMatcher for container image pattern matching
- [ ] **T2.3**: Implement ResourceRequirementsMatcher (CPU, memory comparison)
- [ ] **T2.4**: Implement AnnotationMatcher with regex support
- [ ] **T2.5**: Implement LabelMatcher
- [ ] **T2.6**: Implement NamespaceSelectorMatcher
- [ ] **T2.7**: Implement NodeSelectorMatcher
- [ ] **T2.8**: Implement TaintMatcher
- [ ] **T2.9**: Implement composite matcher that evaluates all criteria with AND logic
- [ ] **T2.10**: Write unit tests for all matchers (target: >85% coverage)

## Priority 3: Pod Classification Logic
- [ ] **T3.1**: Create PodClassifier struct with Pod evaluation logic
- [ ] **T3.2**: Implement Pod listing and filtering logic
- [ ] **T3.3**: Implement Pod matching against ClassificationRule criteria
- [ ] **T3.4**: Create PodLabelPatcher for non-destructive label application
- [ ] **T3.5**: Implement label ownership tracking (annotation-based)
- [ ] **T3.6**: Implement conflict detection for Pod labels
- [ ] **T3.7**: Handle conflict resolution strategies (Overwrite, Merge, Ignore, Report)
- [ ] **T3.8**: Implement Pod reconciliation in ClassificationRuleReconciler
- [ ] **T3.9**: Add Pod event handlers to trigger reconciliation
- [ ] **T3.10**: Write unit tests for PodClassifier
- [ ] **T3.11**: Write integration tests for Pod labeling

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

## Priority 5: Status & Observability
- [ ] **T5.1**: Implement status update logic in reconciler
- [ ] **T5.2**: Populate matchedResourceCount during reconciliation
- [ ] **T5.3**: Update lastReconciled timestamp
- [ ] **T5.4**: Implement conflict tracking in status
- [ ] **T5.5**: Implement error tracking in status
- [ ] **T5.6**: Add Kubernetes conditions (Available, Progressing, Degraded)
- [ ] **T5.7**: Create Prometheus metrics for rules processed
- [ ] **T5.8**: Create Prometheus metrics for resources labeled
- [ ] **T5.9**: Create Prometheus metrics for conflicts detected
- [ ] **T5.10**: Create Prometheus metrics for reconciliation duration
- [ ] **T5.11**: Create Prometheus metrics for errors
- [ ] **T5.12**: Implement proper logging throughout reconciliation

## Priority 6: RBAC & Security
- [ ] **T6.1**: Review and update ClusterRole permissions
- [ ] **T6.2**: Add permissions for Pod watching/patching
- [ ] **T6.3**: Add permissions for Node watching/patching
- [ ] **T6.4**: Add permissions for ServiceMonitor (Prometheus)
- [ ] **T6.5**: Update ServiceAccount configuration
- [ ] **T6.6**: Configure RoleBinding for RBAC
- [ ] **T6.7**: Review network policies
- [ ] **T6.8**: Ensure metrics endpoint is properly secured
- [ ] **T6.9**: Test Pod Security Policy (Restricted) enforcement

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


MILESTONE 1:
Scaffold project and CRD
Implement Pod labeling controller
Apply labels based on simple metadata rules
Status reporting
Test in local cluster