# Bulk package-level lead audit

Suspect files: **52**

| path | package_leads | action | check_leads_pool |
|------|---------------|--------|------------------|
| `pkg/cli/example_test.go` | 175 | strip_aggregate_keep_prefix | True |
| `pkg/fileops/comparison_test.go` | 173 | strip_aggregate_keep_prefix | True |
| `pkg/fileops/exclusion_test.go` | 173 | strip_aggregate_keep_prefix | True |
| `pkg/config/config_test.go` | 165 | strip_aggregate_keep_prefix | True |
| `pkg/resources/examples_test.go` | 156 | strip_aggregate_keep_prefix | True |
| `pkg/git/git_test.go` | 155 | strip_aggregate_keep_prefix | True |
| `pkg/resources/resources_test.go` | 150 | strip_aggregate_keep_prefix | True |
| `pkg/testutil/integration_demo_test.go` | 149 | strip_all | True |
| `pkg/testutil/testutil_test.go` | 149 | strip_all | True |
| `test/metrics/compliance_rate_test.go` | 145 | strip_all | False |
| `test/metrics/decision_metrics_test.go` | 145 | strip_all | False |
| `test/metrics/goal_alignment_test.go` | 145 | strip_all | False |
| `comparison_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `config_bench_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `config_grouping_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `config_integration_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `diff_command_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `exclude_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `file_stats_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `formatter_adapter_simple_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `formatter_stats_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `formatter_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `git_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `inc_diff_integration_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `incremental_duplicate_prevention_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `internal/testutil/context_test.go` | 144 | strip_all | False |
| `internal/testutil/corruption_test.go` | 144 | strip_all | False |
| `internal/testutil/diskspace_test.go` | 144 | strip_all | False |
| `internal/testutil/errorinjection_test.go` | 144 | strip_all | False |
| `internal/testutil/permissions_test.go` | 144 | strip_all | False |
| `internal/testutil/scenarios_test.go` | 144 | strip_all | False |
| `internal/validation/realtime_validator_test.go` | 144 | strip_all | False |
| `pkg/cli/cli_test.go` | 144 | strip_all | True |
| `pkg/formatter/ai_first_formatter_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `pkg/formatter/formatter_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `pkg/formatter/placeholder_replace_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `pkg/processing/processor_test.go` | 144 | strip_all | True |
| `refactoring_validation_test.go` | 144 | strip_aggregate_keep_prefix | True |
| `test/performance/validation_performance_test.go` | 144 | strip_all | False |
| `test/scenarios/decision_scenarios_test.go` | 144 | strip_all | False |
| `tools/coverage_differential_test.go` | 144 | strip_all | False |
| `adapter_formatters_test.go` | 138 | strip_aggregate_keep_prefix | True |
| `backup_test.go` | 138 | strip_aggregate_keep_prefix | True |
| `formatter_adapter_test.go` | 138 | strip_aggregate_keep_prefix | True |
| `formatter_list_safety_test.go` | 138 | strip_aggregate_keep_prefix | True |
| `archive_test.go` | 135 | strip_aggregate_keep_prefix | True |
| `config_test.go` | 129 | strip_aggregate_keep_prefix | True |
| `errors_test.go` | 126 | strip_aggregate_keep_prefix | True |
| `pkg/errors/errors_test.go` | 126 | strip_all | True |
| `main_cli_composition_test.go` | 125 | strip_aggregate_keep_prefix | True |
| `main_test.go` | 125 | strip_aggregate_keep_prefix | True |
| `scripts/testdata/lead_hygiene/pkg_config_prefix.go` | 2 | keep | False |
