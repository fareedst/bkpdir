# Semantic Token Validation Report

**Validation Date**: 2026-03-29 23:12:31
**Project**: BkpDir
**Validation Script**: validate-semantic-tokens.sh

**Canonical registry**: `tied/semantic-tokens.yaml` (indexes and detail YAML under `tied/`). **`project-tokens.yaml`** is a short pointer only, not a token matrix.

## Validation Summary

### Statistics
- **Total files scanned**: 463
- **Files with semantic tokens**: 335
- **Files with Unicode icons**: 7
- **Total semantic tokens found**: 9088
- **Total Unicode icons found**: 26

### Semantic Token Patterns

#### Priority Tokens
- `\[CRITICAL\]`
- `\[HIGH\]`
- `\[MEDIUM\]`
- `\[LOW\]`

#### Action Tokens
- `\[ACTION:core-functionality\]`
- `\[ACTION:format-processing\]`
- `\[ACTION:discovery\]`
- `\[ACTION:maintenance\]`
- `\[ACTION:validation\]`
- `\[ACTION:migration\]`

#### TIED traceability tokens (REQ/ARCH/IMPL) [REQ-DOC_015]
- `\[REQ[:\-][A-Z0-9_]+\]` (REQ/ARCH/IMPL tokens)
- `\[ARCH[:\-][A-Z0-9_]+\]` (REQ/ARCH/IMPL tokens)
- `\[IMPL[:\-][A-Z0-9_]+\]` (REQ/ARCH/IMPL tokens)

### Unicode Icons to Replace
- `⭐`
- `🔺`
- `🔶`
- `🔻`
- `🔧`
- `📝`
- `🔍`
- `🛠️`
- `🛡️`
- `🔄`

## Detailed File Analysis

### Files with Semantic Tokens
```
/Users/fareed/Documents/dev/go/bkpdir/project-tokens.yaml
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/ui/config.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/generator/generator.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/templates/manager.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/create.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/config.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/version.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/process.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/root.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/default.yml
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/example.yml
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/examples/basic/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/realtime-validator/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/types.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer_test.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/ai-validation/main.go
/Users/fareed/Documents/dev/go/bkpdir/config.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_adapter.go
/Users/fareed/Documents/dev/go/bkpdir/tools/coverage-differential.go
/Users/fareed/Documents/dev/go/bkpdir/tools/coverage_differential_test.go
/Users/fareed/Documents/dev/go/bkpdir/errors_test.go
/Users/fareed/Documents/dev/go/bkpdir/backup.go
/Users/fareed/Documents/dev/go/bkpdir/binary-installation-instructions-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/formatter.go
/Users/fareed/Documents/dev/go/bkpdir/refactoring_validation_test.go
/Users/fareed/Documents/dev/go/bkpdir/archive_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/decision_metrics_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/metrics.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/compliance_rate_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/goal_alignment_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/scenarios/decision_scenarios_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/performance/validation_performance_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/performance/performance_test_integration.go
/Users/fareed/Documents/dev/go/bkpdir/config_bench_test.go
/Users/fareed/Documents/dev/go/bkpdir/CHANGELOG.md
/Users/fareed/Documents/dev/go/bkpdir/exclude_test.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_stats_test.go
/Users/fareed/Documents/dev/go/bkpdir/adapter_formatters_test.go
/Users/fareed/Documents/dev/go/bkpdir/working-plan-package-interdependency-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenarios.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/errorinjection_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/corruption_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/context_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenarios_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/permissions.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/context.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/errorinjection.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/corruption.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenario_helpers.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/diskspace.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/permissions_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/diskspace_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_validation.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/semantic_token_validator.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_error_formatter.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/doc008_validator.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/compliance_monitor.go
/Users/fareed/Documents/dev/go/bkpdir/config_impl.go
/Users/fareed/Documents/dev/go/bkpdir/docs/context/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/icon-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-context-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-troubleshooting.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-requirements.md
/Users/fareed/Documents/dev/go/bkpdir/docs/migration-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/images/stdd-visualization-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/images/stdd-timeline.md
/Users/fareed/Documents/dev/go/bkpdir/docs/package-reference.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inspection-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/user/specification.md
/Users/fareed/Documents/dev/go/bkpdir/docs/integration-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/index.md
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/config.yaml
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/example-code/main.go
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/pipeline-example/example-code/main.go
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-examples.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/2025-09-18.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/2025-07-14.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/2026-01-08.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/next-task.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/features.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/2026-01-14.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/advanced-patterns.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/troubleshooting.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/getting-started.md
/Users/fareed/Documents/dev/go/bkpdir/docs/data/stdd-trace-schema.md
/Users/fareed/Documents/dev/go/bkpdir/docs/data/stdd-validation-log.md
/Users/fareed/Documents/dev/go/bkpdir/docs/package-interdependency-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/inc_diff_integration_test.go
/Users/fareed/Documents/dev/go/bkpdir/tied/requirements.md
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CICD_PIPELINE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-PROJECT_STRUCTURE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-DIRECTORY_COMPARISON.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-ARCHIVE_FORMAT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-INCREMENTAL_DUPLICATE_PREVENTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-ERROR_HANDLING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-OUTPUT_FORMATTING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-DEPLOYMENT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-TOKEN_SYSTEM.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-EXCLUDE_MERGE_FIX.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-RESOURCE_MANAGEMENT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CFG_005.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-TEST_EXCLUDE_MERGE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CUSTOMIZABLE_FORMAT_STRINGS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-AI_DOCUMENTATION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-BUILD_DISTRIBUTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-DIFF_COMMAND.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-AUTO_DETECTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-FILE_OPERATIONS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-PROCESSING_PATTERNS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-PERFORMANCE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-LIST_LIMIT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CLI_COMMANDS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-SYSTEM_COMPONENTS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CONTEXT_SUPPORT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-GIT_INTEGRATION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-LANGUAGE_SELECTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CFG_006.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CONFIG_SYSTEM.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CONFIG_OUTPUT_GROUPING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-PERF_VALIDATION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-SECURITY.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CLI_FRAMEWORK.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-EXCLUSION_PATTERNS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-TESTING_STRATEGY.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-CODE_ORGANIZATION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-PACKAGE_EXTRACTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-EXTENSIBILITY.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions/ARCH-FILE_STATISTICS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TOKEN_MIGRATION_COMPLETE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-RESOURCE_MANAGER.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-EXTRACTION_PRINCIPLES.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-FILE_STATISTICS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-DEPENDENCY_MGMT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-ATOMIC_OPS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-DIFF_COMMAND.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CLI_FRAMEWORK.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-ZIP_FORMAT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CFG_INHERITANCE_PATH_RESOLUTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CONFIG_OUTPUT_GROUPING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-STRUCTURED_ERRORS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CFG_MIXED_SEQUENTIAL_INHERITANCE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CONFIG_DISPLAY_FLATTENING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TIED_FILES.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-BACKWARD_COMPAT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CODE_STYLE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CFG_HIERARCHY_PRESERVATION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-EXTRACT_008_DOC_MIGRATION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CFG_PRECEDENCE_FIX.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-GIT_CLI.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CUSTOMIZABLE_FORMAT_STRINGS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CFG_QUOTED_KEY_PREFIX.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TEST_PREPEND_ORDERING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CFG_006.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TEST_DEFAULT_STRATEGY_EDGES.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-EXCLUDE_MERGE_FIX.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-PACKAGE_EXTRACTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-DATA_MODELS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CONFIG_SCHEMA_FLEX.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-EXCLUSION_PATTERNS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-LIST_FORMAT_SAFETY.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TEST_CFG_005_P1.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-INTERFACE_DRIVEN.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TOKEN_SYSTEM.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-PROCESSING_PATTERNS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-MODULE_VALIDATION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-DOC_ENHANCEMENT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TRACEABILITY.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TESTING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TOKEN_COVERAGE_AUDIT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TEST_UNICODE_HANDLING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-STDD_VIS_ASSETS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TESTING_COMPLEXITY.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-EXTRACTION_CHALLENGES.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-LIST_LIMIT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TEST_EMPTY_STRING_HANDLING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-FILE_OPERATIONS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CFG_MERGE_BEHAVIOR_REGISTRY.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TEST_COVERAGE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-REFACTOR_PREP.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-FILE_STATISTICS_TEMPLATE_FIX.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CFG_MIXED_MODE_MERGE_FIX.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-LARGE_FILE_DECOMP.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-LARGE_FILE_CHALLENGE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-LAYERED_EXTRACTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-DELAYED_OUTPUT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CONTEXT_OPS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-INCREMENTAL_DUPLICATE_PREVENTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-MCP_FEEDBACK_TOOLS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-SEMANTIC_CROSS_REF.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-ZERO_BREAKING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CONFIG_STRUCT.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-AUTO_DETECTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-STDD_VIS_DATA_PIPELINE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-GIT_DIRTY_CONFIG.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-INTERFACE_FIRST.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-DUAL_FORMATTING.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CFG_MERGE_PREPEND_PRECEDENCE_FIX.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-CONFIGURABLE_STRINGS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-TEST_EXCLUDE_MERGE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions/IMPL-DIRECTORY_COMPARISON.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/processes.md
/Users/fareed/Documents/dev/go/bkpdir/tied/requirements.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/detail-files-schema.md
/Users/fareed/Documents/dev/go/bkpdir/tied/ai-assistant-compliance.md
/Users/fareed/Documents/dev/go/bkpdir/tied/methodology/implementation-decisions/IMPL-TIED_FILES.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/methodology/implementation-decisions/IMPL-MCP_FEEDBACK_TOOLS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/methodology/implementation-decisions.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/semantic-tokens.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/requirements/REQ-TEST_EXCLUDE_MERGE.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/requirements/REQ-DIFF_COMMAND.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/requirements/REQ-INCREMENTAL_DUPLICATE_PREVENTION.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/requirements/REQ-DOC_015.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/requirements/REQ-STDD_VIS.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/docs/pseudocode-writing-and-validation.md
/Users/fareed/Documents/dev/go/bkpdir/tied/docs/agent-req-implementation-checklist.md
/Users/fareed/Documents/dev/go/bkpdir/tied/docs/adding-tied-mcp-and-invoking-passes.md
/Users/fareed/Documents/dev/go/bkpdir/tied/docs/agent-req-implementation-checklist.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/docs/impl-code-test-linkage.md
/Users/fareed/Documents/dev/go/bkpdir/tied/docs/using-tied-without-mcp.md
/Users/fareed/Documents/dev/go/bkpdir/tied/docs/methodology-diagrams.md
/Users/fareed/Documents/dev/go/bkpdir/tied/docs/tied-first-implementation-procedure.md
/Users/fareed/Documents/dev/go/bkpdir/tied/impl-code-test-linkage.md
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions.md
/Users/fareed/Documents/dev/go/bkpdir/tied/tasks.md
/Users/fareed/Documents/dev/go/bkpdir/tied/implementation-decisions.md
/Users/fareed/Documents/dev/go/bkpdir/tied/semantic-tokens.md
/Users/fareed/Documents/dev/go/bkpdir/tied/merge-strategy-conflict-resolution.md
/Users/fareed/Documents/dev/go/bkpdir/tied/architecture-decisions.yaml
/Users/fareed/Documents/dev/go/bkpdir/tied/ai-principles.md
/Users/fareed/Documents/dev/go/bkpdir/ai_formatter_adapter.go
/Users/fareed/Documents/dev/go/bkpdir/README.md
/Users/fareed/Documents/dev/go/bkpdir/archive.go
/Users/fareed/Documents/dev/go/bkpdir/diff_command_test.go
/Users/fareed/Documents/dev/go/bkpdir/file_stats_test.go
/Users/fareed/Documents/dev/go/bkpdir/task-analysis-and-plan.md
/Users/fareed/Documents/dev/go/bkpdir/config_grouping_test.go
/Users/fareed/Documents/dev/go/bkpdir/exclude.go
/Users/fareed/Documents/dev/go/bkpdir/config_adapter.go
/Users/fareed/Documents/dev/go/bkpdir/file_stats.go
/Users/fareed/Documents/dev/go/bkpdir/comparison.go
/Users/fareed/Documents/dev/go/bkpdir/refactoring-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/config_interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_list_safety_test.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_adapter_simple_test.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_adapter_simple.go
/Users/fareed/Documents/dev/go/bkpdir/AGENTS.md
/Users/fareed/Documents/dev/go/bkpdir/config_test.go
/Users/fareed/Documents/dev/go/bkpdir/config_integration_test.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_test.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_adapter_test.go
/Users/fareed/Documents/dev/go/bkpdir/SPECIFICATION_AUDIT_REPORT.md
/Users/fareed/Documents/dev/go/bkpdir/ai-principles.md
/Users/fareed/Documents/dev/go/bkpdir/main.go
/Users/fareed/Documents/dev/go/bkpdir/backup_test.go
/Users/fareed/Documents/dev/go/bkpdir/comparison_test.go
/Users/fareed/Documents/dev/go/bkpdir/conversation-template.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/discovery.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/merge_strategies.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/loader.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/utils.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/inheritance.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/config_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/concurrent.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/processor_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/processor.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/naming.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/pipeline.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/doc.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/fileops.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/traversal.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/exclusion.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/comparison.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/validation.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/atomic.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/examples_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/resources_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/context.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/version.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/example_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/types.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/flags.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/context.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/builder.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/dryrun.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/cli_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/ai_core_formatter.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/patterns.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/ai_first_formatter_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/formatter.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/collector.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/ai_output_manager.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/filestats.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/placeholder_replace_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/template.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/placeholder_replace.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/ai_first_formatter.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/ai_pattern_extractor.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/ai_first_interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/AI_FIRST_REFACTORING_SUMMARY.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/formatter_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/classification.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/errors_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/handlers.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/provider.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/config.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/testutil_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/scenarios.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/assertions.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/filesystem.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/doc.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/cli.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/fixtures.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/integration_demo_test.go
/Users/fareed/Documents/dev/go/bkpdir/incremental_duplicate_prevention_test.go
/Users/fareed/Documents/dev/go/bkpdir/main_test.go
/Users/fareed/Documents/dev/go/bkpdir/errors.go
```

### Files with Remaining Unicode Icons
```
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/templates/manager.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/create.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/realtime-validator/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/ai-validation/main.go
/Users/fareed/Documents/dev/go/bkpdir/refactoring_validation_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/decision_metrics_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/compliance_rate_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/performance/validation_performance_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/performance/performance_test_integration.go
/Users/fareed/Documents/dev/go/bkpdir/working-plan-package-interdependency-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator_test.go
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/icon-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-context-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-framework-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-troubleshooting.md
/Users/fareed/Documents/dev/go/bkpdir/docs/migration-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/package-reference.md
/Users/fareed/Documents/dev/go/bkpdir/docs/format-strings-reference.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/2025-09-18.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/next-task.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/features.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/advanced-patterns.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/getting-started.md
/Users/fareed/Documents/dev/go/bkpdir/tied/ai-assistant-compliance.md
/Users/fareed/Documents/dev/go/bkpdir/tied/tasks.md
/Users/fareed/Documents/dev/go/bkpdir/example-custom-formats.yml
/Users/fareed/Documents/dev/go/bkpdir/refactoring-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/config_test.go
/Users/fareed/Documents/dev/go/bkpdir/SPECIFICATION_AUDIT_REPORT.md
/Users/fareed/Documents/dev/go/bkpdir/conversation-template.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/AI_QUICK_REFERENCE.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/AI_FIRST_REFACTORING_SUMMARY.md
```

## Validation Results

### Success Criteria
- ✅ **Semantic tokens present**: 335 files
- ✅ **Total semantic tokens**: 9088 tokens
- ❌ **Unicode icons remaining**: 7 files
- ❌ **Total Unicode icons**: 26 icons

### Recommendations

1. **Token Standardization**: Ensure all new code uses semantic tokens
2. **Documentation Update**: Update documentation to reference semantic tokens
3. **AI Assistant Integration**: Update AI assistant protocols to use semantic tokens
4. **Validation Automation**: Integrate validation into CI/CD pipeline
5. **Continuous Validation**: Add validation to development workflow

## Token Usage Guidelines

### Priority Tokens
- Use `[CRITICAL]` for essential system operations
- Use `[HIGH]` for important features
- Use `[MEDIUM]` for standard features
- Use `[LOW]` for minor improvements

### Action Tokens
- Use `[ACTION:core-functionality]` for essential operations
- Use `[ACTION:format-processing]` for text formatting
- Use `[ACTION-discovery]` for file system operations
- Use `[ACTION-maintenance]` for code cleanup
- Use `[ACTION-validation]` for error checking
- Use `[ACTION-migration]` for system migrations

## Compliance Status

**Overall Status**: ❌ NON-COMPLIANT

**Migration Progress**: 72%

**Remaining Work**: 7 files need migration
