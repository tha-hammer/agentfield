# Silmari Rebrand Manifest

## Summary
This manifest records the completed Silmari first-party rebrand on `integration/silmari-rebrand-agentfield`, the compatibility-sensitive AgentField-family identifiers intentionally preserved in place, and the verification evidence for the final QA matrix.


## Audited Files
| Path | Action | Verification |
|---|---|---|
| .github/BRANCH_PROTECTION.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/ISSUE_TEMPLATE/community-project.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/ISSUE_TEMPLATE/feature_request.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/ISSUE_TEMPLATE/question.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/ISSUE_TEMPLATE/task.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/pull_request_template.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/workflows/contributor-reminders.yml | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/workflows/control-plane.yml | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/workflows/docker.yml | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/workflows/functional-tests.yml | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/workflows/memory-metrics.yml | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/workflows/release.yml | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .github/workflows/update-download-stats.yml | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| .gitignore | audited-no-change | visual review |
| CLAUDE.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| CODE_OF_CONDUCT.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| SECURITY.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| SUPPORT.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| assets/utm-links.csv | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| control-plane/README.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | rebranded-with-preserved-identifiers | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | rebranded-with-preserved-identifiers | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/anti-patterns.md | rebranded | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/capability-playbook.md | rebranded | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/cli-toolkit.md | rebranded-with-preserved-identifiers | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/examples-map.md | rebranded | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | rebranded-with-preserved-identifiers | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/memory-events.md | rebranded | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/model-selection.md | rebranded | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/patterns-emerge.md | rebranded | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | audited-no-change | mirror parity; visual review |
| control-plane/internal/skillkit/skill_data/agentfield/references/project-claude-template.md | rebranded-with-preserved-identifiers | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | rebranded-with-preserved-identifiers | mirror parity |
| control-plane/internal/skillkit/skill_data/agentfield/references/triggers.md | audited-no-change | mirror parity; visual review |
| control-plane/internal/skillkit/skill_data/agentfield/references/verification.md | rebranded | mirror parity |
| control-plane/internal/skillkit/skillkit_test.go | rebranded | go test ./internal/skillkit |
| control-plane/internal/templates/docker/.env.example.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | rebranded-with-preserved-identifiers | go test ./internal/templates/... |
| control-plane/internal/templates/go/.env.example.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/go/README.md.tmpl | rebranded-with-preserved-identifiers | go test ./internal/templates/... |
| control-plane/internal/templates/go/go.mod.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/go/main.go.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/go/reasoners.go.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/python/.env.example.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/python/README.md.tmpl | rebranded-with-preserved-identifiers | go test ./internal/templates/... |
| control-plane/internal/templates/python/main.py.tmpl | rebranded-with-preserved-identifiers | go test ./internal/templates/... |
| control-plane/internal/templates/python/reasoners.py.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/python/requirements.txt.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/templates.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/templates_test.go | rebranded-with-preserved-identifiers | go test ./internal/templates/... |
| control-plane/internal/templates/typescript/.env.example.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/typescript/README.md.tmpl | rebranded-with-preserved-identifiers | go test ./internal/templates/... |
| control-plane/internal/templates/typescript/main.ts.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/typescript/package.json.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/internal/templates/typescript/reasoners.ts.tmpl | excluded-runtime-compatibility | go test ./internal/templates/... |
| control-plane/migrations/README.md | rebranded | `python3 -m pytest tests/test_docs_specs_branding.py` |
| control-plane/scripts/README.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| control-plane/tools/perf/README.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| control-plane/web/client/.env | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/.env.example | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/.env.production | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/.eslint-suppressions.json | audited-no-change | `cd control-plane/web/client && npm run lint` |
| control-plane/web/client/package.json | audited-no-change | `cd control-plane/web/client && npm run lint` |
| control-plane/web/client/src/components/AdminTokenPrompt.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/AppLayout.tsx | rebranded | `cd control-plane/web/client && npm run test`; visual review |
| control-plane/web/client/src/components/AppSidebar.tsx | rebranded | `cd control-plane/web/client && npm run test`; visual review |
| control-plane/web/client/src/components/AuthGuard.tsx | rebranded | `cd control-plane/web/client && npm run test`; visual review |
| control-plane/web/client/src/components/HealthStrip.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/execution/TechnicalDetailsPanel.tsx | rebranded-with-preserved-identifiers | `cd control-plane/web/client && npm run test`; visual review |
| control-plane/web/client/src/components/forms/ConfigField.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/forms/ConfigurationForm.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/forms/ConfigurationWizard.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/forms/EnvironmentVariableForm.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/nodes/EnhancedNodeDetailHeader.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/nodes/EnhancedNodesHeader.tsx | rebranded | `cd control-plane/web/client && npm run test`; visual review |
| control-plane/web/client/src/components/reasoners/EmptyReasonersState.tsx | rebranded | `cd control-plane/web/client && npm run test`; visual review |
| control-plane/web/client/src/components/status/StatusBadge.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/StatusRefreshButton.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/StatusRefreshButton.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/UnifiedStatusIndicator.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/UnifiedStatusIndicator.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/index.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/ui/data-formatters.tsx | rebranded | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/config/navigation.ts | rebranded-with-preserved-identifiers | `cd control-plane/web/client && npm run test`; visual review |
| control-plane/web/client/src/contexts/ModeContext.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/hooks/queries/useAgents.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/pages/AccessManagementPage.tsx | rebranded-with-preserved-identifiers | `cd control-plane/web/client && npm run test`; visual review |
| control-plane/web/client/src/pages/AgentsPage.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/pages/NewDashboardPage.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | rebranded-with-preserved-identifiers | `cd control-plane/web/client && npm run test`; visual review |
| control-plane/web/client/src/services/api.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/services/configurationApi.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/components/AppLayout.test.tsx | rebranded | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/components/AppSidebar.test.tsx | rebranded | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/components/AuthGuard.test.tsx | rebranded | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/components/ConfigurationWizard.test.tsx | rebranded | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/components/execution/ExecutionComponents.test.tsx | rebranded-with-preserved-identifiers | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/components/forms/configuration-form.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/components/nodes/EnhancedNodesHeader.test.tsx | rebranded | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/components/reasoners/EmptyReasonersState.test.tsx | rebranded | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/components/status/StatusBadge.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/config/navigation.test.ts | rebranded-with-preserved-identifiers | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/pages/AccessManagementPage.test.tsx | rebranded-with-preserved-identifiers | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/pages/AgentsPage.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/pages/NewSettingsPage.restored.test.tsx | rebranded-with-preserved-identifiers | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/services/didApi.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/services/identityApi.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/services/observabilityWebhookApi.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/ui/notification.test.tsx | rebranded | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/test/utils/formattingUtils.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/utils/schemaUtils.test.ts | rebranded | `cd control-plane/web/client && npm run test` |
| control-plane/web/client/src/utils/lifecycle-status.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/utils/lifecycle-status.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/utils/node-status.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/utils/node-status.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/ui_embed.go | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| deployments/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/docker/Dockerfile.control-plane | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/docker/Dockerfile.demo-go-agent | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/docker/Dockerfile.demo-python-agent | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/docker/Dockerfile.go-agent | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/docker/Dockerfile.python-agent | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/docker/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/docker/docker-compose.yml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/Chart.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/NOTES.txt | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/_helpers.tpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/api-auth-secret.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/control-plane-deployment.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/control-plane-ingress.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/control-plane-pvc.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/control-plane-service.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/demo-agent-deployment.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/demo-agent-service.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/demo-python-agent-configmap.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/demo-python-agent-deployment.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/demo-python-agent-service.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/postgres-secret.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/postgres-service.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/templates/postgres-statefulset.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/helm/agentfield/values.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/base/control-plane-deployment.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/base/control-plane-pvc.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/base/control-plane-service.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/base/kustomization.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/local-demo/demo-agent-deployment.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/local-demo/demo-agent-service.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/postgres-demo/demo-agent-deployment.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/postgres-demo/demo-agent-service.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/postgres-demo/patch-control-plane-postgres.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/postgres-demo/postgres-secret.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/postgres-demo/postgres-service.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/postgres-demo/postgres-statefulset.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/python-demo/demo-python-agent-configmap.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/python-demo/demo-python-agent-deployment.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/kubernetes/overlays/python-demo/demo-python-agent-service.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| deployments/railway/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| docs/ARCHITECTURE.md | rebranded | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/CONTRIBUTING.md | rebranded | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/COVERAGE.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/DEVELOPMENT.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/ENVIRONMENT_VARIABLES.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/RELEASE.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/api/AGENT_NODE_LOGS.md | audited-no-change | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/design/execution-observability-rfc.md | rebranded | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/design/harness-v2-design.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| docs/silmari-rebrand-manifest.md | rebranded-with-preserved-identifiers | manual review; targeted SDK metadata branding pytest |
| examples/README.md | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/benchmarks/100k-scale/README.md | rebranded | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/benchmarks/100k-scale/analyze.py | rebranded-with-preserved-identifiers | `python3 examples/benchmarks/100k-scale/analyze.py` |
| examples/benchmarks/100k-scale/crewai-bench/benchmark.py | rebranded | manual review |
| examples/benchmarks/100k-scale/go-bench/go.mod | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/go-bench/main.go | rebranded-with-preserved-identifiers | manual review |
| examples/benchmarks/100k-scale/langchain-bench/benchmark.py | rebranded | manual review |
| examples/benchmarks/100k-scale/python-bench/benchmark.py | rebranded-with-preserved-identifiers | manual review |
| examples/benchmarks/100k-scale/requirements.txt | rebranded-with-preserved-identifiers | manual review |
| examples/benchmarks/100k-scale/results/AgentField_Go.json | excluded-historical | visual review |
| examples/benchmarks/100k-scale/results/AgentField_Python.json | excluded-historical | visual review |
| examples/benchmarks/100k-scale/results/AgentField_TypeScript.json | excluded-historical | visual review |
| examples/benchmarks/100k-scale/results/benchmark_summary.png | rebranded | visual review |
| examples/benchmarks/100k-scale/run_benchmarks.sh | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/benchmarks/100k-scale/ts-bench/benchmark.ts | rebranded-with-preserved-identifiers | manual review |
| examples/benchmarks/100k-scale/ts-bench/package-lock.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/ts-bench/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/e2e_resilience_tests/README.md | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/e2e_resilience_tests/agent_flaky.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/e2e_resilience_tests/agent_healthy.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/e2e_resilience_tests/agent_slow.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/e2e_resilience_tests/run_tests.sh | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_agent_nodes/cmd/multi_version/main.go | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/go_agent_nodes/cmd/serverless/main.go | rebranded-with-preserved-identifiers | manual review |
| examples/go_agent_nodes/go.mod | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_agent_nodes/main.go | rebranded-with-preserved-identifiers | manual review |
| examples/go_harness_demo/go.mod | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_harness_demo/main.go | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/agentic_rag/README.md | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/agentic_rag/main.py | rebranded-with-preserved-identifiers | manual review |
| examples/python_agent_nodes/agentic_rag/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/deep_research/README.md | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/deep_research/main.py | rebranded-with-preserved-identifiers | manual review |
| examples/python_agent_nodes/deep_research/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/deep_research/routers/planning.py | rebranded-with-preserved-identifiers | manual review |
| examples/python_agent_nodes/deep_research/routers/research.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/docker_hello_world/main.py | rebranded-with-preserved-identifiers | manual review |
| examples/python_agent_nodes/documentation_chatbot/README.md | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/documentation_chatbot/install.sh | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/main.py | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/documentation_chatbot/pipeline_utils.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/product_context.py | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/documentation_chatbot/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/routers/ingestion.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/routers/qa.py | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/documentation_chatbot/routers/query_planning.py | rebranded-with-preserved-identifiers | manual review |
| examples/python_agent_nodes/documentation_chatbot/routers/retrieval.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/hello_world/main.py | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/hello_world_rag/README.md | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/hello_world_rag/main.py | rebranded-with-preserved-identifiers | manual review |
| examples/python_agent_nodes/hello_world_rag/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/hello_world_rag/test_pydantic_skill.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/image_generation_hello_world/README.md | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/image_generation_hello_world/main.py | rebranded-with-preserved-identifiers | manual review |
| examples/python_agent_nodes/image_generation_hello_world/requirements.txt | rebranded-with-preserved-identifiers | manual review |
| examples/python_agent_nodes/multi_version/main.py | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/permission_agent_a/main.py | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/permission_agent_b/README.md | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/permission_agent_b/main.py | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/rag_evaluation/Dockerfile | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/README.md | audited-no-change | manual review |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | excluded-runtime-compatibility | manual review |
| examples/python_agent_nodes/rag_evaluation/main.py | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/rag_evaluation/rag_eval_client.py | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/rag_evaluation/reasoners/__init__.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/ui/app/page.tsx | rebranded | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/rag_evaluation/ui/package.json | rebranded | `npm run lint` |
| examples/python_agent_nodes/rag_evaluation/ui/public/silmari-icon-dark.svg | rebranded | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/rag_evaluation/ui/public/silmari-icon-light.svg | rebranded | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/rag_evaluation/ui/public/silmari-logo-dark.svg | rebranded | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/rag_evaluation/ui/public/silmari-logo-light.svg | rebranded | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/serverless_hello/main.py | rebranded-with-preserved-identifiers | manual review |
| examples/python_agent_nodes/simulation_engine/README.md | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/simulation_engine/example.ipynb | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/main.py | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/python_agent_nodes/simulation_engine/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/aggregation.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/decision.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/entity.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/scenario.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/simulation.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/tool_calling/orchestrator.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/tool_calling/worker.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/waiting_state/main.py | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/tests/test_silmari_branding.py | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/triggers-demo/Dockerfile | rebranded-with-preserved-identifiers | manual review |
| examples/triggers-demo/README.md | rebranded | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/triggers-demo/agent.py | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/triggers-demo/docker-compose.yml | rebranded-with-preserved-identifiers | manual review |
| examples/triggers-demo/scripts/fire-events.sh | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/ts-node-examples/discovery-memory/main.ts | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/ts-node-examples/init-example/main.ts | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/ts-node-examples/init-example/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/init-example/reasoners.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/multi-version/main.ts | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/ts-node-examples/package-lock.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/permission-agent-a/main.ts | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/ts-node-examples/permission-agent-b/main.ts | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/ts-node-examples/serverless-hello/main.ts | rebranded-with-preserved-identifiers | manual review |
| examples/ts-node-examples/simulation/main.ts | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/ts-node-examples/simulation/routers/aggregation.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/routers/decision.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/routers/entity.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/routers/scenario.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/routers/simulation.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/verifiable-credentials/README.md | rebranded-with-preserved-identifiers | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/ts-node-examples/verifiable-credentials/main.ts | rebranded-with-preserved-identifiers | manual review |
| examples/ts-node-examples/verifiable-credentials/reasoners.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/waiting-state/main.ts | audited-no-change | `python3 -m unittest discover -s examples/tests -p 'test_*.py'` |
| examples/ts-node-examples/waiting-state/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/orchestrator.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/package-lock.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/test_orchestrator.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/worker.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| scripts/check-silmari-rebrand.sh | rebranded | ./scripts/check-silmari-rebrand.sh; python3 -m pytest tests/test_check_silmari_rebrand.py |
| scripts/collect_silmari_rebrand_inventory.py | rebranded | python3 -m pytest tests/test_collect_silmari_rebrand_inventory.py |
| sdk/go/README.md | rebranded-with-preserved-identifiers | scoped SDK brand and compatibility greps |
| sdk/go/agent/agent.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/agent_accepts_webhook_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/agent_did.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/agent_did_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/agent_lifecycle.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/agent_lifecycle_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/agent_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/agent_trigger_origin_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/cli.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/cli_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/did_async_additional_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/discovery.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/harness.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/harness_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/middleware_additional_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/note.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/process_logs.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/process_logs_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/registration_integration_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/router_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/agent/workflow_event_additional_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/ai/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/ai/tool_calling.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/ai/tool_calling_additional_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/ai/tool_calling_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/client/client.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/client/client_invariant_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/client/client_low_coverage_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/client/client_redirect_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/client/client_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/did/types.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/go.mod | audited-no-change | direct Go module-path grep; targeted SDK metadata branding pytest |
| sdk/go/harness/codex_integration_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/harness/coverage_branches_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/harness/gemini_integration_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/harness/opencode.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/harness/runner.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/types/types.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/MANIFEST.in | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/README.md | rebranded-with-preserved-identifiers | scoped SDK brand and compatibility greps |
| sdk/python/agentfield.egg-info/PKG-INFO | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield.egg-info/SOURCES.txt | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield.egg-info/top_level.txt | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent_ai.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent_cli.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent_discovery.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent_field_handler.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent_pause.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent_server.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent_serverless.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent_vc.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/agent_workflow.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/async_config.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/async_execution_manager.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/cancel.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/client.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/connection_manager.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/decorators.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/did_auth.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/did_manager.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/exceptions.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/execution_context.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/fixtures/triggers/github.json | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/harness/__init__.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/harness/_runner.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/harness/providers/__init__.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/harness/providers/_base.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/harness/providers/_factory.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/harness/providers/claude.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/harness/providers/codex.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/harness/providers/gemini.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/harness/providers/opencode.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/http_connection_manager.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/logger.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/media_providers.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/media_router.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/memory.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/memory_events.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/multimodal_response.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/node_logs.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/pydantic_utils.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/rate_limiter.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/status.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/testing.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/tool_calling.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/triggers.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/types.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/vc_generator.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/verification.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/agentfield/vision.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/pyproject.toml | rebranded-with-preserved-identifiers | scoped SDK brand and compatibility greps |
| sdk/python/requirements-dev.txt | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/conftest.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/debug_complex_json.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/helpers.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/integration/conftest.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/integration/test_agentfield_end_to_end.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_accepts_webhook.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_ai.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_ai_comprehensive.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_ai_coverage_additions.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_ai_deadlock_recovery.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_ai_final90.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_bigfiles_final90.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_call.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_cli.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_core.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_coverage_additions.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_field_handler.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_graceful_shutdown.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_helpers.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_instance_id.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_integration.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_lifecycle_invariants.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_networking.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_registry.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_resilience.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_server.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_server_extended.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_serverless.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_utils.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_workflow.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_workflow_extended.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_agent_workflow_registration.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_ai_config.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_approval.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_async_config.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_async_execution.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_async_execution_manager_comprehensive.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_async_execution_manager_final90.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_async_execution_manager_paths.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_cancel.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_client.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_client_auth.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_client_bigfiles_coverage.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_client_coverage_additions.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_client_execution_paths.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_client_execution_vc_payload.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_client_laser_push.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_client_lifecycle.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_client_unit.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_connection_manager.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_connection_manager_invariants.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_cost_tracker.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_decorator_code_origin.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_decorators.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_did_auth.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_did_auth_invariants.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_did_manager.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_did_manager_error_paths.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_exceptions.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_execution_context_core.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_execution_context_coverage_additions.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_execution_context_parent_vc.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_execution_logger.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_execution_state.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_execution_state_invariants.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_agent_wiring.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_ai_schema_repair.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_cli.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_cost_estimation.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_factory.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_functional.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_provider_claude.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_provider_codex.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_provider_gemini.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_provider_opencode.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_runner.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_schema.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_harness_types.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_http_connection_manager.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_image_config.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_invariants.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_litellm_adapters.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_media_integration.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_media_providers.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_media_providers_additional.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_media_router.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_memory_bigfiles_coverage.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_memory_client_core.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_memory_coverage_additions.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_memory_events.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_memory_events_additional.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_memory_flow_core.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_memory_invariants.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_memory_performance.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_multimodal.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_multimodal_response.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_multimodal_response_additional.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_multimodal_response_cost.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_node_logs.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_openrouter_audio.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_openrouter_video.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_pydantic_utils.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_rate_limiter_core.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_reasoner_path_normalization.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_result_cache.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_result_cache_bigfiles_coverage.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_router.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_run_cli_env.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_sdk_metadata_branding.py | rebranded-with-preserved-identifiers | targeted SDK metadata branding pytest |
| sdk/python/tests/test_simulate_trigger.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_skill_pydantic_models.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_status_utils.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_tool_calling.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_tool_calling_error_paths.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_trigger_context.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_trigger_param_binding.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_types.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_utils.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_vc_generator.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_vc_generator_error_paths.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_verification.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_video_output.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_vision.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/tests/test_workflow_parent_child.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/uv.lock | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/README.md | rebranded-with-preserved-identifiers | scoped SDK brand and compatibility greps |
| sdk/typescript/package-lock.json | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/package.json | rebranded-with-preserved-identifiers | scoped SDK brand and compatibility greps |
| sdk/typescript/scripts/did-smoke.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/src/agent/Agent.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/src/agent/cancel.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/src/agent/processLogs.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/src/ai/ToolCalling.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/src/approval/ApprovalClient.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/src/harness/providers/opencode.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/src/status/ExecutionStatus.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/src/verification/LocalVerifier.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/ai_multimodal_response_extra.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/cancel.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/did_client_methods.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/harness_functional.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/harness_runner.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/harness_schema.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/integration.e2e.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/memory_performance.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/multimodal.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/tests/process_logs.test.ts | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/SKILL.md | rebranded-with-preserved-identifiers | visible brand check; mirror parity |
| skills/agentfield/commands/agentfield.md | rebranded-with-preserved-identifiers | visible brand check; mirror parity |
| skills/agentfield/references/anti-patterns.md | rebranded | visible brand check |
| skills/agentfield/references/capability-playbook.md | rebranded | visible brand check |
| skills/agentfield/references/cli-toolkit.md | rebranded-with-preserved-identifiers | visible brand check |
| skills/agentfield/references/examples-map.md | rebranded | visible brand check |
| skills/agentfield/references/live-docs.md | rebranded-with-preserved-identifiers | visible brand check |
| skills/agentfield/references/memory-events.md | rebranded | visible brand check |
| skills/agentfield/references/model-selection.md | rebranded | visible brand check |
| skills/agentfield/references/patterns-emerge.md | rebranded | visible brand check |
| skills/agentfield/references/primitives-snapshot.md | audited-no-change | visual review |
| skills/agentfield/references/project-claude-template.md | rebranded-with-preserved-identifiers | visible brand check |
| skills/agentfield/references/scaffold-recipe.md | rebranded-with-preserved-identifiers | visible brand check |
| skills/agentfield/references/triggers.md | audited-no-change | visual review |
| skills/agentfield/references/verification.md | rebranded | visible brand check |
| specs/README.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| specs/agentplane-ui-api-spec.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| specs/architecture-overview.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| specs/control-plane.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| specs/data-flow.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| specs/deployment.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| specs/sdk-go.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| specs/sdk-python.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| specs/security.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| specs/viewer.html | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/web-ui.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| tests/functional/.env.example | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/README.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| tests/functional/agents/call_chain_agents.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/agents/docs_quick_start_agent.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/agents/memory_agent.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/agents/memory_events_agent.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/agents/memory_events_decorator_agent.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/agents/quick_start_agent.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/agents/router_prefix_agent.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/agents/scoping_agent.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/conftest.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/docker/Dockerfile.log-demo-node | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/docker/Dockerfile.test-runner | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/docker/LOG_DEMO.md | rebranded-with-preserved-identifiers | `python3 -m pytest tests/test_docs_specs_branding.py` |
| tests/functional/docker/agentfield-test.yaml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/docker/docker-compose.local.yml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/docker/docker-compose.log-demo.yml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/docker/docker-compose.postgres.yml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/docker/log-demo-node/log-demo.mjs | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/docker/wait-for-services.sh | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/go_agents/cmd/discovery/main.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/go_agents/cmd/hello/main.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/go_agents/cmd/serverless/main.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/go_agents/go.mod | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/tests/test_app_call.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/tests/test_go_sdk_cli.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/tests/test_go_sdk_discovery.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/tests/test_serverless_agents.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/tests/test_ui_node_logs_proxy.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/tests/test_vc_cli.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/tests/test_waiting_state.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/ts_agents/echo-agent.mjs | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/ts_agents/serverless-agent.mjs | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/utils/agent_server.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/utils/naming.py | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/test_check_silmari_rebrand.py | rebranded-with-preserved-identifiers | ./scripts/check-silmari-rebrand.sh; python3 -m pytest tests/test_check_silmari_rebrand.py |
| tests/test_collect_silmari_rebrand_inventory.py | rebranded | python3 -m pytest tests/test_collect_silmari_rebrand_inventory.py |
| tests/test_docs_specs_branding.py | rebranded-with-preserved-identifiers | ./scripts/check-silmari-rebrand.sh; python3 -m pytest tests/test_docs_specs_branding.py |
| thoughts/searchable/shared/plans/2026-06-29-silmari-rebrand-baseline-inventory.json | rebranded | python3 -m pytest tests/test_collect_silmari_rebrand_inventory.py |

## Preserved Non-Silmari Identifiers
| Path | Identifier | Category | Reason |
|---|---|---|---|
| .github/BRANCH_PROTECTION.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| .github/ISSUE_TEMPLATE/community-project.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| .github/ISSUE_TEMPLATE/community-project.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| .github/ISSUE_TEMPLATE/feature_request.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| .github/ISSUE_TEMPLATE/question.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| .github/ISSUE_TEMPLATE/task.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| .github/ISSUE_TEMPLATE/task.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| .github/pull_request_template.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| .github/workflows/contributor-reminders.yml | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| .github/workflows/control-plane.yml | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| .github/workflows/docker.yml | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| .github/workflows/functional-tests.yml | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| .github/workflows/memory-metrics.yml | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| .github/workflows/release.yml | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| .github/workflows/release.yml | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| .github/workflows/release.yml | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| .github/workflows/update-download-stats.yml | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| .github/workflows/update-download-stats.yml | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| CLAUDE.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| CLAUDE.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| CLAUDE.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| CLAUDE.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| CODE_OF_CONDUCT.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| CODE_OF_CONDUCT.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| README.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| SECURITY.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| SECURITY.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| SUPPORT.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| SUPPORT.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| assets/utm-links.csv | agentfield.ai | published-link-target | All target URLs in this CSV keep the published AgentField domain because the repo does not provide a verified Silmari replacement URL. |
| assets/utm-links.csv | agentfield.ai | published-link-target | Published UTM inventory targets remain on agentfield.ai until verified Silmari replacements exist; covers all occurrences in this file. |
| control-plane/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/README.md | `AGENTFIELD*` | env-var | Control-plane configuration examples retain the stable env var namespace. |
| control-plane/README.md | `agentfield.ai*` | published-link-target | The install script remains published at the legacy domain referenced by this README. |
| control-plane/README.md | `agentfield/control-plane*` | docker-image-or-volume | Docker quick-start examples retain the published image name. |
| control-plane/README.md | `cmd/agentfield-server*` | cli-command | Source-run instructions retain the current standalone server command path. |
| control-plane/README.md | `config/agentfield.yaml` | yaml-config-path | The default YAML config path remains stable and is documented as a peer surface to env vars. |
| control-plane/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | `agentfield` skill | skill-or-repo-slug | The embedded mirror keeps the stable skill name in user instructions. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | agentfield.ai | published-link-target | The embedded mirror keeps the same published docs host references as the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | aliases: [agentfield-multi-reasoner-builder] | skill-or-repo-slug | The embedded mirror keeps the published alias stable for installer compatibility. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | name: agentfield | skill-or-repo-slug | The embedded mirror must preserve the shipped skill frontmatter slug exactly. |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | /agentfield | cli-command | The embedded command mirror preserves the derived slash command surface. |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | `agentfield` skill | skill-or-repo-slug | The embedded command mirror preserves the stable skill slug in help text. |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | https://agentfield.ai/llms.txt | published-link-target | The embedded command mirror preserves the current live-docs endpoint. |
| control-plane/internal/skillkit/skill_data/agentfield/references/anti-patterns.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/capability-playbook.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/cli-toolkit.md | AGENTFIELD_HOME | env-var | Skill CLI docs keep AGENTFIELD_HOME stable as the compatibility env var for the canonical skill store path; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/cli-toolkit.md | AGENTFIELD_HOME | env-var | The embedded mirror preserves the documented skill-store env var override. |
| control-plane/internal/skillkit/skill_data/agentfield/references/cli-toolkit.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/examples-map.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | AGENTFIELD_HOME | env-var | Skill live-doc guidance keeps AGENTFIELD_HOME stable as the compatibility cache-directory env var; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | AGENTFIELD_HOME | env-var | The embedded mirror preserves the cache-path env var guidance from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | agentfield.ai | published-link-target | The embedded mirror preserves the published docs host references from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/memory-events.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/model-selection.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/patterns-emerge.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | AGENTFIELD_SERVER | env-var | The embedded mirror preserves the current server env var example from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | agentfield.ai | published-link-target | The embedded mirror preserves the offline fallback docs host references from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | sdk/python/agentfield/agent.py | import-module-path | The embedded mirror preserves the existing Python package path reference from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/project-claude-template.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/project-claude-template.md | `agentfield` skill | skill-or-repo-slug | The embedded mirror preserves the stable skill slug in generated handoff guidance. |
| control-plane/internal/skillkit/skill_data/agentfield/references/project-claude-template.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/project-claude-template.md | agentfield/control-plane:latest | docker-image-or-volume | The embedded mirror preserves the published control-plane image name from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | AGENTFIELD_HTTP_ADDR | env-var | The embedded scaffold mirror preserves the HTTP bind-address env var name from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | AGENTFIELD_HTTP_PORT | env-var | The embedded scaffold mirror preserves the HTTP port env var name from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | AGENTFIELD_SERVER | env-var | The embedded scaffold mirror preserves the control-plane URL env var name from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | AGENTFIELD_STORAGE_MODE | env-var | The embedded scaffold mirror preserves the storage-mode env var name from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | agentfield | package-name | The embedded scaffold mirror preserves the Python package name mentions from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | agentfield-data | docker-image-or-volume | The embedded scaffold mirror preserves the sample compose volume name from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | agentfield/control-plane:latest | docker-image-or-volume | The embedded scaffold mirror preserves the published control-plane image tag from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | from agentfield import | import-module-path | The embedded scaffold mirror preserves the installed package import path; this row covers all import examples in the file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/triggers.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/triggers.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/triggers.md | agentfield.ai/docs | published-link-target | The embedded trigger mirror preserves the currently published docs host reference from the canonical source. |
| control-plane/internal/skillkit/skill_data/agentfield/references/triggers.md | from agentfield import | import-module-path | The embedded trigger mirror preserves the installed Python import path; this row covers all import examples in the file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/verification.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/docker/.env.example.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/docker/.env.example.tmpl | AGENTFIELD_HTTP_PORT | env-var | The generated Docker environment example keeps the existing control-plane port override variable. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | AGENTFIELD_HTTP_ADDR | env-var | The generated Docker Compose control-plane service keeps the existing HTTP bind address environment variable. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | AGENTFIELD_HTTP_PORT | env-var | The generated Docker Compose scaffold keeps the existing HTTP port override variable. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | AGENTFIELD_SERVER | env-var | The generated agent container keeps the existing control-plane server environment variable contract. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | AGENTFIELD_STORAGE_MODE | env-var | The generated Docker Compose control-plane service keeps the existing storage mode environment variable. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | agentfield-data | docker-image-or-volume | This row covers both `agentfield-data` volume occurrences in the Compose file because the persisted Docker volume name must stay stable. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | agentfield/control-plane:latest | docker-image-or-volume | The generated Compose guidance keeps the published control-plane image name for compatibility. |
| control-plane/internal/templates/go/.env.example.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/go/.env.example.tmpl | AGENTFIELD_CONTROL_PLANE_URL | env-var | The generated Go environment example keeps the existing control-plane URL variable for local setup compatibility. |
| control-plane/internal/templates/go/README.md.tmpl | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/go/README.md.tmpl | af init | cli-command | Generated onboarding keeps the existing `af init` CLI entry point for scaffold compatibility. |
| control-plane/internal/templates/go/README.md.tmpl | af server | cli-command | Generated quick-start instructions must keep the existing `af server` control-plane command. |
| control-plane/internal/templates/go/README.md.tmpl | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/templates/go/README.md.tmpl | http://localhost:8080/api/v1/execute/ | api-path | The generated Go README keeps the published execute endpoint path in curl examples for backward-compatible local verification. |
| control-plane/internal/templates/go/README.md.tmpl | https://agentfield.ai/docs/learn | published-link-target | The README link target stays on the verified legacy docs domain until a Silmari docs URL exists. |
| control-plane/internal/templates/go/README.md.tmpl | https://agentfield.ai/docs/reference/sdks/go | published-link-target | The Go SDK reference target stays on the published legacy docs domain for backward-compatible links. |
| control-plane/internal/templates/go/go.mod.tmpl | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/go.mod.tmpl | agentfield | go-module-path | Go module or import compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/go.mod.tmpl | github.com/Agent-Field/agentfield/sdk/go | go-module-path | The generated Go module file must keep the published SDK module path and version requirement. |
| control-plane/internal/templates/go/main.go.tmpl | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/main.go.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/main.go.tmpl | github.com/Agent-Field/agentfield/sdk/go | go-module-path | This row covers both Go SDK imports in the generated main program because the published module path is a compatibility surface. |
| control-plane/internal/templates/go/reasoners.go.tmpl | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/reasoners.go.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/reasoners.go.tmpl | github.com/Agent-Field/agentfield/sdk/go | go-module-path | This row covers both Go SDK imports in the generated reasoner file because the published module path must stay stable. |
| control-plane/internal/templates/python/.env.example.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/python/.env.example.tmpl | AGENTFIELD_CONTROL_PLANE_URL | env-var | The generated Python environment example keeps the existing control-plane URL variable for local setup compatibility. |
| control-plane/internal/templates/python/README.md.tmpl | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/python/README.md.tmpl | af init | cli-command | Generated onboarding keeps the existing `af init` CLI entry point for scaffold compatibility. |
| control-plane/internal/templates/python/README.md.tmpl | af server | cli-command | Generated quick-start instructions must keep the existing `af server` control-plane command. |
| control-plane/internal/templates/python/README.md.tmpl | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/templates/python/README.md.tmpl | http://localhost:8080/api/v1/execute/ | api-path | The generated Python README keeps the published execute endpoint path in curl examples for backward-compatible local verification. |
| control-plane/internal/templates/python/README.md.tmpl | https://agentfield.ai/docs/learn | published-link-target | The README link target stays on the verified legacy docs domain until a Silmari docs URL exists. |
| control-plane/internal/templates/python/README.md.tmpl | https://agentfield.ai/docs/reference/sdks/python | published-link-target | The Python SDK reference target stays on the published legacy docs domain for backward-compatible links. |
| control-plane/internal/templates/python/main.py.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/python/main.py.tmpl | AGENTFIELD_SERVER | env-var | The generated Python entrypoint keeps the existing control-plane server environment variable contract. |
| control-plane/internal/templates/python/main.py.tmpl | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/python/main.py.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/python/main.py.tmpl | from agentfield import Agent, AIConfig | import-module-path | The generated Python entrypoint must import the published `agentfield` SDK module to run. |
| control-plane/internal/templates/python/reasoners.py.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/python/reasoners.py.tmpl | from agentfield import AgentRouter | import-module-path | The generated Python reasoner router must import the published `agentfield` SDK module. |
| control-plane/internal/templates/python/requirements.txt.tmpl | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| control-plane/internal/templates/python/requirements.txt.tmpl | agentfield | package-name | The generated Python dependency must keep the published package name from PyPI. |
| control-plane/internal/templates/templates.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/templates_test.go | @agentfield/sdk | package-name | This row covers all npm package assertions in the template tests because the published package name must stay stable. |
| control-plane/internal/templates/templates_test.go | AGENTFIELD | env-var | Template regression tests intentionally assert legacy AGENTFIELD env vars remain stable; covers all occurrences in this file. |
| control-plane/internal/templates/templates_test.go | AGENTFIELD_CONTROL_PLANE_URL | env-var | This row covers the Python and Go env example assertions that preserve the existing control-plane URL variable. |
| control-plane/internal/templates/templates_test.go | AGENTFIELD_HTTP_ADDR | env-var | The template-test assertions keep the existing Docker HTTP bind variable to verify Compose compatibility. |
| control-plane/internal/templates/templates_test.go | AGENTFIELD_HTTP_PORT | env-var | This row covers both Docker port override assertions in the template tests because the existing env var contract must stay stable. |
| control-plane/internal/templates/templates_test.go | AGENTFIELD_SERVER | env-var | This row covers the Python and Docker compatibility assertions that keep the existing server environment variable contract stable. |
| control-plane/internal/templates/templates_test.go | AGENTFIELD_STORAGE_MODE | env-var | The template-test assertions keep the existing Docker storage mode variable to verify Compose compatibility. |
| control-plane/internal/templates/templates_test.go | AGENTFIELD_URL | env-var | This row covers the TypeScript runtime and env example assertions that preserve the existing control-plane URL variable. |
| control-plane/internal/templates/templates_test.go | AgentField | test-fixture | Template regression tests intentionally mention AgentField when asserting removed prose and compatibility fixtures; covers all occurrences in this file. |
| control-plane/internal/templates/templates_test.go | AgentField | test-fixture | This row covers both README branding `assertNotContainsAny` fixtures that keep the legacy product name only to prove generated output no longer emits it. |
| control-plane/internal/templates/templates_test.go | AgentField agent node | test-fixture | The generated source-guidance test keeps this legacy docstring fragment only as a negative fixture to prove the old branding is gone. |
| control-plane/internal/templates/templates_test.go | AgentFieldURL: | test-fixture | The Go compatibility assertion keeps the legacy SDK config field name only to verify generated code stays byte-compatible. |
| control-plane/internal/templates/templates_test.go | agentfield | package-name | This row covers the generated Python dependency assertion that keeps the published PyPI package name stable. |
| control-plane/internal/templates/templates_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/internal/templates/templates_test.go | agentfield skill's scaffold-recipe | test-fixture | The generated source-guidance test keeps this legacy scaffold wording only as a negative fixture to prove the old branding is gone. |
| control-plane/internal/templates/templates_test.go | agentfield-data | docker-image-or-volume | This row covers the preserved Docker volume assertions in the template tests because the existing volume name must stay stable. |
| control-plane/internal/templates/templates_test.go | agentfield.ai | published-link-target | Template regression tests intentionally assert legacy published docs URLs remain on agentfield.ai; covers all occurrences in this file. |
| control-plane/internal/templates/templates_test.go | agentfield/control-plane:latest | docker-image-or-volume | This row covers both the fixture value and Compose compatibility assertions that keep the published control-plane image name stable. |
| control-plane/internal/templates/templates_test.go | from agentfield import Agent, AIConfig | import-module-path | The template-test assertions intentionally keep the published Python SDK import to verify generated output remains compatible. |
| control-plane/internal/templates/templates_test.go | from agentfield import AgentRouter | import-module-path | The template-test assertions intentionally keep the published Python router import to verify generated output remains compatible. |
| control-plane/internal/templates/templates_test.go | github.com/Agent-Field/agentfield/sdk/go | go-module-path | Template regression tests intentionally assert the published Go SDK module path remains stable; covers all occurrences in this file. |
| control-plane/internal/templates/templates_test.go | github.com/Agent-Field/agentfield/sdk/go | go-module-path | This row covers all Go SDK module-path assertions in the template tests because the published module path must stay stable. |
| control-plane/internal/templates/templates_test.go | http://localhost:8080/api/v1/execute/ | api-path | This row covers all README and compatibility assertions that preserve the published execute endpoint path for backward-compatible template output. |
| control-plane/internal/templates/templates_test.go | https://agentfield.ai/docs/learn | published-link-target | This row covers all template-test assertions that keep the verified legacy docs domain until a Silmari replacement URL exists. |
| control-plane/internal/templates/templates_test.go | https://agentfield.ai/docs/reference/sdks/go | published-link-target | The template-test assertions intentionally keep the published Go SDK reference URL until a Silmari replacement URL exists. |
| control-plane/internal/templates/templates_test.go | https://agentfield.ai/docs/reference/sdks/python | published-link-target | The template-test assertions intentionally keep the published Python SDK reference URL until a Silmari replacement URL exists. |
| control-plane/internal/templates/templates_test.go | https://agentfield.ai/docs/reference/sdks/typescript | published-link-target | The template-test assertions intentionally keep the published TypeScript SDK reference URL until a Silmari replacement URL exists. |
| control-plane/internal/templates/typescript/.env.example.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/.env.example.tmpl | AGENTFIELD_URL | env-var | The generated TypeScript environment example keeps the existing control-plane URL variable for local setup compatibility. |
| control-plane/internal/templates/typescript/README.md.tmpl | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/README.md.tmpl | af init | cli-command | Generated onboarding keeps the existing `af init` CLI entry point for scaffold compatibility. |
| control-plane/internal/templates/typescript/README.md.tmpl | af server | cli-command | Generated quick-start instructions must keep the existing `af server` control-plane command. |
| control-plane/internal/templates/typescript/README.md.tmpl | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/README.md.tmpl | http://localhost:8080/api/v1/execute/ | api-path | The generated TypeScript README keeps the published execute endpoint path in curl examples for backward-compatible local verification. |
| control-plane/internal/templates/typescript/README.md.tmpl | https://agentfield.ai/docs/learn | published-link-target | The README link target stays on the verified legacy docs domain until a Silmari docs URL exists. |
| control-plane/internal/templates/typescript/README.md.tmpl | https://agentfield.ai/docs/reference/sdks/typescript | published-link-target | The TypeScript SDK reference target stays on the published legacy docs domain for backward-compatible links. |
| control-plane/internal/templates/typescript/main.ts.tmpl | @agentfield/sdk | package-name | The generated TypeScript entrypoint must keep the published npm package name for runtime compatibility. |
| control-plane/internal/templates/typescript/main.ts.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/main.ts.tmpl | AGENTFIELD_URL | env-var | The generated TypeScript entrypoint keeps the existing control-plane URL environment variable contract. |
| control-plane/internal/templates/typescript/main.ts.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/package.json.tmpl | @agentfield/sdk | package-name | The generated TypeScript package manifest must keep the published npm package dependency name. |
| control-plane/internal/templates/typescript/package.json.tmpl | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/reasoners.ts.tmpl | @agentfield/sdk | package-name | The generated TypeScript reasoner router must keep the published npm package import path. |
| control-plane/internal/templates/typescript/reasoners.ts.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/migrations/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/migrations/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/scripts/README.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| control-plane/scripts/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/scripts/README.md | `Agent-Field/agentfield*` | skill-or-repo-slug | Install and release URLs retain the current GitHub slug. |
| control-plane/scripts/README.md | `agentfield-*` | cli-command | Release script output names retain the published binary and artifact stems. |
| control-plane/scripts/README.md | `agentfield` | cli-command | Shell examples still use the published CLI binary name. |
| control-plane/scripts/README.md | `~/.agentfield*` | cli-command | Installer-managed directories remain stable for existing user setups. |
| control-plane/scripts/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/tools/perf/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/tools/perf/README.md | `agentfield-*` | cli-command | Built binary names remain stable for the perf workflow. |
| control-plane/tools/perf/README.md | `agentfield-server` | cli-command | The Docker perf guide keeps the published server binary name. |
| control-plane/tools/perf/README.md | `cmd/agentfield-server*` | cli-command | Perf instructions retain the current standalone server command path. |
| control-plane/tools/perf/README.md | `config/agentfield.yaml` | yaml-config-path | The perf doc references the stable default YAML config path. |
| control-plane/tools/perf/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/web/client/.env | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/.env.example | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/.env.production | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/.eslint-suppressions.json | src/types/agentfield.ts | import-module-path | ESLint suppressions keep the stable `src/types/agentfield.ts` module path until runtime type imports are renamed separately; covers all occurrences in this file. |
| control-plane/web/client/src/components/AdminTokenPrompt.tsx | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/web/client/src/components/AppLayout.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/AppSidebar.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/AuthGuard.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/HealthStrip.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/execution/TechnicalDetailsPanel.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/execution/TechnicalDetailsPanel.tsx | agentfield_request_id | json-field | Execution details still read the stable API field directly; only the visible label changed to `Request ID`. |
| control-plane/web/client/src/components/forms/ConfigField.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/forms/ConfigurationForm.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/forms/ConfigurationWizard.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/forms/EnvironmentVariableForm.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/nodes/EnhancedNodeDetailHeader.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/nodes/EnhancedNodesHeader.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/reasoners/EmptyReasonersState.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/status/StatusBadge.tsx | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/components/status/StatusRefreshButton.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/components/status/StatusRefreshButton.tsx | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/components/status/UnifiedStatusIndicator.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/components/status/UnifiedStatusIndicator.tsx | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/components/status/index.ts | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/components/ui/data-formatters.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/config/navigation.ts | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| control-plane/web/client/src/config/navigation.ts | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/config/navigation.ts | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/web/client/src/config/navigation.ts | https://agentfield.ai/docs | published-link-target | No verified Silmari docs URL exists in this repo, so the visible label was rebranded while the published docs target stayed unchanged. |
| control-plane/web/client/src/config/navigation.ts | https://github.com/Agent-Field/agentfield/ | skill-or-repo-slug | The OSS repository slug remains the stable distribution target for the embedded UI resource link. |
| control-plane/web/client/src/contexts/ModeContext.tsx | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/hooks/queries/useAgents.ts | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/pages/AccessManagementPage.tsx | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/web/client/src/pages/AccessManagementPage.tsx | AGENTFIELD_AUTHORIZATION_ADMIN_TOKEN | env-var | Disabled-state guidance and browser-admin-token help text both reference the exact compatibility env var in this file. |
| control-plane/web/client/src/pages/AccessManagementPage.tsx | AGENTFIELD_AUTHORIZATION_ENABLED | env-var | Authorization setup guidance keeps the exact compatibility env var in the disabled-state example. |
| control-plane/web/client/src/pages/AccessManagementPage.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/pages/AccessManagementPage.tsx | config/agentfield.yaml | yaml-config-path | The access-management guidance points operators at the stable YAML config path for VC authorization setup. |
| control-plane/web/client/src/pages/AgentsPage.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/pages/NewDashboardPage.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | /api/v1/did/agentfield-server | api-path | The settings page still fetches the server DID from the stable compatibility endpoint. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | AGENTFIELD_NODE_LOG_MAX_TAIL_LINES | env-var | Node log proxy guidance keeps the exact max-tail-lines override env var. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | AGENTFIELD_NODE_LOG_PROXY_* | env-var | Node log proxy guidance keeps the exact override env-var family for backward-compatible operations. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | AGENTFIELD_SERVER | env-var | General settings copy and export example keep the exact compatibility env var for agent bootstrap. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | agentfield_server_did | json-field | DID bootstrap still reads the existing response field name in this file. |
| control-plane/web/client/src/services/api.ts | agentfield | api-path | Runtime API compatibility keeps agentfield stable for existing request and response contracts; covers all occurrences in this file. |
| control-plane/web/client/src/services/configurationApi.ts | agentfield | api-path | Runtime API compatibility keeps agentfield stable for existing request and response contracts; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/AppLayout.test.tsx | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/AppSidebar.test.tsx | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/ConfigurationWizard.test.tsx | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/ConfigurationWizard.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/execution/ExecutionComponents.test.tsx | agentfield_request_id | test-fixture | Test data continues to pin the legacy JSON field name while the visible technical-details label is rebranded. |
| control-plane/web/client/src/test/components/forms/configuration-form.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/status/StatusBadge.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/config/navigation.test.ts | https://agentfield.ai/docs | published-link-target | Navigation regression tests intentionally assert the published docs URL target remains on agentfield.ai; covers all occurrences in this file. |
| control-plane/web/client/src/test/config/navigation.test.ts | https://agentfield.ai/docs | test-fixture | Test assertions pin the preserved docs target behind the rebranded `Silmari Docs` label. |
| control-plane/web/client/src/test/config/navigation.test.ts | https://github.com/Agent-Field/agentfield/ | skill-or-repo-slug | Navigation regression tests intentionally assert the published GitHub repository slug remains Agent-Field/agentfield; covers all occurrences in this file. |
| control-plane/web/client/src/test/config/navigation.test.ts | https://github.com/Agent-Field/agentfield/ | test-fixture | Test assertions pin the preserved GitHub repo target behind the rebranded `Silmari GitHub` label. |
| control-plane/web/client/src/test/pages/AccessManagementPage.test.tsx | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/web/client/src/test/pages/AccessManagementPage.test.tsx | AGENTFIELD_AUTHORIZATION_ENABLED | test-fixture | Test coverage keeps the exact compatibility env var visible in disabled-state setup guidance. |
| control-plane/web/client/src/test/pages/AccessManagementPage.test.tsx | config/agentfield.yaml | test-fixture | Test coverage pins the stable YAML config path mentioned in setup guidance. |
| control-plane/web/client/src/test/pages/AccessManagementPage.test.tsx | config/agentfield.yaml | yaml-config-path | UI regression tests intentionally assert the legacy-compatible YAML config path remains visible in help copy; covers all occurrences in this file. |
| control-plane/web/client/src/test/pages/AgentsPage.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/pages/NewSettingsPage.restored.test.tsx | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/pages/NewSettingsPage.restored.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/pages/NewSettingsPage.restored.test.tsx | agentfield_server_did | test-fixture | Restored settings coverage keeps the DID response field name unchanged. |
| control-plane/web/client/src/test/pages/NewSettingsPage.restored.test.tsx | did:web:agentfield.example.test | test-fixture | Restored settings coverage keeps the legacy DID fixture value while asserting rebranded visible copy. |
| control-plane/web/client/src/test/services/didApi.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/services/identityApi.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/services/observabilityWebhookApi.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/ui/notification.test.tsx | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/utils/formattingUtils.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/utils/schemaUtils.test.ts | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/utils/lifecycle-status.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/utils/lifecycle-status.ts | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/utils/node-status.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/utils/node-status.ts | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/ui_embed.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/README.md | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| deployments/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| deployments/docker/Dockerfile.control-plane | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| deployments/docker/Dockerfile.demo-go-agent | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| deployments/docker/Dockerfile.demo-python-agent | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/docker/Dockerfile.go-agent | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/docker/Dockerfile.python-agent | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/docker/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/docker/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/docker/README.md | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| deployments/docker/docker-compose.yml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/docker/docker-compose.yml | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| deployments/docker/docker-compose.yml | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| deployments/helm/agentfield/Chart.yaml | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/helm/agentfield/Chart.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/helm/agentfield/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/helm/agentfield/README.md | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/NOTES.txt | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/NOTES.txt | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/_helpers.tpl | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/api-auth-secret.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/control-plane-deployment.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/control-plane-deployment.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/control-plane-ingress.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/control-plane-pvc.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/control-plane-service.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/demo-agent-deployment.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/demo-agent-deployment.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/demo-agent-service.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/demo-python-agent-configmap.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/demo-python-agent-configmap.yaml | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/demo-python-agent-configmap.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/demo-python-agent-deployment.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/demo-python-agent-deployment.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/demo-python-agent-service.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/postgres-secret.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/postgres-service.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/templates/postgres-statefulset.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/helm/agentfield/values.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/helm/agentfield/values.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/kubernetes/README.md | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/base/control-plane-deployment.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/kubernetes/base/control-plane-deployment.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/base/control-plane-pvc.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/base/control-plane-service.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/base/kustomization.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/local-demo/demo-agent-deployment.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/local-demo/demo-agent-deployment.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/local-demo/demo-agent-service.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/postgres-demo/demo-agent-deployment.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/postgres-demo/demo-agent-deployment.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/postgres-demo/demo-agent-service.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/postgres-demo/patch-control-plane-postgres.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/postgres-demo/patch-control-plane-postgres.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/postgres-demo/postgres-secret.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/postgres-demo/postgres-service.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/postgres-demo/postgres-statefulset.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/python-demo/demo-python-agent-configmap.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/python-demo/demo-python-agent-configmap.yaml | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/kubernetes/overlays/python-demo/demo-python-agent-configmap.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/python-demo/demo-python-agent-deployment.yaml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/python-demo/demo-python-agent-deployment.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/kubernetes/overlays/python-demo/demo-python-agent-service.yaml | agentfield | helm-chart-or-k8s-name | Helm chart and Kubernetes naming compatibility keeps agentfield stable for deployment examples; covers all occurrences in this file. |
| deployments/railway/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| deployments/railway/README.md | Agent-Field | historical-record | Legacy branch snapshot intentionally keeps Agent-Field visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/railway/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| deployments/railway/README.md | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| deployments/railway/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| docs/ARCHITECTURE.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/CONTRIBUTING.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/COVERAGE.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/COVERAGE.md | `sdk/python/agentfield*` | import-module-path | Coverage commands still point at the published Python package path in this file. |
| docs/COVERAGE.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/DEVELOPMENT.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| docs/DEVELOPMENT.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| docs/DEVELOPMENT.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/DEVELOPMENT.md | `AGENTFIELD*` | env-var | Development setup and migration commands in this file must keep the existing env var namespace. |
| docs/DEVELOPMENT.md | `Agent-Field/agentfield*` | skill-or-repo-slug | The repository slug and clone URL remain the current canonical checkout target. |
| docs/DEVELOPMENT.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/ENVIRONMENT_VARIABLES.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| docs/ENVIRONMENT_VARIABLES.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/ENVIRONMENT_VARIABLES.md | `/etc/agentfield/config/agentfield.yaml` | yaml-config-path | The container-mounted YAML config path stays aligned with the existing compatibility layout. |
| docs/ENVIRONMENT_VARIABLES.md | `AGENTFIELD*` | env-var | The control-plane and SDK environment variable namespace is intentionally unchanged in this reference doc. |
| docs/ENVIRONMENT_VARIABLES.md | `agentfield.ai*` | published-link-target | The hosted telemetry endpoint still uses the published legacy domain. |
| docs/ENVIRONMENT_VARIABLES.md | `config/agentfield.yaml` | yaml-config-path | The default YAML config path remains stable and is documented as a first-class compatibility surface. |
| docs/ENVIRONMENT_VARIABLES.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/ENVIRONMENT_VARIABLES.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| docs/RELEASE.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| docs/RELEASE.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/RELEASE.md | `@agentfield/sdk*` | package-name | npm install examples retain the published package name and channel tags. |
| docs/RELEASE.md | `Agent-Field/agentfield*` | skill-or-repo-slug | GitHub release workflow, issue, and raw-install URLs retain the current repo slug. |
| docs/RELEASE.md | `agentfield-*` | cli-command | Release artifacts and binary filenames remain stable for installers and release consumers. |
| docs/RELEASE.md | `agentfield.ai*` | published-link-target | Install and uninstall scripts still live at the published legacy domain. |
| docs/RELEASE.md | `agentfield/control-plane*` | docker-image-or-volume | Docker release instructions retain the published image name and tags. |
| docs/RELEASE.md | `agentfield` | package-name | Python installation examples retain the published package name. |
| docs/RELEASE.md | `~/.agentfield*` | cli-command | Install locations remain stable for the current shell installer and user workflows. |
| docs/RELEASE.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/RELEASE.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | `AGENTFIELD*` | env-var | Authorization env vars remain stable compatibility surfaces in the examples. |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | `config/agentfield.yaml` | yaml-config-path | Authorization examples intentionally use the stable YAML config path. |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/api/AGENT_NODE_LOGS.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| docs/api/AGENT_NODE_LOGS.md | `/agentfield/v1/logs` | api-path | The node log endpoint path is a runtime HTTP contract between agents and the control plane. |
| docs/api/AGENT_NODE_LOGS.md | `AGENTFIELD*` | env-var | Node log authentication and toggle env vars remain stable compatibility surfaces. |
| docs/api/AGENT_NODE_LOGS.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/design/execution-observability-rfc.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/design/harness-v2-design.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/design/harness-v2-design.md | `.agentfield_output.json` | test-fixture | The harness output filename stays stable so structured-output tooling can read the response file reliably. |
| docs/design/harness-v2-design.md | `.agentfield_schema.json` | test-fixture | The harness schema filename stays stable so large-schema runs can hand off the compatibility file predictably. |
| docs/design/harness-v2-design.md | `@agentfield/sdk*` | package-name | TypeScript examples use the published npm package name. |
| docs/design/harness-v2-design.md | `agentfield*` | import-module-path | Python import examples and code references use the published package path. |
| docs/design/harness-v2-design.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| examples/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| examples/README.md | https://agentfield.ai/docs/learn/examples | published-link-target | The examples index keeps the current published docs target until a Silmari docs URL exists in the repo. |
| examples/benchmarks/100k-scale/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/analyze.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/analyze.py | AgentField | historical-record | The benchmark loader keeps the historical `AgentField` filename prefix so checked-in result discovery still finds preserved artifacts. |
| examples/benchmarks/100k-scale/analyze.py | AgentField_Go | historical-record | All occurrences of this historical Go result key stay aligned with the checked-in `AgentField_Go.json` artifact and plotting schema in this file. |
| examples/benchmarks/100k-scale/analyze.py | AgentField_Python | historical-record | All occurrences of this historical Python result key stay aligned with the checked-in `AgentField_Python.json` artifact and plotting schema in this file. |
| examples/benchmarks/100k-scale/analyze.py | AgentField_TypeScript | historical-record | All occurrences of this historical TypeScript result key stay aligned with the checked-in `AgentField_TypeScript.json` artifact and plotting schema in this file. |
| examples/benchmarks/100k-scale/crewai-bench/benchmark.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/go.mod | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/go.mod | agentfield | go-module-path | Go module or import compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/main.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/main.go | AgentField | historical-record | The JSON `framework` field stays aligned with the checked-in benchmark artifact schema. |
| examples/benchmarks/100k-scale/go-bench/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/main.go | github.com/Agent-Field/agentfield/sdk/go/agent | import-module-path | The Go benchmark keeps the published SDK import path, which preserves both the repo slug and module name in this file. |
| examples/benchmarks/100k-scale/langchain-bench/benchmark.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/python-bench/benchmark.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/python-bench/benchmark.py | AgentField | historical-record | The JSON `framework` field stays aligned with the checked-in benchmark artifact schema. |
| examples/benchmarks/100k-scale/python-bench/benchmark.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/python-bench/benchmark.py | agentfield_server | json-field | All constructor kwargs in this benchmark keep the legacy `agentfield_server` keyword alias so the harness still exercises the Python SDK runtime compatibility path. |
| examples/benchmarks/100k-scale/python-bench/benchmark.py | from agentfield import Agent | import-module-path | The Python benchmark keeps the published SDK import path so the harness still runs against the packaged module name. |
| examples/benchmarks/100k-scale/requirements.txt | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/requirements.txt | agentfield | package-name | The benchmark environment installs the published PyPI package name. |
| examples/benchmarks/100k-scale/results/AgentField_Go.json | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/results/AgentField_Go.json | AgentField | historical-record | This checked-in benchmark result is preserved as a historical artifact. |
| examples/benchmarks/100k-scale/results/AgentField_Go.json | AgentField_Go.json | historical-record | The checked-in Go benchmark artifact keeps its legacy filename so historical comparisons still reference the preserved result file. |
| examples/benchmarks/100k-scale/results/AgentField_Python.json | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/results/AgentField_Python.json | AgentField | historical-record | This checked-in benchmark result is preserved as a historical artifact. |
| examples/benchmarks/100k-scale/results/AgentField_Python.json | AgentField_Python.json | historical-record | The checked-in Python benchmark artifact keeps its legacy filename so historical comparisons still reference the preserved result file. |
| examples/benchmarks/100k-scale/results/AgentField_TypeScript.json | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/results/AgentField_TypeScript.json | AgentField | historical-record | This checked-in benchmark result is preserved as a historical artifact. |
| examples/benchmarks/100k-scale/results/AgentField_TypeScript.json | AgentField_TypeScript.json | historical-record | The checked-in TypeScript benchmark artifact keeps its legacy filename so historical comparisons still reference the preserved result file. |
| examples/benchmarks/100k-scale/run_benchmarks.sh | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/run_benchmarks.sh | AgentField_Go.json | historical-record | Both shell references keep the historical Go benchmark JSON filename consumed by the checked-in results and visualization flow. |
| examples/benchmarks/100k-scale/run_benchmarks.sh | AgentField_Python.json | historical-record | Both shell references keep the historical Python benchmark JSON filename consumed by the checked-in results and visualization flow. |
| examples/benchmarks/100k-scale/ts-bench/benchmark.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/ts-bench/benchmark.ts | AgentField | historical-record | The JSON `framework` field stays aligned with the checked-in benchmark artifact schema. |
| examples/benchmarks/100k-scale/ts-bench/package-lock.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/ts-bench/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_EXECUTION_MAX_RETRIES | env-var | The configuration table keeps the execution retry-count environment variable expected by the control plane. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_EXECUTION_RETRY_BACKOFF | env-var | The configuration table keeps the execution retry-backoff environment variable expected by the control plane. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_LLM_HEALTH_CHECK_INTERVAL | env-var | Manual resilience steps keep the health-check interval environment variable expected by the control plane. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_LLM_HEALTH_CHECK_TIMEOUT | env-var | The configuration table keeps the health-check timeout environment variable expected by the control plane. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_LLM_HEALTH_ENABLED | env-var | Manual resilience steps keep the health-monitoring feature flag name that the control plane reads today. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_LLM_HEALTH_ENDPOINT | env-var | Manual resilience steps keep the probe URL environment variable expected by the control plane health checker. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_LLM_HEALTH_ENDPOINT_NAME | env-var | Manual resilience steps keep the endpoint label environment variable expected by the control plane health checker. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_LLM_HEALTH_FAILURE_THRESHOLD | env-var | Manual resilience steps and the configuration table keep the failure-threshold environment variable expected by the control plane. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_LLM_HEALTH_RECOVERY_TIMEOUT | env-var | Manual resilience steps and the configuration table keep the recovery-timeout environment variable expected by the control plane. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD_MAX_CONCURRENT_PER_AGENT | env-var | The concurrency example keeps the per-agent concurrency environment variable expected by the control plane. |
| examples/e2e_resilience_tests/README.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| examples/e2e_resilience_tests/README.md | Agent-Field/agentfield | skill-or-repo-slug | The issue links keep the current GitHub org and repository slug because the repo has not moved to a Silmari namespace. |
| examples/e2e_resilience_tests/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/e2e_resilience_tests/README.md | agentfield-server | cli-command | Manual resilience commands keep the existing control-plane binary name that local operators run today. |
| examples/e2e_resilience_tests/agent_flaky.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_flaky.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_healthy.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_healthy.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_slow.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_slow.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/e2e_resilience_tests/run_tests.sh | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/run_tests.sh | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/multi_version/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/multi_version/main.go | AGENTFIELD_TOKEN | env-var | The multi-version Go example keeps the bearer-token environment variable expected by the published SDK config. |
| examples/go_agent_nodes/cmd/multi_version/main.go | AGENTFIELD_URL | env-var | The multi-version Go example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/go_agent_nodes/cmd/multi_version/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/multi_version/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/multi_version/main.go | github.com/Agent-Field/agentfield/sdk/go/agent | import-module-path | The multi-version Go example keeps the published Go SDK module path used by the sample agent. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | AGENTFIELD_INTERNAL_TOKEN | env-var | The permission-agent Go example keeps the internal authorization token environment variable expected by the published SDK config. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | AGENTFIELD_TOKEN | env-var | The permission-agent Go example keeps the bearer-token environment variable expected by the published SDK config. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | AGENTFIELD_URL | env-var | The permission-agent Go example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | github.com/Agent-Field/agentfield/sdk/go/agent | import-module-path | The permission-agent Go example keeps the published Go SDK module path used by the sample agent. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | AGENTFIELD_INTERNAL_TOKEN | env-var | The protected-agent Go example keeps the internal authorization token environment variable expected by the published SDK config. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | AGENTFIELD_TOKEN | env-var | The protected-agent Go example keeps the bearer-token environment variable expected by the published SDK config. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | AGENTFIELD_URL | env-var | The protected-agent Go example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | github.com/Agent-Field/agentfield/sdk/go/agent | import-module-path | The protected-agent Go example keeps the published Go SDK module path used by the sample agent. |
| examples/go_agent_nodes/cmd/serverless/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/serverless/main.go | AGENTFIELD_TOKEN | env-var | The serverless Go example keeps the bearer-token environment variable expected by the published SDK config. |
| examples/go_agent_nodes/cmd/serverless/main.go | AGENTFIELD_URL | env-var | The serverless Go example keeps the compatibility environment variable for the control-plane URL. |
| examples/go_agent_nodes/cmd/serverless/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/serverless/main.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/serverless/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/serverless/main.go | github.com/Agent-Field/agentfield/sdk/go/agent | import-module-path | The serverless Go example keeps the published SDK import path. |
| examples/go_agent_nodes/go.mod | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/go.mod | agentfield | go-module-path | Go module or import compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/main.go | AGENTFIELD_AUTHORIZATION_INTERNAL_TOKEN | env-var | The Go example keeps the internal authorization token environment variable expected by the published SDK config. |
| examples/go_agent_nodes/main.go | AGENTFIELD_TOKEN | env-var | The Go example keeps the bearer-token environment variable expected by the published SDK config. |
| examples/go_agent_nodes/main.go | AGENTFIELD_URL | env-var | The Go example keeps the compatibility environment variable for server mode. |
| examples/go_agent_nodes/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/main.go | github.com/Agent-Field/agentfield/sdk/go/agent | import-module-path | The Go example imports the published Go SDK module path. |
| examples/go_harness_demo/go.mod | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_harness_demo/go.mod | agentfield | go-module-path | Go module or import compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_harness_demo/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_harness_demo/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_harness_demo/main.go | github.com/Agent-Field/agentfield/sdk/go/harness | import-module-path | The Go harness demo keeps the published Go harness module path used by the sample code. |
| examples/python_agent_nodes/agentic_rag/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/agentic_rag/README.md | AGENTFIELD_SERVER | env-var | The setup snippet keeps the compatibility environment variable for the control-plane URL used by the example runtime. |
| examples/python_agent_nodes/agentic_rag/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/agentic_rag/main.py | AGENTFIELD_SERVER | env-var | The agentic RAG example keeps the existing control-plane environment variable used by the runtime setup. |
| examples/python_agent_nodes/agentic_rag/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/agentic_rag/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/agentic_rag/main.py | agentfield_server | json-field | The agentic RAG example keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/agentic_rag/main.py | from agentfield import Agent, AIConfig | import-module-path | The agentic RAG example keeps the published Python SDK import path for all agent and AI config imports in this file. |
| examples/python_agent_nodes/agentic_rag/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/README.md | AGENTFIELD_SERVER | env-var | Setup instructions keep the control-plane environment variable used by the SDK. |
| examples/python_agent_nodes/deep_research/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/main.py | AGENTFIELD_SERVER | env-var | The deep research example keeps the existing environment variable name so its startup command remains backward compatible. |
| examples/python_agent_nodes/deep_research/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/main.py | agentfield_server | json-field | All occurrences of the legacy `agentfield_server` keyword alias stay in this example because the Python SDK still accepts that runtime argument. |
| examples/python_agent_nodes/deep_research/main.py | from agentfield import AIConfig, Agent | import-module-path | The example imports the published Python SDK path. |
| examples/python_agent_nodes/deep_research/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/routers/planning.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/routers/planning.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/routers/planning.py | from agentfield import AgentRouter | import-module-path | The planning router keeps the published Python SDK router import path used by the example agent. |
| examples/python_agent_nodes/deep_research/routers/research.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/docker_hello_world/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/docker_hello_world/main.py | AGENTFIELD_URL | env-var | The Docker hello world example keeps the existing control-plane URL environment variable for container runs. |
| examples/python_agent_nodes/docker_hello_world/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/docker_hello_world/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/docker_hello_world/main.py | agentfield_server | json-field | The Docker hello world example keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/docker_hello_world/main.py | from agentfield import Agent | import-module-path | The Docker hello world example keeps the published Python SDK import path used by the agent bootstrap. |
| examples/python_agent_nodes/documentation_chatbot/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/README.md | AGENTFIELD_SERVER | env-var | The chatbot setup keeps the compatibility environment variable for the control-plane URL. |
| examples/python_agent_nodes/documentation_chatbot/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/install.sh | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/install.sh | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/main.py | AGENTFIELD_API_KEY | env-var | The documentation chatbot keeps the compatibility API key environment variable used by the SDK client. |
| examples/python_agent_nodes/documentation_chatbot/main.py | AGENTFIELD_SERVER | env-var | The documentation chatbot keeps the compatibility environment variable for the control-plane URL. |
| examples/python_agent_nodes/documentation_chatbot/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/main.py | agentfield_server | json-field | The documentation chatbot keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/documentation_chatbot/main.py | from agentfield | import-module-path | The documentation chatbot keeps the published Python SDK module prefix for both its main agent import and logger helper import. |
| examples/python_agent_nodes/documentation_chatbot/pipeline_utils.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/product_context.py | @agentfield/sdk | package-name | Search guidance needs the published npm package name for accurate install answers. |
| examples/python_agent_nodes/documentation_chatbot/product_context.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/product_context.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/product_context.py | agentfield | package-name | Search guidance keeps the published PyPI package name in `pip install agentfield`. |
| examples/python_agent_nodes/documentation_chatbot/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/ingestion.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/ingestion.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/qa.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/qa.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/qa.py | config/agentfield.yaml | yaml-config-path | The answer template keeps the control-plane YAML path that remains a first-class compatibility surface. |
| examples/python_agent_nodes/documentation_chatbot/routers/qa.py | from agentfield | import-module-path | The QA router keeps the published Python SDK module prefix for both its router import and logger import in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/qa.py | npm install -g agentfield | cli-command | The QA example keeps the installed CLI command users run today. |
| examples/python_agent_nodes/documentation_chatbot/routers/query_planning.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/query_planning.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/query_planning.py | from agentfield import AgentRouter | import-module-path | The query planning router keeps the published Python SDK router import path used by the documentation chatbot example. |
| examples/python_agent_nodes/documentation_chatbot/routers/retrieval.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world/main.py | AGENTFIELD_URL | env-var | The example keeps the existing control-plane environment variable. |
| examples/python_agent_nodes/hello_world/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world/main.py | agentfield_server | json-field | The hello world example keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/hello_world/main.py | from agentfield | import-module-path | The hello world example keeps the published Python SDK module prefix for both of its agent imports in this file. |
| examples/python_agent_nodes/hello_world_rag/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/README.md | agentfield call | cli-command | CLI examples keep the installed command name that existing users already have on PATH. |
| examples/python_agent_nodes/hello_world_rag/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/main.py | AGENTFIELD_SERVER | env-var | The RAG example keeps the established environment variable name for the control-plane URL shown in the runnable snippet. |
| examples/python_agent_nodes/hello_world_rag/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/main.py | agentfield_server | json-field | All occurrences of the legacy `agentfield_server` keyword alias stay in this example because the Python SDK still accepts that runtime argument. |
| examples/python_agent_nodes/hello_world_rag/main.py | from agentfield | import-module-path | The RAG example keeps the published Python SDK module prefix for both its agent import and logger import in this file. |
| examples/python_agent_nodes/hello_world_rag/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/test_pydantic_skill.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/README.md | https://agentfield.ai/docs/learn | published-link-target | The README keeps the current published docs URL until the repo exposes a Silmari replacement. |
| examples/python_agent_nodes/image_generation_hello_world/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/main.py | agentfield_server | json-field | The image generation example keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/image_generation_hello_world/main.py | from agentfield | import-module-path | The image generation example keeps the published Python SDK module prefix for both `AIConfig` and `Agent` imports in this file. |
| examples/python_agent_nodes/image_generation_hello_world/requirements.txt | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/requirements.txt | agentfield | package-name | The example installs the published PyPI package name. |
| examples/python_agent_nodes/multi_version/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/multi_version/main.py | AGENTFIELD_URL | env-var | The multi-version Python example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/python_agent_nodes/multi_version/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/multi_version/main.py | agentfield_server | json-field | The multi-version Python example keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/multi_version/main.py | from agentfield import Agent | import-module-path | The multi-version Python example keeps the published Python SDK import path used by the sample agent. |
| examples/python_agent_nodes/permission_agent_a/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_a/main.py | AGENTFIELD_URL | env-var | The permission-agent Python example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/python_agent_nodes/permission_agent_a/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_a/main.py | agentfield_server | json-field | The permission-agent Python example keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/permission_agent_a/main.py | from agentfield import Agent | import-module-path | The permission-agent Python example keeps the published Python SDK import path used by the sample agent. |
| examples/python_agent_nodes/permission_agent_b/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_b/README.md | control-plane/config/agentfield.yaml | yaml-config-path | The protection rule walkthrough keeps the control-plane YAML file path that remains a first-class compatibility surface. |
| examples/python_agent_nodes/permission_agent_b/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_b/main.py | AGENTFIELD_URL | env-var | The protected-agent Python example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/python_agent_nodes/permission_agent_b/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_b/main.py | agentfield_server | json-field | The protected-agent Python example keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/permission_agent_b/main.py | from agentfield import Agent | import-module-path | The protected-agent Python example keeps the published Python SDK import path used by the sample agent. |
| examples/python_agent_nodes/rag_evaluation/Dockerfile | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/README.md | AGENTFIELD_SERVER | env-var | The RAG evaluation README keeps the compatibility environment variable shown in the configuration table. |
| examples/python_agent_nodes/rag_evaluation/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/README.md | https://agentfield.ai/docs/learn/examples | published-link-target | The RAG evaluation README keeps the current published docs target until a Silmari URL is available. |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | AGENTFIELD_HTTP_ADDR | env-var | The compose example keeps the control-plane listen-address environment variable used by the shipped local stack. |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | AGENTFIELD_SERVER | env-var | The compose example keeps the existing environment variable contract for the agent container. |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | AGENTFIELD_STORAGE_MODE | env-var | The compose example keeps the control-plane storage mode environment variable used by the shipped local stack. |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | agentfield-data | docker-image-or-volume | The compose example keeps the named volume stable so local benchmark and demo data mounts still work. |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | agentfield/control-plane:latest | docker-image-or-volume | The compose example keeps the published control-plane image name. |
| examples/python_agent_nodes/rag_evaluation/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/main.py | AGENTFIELD_SERVER | env-var | The RAG evaluation entrypoint keeps the compatibility environment variable for the control-plane URL. |
| examples/python_agent_nodes/rag_evaluation/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/main.py | agentfield_server | json-field | The RAG evaluation entrypoint keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/rag_evaluation/main.py | from agentfield import Agent, AIConfig | import-module-path | The RAG evaluation entrypoint keeps the published Python SDK import path used by the sample agent. |
| examples/python_agent_nodes/rag_evaluation/rag_eval_client.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/rag_eval_client.py | agentfield_server | json-field | The RAG evaluation client still accepts the legacy `agentfield_server` keyword alias so older sample code keeps working at runtime. |
| examples/python_agent_nodes/rag_evaluation/reasoners/__init__.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/app/page.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/app/page.tsx | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx | https://github.com/Agent-Field/agentfield | skill-or-repo-slug | The UI footer still links to the current public repository slug. |
| examples/python_agent_nodes/serverless_hello/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/serverless_hello/main.py | AGENTFIELD_URL | env-var | The serverless hello world example keeps the compatibility control-plane URL environment variable used in local smoke tests. |
| examples/python_agent_nodes/serverless_hello/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/serverless_hello/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/serverless_hello/main.py | agentfield_server | json-field | The serverless hello world example keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/serverless_hello/main.py | from agentfield | import-module-path | The serverless hello world example keeps the published Python SDK module prefix for both the main agent import and async config import. |
| examples/python_agent_nodes/simulation_engine/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/README.md | AGENTFIELD_SERVER | env-var | The simulation engine setup keeps the compatibility environment variable for the control-plane URL. |
| examples/python_agent_nodes/simulation_engine/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/README.md | pip install agentfield | package-name | The simulation engine setup keeps the published PyPI package name in the install command. |
| examples/python_agent_nodes/simulation_engine/example.ipynb | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/example.ipynb | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/main.py | AGENTFIELD_SERVER | env-var | The simulation engine keeps the compatibility environment variable for the control-plane URL. |
| examples/python_agent_nodes/simulation_engine/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/main.py | agentfield_server | json-field | The simulation engine keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/simulation_engine/main.py | from agentfield import AIConfig, Agent | import-module-path | The simulation engine keeps the published Python SDK import path for its agent and AI config imports. |
| examples/python_agent_nodes/simulation_engine/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/routers/aggregation.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/routers/decision.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/routers/entity.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/routers/scenario.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/routers/simulation.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/tool_calling/orchestrator.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/tool_calling/orchestrator.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/tool_calling/worker.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/tool_calling/worker.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/waiting_state/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/waiting_state/main.py | AGENTFIELD_URL | env-var | The waiting-state Python example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/python_agent_nodes/waiting_state/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/waiting_state/main.py | agentfield_server | json-field | The waiting-state Python example keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/python_agent_nodes/waiting_state/main.py | from agentfield import Agent, AIConfig, ApprovalResult | import-module-path | The waiting-state Python example keeps the published Python SDK import path used by the sample agent. |
| examples/tests/test_silmari_branding.py | AGENTFIELD | test-fixture | Example branding regression tests intentionally keep legacy brand tokens as fixture data for compatibility assertions; covers all occurrences in this file. |
| examples/tests/test_silmari_branding.py | AGENTFIELD | test-fixture | The shared regex fixture intentionally keeps the uppercase environment-variable form so missed compatibility identifiers still fail the branding tests. |
| examples/tests/test_silmari_branding.py | Agent-Field | test-fixture | Example branding regression tests intentionally keep legacy brand tokens as fixture data for compatibility assertions; covers all occurrences in this file. |
| examples/tests/test_silmari_branding.py | Agent-Field | test-fixture | The shared regex fixture intentionally keeps the GitHub slug variant so preserved repository-link coverage remains under test. |
| examples/tests/test_silmari_branding.py | AgentField | test-fixture | Camel-case legacy brand strings in this unittest intentionally verify that visible example copy and historical benchmark labels have been rebranded or explicitly manifested. |
| examples/tests/test_silmari_branding.py | AgentField | test-fixture | Example branding regression tests intentionally keep legacy brand tokens as fixture data for compatibility assertions; covers all occurrences in this file. |
| examples/tests/test_silmari_branding.py | AgentField_Go | test-fixture | This unittest fixture keeps the exact historical benchmark key so the manifest audit proves underscore-delimited legacy identifiers stay covered. |
| examples/tests/test_silmari_branding.py | AgentField_Go.json | test-fixture | This unittest fixture keeps the exact historical benchmark filename so the manifest audit proves underscore-delimited legacy identifiers stay covered. |
| examples/tests/test_silmari_branding.py | AgentField_Python | test-fixture | This unittest fixture keeps the exact historical benchmark key so the manifest audit proves underscore-delimited legacy identifiers stay covered. |
| examples/tests/test_silmari_branding.py | AgentField_Python.json | test-fixture | This unittest fixture keeps the exact historical benchmark filename so the manifest audit proves underscore-delimited legacy identifiers stay covered. |
| examples/tests/test_silmari_branding.py | AgentField_TypeScript | test-fixture | This unittest fixture keeps the exact historical benchmark key so the manifest audit proves underscore-delimited legacy identifiers stay covered. |
| examples/tests/test_silmari_branding.py | AgentField_TypeScript.json | test-fixture | This unittest fixture keeps the exact historical benchmark filename so the manifest audit proves underscore-delimited legacy identifiers stay covered. |
| examples/tests/test_silmari_branding.py | AgentPlane | test-fixture | Example branding regression tests intentionally keep legacy AgentPlane fixture strings for scanner coverage; covers all occurrences in this file. |
| examples/tests/test_silmari_branding.py | AgentPlane | test-fixture | The shared regex fixture intentionally keeps the older alias so the branding checks still catch that legacy naming if it resurfaces. |
| examples/tests/test_silmari_branding.py | Agent Plane | test-fixture | The shared regex fixture intentionally keeps the spaced legacy alias so the branding checks still catch that wording if it resurfaces. |
| examples/tests/test_silmari_branding.py | agentfield | test-fixture | Example branding regression tests intentionally keep legacy brand tokens as fixture data for compatibility assertions; covers all occurrences in this file. |
| examples/tests/test_silmari_branding.py | agentfield | test-fixture | Lowercase legacy tokens in regexes, helper edge cases, and asset-prefix assertions intentionally exercise the branding scanner and compatibility coverage checks. |
| examples/tests/test_silmari_branding.py | agentfield.ai | test-fixture | Example branding regression tests intentionally keep the legacy published-domain token as fixture data for compatibility assertions; covers all occurrences in this file. |
| examples/tests/test_silmari_branding.py | agentfield.ai | test-fixture | The shared regex fixture intentionally keeps the published legacy domain so preservation rules for public documentation links stay covered by the tests. |
| examples/tests/test_silmari_branding.py | agentfield_server | test-fixture | This unittest fixture keeps the exact legacy keyword so the manifest audit proves underscore-delimited compatibility identifiers stay covered. |
| examples/tests/test_silmari_branding.py | agentplane | test-fixture | Example branding regression tests intentionally keep legacy agentplane fixture strings for scanner coverage; covers all occurrences in this file. |
| examples/tests/test_silmari_branding.py | agentplane | test-fixture | The lowercase legacy alias remains in the regex fixture so the scanner tests continue to catch both case variants. |
| examples/tests/test_silmari_branding.py | agent plane | test-fixture | The shared regex fixture intentionally keeps the lowercase spaced legacy alias so the branding checks still catch that wording if it resurfaces. |
| examples/triggers-demo/Dockerfile | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/triggers-demo/Dockerfile | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| examples/triggers-demo/Dockerfile | agentfield/triggers-demo-agent | docker-image-or-volume | The example build tag stays aligned with the existing demo image naming used in local instructions. |
| examples/triggers-demo/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/triggers-demo/agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/triggers-demo/agent.py | AGENTFIELD_URL | env-var | The demo keeps the existing control-plane environment variable expected by the example runtime. |
| examples/triggers-demo/agent.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/triggers-demo/agent.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/triggers-demo/agent.py | agentfield_server | json-field | The triggers demo keeps the legacy `agentfield_server` keyword alias that the Python SDK still accepts at runtime. |
| examples/triggers-demo/agent.py | from agentfield import ( | import-module-path | The trigger demo imports the published Python SDK module path. |
| examples/triggers-demo/docker-compose.yml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/triggers-demo/docker-compose.yml | AGENTFIELD_FEATURES_DID_ENABLED | env-var | The compose file keeps the DID feature-flag environment variable used by the shipped demo container. |
| examples/triggers-demo/docker-compose.yml | AGENTFIELD_HOME | env-var | The compose file must keep the control-plane environment variable names used by the shipped demo container. |
| examples/triggers-demo/docker-compose.yml | AGENTFIELD_HTTP_ADDR | env-var | The compose file keeps the control-plane listen-address environment variable used by the shipped demo container. |
| examples/triggers-demo/docker-compose.yml | AGENTFIELD_MODE | env-var | The compose file keeps the control-plane mode environment variable used by the shipped demo container. |
| examples/triggers-demo/docker-compose.yml | AGENTFIELD_STORAGE_MODE | env-var | The compose file keeps the control-plane storage mode environment variable used by the shipped demo container. |
| examples/triggers-demo/docker-compose.yml | AGENTFIELD_URL | env-var | The demo agent service keeps the compatibility control-plane URL environment variable passed by compose. |
| examples/triggers-demo/docker-compose.yml | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/triggers-demo/docker-compose.yml | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| examples/triggers-demo/docker-compose.yml | agentfield-data | docker-image-or-volume | The named volume remains stable so existing local demo data mounts continue to work. |
| examples/triggers-demo/scripts/fire-events.sh | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/triggers-demo/scripts/fire-events.sh | AGENTFIELD_URL | env-var | The helper script keeps the operator override for the control-plane endpoint. |
| examples/triggers-demo/scripts/fire-events.sh | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/discovery-memory/main.ts | @agentfield/sdk | import-module-path | The discovery-memory TypeScript example keeps the published npm import path used by the sample agent and router. |
| examples/ts-node-examples/discovery-memory/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/discovery-memory/main.ts | AGENTFIELD_URL | env-var | The discovery-memory TypeScript example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/ts-node-examples/discovery-memory/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/init-example/main.ts | @agentfield/sdk | import-module-path | The init-example TypeScript sample keeps the published npm import path used by the sample agent. |
| examples/ts-node-examples/init-example/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/init-example/main.ts | AGENTFIELD_API_KEY | env-var | The init-example TypeScript sample keeps the compatibility API key environment variable used by the SDK client. |
| examples/ts-node-examples/init-example/main.ts | AGENTFIELD_URL | env-var | The init-example TypeScript sample keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/ts-node-examples/init-example/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/init-example/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/ts-node-examples/init-example/reasoners.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/multi-version/main.ts | @agentfield/sdk | import-module-path | The multi-version TypeScript example keeps the published npm import path used by the sample agent. |
| examples/ts-node-examples/multi-version/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/multi-version/main.ts | AGENTFIELD_URL | env-var | The multi-version TypeScript example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/ts-node-examples/multi-version/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/package-lock.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/ts-node-examples/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/ts-node-examples/permission-agent-a/main.ts | @agentfield/sdk | import-module-path | The permission-agent TypeScript example keeps the published npm import path used by the sample agent. |
| examples/ts-node-examples/permission-agent-a/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/permission-agent-a/main.ts | AGENTFIELD_URL | env-var | The permission-agent TypeScript example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/ts-node-examples/permission-agent-a/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/permission-agent-b/main.ts | @agentfield/sdk | import-module-path | The protected-agent TypeScript example keeps the published npm import path used by the sample agent. |
| examples/ts-node-examples/permission-agent-b/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/permission-agent-b/main.ts | AGENTFIELD_URL | env-var | The protected-agent TypeScript example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/ts-node-examples/permission-agent-b/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/serverless-hello/main.ts | @agentfield/sdk | import-module-path | The TypeScript example keeps the published npm import path. |
| examples/ts-node-examples/serverless-hello/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/serverless-hello/main.ts | AGENTFIELD_URL | env-var | The TypeScript serverless example keeps the compatibility environment variable for the control-plane URL. |
| examples/ts-node-examples/serverless-hello/main.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/serverless-hello/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/main.ts | @agentfield/sdk | import-module-path | The simulation TypeScript example keeps the published npm import path used by the sample agent. |
| examples/ts-node-examples/simulation/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/aggregation.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/decision.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/entity.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/scenario.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/simulation.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/README.md | agentfield.yaml | yaml-config-path | The README keeps the compatibility YAML filename shown in the DID setup snippet for existing control-plane configuration flows. |
| examples/ts-node-examples/verifiable-credentials/README.md | agentfield/examples/ts-node-examples | skill-or-repo-slug | The quick-start command keeps the current repository checkout path so the example still matches the `agentfield` repo slug and folder layout. |
| examples/ts-node-examples/verifiable-credentials/README.md | https://agentfield.ai/credentials/v1 | published-link-target | The sample VC context URL remains the published context endpoint until a Silmari replacement exists. |
| examples/ts-node-examples/verifiable-credentials/main.ts | @agentfield/sdk | import-module-path | The VC example keeps the published TypeScript SDK import path. |
| examples/ts-node-examples/verifiable-credentials/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/main.ts | AGENTFIELD_URL | env-var | The VC example keeps the control-plane URL environment variable used by the SDK. |
| examples/ts-node-examples/verifiable-credentials/main.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/reasoners.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/waiting-state/main.ts | @agentfield/sdk | import-module-path | The waiting-state TypeScript example keeps the published npm import path used by the sample agent and approval client. |
| examples/ts-node-examples/waiting-state/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/waiting-state/main.ts | AGENTFIELD_API_KEY | env-var | The waiting-state TypeScript example keeps the compatibility API key environment variable used by the SDK clients. |
| examples/ts-node-examples/waiting-state/main.ts | AGENTFIELD_URL | env-var | The waiting-state TypeScript example keeps the compatibility control-plane URL environment variable used by the sample agent. |
| examples/ts-node-examples/waiting-state/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/waiting-state/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/ts_agent_nodes/tool_calling/orchestrator.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts_agent_nodes/tool_calling/orchestrator.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts_agent_nodes/tool_calling/package-lock.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/ts_agent_nodes/tool_calling/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/ts_agent_nodes/tool_calling/test_orchestrator.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts_agent_nodes/tool_calling/test_orchestrator.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts_agent_nodes/tool_calling/worker.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts_agent_nodes/tool_calling/worker.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/README.md | Agent-Field | import-module-path | Import or module path compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/go/README.md | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/README.md | github.com/Agent-Field/agentfield/sdk/go | go-module-path | The published Go module path stays stable across installation and import examples in this README. |
| sdk/go/agent/agent.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/go/agent/agent.go | Agent-Field | import-module-path | Import or module path compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/agent.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/go/agent/agent.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/agent_accepts_webhook_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/agent_accepts_webhook_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/agent_did.go | Agent-Field | import-module-path | Import or module path compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/agent_did.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/agent_did_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/agent_did_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/agent_lifecycle.go | Agent-Field | import-module-path | Import or module path compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/agent_lifecycle.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/go/agent/agent_lifecycle.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/agent_lifecycle_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/agent_lifecycle_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/agent_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/agent_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/agent_trigger_origin_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/agent_trigger_origin_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/cli.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/go/agent/cli_test.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/go/agent/cli_test.go | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/did_async_additional_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/did_async_additional_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/discovery.go | Agent-Field | import-module-path | Import or module path compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/discovery.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/harness.go | Agent-Field | import-module-path | Import or module path compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/harness.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/agent/harness_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/harness_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/middleware_additional_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/note.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/go/agent/process_logs.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/go/agent/process_logs_test.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/go/agent/process_logs_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/registration_integration_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/registration_integration_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/router_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/router_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/workflow_event_additional_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/agent/workflow_event_additional_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/ai/README.md | Agent-Field | import-module-path | Import or module path compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/ai/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/go/ai/README.md | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/ai/tool_calling.go | Agent-Field | import-module-path | Import or module path compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/ai/tool_calling.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/ai/tool_calling_additional_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/ai/tool_calling_additional_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/ai/tool_calling_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/ai/tool_calling_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/client/client.go | Agent-Field | import-module-path | Import or module path compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/client/client.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/go/client/client.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/client/client_invariant_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/client/client_low_coverage_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/client/client_low_coverage_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/client/client_redirect_test.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/go/client/client_test.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/client/client_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/did/types.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/go/go.mod | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| sdk/go/go.mod | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/go.mod | module github.com/Agent-Field/agentfield/sdk/go | go-module-path | The authoritative Go module declaration remains unchanged so existing consumers keep the published module path. |
| sdk/go/harness/codex_integration_test.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/go/harness/coverage_branches_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/harness/gemini_integration_test.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/go/harness/opencode.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/harness/runner.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/types/types.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/MANIFEST.in | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| sdk/python/README.md | ./scripts/run_pytest.sh --cov=agentfield --cov-report=term-missing | package-name | The README coverage command keeps the published Python package name in the pytest targets. |
| sdk/python/README.md | Agent-Field | historical-record | Legacy branch snapshot intentionally keeps Agent-Field visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/README.md | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/README.md | from agentfield import Agent | import-module-path | The published Python import root stays stable across the README code samples. |
| sdk/python/README.md | https://github.com/Agent-Field/agentfield.git | skill-or-repo-slug | The public repository slug remains stable for local SDK development instructions. |
| sdk/python/README.md | pip install agentfield | package-name | The published PyPI package name remains stable for SDK installation instructions. |
| sdk/python/agentfield.egg-info/PKG-INFO | Agent-Field | historical-record | Legacy branch snapshot intentionally keeps Agent-Field visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield.egg-info/PKG-INFO | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield.egg-info/PKG-INFO | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield.egg-info/SOURCES.txt | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield.egg-info/top_level.txt | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/agentfield/agent.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/agent.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent_ai.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/agent_ai.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent_cli.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/agent_cli.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent_discovery.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/agent_discovery.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent_field_handler.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/agent_field_handler.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent_pause.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent_server.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/agentfield/agent_server.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/agent_server.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent_serverless.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/agent_serverless.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent_vc.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/agent_workflow.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/agent_workflow.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/async_config.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/agentfield/async_config.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/async_execution_manager.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/cancel.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/cancel.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/client.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/client.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/connection_manager.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/connection_manager.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/decorators.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/decorators.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/did_auth.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/did_manager.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/exceptions.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/execution_context.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/fixtures/triggers/github.json | Agent-Field | historical-record | Legacy branch snapshot intentionally keeps Agent-Field visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/fixtures/triggers/github.json | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/harness/__init__.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/harness/_runner.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/harness/providers/__init__.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/harness/providers/_base.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/harness/providers/_factory.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/harness/providers/claude.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/harness/providers/codex.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/harness/providers/gemini.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/harness/providers/opencode.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/agentfield/harness/providers/opencode.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/http_connection_manager.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/logger.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/agentfield/logger.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/logger.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/media_providers.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/media_providers.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/media_router.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/memory.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/memory.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/memory_events.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/multimodal_response.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/node_logs.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/agentfield/node_logs.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/pydantic_utils.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/pydantic_utils.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/rate_limiter.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/status.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/testing.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/tool_calling.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/tool_calling.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/triggers.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/triggers.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/agentfield/types.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/vc_generator.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/verification.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/agentfield/vision.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/agentfield/vision.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/python/pyproject.toml | Agent-Field | historical-record | Legacy branch snapshot intentionally keeps Agent-Field visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/pyproject.toml | Agent-Field | skill-or-repo-slug | The GitHub organization slug stays stable in project URLs and covers all occurrences in this file. |
| sdk/python/pyproject.toml | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/pyproject.toml | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| sdk/python/pyproject.toml | agentfield | package-name | The published PyPI package name, discovery keyword, setuptools package roots, and coverage module targets stay legacy-compatible and covers all occurrences in this file. |
| sdk/python/requirements-dev.txt | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/tests/conftest.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/conftest.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/debug_complex_json.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/helpers.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/helpers.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/integration/conftest.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/tests/integration/conftest.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/integration/conftest.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/integration/test_agentfield_end_to_end.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/integration/test_agentfield_end_to_end.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_accepts_webhook.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_ai.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_ai_comprehensive.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_ai_coverage_additions.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_ai_deadlock_recovery.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_ai_final90.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_bigfiles_final90.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_call.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_cli.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_core.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_coverage_additions.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_field_handler.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_graceful_shutdown.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_helpers.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_instance_id.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_integration.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/tests/test_agent_integration.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_integration.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_lifecycle_invariants.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_networking.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_registry.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_resilience.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_server.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/tests/test_agent_server.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_server_extended.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_serverless.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_utils.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_workflow.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_workflow_extended.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_agent_workflow_registration.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_ai_config.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_approval.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_async_config.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/tests/test_async_config.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_async_execution.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_async_execution_manager_comprehensive.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_async_execution_manager_final90.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_async_execution_manager_paths.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_cancel.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_cancel.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_client.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_client_auth.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_client_bigfiles_coverage.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_client_coverage_additions.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_client_execution_paths.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_client_execution_vc_payload.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_client_laser_push.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_client_lifecycle.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_client_unit.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_connection_manager.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_connection_manager_invariants.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_cost_tracker.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_decorator_code_origin.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_decorators.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_did_auth.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_did_auth_invariants.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_did_manager.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_did_manager_error_paths.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_exceptions.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_exceptions.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_execution_context_core.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_execution_context_coverage_additions.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_execution_context_parent_vc.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_execution_logger.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/tests/test_execution_logger.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_execution_state.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_execution_state_invariants.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_agent_wiring.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_ai_schema_repair.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_cli.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/tests/test_harness_cli.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_cost_estimation.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_factory.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_functional.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_provider_claude.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_provider_codex.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_provider_gemini.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_provider_opencode.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_runner.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_schema.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_harness_types.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_http_connection_manager.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_image_config.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_invariants.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_invariants.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_litellm_adapters.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_media_integration.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_media_providers.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_media_providers_additional.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_media_router.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_memory_bigfiles_coverage.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_memory_client_core.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_memory_coverage_additions.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_memory_events.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_memory_events_additional.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_memory_flow_core.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_memory_invariants.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_memory_performance.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_memory_performance.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_multimodal.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_multimodal_response.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_multimodal_response_additional.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_multimodal_response_cost.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_node_logs.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/tests/test_node_logs.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_openrouter_audio.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_openrouter_video.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_pydantic_utils.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_rate_limiter_core.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_reasoner_path_normalization.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_result_cache.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_result_cache_bigfiles_coverage.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_router.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_run_cli_env.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/python/tests/test_run_cli_env.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_sdk_metadata_branding.py | @agentfield/sdk | test-fixture | SDK metadata regression tests intentionally keep @agentfield/sdk fixture strings to assert published package compatibility; covers all occurrences in this file. |
| sdk/python/tests/test_sdk_metadata_branding.py | @agentfield/sdk | test-fixture | The regression suite intentionally keeps the published npm package name stable across install and import fixture assertions in this file and covers all occurrences in this file. |
| sdk/python/tests/test_sdk_metadata_branding.py | AgentField | test-fixture | SDK metadata regression tests intentionally keep AgentField fixture strings to assert the rebrand and compatibility boundaries; covers all occurrences in this file. |
| sdk/python/tests/test_sdk_metadata_branding.py | AgentField | test-fixture | The regression suite keeps the legacy brand token in a negative assertion so README prose regressions still fail when `AgentField` reappears. |
| sdk/python/tests/test_sdk_metadata_branding.py | agentfield | test-fixture | SDK metadata regression tests intentionally keep agentfield fixture strings to assert the rebrand and compatibility boundaries; covers all occurrences in this file. |
| sdk/python/tests/test_sdk_metadata_branding.py | agentfield | test-fixture | The regression suite intentionally keeps the published Python package name, discovery keyword, and import-root fixtures stable and covers all bare `agentfield` occurrences in this file. |
| sdk/python/tests/test_sdk_metadata_branding.py | github.com/Agent-Field/agentfield/sdk/go | test-fixture | SDK metadata regression tests intentionally keep the published Go SDK module path fixture stable; covers all occurrences in this file. |
| sdk/python/tests/test_sdk_metadata_branding.py | github.com/Agent-Field/agentfield/sdk/go | test-fixture | The regression suite intentionally keeps the published Go module path stable across install and import fixture assertions in this file and covers all Go SDK module-path occurrences in this file. |
| sdk/python/tests/test_sdk_metadata_branding.py | https://github.com/Agent-Field/agentfield | test-fixture | SDK metadata regression tests intentionally keep the published GitHub repository URL fixture stable; covers all occurrences in this file. |
| sdk/python/tests/test_sdk_metadata_branding.py | https://github.com/Agent-Field/agentfield | test-fixture | The regression suite intentionally keeps the published GitHub slug stable across homepage, docs, repository, and issues URL fixtures in this file and covers all `https://github.com/Agent-Field/agentfield...` occurrences in this file. |
| sdk/python/tests/test_simulate_trigger.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_skill_pydantic_models.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_status_utils.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_tool_calling.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_tool_calling_error_paths.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_trigger_context.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_trigger_param_binding.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_types.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_utils.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_vc_generator.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_vc_generator_error_paths.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_verification.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_video_output.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_vision.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/tests/test_workflow_parent_child.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/python/uv.lock | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| sdk/typescript/README.md | @agentfield/sdk | package-name | The published npm package name remains stable across installation and import examples in this README. |
| sdk/typescript/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/typescript/README.md | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/typescript/package-lock.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| sdk/typescript/package.json | "agentfield" | package-name | The legacy discovery keyword remains in package metadata so npm users can still find the stable package name. |
| sdk/typescript/package.json | @agentfield/sdk | package-name | The published npm package name remains stable in package metadata. |
| sdk/typescript/package.json | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| sdk/typescript/package.json | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/typescript/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| sdk/typescript/package.json | https://github.com/Agent-Field/agentfield | skill-or-repo-slug | The repository, homepage, and issues URLs keep the stable GitHub slug for published package metadata. |
| sdk/typescript/scripts/did-smoke.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/typescript/src/agent/Agent.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/typescript/src/agent/Agent.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/typescript/src/agent/cancel.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/typescript/src/agent/processLogs.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/typescript/src/agent/processLogs.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/typescript/src/ai/ToolCalling.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/typescript/src/approval/ApprovalClient.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/typescript/src/harness/providers/opencode.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/typescript/src/status/ExecutionStatus.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/typescript/src/verification/LocalVerifier.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/typescript/tests/ai_multimodal_response_extra.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/cancel.test.ts | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/cancel.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/did_client_methods.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/harness_functional.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/harness_runner.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/harness_schema.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/integration.e2e.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/memory_performance.test.ts | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/multimodal.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/typescript/tests/process_logs.test.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/typescript/tests/process_logs.test.ts | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| skills/agentfield/SKILL.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/SKILL.md | `agentfield` skill | skill-or-repo-slug | User instructions must keep the existing skill name; this row covers all skill-name mentions in the file. |
| skills/agentfield/SKILL.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/SKILL.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/SKILL.md | agentfield.ai | published-link-target | The live SDK docs still publish at `agentfield.ai`; this row covers every docs host reference in the file. |
| skills/agentfield/SKILL.md | aliases: [agentfield-multi-reasoner-builder] | skill-or-repo-slug | The published alias remains stable for existing skill references and install surfaces. |
| skills/agentfield/SKILL.md | name: agentfield | skill-or-repo-slug | The shipped skill frontmatter slug must stay `agentfield` for catalog lookup and installer compatibility. |
| skills/agentfield/commands/agentfield.md | /agentfield | cli-command | The slash command is derived from `commands/agentfield.md` and must remain `/agentfield` for existing agent integrations. |
| skills/agentfield/commands/agentfield.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/commands/agentfield.md | `agentfield` skill | skill-or-repo-slug | Command help still refers to the stable skill slug that users invoke today. |
| skills/agentfield/commands/agentfield.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/commands/agentfield.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/commands/agentfield.md | https://agentfield.ai/llms.txt | published-link-target | The hard-gate live-docs endpoint has no verified Silmari replacement URL. |
| skills/agentfield/references/anti-patterns.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/capability-playbook.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/cli-toolkit.md | AGENTFIELD_HOME | env-var | Skill CLI docs keep AGENTFIELD_HOME stable as the compatibility env var for the canonical skill store path; covers all occurrences in this file. |
| skills/agentfield/references/cli-toolkit.md | AGENTFIELD_HOME | env-var | The canonical skill-store path is described through the existing `AGENTFIELD_HOME` override surface. |
| skills/agentfield/references/cli-toolkit.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/examples-map.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/live-docs.md | AGENTFIELD_HOME | env-var | Cache-path guidance now routes through `AGENTFIELD_HOME`, which remains the backward-compatible env var surface. |
| skills/agentfield/references/live-docs.md | AGENTFIELD_HOME | env-var | Skill live-doc guidance keeps AGENTFIELD_HOME stable as the compatibility cache-directory env var; covers all occurrences in this file. |
| skills/agentfield/references/live-docs.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/live-docs.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/live-docs.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/references/live-docs.md | agentfield.ai | published-link-target | The published docs host remains `agentfield.ai`; this row covers all host and URL references in the file. |
| skills/agentfield/references/memory-events.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/model-selection.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/patterns-emerge.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/primitives-snapshot.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| skills/agentfield/references/primitives-snapshot.md | AGENTFIELD_SERVER | env-var | The runtime env var example keeps the current server address knob unchanged. |
| skills/agentfield/references/primitives-snapshot.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/primitives-snapshot.md | agentfield.ai | published-link-target | Offline fallback guidance still points at the published `agentfield.ai` docs corpus. |
| skills/agentfield/references/primitives-snapshot.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/references/primitives-snapshot.md | sdk/python/agentfield/agent.py | import-module-path | The offline snapshot points to the existing Python package path for direct source inspection. |
| skills/agentfield/references/project-claude-template.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/project-claude-template.md | `agentfield` skill | skill-or-repo-slug | Generated project handoff text must keep the skill slug users rerun later. |
| skills/agentfield/references/project-claude-template.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/project-claude-template.md | agentfield/control-plane:latest | docker-image-or-volume | Generated runtime instructions must keep the published control-plane image name. |
| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD_HTTP_ADDR | env-var | Docker compose examples keep the current HTTP bind-address env var name. |
| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD_HTTP_PORT | env-var | Port override examples keep the current HTTP port env var name. |
| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD_SERVER | env-var | Agent container examples keep the current control-plane URL env var name. |
| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD_STORAGE_MODE | env-var | Docker compose examples keep the current storage-mode env var name. |
| skills/agentfield/references/scaffold-recipe.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/scaffold-recipe.md | agentfield | package-name | Requirements guidance keeps the installed Python package name unchanged; this row covers package-name mentions in the file. |
| skills/agentfield/references/scaffold-recipe.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/scaffold-recipe.md | agentfield-data | docker-image-or-volume | The sample compose volume name stays unchanged for backward-compatible local data mounts. |
| skills/agentfield/references/scaffold-recipe.md | agentfield/control-plane:latest | docker-image-or-volume | Compose examples keep the published control-plane image tag unchanged. |
| skills/agentfield/references/scaffold-recipe.md | from agentfield import | import-module-path | Python scaffold snippets retain the installed package import path; this row covers all import examples in the file. |
| skills/agentfield/references/triggers.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/triggers.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/references/triggers.md | agentfield.ai/docs | published-link-target | Trigger guidance points at the currently published docs site because no verified Silmari docs host exists yet. |
| skills/agentfield/references/triggers.md | from agentfield import | import-module-path | Trigger examples retain the installed Python import path; this row covers all import examples in the file. |
| skills/agentfield/references/verification.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/README.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| specs/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/README.md | `agentfield.*` | import-module-path | The generic Python import-path notation stays aligned with the published package namespace. |
| specs/README.md | `cmd/agentfield-server*` | cli-command | The standalone server entrypoint path remains tied to the current binary name. |
| specs/README.md | `control-plane/cmd/agentfield*` | cli-command | The CLI entrypoint path remains tied to the current command layout. |
| specs/README.md | `github.com/Agent-Field/agentfield/control-plane` | go-module-path | Go code references retain the repository module path. |
| specs/README.md | `sdk/python/agentfield*` | import-module-path | Python SDK code references use the published package path. |
| specs/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/agentplane-ui-api-spec.md | AgentPlane | historical-record | Legacy branch snapshot intentionally keeps AgentPlane visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/agentplane-ui-api-spec.md | `/api/v1/did/agentfield-server` | api-path | This DID endpoint path is an existing compatibility surface for the UI API spec. |
| specs/agentplane-ui-api-spec.md | `agentplane-ui-api-worklist.md` | historical-record | The linked planning note keeps its historical filename even though the visible spec branding was updated. |
| specs/agentplane-ui-api-spec.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/agentplane-ui-api-spec.md | agentplane | historical-record | Legacy branch snapshot intentionally keeps agentplane visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/architecture-overview.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| specs/architecture-overview.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/architecture-overview.md | `AGENTFIELD*` | env-var | Configuration precedence examples intentionally keep the stable env var namespace. |
| specs/architecture-overview.md | `cmd/agentfield-server*` | cli-command | Server entrypoint references retain the current binary name and source path. |
| specs/architecture-overview.md | `config/agentfield.yaml` | yaml-config-path | The default YAML config path remains a first-class compatibility surface in this architecture doc. |
| specs/architecture-overview.md | `sdk/python/agentfield*` | import-module-path | Python SDK path references retain the published package layout. |
| specs/architecture-overview.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/control-plane.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| specs/control-plane.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/control-plane.md | `AGENTFIELD*` | env-var | Control-plane configuration examples retain the stable env var namespace. |
| specs/control-plane.md | `agentfield-server` | cli-command | The standalone server binary name remains stable in this control-plane spec. |
| specs/control-plane.md | `cmd/agentfield-server*` | cli-command | Source entrypoint references retain the current server command path. |
| specs/control-plane.md | `config/agentfield.yaml` | yaml-config-path | The default YAML config path remains stable and is documented alongside env vars. |
| specs/control-plane.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/data-flow.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/data-flow.md | `sdk/python/agentfield*` | import-module-path | Python client and memory code references retain the published package path. |
| specs/data-flow.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/deployment.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| specs/deployment.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/deployment.md | `AGENTFIELD*` | env-var | Deployment docs intentionally keep the stable env var namespace. |
| specs/deployment.md | `agentfield-server` | cli-command | Deployment examples retain the standalone server binary name. |
| specs/deployment.md | `config/agentfield.yaml` | yaml-config-path | The default YAML config path remains stable in deployment guidance. |
| specs/deployment.md | `sdk/python/agentfield*` | import-module-path | The serverless adapter reference keeps the published Python package path. |
| specs/deployment.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/sdk-go.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| specs/sdk-go.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/sdk-go.md | `github.com/Agent-Field/agentfield/sdk/go` | go-module-path | The Go SDK module path is a published compatibility surface. |
| specs/sdk-go.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/sdk-python.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/sdk-python.md | `AgentFieldClient` | import-module-path | The public Python SDK client class name remains stable so the documented import surface matches the shipped package. |
| specs/sdk-python.md | `AgentFieldHandler` | import-module-path | The public Python SDK handler class name remains stable so the documented import surface matches the shipped package. |
| specs/sdk-python.md | `agentfield*` | import-module-path | Python import examples and coverage commands intentionally use the published package namespace. |
| specs/sdk-python.md | `sdk/python/agentfield*` | import-module-path | The Python SDK root and code references retain the published package path. |
| specs/sdk-python.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/security.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/security.md | `sdk/python/agentfield*` | import-module-path | Security code references retain the published Python package path. |
| specs/security.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/viewer.html | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| specs/viewer.html | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/viewer.html | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/viewer.html | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| specs/web-ui.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/web-ui.md | `cmd/agentfield-server*` | cli-command | Web UI development instructions retain the current server command path. |
| specs/web-ui.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| tests/functional/.env.example | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/.env.example | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/README.md | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/README.md | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/README.md | `AGENTFIELD*` | env-var | Functional test configuration retains the stable env var namespace. |
| tests/functional/README.md | `Agent-Field/agentfield*` | skill-or-repo-slug | The documentation link still targets the current repository slug. |
| tests/functional/README.md | `agentfield-test.yaml` | test-fixture | The functional test harness and compose files depend on the existing fixture filename. |
| tests/functional/README.md | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/agents/call_chain_agents.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/agents/call_chain_agents.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/agents/docs_quick_start_agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/agents/docs_quick_start_agent.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/agents/memory_agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/agents/memory_agent.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/agents/memory_events_agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/agents/memory_events_agent.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/agents/memory_events_decorator_agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/agents/memory_events_decorator_agent.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/agents/quick_start_agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/agents/quick_start_agent.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/agents/quick_start_agent.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/agents/router_prefix_agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/agents/router_prefix_agent.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/agents/scoping_agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/agents/scoping_agent.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/conftest.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/conftest.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/conftest.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/Dockerfile.log-demo-node | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/Dockerfile.test-runner | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/Dockerfile.test-runner | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/LOG_DEMO.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/docker/LOG_DEMO.md | `/agentfield/v1/logs` | api-path | The log proxy demo references the live agent log endpoint path. |
| tests/functional/docker/LOG_DEMO.md | `/tmp/agentfield-log-demo` | test-fixture | The host-mode demo scripts and Makefile still use this writable compatibility data path by default. |
| tests/functional/docker/LOG_DEMO.md | `AGENTFIELD*` | env-var | The host-mode log demo keeps its stable internal-auth and data-dir env var names. |
| tests/functional/docker/LOG_DEMO.md | `agentfield-server` | cli-command | The native log demo still builds the published standalone server binary. |
| tests/functional/docker/LOG_DEMO.md | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/agentfield-test.yaml | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/docker-compose.local.yml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/docker/docker-compose.local.yml | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/docker-compose.log-demo.yml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/docker/docker-compose.log-demo.yml | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/docker-compose.postgres.yml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/docker/docker-compose.postgres.yml | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/log-demo-node/log-demo.mjs | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/docker/log-demo-node/log-demo.mjs | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/docker/wait-for-services.sh | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/docker/wait-for-services.sh | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/discovery/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/discovery/main.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/discovery/main.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/hello/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/hello/main.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/hello/main.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/serverless/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/serverless/main.go | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/serverless/main.go | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/go_agents/cmd/serverless/main.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/go_agents/go.mod | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/go_agents/go.mod | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/tests/test_app_call.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/tests/test_app_call.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/tests/test_go_sdk_cli.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/tests/test_go_sdk_discovery.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/tests/test_serverless_agents.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/tests/test_serverless_agents.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/tests/test_serverless_agents.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/tests/test_ui_node_logs_proxy.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/tests/test_vc_cli.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/tests/test_vc_cli.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/tests/test_waiting_state.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/ts_agents/echo-agent.mjs | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/ts_agents/echo-agent.mjs | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/ts_agents/serverless-agent.mjs | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/ts_agents/serverless-agent.mjs | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/ts_agents/serverless-agent.mjs | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/utils/agent_server.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/utils/agent_server.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/utils/naming.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/test_check_silmari_rebrand.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/test_check_silmari_rebrand.py | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/test_check_silmari_rebrand.py | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/test_check_silmari_rebrand.py | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/test_docs_specs_branding.py | AGENTFIELD | test-fixture | Docs/specs branding regression tests intentionally keep legacy brand tokens as fixture data for compatibility assertions; covers all occurrences in this file. |
| tests/test_docs_specs_branding.py | Agent Plane | test-fixture | Docs/specs branding regression tests intentionally keep legacy Agent Plane fixture strings for scanner coverage; covers all occurrences in this file. |
| tests/test_docs_specs_branding.py | Agent-Field | test-fixture | Docs/specs branding regression tests intentionally keep legacy brand tokens as fixture data for compatibility assertions; covers all occurrences in this file. |
| tests/test_docs_specs_branding.py | AgentField | test-fixture | Docs/specs branding regression tests intentionally keep legacy brand tokens as fixture data for compatibility assertions; covers all occurrences in this file. |
| tests/test_docs_specs_branding.py | AgentPlane | test-fixture | Docs/specs branding regression tests intentionally keep legacy AgentPlane fixture strings for scanner coverage; covers all occurrences in this file. |
| tests/test_docs_specs_branding.py | agent plane | test-fixture | Docs/specs branding regression tests intentionally keep legacy agent plane fixture strings for scanner coverage; covers all occurrences in this file. |
| tests/test_docs_specs_branding.py | agentfield | test-fixture | Docs/specs branding regression tests intentionally keep legacy brand tokens as fixture data for compatibility assertions; covers all occurrences in this file. |
| tests/test_docs_specs_branding.py | agentfield.ai | test-fixture | Docs/specs branding regression tests intentionally keep the legacy published-domain token as fixture data for compatibility assertions; covers all occurrences in this file. |
| tests/test_docs_specs_branding.py | agentplane | test-fixture | Docs/specs branding regression tests intentionally keep legacy agentplane fixture strings for scanner coverage; covers all occurrences in this file. |

## CodeCleanup Passes
| Pass | Result | Notes |
|---|---|---|
| Brand surface pass | pass | README, docs/specs, templates, skills, examples, deployments, and admin UI checks now show Silmari as the visible first-party brand while manifested compatibility identifiers remain intact. |
| Compatibility pass | pass | Preserved rows cover package/import names, module paths, YAML and env configuration surfaces, API paths, published domains, and historical or test-fixture identifiers that must stay stable. |
| Mirror and generated-template pass | pass | The embedded skill mirror is byte-identical with `skills/agentfield`, and targeted `control-plane/internal/templates` plus `control-plane/internal/skillkit` Go tests passed. |
| Link and asset pass | pass | Published external links resolved under the fallback Markdown link check; the only non-zero result was the expected `http://localhost:8080` local-demo target with no server running. |
| Formatting and lint pass | pass | Web and SDK TypeScript lint checks exited cleanly; the web client reported only pre-existing warnings under the checked-in suppressions baseline. |
| Verification pass | fail | Scanner, manifest, docs/specs, examples, web, and SDK checks are recorded below; the remaining blockers are the environment-specific `control-plane` and aggregate Go test runs, which require capabilities this runner does not provide. |


## Verification Commands
| Command | Working Directory | Exit Code | Result |
|---|---|---|---|
| `lychee --version` | `.` | 127 | `lychee` is not installed in PATH on this runner, so the fallback Markdown link check was used instead. |
| `npx --yes markdown-link-check README.md docs/**/*.md specs/*.md examples/**/*.md deployments/**/*.md skills/**/*.md` | `.` | 1 | Fallback link check examined the scoped Markdown corpus and only flagged `http://localhost:8080`, which is an expected local-demo target when no local server is running. |
| `git diff --check` | `.` | 0 | No whitespace or patch-formatting issues remain in the rebrand worktree changes. |
| `./scripts/sync-embedded-skills.sh --check` | `.` | 0 | Embedded skills are in sync with the canonical `skills/agentfield` source tree. |
| `diff -qr skills/agentfield control-plane/internal/skillkit/skill_data/agentfield` | `.` | 0 | The canonical skill source and embedded mirror are byte-identical. |
| `rg -n -e 'AgentField' -e 'agentfield' -e 'AgentPlane' -e 'agentplane' -e 'Agent Plane' -e 'agent plane' docs specs control-plane/README.md control-plane/scripts/README.md control-plane/tools/perf/README.md control-plane/migrations/README.md tests/functional/README.md tests/functional/docker/LOG_DEMO.md --glob '*.md'` | `.` | 0 | Scoped Markdown review shows only manifested compatibility tokens remain, including `.agentfield_output.json`, `.agentfield_schema.json`, `AgentFieldHandler`, `AgentFieldClient`, and the historical `agentplane-ui-api-worklist.md` reference. |
| `rg -n -e '\\.agentfield_output\\.json' -e '\\.agentfield_schema\\.json' -e 'AgentFieldHandler' -e 'AgentFieldClient' docs/design/harness-v2-design.md specs/sdk-python.md` | `.` | 0 | Passed: embedded compatibility identifiers are limited to the manifested harness filenames and public Python SDK class names. |
| `rg -n -e 'config/agentfield.yaml' -e 'AGENTFIELD_CONFIG_FILE' docs specs control-plane/README.md control-plane/scripts/README.md control-plane/tools/perf/README.md control-plane/migrations/README.md tests/functional/README.md tests/functional/docker/LOG_DEMO.md` | `.` | 0 | Passed: YAML config review confirmed `config/agentfield.yaml` and `AGENTFIELD_CONFIG_FILE` remain documented as first-class surfaces alongside the rebranded product copy. |
| `./scripts/check-silmari-rebrand.sh` | `.` | 0 | Passed: scanned 1336 files, audited 7 changed files, validated 1327 preserved identifier rows, and confirmed skill mirror parity. |
| `python3 -m pytest tests/test_check_silmari_rebrand.py tests/test_collect_silmari_rebrand_inventory.py` | `.` | 0 | Passed: 64 scanner and manifest inventory regression tests. |
| `python3 -m pytest tests/test_docs_specs_branding.py` | `.` | 0 | Passed: 27 docs/spec manifest regression tests, including compatibility identifier and YAML-config assertions. |
| `python3 -m pytest examples/tests/test_silmari_branding.py` | `.` | 0 | Passed: 18 example-surface branding regression tests. |
| `go test ./internal/templates/... ./internal/skillkit/...` | `control-plane` | 0 | Passed targeted rebrand coverage for generated templates and the embedded skill mirror packages. |
| `npm run lint` | `control-plane/web/client` | 0 | Passed with the checked-in suppressions baseline; only pre-existing warnings remained in unrelated files. |
| `npm run test` | `control-plane/web/client` | 0 | Passed: 131 test files and 662 tests after stabilizing the async Workflow DAG node-count assertion. |
| `npm run build` | `control-plane/web/client` | 0 | Production build succeeded; Vite only reported the existing browser-data and chunk-size warnings. |
| `go test ./...` | `control-plane` | 1 | Environment-specific failure: this runner has `CGO_ENABLED=0` and no `gcc` or `cc`, so SQLite-backed packages fail immediately; inherited `AGENTFIELD_SERVER{,_URL}` values also trip localhost-sensitive service tests. |
| `go test ./...` | `sdk/go` | 0 | Passed all 6 Go SDK packages. |
| `python3 -m pytest` | `sdk/python` | 0 | Passed: 1619 tests, 4 skipped, 12 deselected, and 27 warnings. |
| `npm run lint` | `sdk/typescript` | 0 | Passed `tsc --noEmit`. |
| `CI=1 npm run test:core` | `sdk/typescript` | 0 | Passed: 65 test files and 609 tests. |
| `./scripts/test-all.sh` | `.` | 1 | Aggregate verification stops in the `control-plane` Go suite for the same runner constraints: CGO-backed SQLite is unavailable, root-level `AGENTFIELD_SERVER{,_URL}` overrides break localhost-sensitive tests, and root can bypass the registry write-permission failure expected by `internal/core/services`. |

## Deferred Or Excluded
- docs/silmari-rebrand-manifest.md remains validation input and is excluded from old-brand match scanning.
- Verification in this runner still depends on environment capabilities outside the rebrand diff: `lychee` is absent, `control-plane/web/client` and `sdk/typescript` required local `npm ci`, and full `control-plane` Go tests cannot enable CGO-backed SQLite because no system C compiler is installed.
