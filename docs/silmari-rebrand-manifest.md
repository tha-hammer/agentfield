# Silmari Rebrand Manifest

## Summary
This manifest reconciles the current integration branch against the scanner contract and records every changed file plus every scanned file that still contains legacy AgentField-family identifiers.
The branch does not yet contain the full stacked surface rebrand history, so remaining visible legacy tokens are captured here as audited legacy records rather than hidden by scanner exclusions.

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
| CLAUDE.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| CODE_OF_CONDUCT.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| SECURITY.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| SUPPORT.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/anti-patterns.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/capability-playbook.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/cli-toolkit.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/examples-map.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/memory-events.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/model-selection.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/patterns-emerge.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/project-claude-template.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/triggers.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/skillkit/skill_data/agentfield/references/verification.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/docker/.env.example.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/go/.env.example.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/go/README.md.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/go/go.mod.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/go/main.go.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/go/reasoners.go.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/python/.env.example.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/python/README.md.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/python/main.py.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/python/reasoners.py.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/python/requirements.txt.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/templates.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/templates_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/typescript/.env.example.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/typescript/README.md.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/typescript/main.ts.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/typescript/package.json.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/internal/templates/typescript/reasoners.ts.tmpl | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| control-plane/migrations/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/scripts/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/tools/perf/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/.env | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/.env.example | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/.env.production | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/AdminTokenPrompt.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/AppLayout.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/AppSidebar.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/AuthGuard.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/HealthStrip.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/execution/TechnicalDetailsPanel.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/forms/ConfigField.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/forms/ConfigurationForm.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/forms/ConfigurationWizard.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/forms/EnvironmentVariableForm.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/nodes/EnhancedNodeDetailHeader.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/nodes/EnhancedNodesHeader.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/reasoners/EmptyReasonersState.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/StatusBadge.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/StatusRefreshButton.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/StatusRefreshButton.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/UnifiedStatusIndicator.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/UnifiedStatusIndicator.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/status/index.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/components/ui/data-formatters.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/config/navigation.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/contexts/ModeContext.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/hooks/queries/useAgents.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/pages/AccessManagementPage.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/pages/AgentsPage.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/pages/NewDashboardPage.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/services/api.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/services/configurationApi.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/components/AppLayout.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/components/AppSidebar.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/components/ConfigurationWizard.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/components/forms/configuration-form.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/components/status/StatusBadge.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/pages/AccessManagementPage.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/pages/AgentsPage.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/pages/NewSettingsPage.restored.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/services/didApi.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/services/identityApi.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/services/observabilityWebhookApi.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/ui/notification.test.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/utils/formattingUtils.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| control-plane/web/client/src/test/utils/schemaUtils.test.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
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
| docs/ARCHITECTURE.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/CONTRIBUTING.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/COVERAGE.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/DEVELOPMENT.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/ENVIRONMENT_VARIABLES.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/RELEASE.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/api/AGENT_NODE_LOGS.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/design/execution-observability-rfc.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/design/harness-v2-design.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| docs/silmari-rebrand-manifest.md | rebranded | ./scripts/check-silmari-rebrand.sh; python3 -m pytest tests/test_check_silmari_rebrand.py |
| examples/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/analyze.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/crewai-bench/benchmark.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/go-bench/go.mod | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/go-bench/main.go | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/langchain-bench/benchmark.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/python-bench/benchmark.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/results/AgentField_Go.json | excluded-historical | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/results/AgentField_Python.json | excluded-historical | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/results/AgentField_TypeScript.json | excluded-historical | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/run_benchmarks.sh | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/ts-bench/benchmark.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/ts-bench/package-lock.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/benchmarks/100k-scale/ts-bench/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/e2e_resilience_tests/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/e2e_resilience_tests/agent_flaky.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/e2e_resilience_tests/agent_healthy.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/e2e_resilience_tests/agent_slow.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/e2e_resilience_tests/run_tests.sh | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_agent_nodes/cmd/multi_version/main.go | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_agent_nodes/cmd/serverless/main.go | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_agent_nodes/go.mod | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_agent_nodes/main.go | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_harness_demo/go.mod | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/go_harness_demo/main.go | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/agentic_rag/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/agentic_rag/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/agentic_rag/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/deep_research/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/deep_research/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/deep_research/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/deep_research/routers/planning.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/deep_research/routers/research.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/docker_hello_world/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/install.sh | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/pipeline_utils.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/product_context.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/routers/ingestion.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/routers/qa.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/routers/query_planning.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/documentation_chatbot/routers/retrieval.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/hello_world/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/hello_world_rag/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/hello_world_rag/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/hello_world_rag/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/hello_world_rag/test_pydantic_skill.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/image_generation_hello_world/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/image_generation_hello_world/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/image_generation_hello_world/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/multi_version/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/permission_agent_a/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/permission_agent_b/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/permission_agent_b/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/Dockerfile | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/rag_eval_client.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/reasoners/__init__.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/ui/app/page.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/ui/public/agentfield-logo-dark.svg | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/rag_evaluation/ui/public/agentfield-logo-light.svg | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/serverless_hello/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/example.ipynb | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/requirements.txt | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/aggregation.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/decision.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/entity.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/scenario.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/simulation_engine/routers/simulation.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/tool_calling/orchestrator.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/tool_calling/worker.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/python_agent_nodes/waiting_state/main.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/triggers-demo/Dockerfile | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/triggers-demo/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/triggers-demo/agent.py | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/triggers-demo/docker-compose.yml | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/triggers-demo/scripts/fire-events.sh | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/discovery-memory/main.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/init-example/main.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/init-example/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/init-example/reasoners.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/multi-version/main.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/package-lock.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/permission-agent-a/main.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/permission-agent-b/main.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/serverless-hello/main.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/main.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/routers/aggregation.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/routers/decision.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/routers/entity.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/routers/scenario.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/simulation/routers/simulation.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/verifiable-credentials/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/verifiable-credentials/main.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/verifiable-credentials/reasoners.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/waiting-state/main.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts-node-examples/waiting-state/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/orchestrator.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/package-lock.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/package.json | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/test_orchestrator.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| examples/ts_agent_nodes/tool_calling/worker.ts | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| scripts/check-silmari-rebrand.sh | rebranded | ./scripts/check-silmari-rebrand.sh; python3 -m pytest tests/test_check_silmari_rebrand.py |
| scripts/collect_silmari_rebrand_inventory.py | rebranded | python3 -m pytest tests/test_collect_silmari_rebrand_inventory.py |
| sdk/go/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
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
| sdk/go/go.mod | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/harness/codex_integration_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/harness/coverage_branches_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/harness/gemini_integration_test.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/harness/opencode.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/harness/runner.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/go/types/types.go | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/MANIFEST.in | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/python/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
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
| sdk/python/pyproject.toml | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
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
| sdk/typescript/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/package-lock.json | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| sdk/typescript/package.json | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
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
| skills/agentfield/SKILL.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/commands/agentfield.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/anti-patterns.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/capability-playbook.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/cli-toolkit.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/examples-map.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/live-docs.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/memory-events.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/model-selection.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/patterns-emerge.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/primitives-snapshot.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/project-claude-template.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/scaffold-recipe.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/triggers.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| skills/agentfield/references/verification.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| specs/README.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/agentplane-ui-api-spec.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/architecture-overview.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/control-plane.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/data-flow.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/deployment.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/sdk-go.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/sdk-python.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/security.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/viewer.html | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| specs/web-ui.md | audited-no-change | ./scripts/check-silmari-rebrand.sh |
| tests/functional/.env.example | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
| tests/functional/README.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
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
| tests/functional/docker/LOG_DEMO.md | excluded-runtime-compatibility | ./scripts/check-silmari-rebrand.sh |
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
| control-plane/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/SKILL.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/commands/agentfield.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/anti-patterns.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/capability-playbook.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/cli-toolkit.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/examples-map.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/live-docs.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/memory-events.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/model-selection.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/patterns-emerge.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/primitives-snapshot.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/project-claude-template.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/project-claude-template.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/scaffold-recipe.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/triggers.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/triggers.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/skillkit/skill_data/agentfield/references/verification.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/docker/.env.example.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/docker/docker-compose.yml.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/.env.example.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/go/README.md.tmpl | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/go/README.md.tmpl | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/templates/go/go.mod.tmpl | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/go.mod.tmpl | agentfield | go-module-path | Go module or import compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/main.go.tmpl | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/main.go.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/reasoners.go.tmpl | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/go/reasoners.go.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/python/.env.example.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/python/README.md.tmpl | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/python/README.md.tmpl | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/templates/python/main.py.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/python/main.py.tmpl | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/python/main.py.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/python/reasoners.py.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/python/requirements.txt.tmpl | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| control-plane/internal/templates/templates.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/templates_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/.env.example.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/README.md.tmpl | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/README.md.tmpl | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/main.ts.tmpl | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/main.ts.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/package.json.tmpl | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| control-plane/internal/templates/typescript/reasoners.ts.tmpl | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| control-plane/migrations/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/migrations/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/scripts/README.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| control-plane/scripts/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/scripts/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/tools/perf/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/tools/perf/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| control-plane/web/client/.env | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/.env.example | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/.env.production | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/AdminTokenPrompt.tsx | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/web/client/src/components/AppLayout.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/AppSidebar.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/AuthGuard.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/HealthStrip.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/components/execution/TechnicalDetailsPanel.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
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
| control-plane/web/client/src/contexts/ModeContext.tsx | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/hooks/queries/useAgents.ts | agentfield | json-field | Serialized field or stored identifier compatibility keeps agentfield stable in this UI surface; covers all occurrences in this file. |
| control-plane/web/client/src/pages/AccessManagementPage.tsx | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/web/client/src/pages/AccessManagementPage.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/pages/AgentsPage.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/pages/NewDashboardPage.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/pages/NewSettingsPage.tsx | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| control-plane/web/client/src/services/api.ts | agentfield | api-path | Runtime API compatibility keeps agentfield stable for existing request and response contracts; covers all occurrences in this file. |
| control-plane/web/client/src/services/configurationApi.ts | agentfield | api-path | Runtime API compatibility keeps agentfield stable for existing request and response contracts; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/AppLayout.test.tsx | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/AppSidebar.test.tsx | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/ConfigurationWizard.test.tsx | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/ConfigurationWizard.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/forms/configuration-form.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/components/status/StatusBadge.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/pages/AccessManagementPage.test.tsx | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| control-plane/web/client/src/test/pages/AgentsPage.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/pages/NewSettingsPage.restored.test.tsx | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| control-plane/web/client/src/test/pages/NewSettingsPage.restored.test.tsx | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
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
| docs/COVERAGE.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/DEVELOPMENT.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| docs/DEVELOPMENT.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| docs/DEVELOPMENT.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/DEVELOPMENT.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/ENVIRONMENT_VARIABLES.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| docs/ENVIRONMENT_VARIABLES.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/ENVIRONMENT_VARIABLES.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/ENVIRONMENT_VARIABLES.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| docs/RELEASE.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| docs/RELEASE.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/RELEASE.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/RELEASE.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/VC_AUTHORIZATION_ARCHITECTURE.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/api/AGENT_NODE_LOGS.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| docs/api/AGENT_NODE_LOGS.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| docs/design/execution-observability-rfc.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/design/harness-v2-design.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| docs/design/harness-v2-design.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| examples/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/analyze.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/crewai-bench/benchmark.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/go.mod | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/go.mod | agentfield | go-module-path | Go module or import compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/main.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/go-bench/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/langchain-bench/benchmark.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/python-bench/benchmark.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/python-bench/benchmark.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/requirements.txt | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/results/AgentField_Go.json | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/results/AgentField_Python.json | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/results/AgentField_TypeScript.json | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/run_benchmarks.sh | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/ts-bench/benchmark.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/ts-bench/package-lock.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/benchmarks/100k-scale/ts-bench/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/e2e_resilience_tests/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/README.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| examples/e2e_resilience_tests/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_flaky.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_flaky.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_healthy.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_healthy.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_slow.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/agent_slow.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/e2e_resilience_tests/run_tests.sh | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/e2e_resilience_tests/run_tests.sh | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/multi_version/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/multi_version/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/multi_version/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_a/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/permission_agent_b/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/serverless/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/serverless/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/serverless/main.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/go_agent_nodes/cmd/serverless/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/go.mod | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/go.mod | agentfield | go-module-path | Go module or import compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/main.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/go_agent_nodes/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_agent_nodes/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_harness_demo/go.mod | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_harness_demo/go.mod | agentfield | go-module-path | Go module or import compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/go_harness_demo/main.go | Agent-Field | go-module-path | Go module or import compatibility keeps Agent-Field stable for existing source references; covers all occurrences in this file. |
| examples/go_harness_demo/main.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/agentic_rag/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/agentic_rag/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/agentic_rag/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/agentic_rag/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/agentic_rag/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/routers/planning.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/routers/planning.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/deep_research/routers/research.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/docker_hello_world/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/docker_hello_world/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/docker_hello_world/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/install.sh | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/install.sh | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/pipeline_utils.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/product_context.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/product_context.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/ingestion.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/ingestion.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/qa.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/qa.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/query_planning.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/query_planning.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/documentation_chatbot/routers/retrieval.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/hello_world_rag/test_pydantic_skill.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/requirements.txt | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/image_generation_hello_world/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/multi_version/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/multi_version/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_a/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_a/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_b/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_b/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/permission_agent_b/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/Dockerfile | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/docker-compose.yml | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/rag_eval_client.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/reasoners/__init__.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/requirements.txt | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/app/page.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/app/page.tsx | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/components/PoweredBy.tsx | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/public/agentfield-logo-dark.svg | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/rag_evaluation/ui/public/agentfield-logo-light.svg | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/serverless_hello/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/serverless_hello/main.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/serverless_hello/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/example.ipynb | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/example.ipynb | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/main.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/python_agent_nodes/simulation_engine/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
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
| examples/python_agent_nodes/waiting_state/main.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/triggers-demo/Dockerfile | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/triggers-demo/Dockerfile | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| examples/triggers-demo/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/triggers-demo/agent.py | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/triggers-demo/agent.py | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/triggers-demo/agent.py | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/triggers-demo/docker-compose.yml | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/triggers-demo/docker-compose.yml | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/triggers-demo/docker-compose.yml | agentfield | docker-image-or-volume | Container image or volume compatibility keeps agentfield stable for existing deployment flows; covers all occurrences in this file. |
| examples/triggers-demo/scripts/fire-events.sh | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/triggers-demo/scripts/fire-events.sh | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/discovery-memory/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/discovery-memory/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/init-example/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/init-example/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/init-example/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/ts-node-examples/init-example/reasoners.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/multi-version/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/multi-version/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/package-lock.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/ts-node-examples/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| examples/ts-node-examples/permission-agent-a/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/permission-agent-a/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/permission-agent-b/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/permission-agent-b/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/serverless-hello/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/serverless-hello/main.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/serverless-hello/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/aggregation.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/decision.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/entity.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/scenario.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/simulation/routers/simulation.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/README.md | agentfield | historical-record | Legacy branch snapshot intentionally keeps agentfield visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/README.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/main.ts | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/main.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/verifiable-credentials/reasoners.ts | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| examples/ts-node-examples/waiting-state/main.ts | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
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
| sdk/go/harness/codex_integration_test.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/go/harness/coverage_branches_test.go | agentfield | test-fixture | Test coverage intentionally keeps agentfield so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| sdk/go/harness/gemini_integration_test.go | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| sdk/go/harness/opencode.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/harness/runner.go | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/go/types/types.go | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/MANIFEST.in | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| sdk/python/README.md | Agent-Field | historical-record | Legacy branch snapshot intentionally keeps Agent-Field visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/README.md | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
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
| sdk/python/pyproject.toml | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/python/pyproject.toml | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
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
| sdk/typescript/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/typescript/README.md | agentfield | import-module-path | Import or module path compatibility keeps agentfield stable for existing source references; covers all occurrences in this file. |
| sdk/typescript/package-lock.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
| sdk/typescript/package.json | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| sdk/typescript/package.json | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| sdk/typescript/package.json | agentfield | package-name | Published package metadata keeps agentfield stable for install and dependency compatibility; covers all occurrences in this file. |
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
| skills/agentfield/SKILL.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/SKILL.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/commands/agentfield.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/commands/agentfield.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/commands/agentfield.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/references/anti-patterns.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/capability-playbook.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/cli-toolkit.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/examples-map.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/live-docs.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/live-docs.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/live-docs.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/references/memory-events.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/model-selection.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/patterns-emerge.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/primitives-snapshot.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| skills/agentfield/references/primitives-snapshot.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/primitives-snapshot.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/references/project-claude-template.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/project-claude-template.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/scaffold-recipe.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| skills/agentfield/references/scaffold-recipe.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| skills/agentfield/references/scaffold-recipe.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/triggers.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| skills/agentfield/references/triggers.md | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| skills/agentfield/references/verification.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/README.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| specs/README.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/README.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/agentplane-ui-api-spec.md | AgentPlane | historical-record | Legacy branch snapshot intentionally keeps AgentPlane visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/agentplane-ui-api-spec.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/agentplane-ui-api-spec.md | agentplane | historical-record | Legacy branch snapshot intentionally keeps agentplane visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/architecture-overview.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| specs/architecture-overview.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/architecture-overview.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/control-plane.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| specs/control-plane.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/control-plane.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/data-flow.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/data-flow.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/deployment.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| specs/deployment.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/deployment.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/sdk-go.md | Agent-Field | skill-or-repo-slug | Repository or skill slug compatibility keeps Agent-Field stable for existing paths and references; covers all occurrences in this file. |
| specs/sdk-go.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/sdk-go.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/sdk-python.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/sdk-python.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/security.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/security.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/viewer.html | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| specs/viewer.html | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/viewer.html | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| specs/viewer.html | agentfield.ai | published-link-target | Published URL target remains agentfield.ai until a verified Silmari replacement exists; covers all occurrences in this file. |
| specs/web-ui.md | AgentField | historical-record | Legacy branch snapshot intentionally keeps AgentField visible while stacked rebrand surfaces are reconciled; covers all occurrences in this file. |
| specs/web-ui.md | agentfield | skill-or-repo-slug | Repository or skill slug compatibility keeps agentfield stable for existing paths and references; covers all occurrences in this file. |
| tests/functional/.env.example | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/.env.example | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/README.md | AGENTFIELD | env-var | Environment variable token remains AGENTFIELD for backward-compatible configuration examples; covers all occurrences in this file. |
| tests/functional/README.md | Agent-Field | test-fixture | Test coverage intentionally keeps Agent-Field so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
| tests/functional/README.md | AgentField | test-fixture | Test coverage intentionally keeps AgentField so legacy compatibility behavior remains asserted; covers all occurrences in this file. |
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

## CodeCleanup Passes
| Pass | Result | Notes |
|---|---|---|
| Brand surface pass | fail | Legacy AgentField prose still appears in this branch outside the manifest and scanner slice, so the audit records those files explicitly. |
| Compatibility pass | pass | The scanner now requires concrete preserved-identifier categories, changed-file coverage, and span-based match coverage. |
| Mirror and generated-template pass | fail | Skill and generated-template surfaces are audited here, but this branch still contains legacy visible copy that later stacked branches must finish replacing. |
| Link and asset pass | fail | Published agentfield.ai link targets and legacy asset labels remain audited in this branch until replacement URLs and assets are verified. |
| Formatting and lint pass | pass | The manifest tables, scanner contract, and focused pytest coverage are all machine-validated in this slice. |
| Verification pass | pass | Scanner, focused pytest coverage, and fallback link-check evidence are recorded in the verification table below. |

## Verification Commands
| Command | Working Directory | Exit Code | Result |
|---|---|---|---|
| lychee --version | . | not-run | lychee is not installed in PATH on this runner, so the primary link checker could not run. |
| npx --yes markdown-link-check docs/silmari-rebrand-manifest.md | . | 0 | Fallback link check ran successfully; the manifest currently contains no hyperlinks. |
| npx --yes markdown-link-check README.md | . | 1 | Fallback link check examined 77 links and only flagged http://localhost:8080 because no local demo server was running during verification. |
| python3 -m pytest tests/test_check_silmari_rebrand.py | . | 0 | Passed 32 scanner contract and property tests covering changed-file coverage, preserved-row validation, and manifest self-scan exclusions. |
| python3 -m pytest tests/test_collect_silmari_rebrand_inventory.py | . | 0 | Passed 24 inventory regression tests after the manifest reconciliation changes. |
| ./scripts/sync-embedded-skills.sh --check | . | 0 | Embedded skills are in sync with the source skill tree. |
| ./scripts/check-silmari-rebrand.sh | . | 0 | Passed with 1333 scanned files, 6 changed files, and 914 preserved identifier rows. |

## Deferred Or Excluded
- docs/silmari-rebrand-manifest.md remains validation input and is excluded from old-brand match scanning.
- Remaining legacy visible copy in this branch is captured as audited legacy inventory because the dependent surface branches are not merged here.
