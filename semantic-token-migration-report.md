# Semantic Token Migration Report

**Migration Date**: 2025-07-13 15:38:20
**Project**: BkpDir
**Migration Script**: migrate-to-semantic-tokens.sh

## Migration Summary

### Token Mapping Applied

#### Priority Tokens
- [CRITICAL] → [CRITICAL]
- [HIGH] → [HIGH]
- [MEDIUM] → [MEDIUM]
- [LOW] → [LOW]

#### Action Tokens
- [ACTION:core-functionality] → [ACTION:core-functionality]
- [ACTION:format-processing] → [ACTION:format-processing]
- [ACTION:discovery] → [ACTION:discovery]
- [ACTION:maintenance] → [ACTION:maintenance]
- [ACTION:validation] → [ACTION:validation]
- [ACTION:migration] → [ACTION:migration]

### Files Processed

```
/Users/fareed/Documents/dev/go/bkpdir/.bkpdir.yml
/Users/fareed/Documents/dev/go/bkpdir/archive_test.go
/Users/fareed/Documents/dev/go/bkpdir/archive.go
/Users/fareed/Documents/dev/go/bkpdir/backup_test.go
/Users/fareed/Documents/dev/go/bkpdir/backup.go
/Users/fareed/Documents/dev/go/bkpdir/binary-installation-instructions-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-implementation-plan.md
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-performance-optimization-plan.md
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-comprehensive-testing-plan.md
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-8-documentation-plan.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/ai-validation/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/config.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/create.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/process.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/root.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/version.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/default.yml
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/example.yml
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/examples/basic/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/realtime-validator/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/generator/generator.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/templates/manager.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/ui/config.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer_test.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/types.go
/Users/fareed/Documents/dev/go/bkpdir/comparison_test.go
/Users/fareed/Documents/dev/go/bkpdir/comparison.go
/Users/fareed/Documents/dev/go/bkpdir/config_adapter.go
/Users/fareed/Documents/dev/go/bkpdir/config_bench_test.go
/Users/fareed/Documents/dev/go/bkpdir/config_impl.go
/Users/fareed/Documents/dev/go/bkpdir/config_interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/config_test.go
/Users/fareed/Documents/dev/go/bkpdir/config.go
/Users/fareed/Documents/dev/go/bkpdir/DEVELOPMENT.md
/Users/fareed/Documents/dev/go/bkpdir/DOC-014-subtask-4-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/ai-first-token-system-proposal.md
/Users/fareed/Documents/dev/go/bkpdir/docs/build-system-fixes-implementation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/config-abstraction.md
/Users/fareed/Documents/dev/go/bkpdir/docs/config-schema-abstraction.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-examples.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inheritance.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inspection-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-troubleshooting.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-compliance.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-optimization-tasks.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-protocol.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-token-protocol.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-code-maintenance-standards.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-decision-framework.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-documentation-templates.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-configuration-refactoring-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/AI-First-Development-Procedure-Complete-Guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-development-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-error-handling-refactoring-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-formatter-refactoring-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-refactoring-priority-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-token-system-comprehensive.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/architecture.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/CFG-005-implementation-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/change-rejection-criteria.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/code-marker-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/context-file-checklist.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/context-file-responsibilities.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/cross-reference-template.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-training-examples.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-troubleshooting.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-usage-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/doc-validation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/enforcement-mechanisms.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-migration-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-syntax-specification.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-traceability.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-change-protocol.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-documentation-standards.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-tracking.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-system-proposal.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-analysis.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-guidelines.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-validation-enforcement.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/immutable.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-decisions.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-status.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/requirements.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/semantic-links-implementation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/source-code-icon-guidelines.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/specification.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/structure-optimization-analysis.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/sync-framework.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/task-completion-enforcement.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/testing.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/unicode-to-semantic-token-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/validation-automation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/coverage-baseline.md
/Users/fareed/Documents/dev/go/bkpdir/docs/coverage-implementation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/delayed-output-usage.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-010-implementation-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-011-implementation-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-012-completion-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-cross-reference-completion-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-quick-start-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-5-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/config.yaml
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/example-code/main.go
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/git-aware-backup/example-code/main.go
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/git-aware-backup/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-final-analysis.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-task-completion-final.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-cli-template-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-remaining-subtasks-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-session-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-1-completion.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-2-completion.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-3-completion.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-5-completion.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-009-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extraction-dependencies.md
/Users/fareed/Documents/dev/go/bkpdir/docs/formatter-decomposition.md
/Users/fareed/Documents/dev/go/bkpdir/docs/integration-documentation-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/integration-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/interface-definitions.md
/Users/fareed/Documents/dev/go/bkpdir/docs/migration-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/milestone-tracking-system.md
/Users/fareed/Documents/dev/go/bkpdir/docs/package-interdependency-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/docs/package-reference.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/add-task.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/adhoc-new-task.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/context-documents.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/continue-task.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/document-completed-task.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/existing-task.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/features.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/fix.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/format-strings.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/import-bkpfile.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/next-task.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/repetition.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-completion-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-phase2f-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-progress-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-implementation-complete.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-requirements.md
/Users/fareed/Documents/dev/go/bkpdir/docs/system-comparison-analysis.md
/Users/fareed/Documents/dev/go/bkpdir/docs/task-config-template-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tiered-performance-testing-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/advanced-patterns.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/getting-started.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/troubleshooting.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask1.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask2.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask3.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask4.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase1-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-critical-tokens-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask1-missing-token-inventory.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-code-layer-progress.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-implementation-progress.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-test-layer-progress.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-creation-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-implementation-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask3-token-implementation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week2-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week3-command-structure-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week4-core-files-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-context-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-framework-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/icon-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/semantic-token-validation-report-20250713_092234.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/semantic-token-validation-report-20250713_092247.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/semantic-token-validation-report-20250713_092636.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/token-traceability-report.md
/Users/fareed/Documents/dev/go/bkpdir/errors_test.go
/Users/fareed/Documents/dev/go/bkpdir/errors.go
/Users/fareed/Documents/dev/go/bkpdir/example-git-config.yml
/Users/fareed/Documents/dev/go/bkpdir/example-inheritance-base.yml
/Users/fareed/Documents/dev/go/bkpdir/example-inheritance-child.yml
/Users/fareed/Documents/dev/go/bkpdir/example-symlink-config.yml
/Users/fareed/Documents/dev/go/bkpdir/exclude_test.go
/Users/fareed/Documents/dev/go/bkpdir/exclude.go
/Users/fareed/Documents/dev/go/bkpdir/file_stats_test.go
/Users/fareed/Documents/dev/go/bkpdir/file_stats.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_adapter.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_test.go
/Users/fareed/Documents/dev/go/bkpdir/formatter.go
/Users/fareed/Documents/dev/go/bkpdir/git_test.go
/Users/fareed/Documents/dev/go/bkpdir/git.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/context_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/context.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/corruption_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/corruption.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/diskspace_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/diskspace.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/errorinjection_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/errorinjection.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/permissions_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/permissions.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenario_helpers.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenarios_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenarios.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_error_formatter.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_validation.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/compliance_monitor.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/decision_checklist_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/decision_checklist.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/doc008_validator.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
/Users/fareed/Documents/dev/go/bkpdir/main_test.go
/Users/fareed/Documents/dev/go/bkpdir/main.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/builder.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/cli_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/context.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/dryrun.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/example_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/flags.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/types.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/version.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/config_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/discovery.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/inheritance.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/loader.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/merge_strategies.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/utils.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/classification.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/errors_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/handlers.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/atomic.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/comparison.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/exclusion.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/fileops.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/traversal.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/validation.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/collector.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/filestats.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/formatter_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/formatter.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/patterns.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/template.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/git/git_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/git/git.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/git/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/concurrent.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/doc.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/naming.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/pipeline.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/processor_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/processor.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/verification.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/context.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/resources_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/assertions.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/cli.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/config.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/doc.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/filesystem.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/fixtures.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/integration_demo_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/provider.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/scenarios.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/testutil_test.go
/Users/fareed/Documents/dev/go/bkpdir/project-tokens.yaml
/Users/fareed/Documents/dev/go/bkpdir/README.md
/Users/fareed/Documents/dev/go/bkpdir/refactoring_validation_test.go
/Users/fareed/Documents/dev/go/bkpdir/refactoring-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/semantic-token-migration-report.md
/Users/fareed/Documents/dev/go/bkpdir/SPECIFICATION_AUDIT_REPORT.md
/Users/fareed/Documents/dev/go/bkpdir/task-analysis-and-plan.md
/Users/fareed/Documents/dev/go/bkpdir/test-dir/.bkpdir.yml
/Users/fareed/Documents/dev/go/bkpdir/test/integration/doc014_integration_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/compliance_rate_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/decision_metrics_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/goal_alignment_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/metrics.go
/Users/fareed/Documents/dev/go/bkpdir/test/performance/performance_test_integration.go
/Users/fareed/Documents/dev/go/bkpdir/test/performance/validation_performance_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/scenarios/decision_scenarios_test.go
/Users/fareed/Documents/dev/go/bkpdir/tools/coverage_differential_test.go
/Users/fareed/Documents/dev/go/bkpdir/tools/coverage-differential.go
/Users/fareed/Documents/dev/go/bkpdir/verify_test.go
/Users/fareed/Documents/dev/go/bkpdir/verify.go
/Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-008.md
/Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-009.md
/Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-010.md
/Users/fareed/Documents/dev/go/bkpdir/working-plan-package-interdependency-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/working-plan-remove-backup-command.md
```

### Migration Log

```
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/project-tokens.yaml
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/project-tokens.yaml
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/project-tokens.yaml
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/ui/config.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/generator/generator.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/templates/manager.go
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/manager.go
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/templates/manager.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/main.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/create.go
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/create.go
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/create.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/config.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/version.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/process.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/root.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/default.yml
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/default.yml
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/default.yml
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/example.yml
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/example.yml
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/example.yml
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/README.md
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/README.md
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/examples/basic/README.md
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/examples/basic/README.md
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/main.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/realtime-validator/main.go
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/main.go
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/realtime-validator/main.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/types.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer.go
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/analyzer.go
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer_test.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/main.go
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/cmd/ai-validation/main.go
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/main.go
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/cmd/ai-validation/main.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/config.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/formatter_adapter.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/tools/coverage-differential.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/tools/coverage_differential_test.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/errors_test.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/backup.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/binary-installation-instructions-working-plan.md
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/binary-installation-instructions-working-plan.md
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/binary-installation-instructions-working-plan.md
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/formatter.go
[2025-07-13 15:38:09] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/refactoring_validation_test.go
[2025-07-13 15:38:09] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/refactoring_validation_test.go
[2025-07-13 15:38:09] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/refactoring_validation_test.go
[2025-07-13 15:38:09] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/archive_test.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-010.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/working-plan-extract-010.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-010.md
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/test/metrics/decision_metrics_test.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/decision_metrics_test.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/test/metrics/decision_metrics_test.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/test/metrics/metrics.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/test/metrics/compliance_rate_test.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/compliance_rate_test.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/test/metrics/compliance_rate_test.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/test/metrics/goal_alignment_test.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/goal_alignment_test.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/test/metrics/goal_alignment_test.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/test/integration/doc014_integration_test.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/doc014_integration_test.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/test/integration/doc014_integration_test.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/test/scenarios/decision_scenarios_test.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/decision_scenarios_test.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/test/scenarios/decision_scenarios_test.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/test/performance/validation_performance_test.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/validation_performance_test.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/test/performance/validation_performance_test.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/test/performance/performance_test_integration.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/performance_test_integration.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/test/performance/performance_test_integration.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/example-inheritance-base.yml
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/config_bench_test.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/exclude_test.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-comprehensive-testing-plan.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/CFG-006-subtask-7-comprehensive-testing-plan.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-comprehensive-testing-plan.md
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/CFG-006-performance-optimization-plan.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/CFG-006-performance-optimization-plan.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-performance-optimization-plan.md
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/working-plan-package-interdependency-mapping.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/working-plan-package-interdependency-mapping.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/working-plan-package-interdependency-mapping.md
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/example-inheritance-child.yml
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenarios.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/errorinjection_test.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/corruption_test.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/context_test.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenarios_test.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/permissions.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/context.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/errorinjection.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/corruption.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenario_helpers.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/diskspace.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/permissions_test.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/diskspace_test.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_validation.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai_validation.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_validation.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/realtime_validator.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_error_formatter.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai_error_formatter.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_error_formatter.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/doc008_validator.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator_test.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/realtime_validator_test.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator_test.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/internal/validation/decision_checklist.go
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/decision_checklist.go
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/decision_checklist.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/compliance_monitor.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/decision_checklist_test.go
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/config_impl.go
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/CFG-006-implementation-plan.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/CFG-006-implementation-plan.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-implementation-plan.md
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-completion-summary.md
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-task-completion-final.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-007-task-completion-final.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-task-completion-final.md
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-working-plan.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/DOC-014-subtask-6-working-plan.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-working-plan.md
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-test-layer-progress.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-phase2-subtask2-test-layer-progress.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-test-layer-progress.md
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask2.md
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-code-layer-progress.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-phase2-subtask2-code-layer-progress.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-code-layer-progress.md
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-010-implementation-summary.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/DOC-010-implementation-summary.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-010-implementation-summary.md
[2025-07-13 15:38:10] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-012-completion-report.md
[2025-07-13 15:38:10] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/DOC-012-completion-report.md
[2025-07-13 15:38:10] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-012-completion-report.md
[2025-07-13 15:38:10] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/interface-definitions.md
[2025-07-13 15:38:11] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/doc-validation.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-training-examples.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/decision-framework-training-examples.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-training-examples.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/specification.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/specification.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/specification.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-optimization-tasks.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-assistant-optimization-tasks.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-optimization-tasks.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-migration-strategy.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/enhanced-token-migration-strategy.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-migration-strategy.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-documentation-templates.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-documentation-templates.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-documentation-templates.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/source-code-icon-guidelines.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/source-code-icon-guidelines.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/source-code-icon-guidelines.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/requirements.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/requirements.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/requirements.md
[2025-07-13 15:38:11] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/task-completion-enforcement.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-token-system-comprehensive.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-first-token-system-comprehensive.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-token-system-comprehensive.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-syntax-specification.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/enhanced-token-syntax-specification.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-syntax-specification.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/architecture.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/architecture.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/architecture.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-formatter-refactoring-plan.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-first-formatter-refactoring-plan.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-formatter-refactoring-plan.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/AI-First-Development-Procedure-Complete-Guide.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/AI-First-Development-Procedure-Complete-Guide.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/AI-First-Development-Procedure-Complete-Guide.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-decision-framework.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-decision-framework.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-decision-framework.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-compliance.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-assistant-compliance.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-compliance.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-status.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/implementation-status.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-status.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/CFG-005-implementation-plan.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/CFG-005-implementation-plan.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/CFG-005-implementation-plan.md
[2025-07-13 15:38:11] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/sync-framework.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-system-proposal.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/icon-system-proposal.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-system-proposal.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-code-maintenance-standards.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-code-maintenance-standards.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-code-maintenance-standards.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-analysis.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/icon-usage-analysis.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-analysis.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-troubleshooting.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/decision-framework-troubleshooting.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-troubleshooting.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/context-file-checklist.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/context-file-checklist.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/context-file-checklist.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-development-strategy.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-first-development-strategy.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-development-strategy.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-token-protocol.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-assistant-token-protocol.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-token-protocol.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/testing.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/testing.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/testing.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-guidelines.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/icon-usage-guidelines.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-guidelines.md
[2025-07-13 15:38:11] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/code-marker-strategy.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-change-protocol.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/feature-change-protocol.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-change-protocol.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/README.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/README.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-configuration-refactoring-plan.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-first-configuration-refactoring-plan.md
[2025-07-13 15:38:11] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-configuration-refactoring-plan.md
[2025-07-13 15:38:11] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/validation-automation.md
[2025-07-13 15:38:11] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/validation-automation.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/validation-automation.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-decisions.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/implementation-decisions.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-decisions.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/structure-optimization-analysis.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/structure-optimization-analysis.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/structure-optimization-analysis.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/immutable.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-protocol.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-assistant-protocol.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-protocol.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-refactoring-priority-summary.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-first-refactoring-priority-summary.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-refactoring-priority-summary.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/change-rejection-criteria.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/cross-reference-template.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/enforcement-mechanisms.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/unicode-to-semantic-token-mapping.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-to-semantic-token-mapping.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/unicode-to-semantic-token-mapping.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/semantic-links-implementation.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-tracking.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/feature-tracking.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-tracking.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-documentation-standards.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-usage-guide.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/decision-framework-usage-guide.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-usage-guide.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/context-file-responsibilities.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-traceability.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-error-handling-refactoring-plan.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-first-error-handling-refactoring-plan.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-error-handling-refactoring-plan.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-validation-enforcement.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/semantic-token-validation-report-20250713_092636.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/token-traceability-report.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/semantic-token-validation-report-20250713_092234.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/semantic-token-validation-report-20250713_092247.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/icon-validation-report.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/icon-validation-report.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/icon-validation-report.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-context-validation-report.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/decision-context-validation-report.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-context-validation-report.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-framework-validation-report.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/decision-framework-validation-report.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-framework-validation-report.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-troubleshooting.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/configuration-troubleshooting.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-troubleshooting.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/coverage-implementation.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask3.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase1-completion-summary.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-phase1-completion-summary.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase1-completion-summary.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/tiered-performance-testing-strategy.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/tiered-performance-testing-strategy.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/tiered-performance-testing-strategy.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-requirements.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/semantic-token-system-requirements.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-requirements.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/migration-guide.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/migration-guide.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/migration-guide.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/package-reference.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/package-reference.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/package-reference.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-completion-summary.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/DOC-014-subtask-6-completion-summary.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-completion-summary.md
[2025-07-13 15:38:12] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/delayed-output-usage.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-3-completion.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-008-subtask-3-completion.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-3-completion.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-session-summary.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-008-session-summary.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-session-summary.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inspection-guide.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/configuration-inspection-guide.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inspection-guide.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week3-command-structure-plan.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-week3-command-structure-plan.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week3-command-structure-plan.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/config-schema-abstraction.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/config-schema-abstraction.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/config-schema-abstraction.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-009-completion-summary.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-009-completion-summary.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-009-completion-summary.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-implementation-progress.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-phase2-subtask2-implementation-progress.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-implementation-progress.md
[2025-07-13 15:38:12] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/milestone-tracking-system.md
[2025-07-13 15:38:12] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/milestone-tracking-system.md
[2025-07-13 15:38:12] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/milestone-tracking-system.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/build-system-fixes-implementation.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/build-system-fixes-implementation.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/build-system-fixes-implementation.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/integration-guide.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/integration-guide.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/integration-guide.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-cross-reference-completion-plan.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/DOC-014-cross-reference-completion-plan.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-cross-reference-completion-plan.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/formatter-decomposition.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/formatter-decomposition.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/formatter-decomposition.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-remaining-subtasks-plan.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-008-remaining-subtasks-plan.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-remaining-subtasks-plan.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-5-completion.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-008-subtask-5-completion.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-5-completion.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/config-abstraction.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-progress-summary.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/semantic-token-migration-progress-summary.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-progress-summary.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-1-completion.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-008-subtask-1-completion.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-1-completion.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week2-completion-summary.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-week2-completion-summary.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week2-completion-summary.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extraction-dependencies.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extraction-dependencies.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extraction-dependencies.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-implementation-complete.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/semantic-token-system-implementation-complete.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-implementation-complete.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-final-analysis.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-007-final-analysis.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-final-analysis.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-implementation-plan.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-phase2-subtask2-token-implementation-plan.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-implementation-plan.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask3-token-implementation.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-phase2-subtask3-token-implementation.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask3-token-implementation.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/coverage-baseline.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-5-working-plan.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/DOC-014-subtask-5-working-plan.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-5-working-plan.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask1-missing-token-inventory.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-phase2-subtask1-missing-token-inventory.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask1-missing-token-inventory.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/config.yaml
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/config.yaml
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/config.yaml
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/README.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/README.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/example-code/main.go
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/examples/README.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/examples/README.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/integration-documentation-working-plan.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/integration-documentation-working-plan.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/integration-documentation-working-plan.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-examples.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/configuration-examples.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-examples.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/ai-first-token-system-proposal.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/ai-first-token-system-proposal.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/ai-first-token-system-proposal.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/task-config-template-working-plan.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/task-config-template-working-plan.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/task-config-template-working-plan.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-cli-template-plan.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-008-cli-template-plan.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-cli-template-plan.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/fix.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/document-completed-task.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/repetition.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/import-bkpfile.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/context-documents.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/next-task.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/next-task.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/next-task.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/adhoc-new-task.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/add-task.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/existing-task.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/features.md
[2025-07-13 15:38:13] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/features.md
[2025-07-13 15:38:13] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/features.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/format-strings.md
[2025-07-13 15:38:13] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/continue-task.md
[2025-07-13 15:38:13] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-2-completion.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-008-subtask-2-completion.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-2-completion.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/advanced-patterns.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/advanced-patterns.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/advanced-patterns.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/troubleshooting.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/troubleshooting.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/troubleshooting.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/getting-started.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/getting-started.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/getting-started.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-completion-report.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/semantic-token-migration-completion-report.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-completion-report.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inheritance.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-working-plan.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/extract-007-working-plan.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-working-plan.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-quick-start-guide.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-011-implementation-summary.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/DOC-011-implementation-summary.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-011-implementation-summary.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask4.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-phase2f-summary.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/semantic-token-migration-phase2f-summary.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-phase2f-summary.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/system-comparison-analysis.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/system-comparison-analysis.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/system-comparison-analysis.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/package-interdependency-mapping.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/package-interdependency-mapping.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/package-interdependency-mapping.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week4-core-files-plan.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-week4-core-files-plan.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week4-core-files-plan.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-creation-strategy.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-phase2-subtask2-token-creation-strategy.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-creation-strategy.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-critical-tokens-completion-summary.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/unicode-migration-phase2-critical-tokens-completion-summary.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-critical-tokens-completion-summary.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask1.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/README.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/archive.go
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-completion-summary.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/CFG-006-subtask-7-completion-summary.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-completion-summary.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-009.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/working-plan-extract-009.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-009.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/file_stats_test.go
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/task-analysis-and-plan.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/task-analysis-and-plan.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/task-analysis-and-plan.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-008.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/exclude.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/config_adapter.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/file_stats.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/test-dir/.bkpdir.yml
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/comparison.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/verify_test.go
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/DOC-014-subtask-4-working-plan.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/DOC-014-subtask-4-working-plan.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/DOC-014-subtask-4-working-plan.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/refactoring-validation-report.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/refactoring-validation-report.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/refactoring-validation-report.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/config_interfaces.go
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/working-plan-remove-backup-command.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/working-plan-remove-backup-command.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/working-plan-remove-backup-command.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/config_test.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/formatter_test.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/DEVELOPMENT.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-8-documentation-plan.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/CFG-006-subtask-8-documentation-plan.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-8-documentation-plan.md
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/SPECIFICATION_AUDIT_REPORT.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/SPECIFICATION_AUDIT_REPORT.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/SPECIFICATION_AUDIT_REPORT.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/main.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/example-symlink-config.yml
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/.bkpdir.yml
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/backup_test.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/verify.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/comparison_test.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/config/discovery.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/config/interfaces.go
[2025-07-13 15:38:14] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/pkg/config/README.md
[2025-07-13 15:38:14] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:14] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/pkg/config/README.md
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/config/merge_strategies.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/config/loader.go
[2025-07-13 15:38:14] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/config/utils.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/config/inheritance.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/config/config_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/verification.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/concurrent.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/processor_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/processor.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/naming.go
[2025-07-13 15:38:15] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/README.md
[2025-07-13 15:38:15] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:15] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/README.md
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/pipeline.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/doc.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/fileops.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/traversal.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/README.md
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/exclusion.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/comparison.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/validation.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/atomic.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/resources/interfaces.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/resources/resources_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/resources/README.md
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/resources/context.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/version.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/example_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/types.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/flags.go
[2025-07-13 15:38:15] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/README.md
[2025-07-13 15:38:15] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:15] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/README.md
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/context.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/builder.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/dryrun.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/cli_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/patterns.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/formatter.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/interfaces.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/collector.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/filestats.go
[2025-07-13 15:38:15] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/README.md
[2025-07-13 15:38:15] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:15] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/README.md
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/template.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/formatter_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/classification.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/errors_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/interfaces.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/handlers.go
[2025-07-13 15:38:15] [INFO] Processing file: /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/README.md
[2025-07-13 15:38:15] [BACKUP] Created backup: /Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/README.md
[2025-07-13 15:38:15] [INFO] No changes needed for /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/README.md
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/provider.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/config.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/interfaces.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/testutil_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/scenarios.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/assertions.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/filesystem.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/README.md
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/doc.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/cli.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/fixtures.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/integration_demo_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/main_test.go
[2025-07-13 15:38:15] [SKIP] No Unicode icons found in /Users/fareed/Documents/dev/go/bkpdir/errors.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/project-tokens.yaml
[2025-07-13 15:38:15] [ERROR] Found       12 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/project-tokens.yaml
[2025-07-13 15:38:15] [SUCCESS] Found       17 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/project-tokens.yaml
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/ui/config.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/generator/generator.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/templates/manager.go
[2025-07-13 15:38:15] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/templates/manager.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
[2025-07-13 15:38:15] [ERROR] Found        6 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
[2025-07-13 15:38:15] [SUCCESS] Found       10 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/main.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/create.go
[2025-07-13 15:38:15] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/create.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/config.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/version.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/process.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/root.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/default.yml
[2025-07-13 15:38:15] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/default.yml
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/example.yml
[2025-07-13 15:38:15] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/example.yml
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/README.md
[2025-07-13 15:38:15] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/README.md
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/examples/basic/README.md
[2025-07-13 15:38:15] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/examples/basic/README.md
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/main.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/realtime-validator/main.go
[2025-07-13 15:38:15] [ERROR] Found        5 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/realtime-validator/main.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/types.go
[2025-07-13 15:38:15] [SUCCESS] Found        8 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/types.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer.go
[2025-07-13 15:38:15] [ERROR] Found        6 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer.go
[2025-07-13 15:38:15] [SUCCESS] Found       10 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer_test.go
[2025-07-13 15:38:15] [SUCCESS] Found        5 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer_test.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
[2025-07-13 15:38:15] [ERROR] Found        9 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
[2025-07-13 15:38:15] [SUCCESS] Found        8 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
[2025-07-13 15:38:15] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/cmd/ai-validation/main.go
[2025-07-13 15:38:15] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/cmd/ai-validation/main.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/config.go
[2025-07-13 15:38:16] [SUCCESS] Found        3 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/config.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/formatter_adapter.go
[2025-07-13 15:38:16] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/formatter_adapter.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/tools/coverage-differential.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/tools/coverage_differential_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/errors_test.go
[2025-07-13 15:38:16] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/errors_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/backup.go
[2025-07-13 15:38:16] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/backup.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/binary-installation-instructions-working-plan.md
[2025-07-13 15:38:16] [ERROR] Found        3 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/binary-installation-instructions-working-plan.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/formatter.go
[2025-07-13 15:38:16] [SUCCESS] Found        4 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/formatter.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/refactoring_validation_test.go
[2025-07-13 15:38:16] [ERROR] Found       13 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/refactoring_validation_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/archive_test.go
[2025-07-13 15:38:16] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/archive_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-010.md
[2025-07-13 15:38:16] [ERROR] Found       26 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-010.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/test/metrics/decision_metrics_test.go
[2025-07-13 15:38:16] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/test/metrics/decision_metrics_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/test/metrics/metrics.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/test/metrics/compliance_rate_test.go
[2025-07-13 15:38:16] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/test/metrics/compliance_rate_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/test/metrics/goal_alignment_test.go
[2025-07-13 15:38:16] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/test/metrics/goal_alignment_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/test/integration/doc014_integration_test.go
[2025-07-13 15:38:16] [ERROR] Found        7 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/test/integration/doc014_integration_test.go
[2025-07-13 15:38:16] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/test/integration/doc014_integration_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/test/scenarios/decision_scenarios_test.go
[2025-07-13 15:38:16] [ERROR] Found        4 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/test/scenarios/decision_scenarios_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/test/performance/validation_performance_test.go
[2025-07-13 15:38:16] [ERROR] Found        4 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/test/performance/validation_performance_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/test/performance/performance_test_integration.go
[2025-07-13 15:38:16] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/test/performance/performance_test_integration.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/example-inheritance-base.yml
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/config_bench_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/exclude_test.go
[2025-07-13 15:38:16] [SUCCESS] Found        1 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/exclude_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-comprehensive-testing-plan.md
[2025-07-13 15:38:16] [ERROR] Found       17 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-comprehensive-testing-plan.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-performance-optimization-plan.md
[2025-07-13 15:38:16] [ERROR] Found       36 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/CFG-006-performance-optimization-plan.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/working-plan-package-interdependency-mapping.md
[2025-07-13 15:38:16] [ERROR] Found       17 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/working-plan-package-interdependency-mapping.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/example-inheritance-child.yml
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenarios.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/errorinjection_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/corruption_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/context_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenarios_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/permissions.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/context.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/errorinjection.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/corruption.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/scenario_helpers.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/diskspace.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/permissions_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/testutil/diskspace_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_validation.go
[2025-07-13 15:38:16] [ERROR] Found        4 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_validation.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
[2025-07-13 15:38:16] [ERROR] Found       11 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
[2025-07-13 15:38:16] [SUCCESS] Found        4 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_error_formatter.go
[2025-07-13 15:38:16] [ERROR] Found       10 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_error_formatter.go
[2025-07-13 15:38:16] [SUCCESS] Found        8 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_error_formatter.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/doc008_validator.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator_test.go
[2025-07-13 15:38:16] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/decision_checklist.go
[2025-07-13 15:38:16] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/internal/validation/decision_checklist.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/compliance_monitor.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/internal/validation/decision_checklist_test.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/config_impl.go
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-implementation-plan.md
[2025-07-13 15:38:16] [ERROR] Found        8 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/CFG-006-implementation-plan.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-completion-summary.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-task-completion-final.md
[2025-07-13 15:38:16] [ERROR] Found        5 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-task-completion-final.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-working-plan.md
[2025-07-13 15:38:16] [ERROR] Found        9 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-working-plan.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-test-layer-progress.md
[2025-07-13 15:38:16] [ERROR] Found        3 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-test-layer-progress.md
[2025-07-13 15:38:16] [SUCCESS] Found        1 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-test-layer-progress.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask2.md
[2025-07-13 15:38:16] [SUCCESS] Found        4 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask2.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-code-layer-progress.md
[2025-07-13 15:38:16] [ERROR] Found        3 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-code-layer-progress.md
[2025-07-13 15:38:16] [SUCCESS] Found        1 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-code-layer-progress.md
[2025-07-13 15:38:16] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-010-implementation-summary.md
[2025-07-13 15:38:16] [ERROR] Found       50 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-010-implementation-summary.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-012-completion-report.md
[2025-07-13 15:38:17] [ERROR] Found        3 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-012-completion-report.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/interface-definitions.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/doc-validation.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-training-examples.md
[2025-07-13 15:38:17] [ERROR] Found       35 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-training-examples.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/specification.md
[2025-07-13 15:38:17] [ERROR] Found       23 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/specification.md
[2025-07-13 15:38:17] [SUCCESS] Found        6 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/specification.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-optimization-tasks.md
[2025-07-13 15:38:17] [ERROR] Found       34 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-optimization-tasks.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-migration-strategy.md
[2025-07-13 15:38:17] [ERROR] Found       42 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-migration-strategy.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-documentation-templates.md
[2025-07-13 15:38:17] [ERROR] Found       95 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-documentation-templates.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/source-code-icon-guidelines.md
[2025-07-13 15:38:17] [ERROR] Found      164 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/source-code-icon-guidelines.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/requirements.md
[2025-07-13 15:38:17] [ERROR] Found        9 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/requirements.md
[2025-07-13 15:38:17] [SUCCESS] Found        1 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/requirements.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/task-completion-enforcement.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-token-system-comprehensive.md
[2025-07-13 15:38:17] [ERROR] Found       68 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-token-system-comprehensive.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-syntax-specification.md
[2025-07-13 15:38:17] [ERROR] Found       70 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-syntax-specification.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/architecture.md
[2025-07-13 15:38:17] [ERROR] Found       49 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/architecture.md
[2025-07-13 15:38:17] [SUCCESS] Found       31 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/architecture.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-formatter-refactoring-plan.md
[2025-07-13 15:38:17] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-formatter-refactoring-plan.md
[2025-07-13 15:38:17] [SUCCESS] Found       27 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-formatter-refactoring-plan.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/AI-First-Development-Procedure-Complete-Guide.md
[2025-07-13 15:38:17] [ERROR] Found       38 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/AI-First-Development-Procedure-Complete-Guide.md
[2025-07-13 15:38:17] [SUCCESS] Found       33 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/AI-First-Development-Procedure-Complete-Guide.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-decision-framework.md
[2025-07-13 15:38:17] [ERROR] Found       60 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-decision-framework.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-compliance.md
[2025-07-13 15:38:17] [ERROR] Found       95 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-compliance.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-status.md
[2025-07-13 15:38:17] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-status.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/CFG-005-implementation-plan.md
[2025-07-13 15:38:17] [ERROR] Found       27 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/CFG-005-implementation-plan.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/sync-framework.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-system-proposal.md
[2025-07-13 15:38:17] [ERROR] Found       71 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-system-proposal.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-code-maintenance-standards.md
[2025-07-13 15:38:17] [ERROR] Found      128 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-code-maintenance-standards.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-analysis.md
[2025-07-13 15:38:17] [ERROR] Found       39 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-analysis.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-troubleshooting.md
[2025-07-13 15:38:17] [ERROR] Found       38 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-troubleshooting.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/context-file-checklist.md
[2025-07-13 15:38:17] [ERROR] Found       23 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/context-file-checklist.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-development-strategy.md
[2025-07-13 15:38:17] [ERROR] Found      106 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-development-strategy.md
[2025-07-13 15:38:17] [SUCCESS] Found       33 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-development-strategy.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-token-protocol.md
[2025-07-13 15:38:17] [ERROR] Found       23 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-token-protocol.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/testing.md
[2025-07-13 15:38:17] [ERROR] Found       40 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/testing.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-guidelines.md
[2025-07-13 15:38:17] [ERROR] Found      109 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-guidelines.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/code-marker-strategy.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-change-protocol.md
[2025-07-13 15:38:17] [ERROR] Found        6 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-change-protocol.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/README.md
[2025-07-13 15:38:17] [ERROR] Found       57 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/README.md
[2025-07-13 15:38:17] [SUCCESS] Found        1 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/README.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-configuration-refactoring-plan.md
[2025-07-13 15:38:17] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-configuration-refactoring-plan.md
[2025-07-13 15:38:17] [SUCCESS] Found       27 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-configuration-refactoring-plan.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/validation-automation.md
[2025-07-13 15:38:17] [ERROR] Found       20 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/validation-automation.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-decisions.md
[2025-07-13 15:38:17] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-decisions.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/structure-optimization-analysis.md
[2025-07-13 15:38:17] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/structure-optimization-analysis.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/immutable.md
[2025-07-13 15:38:17] [SUCCESS] Found        3 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/immutable.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-protocol.md
[2025-07-13 15:38:17] [ERROR] Found      105 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-protocol.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-refactoring-priority-summary.md
[2025-07-13 15:38:17] [ERROR] Found        4 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-refactoring-priority-summary.md
[2025-07-13 15:38:17] [SUCCESS] Found       34 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-refactoring-priority-summary.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/change-rejection-criteria.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/cross-reference-template.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/enforcement-mechanisms.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/unicode-to-semantic-token-mapping.md
[2025-07-13 15:38:17] [ERROR] Found       35 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/unicode-to-semantic-token-mapping.md
[2025-07-13 15:38:17] [SUCCESS] Found        1 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/unicode-to-semantic-token-mapping.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/semantic-links-implementation.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-tracking.md
[2025-07-13 15:38:17] [ERROR] Found      505 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-tracking.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-documentation-standards.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-usage-guide.md
[2025-07-13 15:38:17] [ERROR] Found       45 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-usage-guide.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/context-file-responsibilities.md
[2025-07-13 15:38:17] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-traceability.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-error-handling-refactoring-plan.md
[2025-07-13 15:38:18] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-error-handling-refactoring-plan.md
[2025-07-13 15:38:18] [SUCCESS] Found       33 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-error-handling-refactoring-plan.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-validation-enforcement.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/semantic-token-validation-report-20250713_092636.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/token-traceability-report.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/semantic-token-validation-report-20250713_092234.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/semantic-token-validation-report-20250713_092247.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/icon-validation-report.md
[2025-07-13 15:38:18] [ERROR] Found        6 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/icon-validation-report.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-context-validation-report.md
[2025-07-13 15:38:18] [ERROR] Found        6 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-context-validation-report.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-framework-validation-report.md
[2025-07-13 15:38:18] [ERROR] Found        4 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-framework-validation-report.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-troubleshooting.md
[2025-07-13 15:38:18] [ERROR] Found       10 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-troubleshooting.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/coverage-implementation.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask3.md
[2025-07-13 15:38:18] [SUCCESS] Found        4 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask3.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase1-completion-summary.md
[2025-07-13 15:38:18] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase1-completion-summary.md
[2025-07-13 15:38:18] [SUCCESS] Found        4 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase1-completion-summary.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/tiered-performance-testing-strategy.md
[2025-07-13 15:38:18] [ERROR] Found        3 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/tiered-performance-testing-strategy.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-requirements.md
[2025-07-13 15:38:18] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-requirements.md
[2025-07-13 15:38:18] [SUCCESS] Found        6 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-requirements.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/migration-guide.md
[2025-07-13 15:38:18] [ERROR] Found       28 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/migration-guide.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/package-reference.md
[2025-07-13 15:38:18] [ERROR] Found       27 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/package-reference.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-completion-summary.md
[2025-07-13 15:38:18] [ERROR] Found        8 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-completion-summary.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/delayed-output-usage.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-3-completion.md
[2025-07-13 15:38:18] [ERROR] Found       18 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-3-completion.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-session-summary.md
[2025-07-13 15:38:18] [ERROR] Found        3 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-session-summary.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inspection-guide.md
[2025-07-13 15:38:18] [ERROR] Found        8 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inspection-guide.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week3-command-structure-plan.md
[2025-07-13 15:38:18] [ERROR] Found       25 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week3-command-structure-plan.md
[2025-07-13 15:38:18] [SUCCESS] Found       24 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week3-command-structure-plan.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/config-schema-abstraction.md
[2025-07-13 15:38:18] [ERROR] Found       11 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/config-schema-abstraction.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-009-completion-summary.md
[2025-07-13 15:38:18] [ERROR] Found       13 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-009-completion-summary.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-implementation-progress.md
[2025-07-13 15:38:18] [ERROR] Found       11 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-implementation-progress.md
[2025-07-13 15:38:18] [SUCCESS] Found        4 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-implementation-progress.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/milestone-tracking-system.md
[2025-07-13 15:38:18] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/milestone-tracking-system.md
[2025-07-13 15:38:18] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/milestone-tracking-system.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/build-system-fixes-implementation.md
[2025-07-13 15:38:18] [ERROR] Found       12 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/build-system-fixes-implementation.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/integration-guide.md
[2025-07-13 15:38:18] [ERROR] Found       63 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/integration-guide.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-cross-reference-completion-plan.md
[2025-07-13 15:38:18] [ERROR] Found       18 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-cross-reference-completion-plan.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/formatter-decomposition.md
[2025-07-13 15:38:18] [ERROR] Found       12 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/formatter-decomposition.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-remaining-subtasks-plan.md
[2025-07-13 15:38:18] [ERROR] Found       43 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-remaining-subtasks-plan.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-5-completion.md
[2025-07-13 15:38:18] [ERROR] Found        3 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-5-completion.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/config-abstraction.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-progress-summary.md
[2025-07-13 15:38:18] [ERROR] Found        7 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-progress-summary.md
[2025-07-13 15:38:18] [SUCCESS] Found       16 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-progress-summary.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-1-completion.md
[2025-07-13 15:38:18] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-1-completion.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week2-completion-summary.md
[2025-07-13 15:38:18] [ERROR] Found       11 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week2-completion-summary.md
[2025-07-13 15:38:18] [SUCCESS] Found       18 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week2-completion-summary.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extraction-dependencies.md
[2025-07-13 15:38:18] [ERROR] Found       21 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extraction-dependencies.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-implementation-complete.md
[2025-07-13 15:38:18] [ERROR] Found        4 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-implementation-complete.md
[2025-07-13 15:38:18] [SUCCESS] Found       12 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-implementation-complete.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-final-analysis.md
[2025-07-13 15:38:18] [ERROR] Found       11 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-final-analysis.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-implementation-plan.md
[2025-07-13 15:38:18] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-implementation-plan.md
[2025-07-13 15:38:18] [SUCCESS] Found        7 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-implementation-plan.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask3-token-implementation.md
[2025-07-13 15:38:18] [ERROR] Found        2 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask3-token-implementation.md
[2025-07-13 15:38:18] [SUCCESS] Found       33 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask3-token-implementation.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/coverage-baseline.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-5-working-plan.md
[2025-07-13 15:38:18] [ERROR] Found        8 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-5-working-plan.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask1-missing-token-inventory.md
[2025-07-13 15:38:18] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask1-missing-token-inventory.md
[2025-07-13 15:38:18] [SUCCESS] Found        4 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask1-missing-token-inventory.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/config.yaml
[2025-07-13 15:38:18] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/config.yaml
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/README.md
[2025-07-13 15:38:18] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/README.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/example-code/main.go
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/examples/README.md
[2025-07-13 15:38:18] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/examples/README.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/integration-documentation-working-plan.md
[2025-07-13 15:38:18] [ERROR] Found       18 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/integration-documentation-working-plan.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-examples.md
[2025-07-13 15:38:18] [ERROR] Found       12 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-examples.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/ai-first-token-system-proposal.md
[2025-07-13 15:38:18] [ERROR] Found        5 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/ai-first-token-system-proposal.md
[2025-07-13 15:38:18] [SUCCESS] Found       27 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/ai-first-token-system-proposal.md
[2025-07-13 15:38:18] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/task-config-template-working-plan.md
[2025-07-13 15:38:18] [ERROR] Found        8 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/task-config-template-working-plan.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-cli-template-plan.md
[2025-07-13 15:38:19] [ERROR] Found       15 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-cli-template-plan.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/fix.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/document-completed-task.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/repetition.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/import-bkpfile.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/context-documents.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/next-task.md
[2025-07-13 15:38:19] [ERROR] Found        3 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/next-task.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/adhoc-new-task.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/add-task.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/existing-task.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/features.md
[2025-07-13 15:38:19] [ERROR] Found        8 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/features.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/format-strings.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/prompts/continue-task.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-2-completion.md
[2025-07-13 15:38:19] [ERROR] Found       25 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-2-completion.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/advanced-patterns.md
[2025-07-13 15:38:19] [ERROR] Found       19 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/advanced-patterns.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/troubleshooting.md
[2025-07-13 15:38:19] [ERROR] Found       15 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/troubleshooting.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/getting-started.md
[2025-07-13 15:38:19] [ERROR] Found       11 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/getting-started.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-completion-report.md
[2025-07-13 15:38:19] [ERROR] Found       13 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-completion-report.md
[2025-07-13 15:38:19] [SUCCESS] Found       12 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-completion-report.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inheritance.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-working-plan.md
[2025-07-13 15:38:19] [ERROR] Found        3 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-working-plan.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-quick-start-guide.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-011-implementation-summary.md
[2025-07-13 15:38:19] [ERROR] Found       14 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/DOC-011-implementation-summary.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask4.md
[2025-07-13 15:38:19] [SUCCESS] Found        4 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask4.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-phase2f-summary.md
[2025-07-13 15:38:19] [ERROR] Found       10 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-phase2f-summary.md
[2025-07-13 15:38:19] [SUCCESS] Found       20 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-phase2f-summary.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/system-comparison-analysis.md
[2025-07-13 15:38:19] [ERROR] Found       16 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/system-comparison-analysis.md
[2025-07-13 15:38:19] [SUCCESS] Found        6 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/system-comparison-analysis.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/package-interdependency-mapping.md
[2025-07-13 15:38:19] [ERROR] Found       47 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/package-interdependency-mapping.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week4-core-files-plan.md
[2025-07-13 15:38:19] [ERROR] Found       17 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week4-core-files-plan.md
[2025-07-13 15:38:19] [SUCCESS] Found       23 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week4-core-files-plan.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-creation-strategy.md
[2025-07-13 15:38:19] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-creation-strategy.md
[2025-07-13 15:38:19] [SUCCESS] Found       17 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-creation-strategy.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-critical-tokens-completion-summary.md
[2025-07-13 15:38:19] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-critical-tokens-completion-summary.md
[2025-07-13 15:38:19] [SUCCESS] Found       33 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-critical-tokens-completion-summary.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask1.md
[2025-07-13 15:38:19] [SUCCESS] Found        4 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask1.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/README.md
[2025-07-13 15:38:19] [SUCCESS] Found       15 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/README.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/archive.go
[2025-07-13 15:38:19] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/archive.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-completion-summary.md
[2025-07-13 15:38:19] [ERROR] Found        5 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-completion-summary.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-009.md
[2025-07-13 15:38:19] [ERROR] Found       27 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-009.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/file_stats_test.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/task-analysis-and-plan.md
[2025-07-13 15:38:19] [ERROR] Found        4 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/task-analysis-and-plan.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-008.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/exclude.go
[2025-07-13 15:38:19] [SUCCESS] Found        3 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/exclude.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/config_adapter.go
[2025-07-13 15:38:19] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/config_adapter.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/file_stats.go
[2025-07-13 15:38:19] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/file_stats.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/test-dir/.bkpdir.yml
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/comparison.go
[2025-07-13 15:38:19] [SUCCESS] Found        3 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/comparison.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/verify_test.go
[2025-07-13 15:38:19] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/verify_test.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/DOC-014-subtask-4-working-plan.md
[2025-07-13 15:38:19] [ERROR] Found       23 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/DOC-014-subtask-4-working-plan.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/refactoring-validation-report.md
[2025-07-13 15:38:19] [ERROR] Found       21 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/refactoring-validation-report.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/config_interfaces.go
[2025-07-13 15:38:19] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/config_interfaces.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/working-plan-remove-backup-command.md
[2025-07-13 15:38:19] [ERROR] Found       20 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/working-plan-remove-backup-command.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/config_test.go
[2025-07-13 15:38:19] [SUCCESS] Found        3 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/config_test.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/formatter_test.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/DEVELOPMENT.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-8-documentation-plan.md
[2025-07-13 15:38:19] [ERROR] Found       25 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-8-documentation-plan.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/SPECIFICATION_AUDIT_REPORT.md
[2025-07-13 15:38:19] [ERROR] Found       44 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/SPECIFICATION_AUDIT_REPORT.md
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/main.go
[2025-07-13 15:38:19] [SUCCESS] Found        5 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/main.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/example-symlink-config.yml
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/.bkpdir.yml
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/backup_test.go
[2025-07-13 15:38:19] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/backup_test.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/verify.go
[2025-07-13 15:38:19] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/verify.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/comparison_test.go
[2025-07-13 15:38:19] [SUCCESS] Found        1 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/comparison_test.go
[2025-07-13 15:38:19] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/config/discovery.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/config/interfaces.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/config/README.md
[2025-07-13 15:38:20] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/pkg/config/README.md
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/config/merge_strategies.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/config/loader.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/config/utils.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/config/inheritance.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/config/config_test.go
[2025-07-13 15:38:20] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/pkg/config/config_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/verification.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/concurrent.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/processor_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/processor.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/naming.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/README.md
[2025-07-13 15:38:20] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/README.md
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/pipeline.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/processing/doc.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/fileops.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/traversal.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/README.md
[2025-07-13 15:38:20] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/README.md
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/exclusion.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/comparison.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/validation.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/atomic.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/resources/interfaces.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/resources/resources_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/resources/README.md
[2025-07-13 15:38:20] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/pkg/resources/README.md
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/resources/context.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/version.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/example_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/types.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/flags.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/README.md
[2025-07-13 15:38:20] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/README.md
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/context.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/builder.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/dryrun.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/cli_test.go
[2025-07-13 15:38:20] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/pkg/cli/cli_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/patterns.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/formatter.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/interfaces.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/collector.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/filestats.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/README.md
[2025-07-13 15:38:20] [ERROR] Found       11 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/README.md
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/template.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/formatter_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/classification.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/errors_test.go
[2025-07-13 15:38:20] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/errors_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/interfaces.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/handlers.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/README.md
[2025-07-13 15:38:20] [ERROR] Found        1 remaining Unicode icons in /Users/fareed/Documents/dev/go/bkpdir/pkg/errors/README.md
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/provider.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/config.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/interfaces.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/testutil_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/scenarios.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/assertions.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/filesystem.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/README.md
[2025-07-13 15:38:20] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/README.md
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/doc.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/cli.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/fixtures.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/integration_demo_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/main_test.go
[2025-07-13 15:38:20] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/main_test.go
[2025-07-13 15:38:20] [VALIDATE] Validating migration for /Users/fareed/Documents/dev/go/bkpdir/errors.go
[2025-07-13 15:38:20] [SUCCESS] Found        2 semantic tokens in /Users/fareed/Documents/dev/go/bkpdir/errors.go
```

## Validation Results

### Files with Semantic Tokens
```
/Users/fareed/Documents/dev/go/bkpdir/project-tokens.yaml
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/types.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer_test.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
/Users/fareed/Documents/dev/go/bkpdir/config.go
/Users/fareed/Documents/dev/go/bkpdir/formatter_adapter.go
/Users/fareed/Documents/dev/go/bkpdir/errors_test.go
/Users/fareed/Documents/dev/go/bkpdir/backup.go
/Users/fareed/Documents/dev/go/bkpdir/git.go
/Users/fareed/Documents/dev/go/bkpdir/formatter.go
/Users/fareed/Documents/dev/go/bkpdir/archive_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/integration/doc014_integration_test.go
/Users/fareed/Documents/dev/go/bkpdir/exclude_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_error_formatter.go
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-test-layer-progress.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask2.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-code-layer-progress.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/specification.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/requirements.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/architecture.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-formatter-refactoring-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/AI-First-Development-Procedure-Complete-Guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-development-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-configuration-refactoring-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/immutable.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-refactoring-priority-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/unicode-to-semantic-token-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-error-handling-refactoring-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask3.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase1-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-requirements.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week3-command-structure-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-implementation-progress.md
/Users/fareed/Documents/dev/go/bkpdir/docs/milestone-tracking-system.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-progress-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week2-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-implementation-complete.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-implementation-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask3-token-implementation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask1-missing-token-inventory.md
/Users/fareed/Documents/dev/go/bkpdir/docs/ai-first-token-system-proposal.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-completion-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask4.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-phase2f-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/system-comparison-analysis.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week4-core-files-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-creation-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-critical-tokens-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-audit-report-phase1-subtask1.md
/Users/fareed/Documents/dev/go/bkpdir/README.md
/Users/fareed/Documents/dev/go/bkpdir/archive.go
/Users/fareed/Documents/dev/go/bkpdir/exclude.go
/Users/fareed/Documents/dev/go/bkpdir/config_adapter.go
/Users/fareed/Documents/dev/go/bkpdir/file_stats.go
/Users/fareed/Documents/dev/go/bkpdir/git_test.go
/Users/fareed/Documents/dev/go/bkpdir/comparison.go
/Users/fareed/Documents/dev/go/bkpdir/verify_test.go
/Users/fareed/Documents/dev/go/bkpdir/config_interfaces.go
/Users/fareed/Documents/dev/go/bkpdir/config_test.go
/Users/fareed/Documents/dev/go/bkpdir/main.go
/Users/fareed/Documents/dev/go/bkpdir/backup_test.go
/Users/fareed/Documents/dev/go/bkpdir/verify.go
/Users/fareed/Documents/dev/go/bkpdir/comparison_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/config_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/fileops/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/resources/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/cli_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/errors_test.go
/Users/fareed/Documents/dev/go/bkpdir/pkg/testutil/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/git/git_test.go
/Users/fareed/Documents/dev/go/bkpdir/main_test.go
/Users/fareed/Documents/dev/go/bkpdir/errors.go
```

### Files with Remaining Unicode Icons
```
/Users/fareed/Documents/dev/go/bkpdir/project-tokens.yaml
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/internal/templates/manager.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/scaffolding/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/cmd/create.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/default.yml
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/config/example.yml
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/cli-template/examples/basic/README.md
/Users/fareed/Documents/dev/go/bkpdir/cmd/realtime-validator/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/analyzer.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/token-suggester/main.go
/Users/fareed/Documents/dev/go/bkpdir/cmd/ai-validation/main.go
/Users/fareed/Documents/dev/go/bkpdir/binary-installation-instructions-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/refactoring_validation_test.go
/Users/fareed/Documents/dev/go/bkpdir/example-git-config.yml
/Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-010.md
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/decision_metrics_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/compliance_rate_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/metrics/goal_alignment_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/integration/doc014_integration_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/scenarios/decision_scenarios_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/performance/validation_performance_test.go
/Users/fareed/Documents/dev/go/bkpdir/test/performance/performance_test_integration.go
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-comprehensive-testing-plan.md
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-performance-optimization-plan.md
/Users/fareed/Documents/dev/go/bkpdir/working-plan-package-interdependency-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_validation.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/ai_error_formatter.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/realtime_validator_test.go
/Users/fareed/Documents/dev/go/bkpdir/internal/validation/decision_checklist.go
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-implementation-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-task-completion-final.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-test-layer-progress.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-code-layer-progress.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-010-implementation-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-012-completion-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-training-examples.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/specification.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-optimization-tasks.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-migration-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-documentation-templates.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/source-code-icon-guidelines.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/requirements.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-token-system-comprehensive.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/enhanced-token-syntax-specification.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/architecture.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-formatter-refactoring-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/AI-First-Development-Procedure-Complete-Guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-decision-framework.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-compliance.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-status.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/CFG-005-implementation-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-system-proposal.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-code-maintenance-standards.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-analysis.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-troubleshooting.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/context-file-checklist.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-development-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-token-protocol.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/testing.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/icon-usage-guidelines.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-change-protocol.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-configuration-refactoring-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/validation-automation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/implementation-decisions.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/structure-optimization-analysis.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-assistant-protocol.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-refactoring-priority-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/unicode-to-semantic-token-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/feature-tracking.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/decision-framework-usage-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/context/ai-first-error-handling-refactoring-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/icon-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-context-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/validation-reports/decision-framework-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-troubleshooting.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase1-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tiered-performance-testing-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-requirements.md
/Users/fareed/Documents/dev/go/bkpdir/docs/migration-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/package-reference.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-6-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-3-completion.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-session-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-inspection-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week3-command-structure-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/config-schema-abstraction.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-009-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-implementation-progress.md
/Users/fareed/Documents/dev/go/bkpdir/docs/milestone-tracking-system.md
/Users/fareed/Documents/dev/go/bkpdir/docs/build-system-fixes-implementation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/integration-guide.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-cross-reference-completion-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/formatter-decomposition.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-remaining-subtasks-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-5-completion.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-progress-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-1-completion.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week2-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extraction-dependencies.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-system-implementation-complete.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-final-analysis.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-implementation-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask3-token-implementation.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-014-subtask-5-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask1-missing-token-inventory.md
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/git-aware-backup/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/config.yaml
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/basic-cli-app/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/examples/README.md
/Users/fareed/Documents/dev/go/bkpdir/docs/integration-documentation-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/configuration-examples.md
/Users/fareed/Documents/dev/go/bkpdir/docs/ai-first-token-system-proposal.md
/Users/fareed/Documents/dev/go/bkpdir/docs/task-config-template-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-cli-template-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/next-task.md
/Users/fareed/Documents/dev/go/bkpdir/docs/prompts/features.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-008-subtask-2-completion.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/advanced-patterns.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/troubleshooting.md
/Users/fareed/Documents/dev/go/bkpdir/docs/tutorials/getting-started.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-completion-report.md
/Users/fareed/Documents/dev/go/bkpdir/docs/extract-007-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/DOC-011-implementation-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/semantic-token-migration-phase2f-summary.md
/Users/fareed/Documents/dev/go/bkpdir/docs/system-comparison-analysis.md
/Users/fareed/Documents/dev/go/bkpdir/docs/package-interdependency-mapping.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-week4-core-files-plan.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-subtask2-token-creation-strategy.md
/Users/fareed/Documents/dev/go/bkpdir/docs/unicode-migration-phase2-critical-tokens-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-7-completion-summary.md
/Users/fareed/Documents/dev/go/bkpdir/working-plan-extract-009.md
/Users/fareed/Documents/dev/go/bkpdir/task-analysis-and-plan.md
/Users/fareed/Documents/dev/go/bkpdir/DOC-014-subtask-4-working-plan.md
/Users/fareed/Documents/dev/go/bkpdir/refactoring-validation-report.md
/Users/fareed/Documents/dev/go/bkpdir/working-plan-remove-backup-command.md
/Users/fareed/Documents/dev/go/bkpdir/CFG-006-subtask-8-documentation-plan.md
/Users/fareed/Documents/dev/go/bkpdir/SPECIFICATION_AUDIT_REPORT.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/config/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/processing/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/cli/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/formatter/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/errors/README.md
/Users/fareed/Documents/dev/go/bkpdir/pkg/git/README.md
```

## Next Steps

1. **Review Migration**: Check all migrated files for accuracy
2. **Update Documentation**: Update any documentation referencing old tokens
3. **Update Tests**: Ensure all tests pass with new semantic tokens
4. **Update Validation**: Update any validation scripts to use semantic tokens
5. **AI Assistant Integration**: Update AI assistant protocols to use semantic tokens

## Backup Location

All original files have been backed up to: `/Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809`

## Rollback Instructions

To rollback the migration:

```bash
# Restore from backup
cp -r "/Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809"/* "/Users/fareed/Documents/dev/go/bkpdir"/

# Or restore specific files
cp "/Users/fareed/Documents/dev/go/bkpdir/backup-semantic-tokens-20250713-153809/filename.go" "/Users/fareed/Documents/dev/go/bkpdir/"
```
