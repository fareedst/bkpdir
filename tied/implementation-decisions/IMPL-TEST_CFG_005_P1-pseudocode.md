# [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION]

## Summary contract

Priority-1 configuration inheritance and merge edge-case tests beyond P0 coverage: multi-parent inherit, paths, errors, deep chains, type mismatches, nulls, and whitespace strings.

INPUT: YAML fixture trees under temp dirs, BKPDIR_CONFIG / LoadConfigWithInheritance
OUTPUT: assertions on merged Config fields and expected errors
DATA: inherit directives, merge prefixes, defaults

## TESTMULTIPLEINHERITANCESOURCES

SPEC-ID: IMPL-TEST_CFG_005_P1::TESTMULTIPLEINHERITANCESOURCES
STEP T001: child inheriting parent1 and parent2 merges accumulate exclude_patterns from both parents plus child + prefix while child overrides precedence archive_dir_path

- [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: child inheriting parent1 and parent2 merges accumulate exclude_patterns from both parents plus child + prefix while child overrides precedence archive_dir_path.

PROCEDURE test_multiple_inheritance_sources():
  SETUP:
    CREATE parent1 with exclude_patterns [".git", "node_modules"] and archive_dir_path "/archives/p1"
    CREATE parent2 with exclude_patterns ["tmp"] and archive_dir_path "/archives/p2"
    CREATE child with inherit [parent1, parent2], merge-prefix exclude_patterns "+", and archive_dir_path "/archives/child"
  ACT:
    LOAD child configuration through inheritance resolver
  ASSERT:
    VERIFY merged exclude_patterns equal [".git", "node_modules", "tmp", child additions]
    VERIFY archive_dir_path equals "/archives/child"
    VERIFY no duplicate or reordered parent-first merge regressions occur

## TESTRELATIVEPATHINHERITANCE

SPEC-ID: IMPL-TEST_CFG_005_P1::TESTRELATIVEPATHINHERITANCE
STEP T001: inherit directives using ../base.yml and ./sibling.yml resolve relative to the child config file directory

- [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: inherit directives using ../base.yml and ./sibling.yml resolve relative to the child config file directory.

PROCEDURE test_relative_path_inheritance():
  SETUP:
    CREATE directory tree with child at "configs/env/child.yml"
    CREATE parent files at "../base.yml" and "./sibling.yml" relative to child location
    POPULATE parent files with distinct marker fields for traceable merge verification
  ACT:
    LOAD child configuration using relative inherit directives
  ASSERT:
    VERIFY resolver locates parent files relative to child directory, not process working directory
    VERIFY marker fields from both parents appear in merged result
    VERIFY path resolution is deterministic across repeated runs

## TESTHOMEDIRECTORYEXPANSION

SPEC-ID: IMPL-TEST_CFG_005_P1::TESTHOMEDIRECTORYEXPANSION
STEP T001: tilde in inherit path expands to user home directory and loads the referenced parent config

- [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: tilde in inherit path expands to user home directory and loads the referenced parent config.

PROCEDURE test_home_directory_expansion():
  SETUP:
    CREATE parent config file under simulated user home directory
    CREATE child config whose inherit path starts with "~/" and targets that parent
    DEFINE expected merged field sourced only from home-based parent
  ACT:
    LOAD child configuration with home-prefixed inherit directive
  ASSERT:
    VERIFY "~/" resolves to the active user home directory
    VERIFY parent config is loaded from expanded path
    VERIFY merged output contains expected parent field without fallback path errors

## TESTMISSINGINHERITANCEFILE

SPEC-ID: IMPL-TEST_CFG_005_P1::TESTMISSINGINHERITANCEFILE
STEP T001: missing inherit target returns a clear error without partial silent merge

- [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: missing inherit target returns a clear error without partial silent merge.

PROCEDURE test_missing_inheritance_file():
  SETUP:
    CREATE child config referencing one existing parent and one missing parent path
    DEFINE baseline child field that should not be reported as fully merged on failure
  ACT:
    ATTEMPT to load child configuration
  ASSERT:
    VERIFY load operation returns an explicit "missing inheritance file" error
    VERIFY error identifies the unresolved path for diagnosis
    VERIFY no partial merged config is returned as success

## TESTINVALIDYAMLHANDLING

SPEC-ID: IMPL-TEST_CFG_005_P1::TESTINVALIDYAMLHANDLING
STEP T001: malformed YAML in the inheritance chain surfaces parse error and does not produce a merged config

- [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: malformed YAML in the inheritance chain surfaces parse error and does not produce a merged config.

PROCEDURE test_invalid_yaml_handling():
  SETUP:
    CREATE parent config containing malformed YAML syntax
    CREATE child config that inherits from malformed parent
    PREPARE assertion target for parser diagnostic message
  ACT:
    ATTEMPT to load child configuration
  ASSERT:
    VERIFY parser error is returned and classified as invalid YAML input
    VERIFY diagnostic includes parent file context
    VERIFY merged configuration is not produced on parse failure

## TESTDEEPINHERITANCECHAIN

SPEC-ID: IMPL-TEST_CFG_005_P1::TESTDEEPINHERITANCECHAIN
STEP T001: twelve-level inherit chain still loads and applies child overrides at the leaf

- [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: twelve-level inherit chain still loads and applies child overrides at the leaf.

PROCEDURE test_deep_inheritance_chain():
  SETUP:
    CREATE 12-level inheritance chain where each level adds one traceable field
    CREATE leaf child that overrides one inherited value
    DEFINE expected final map including all inherited additions and leaf override
  ACT:
    LOAD leaf configuration through full chain
  ASSERT:
    VERIFY all 12 parent contributions are present in merged output
    VERIFY leaf override takes precedence over inherited value
    VERIFY load completes within acceptable recursion and cycle-safety limits

## TESTTYPEMISMATCHHANDLING

SPEC-ID: IMPL-TEST_CFG_005_P1::TESTTYPEMISMATCHHANDLING
STEP T001: parent array with child scalar (or incompatible types) errors or rejects merge per CFG-005 rules

- [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: parent array with child scalar (or incompatible types) errors or rejects merge per CFG-005 rules.

PROCEDURE test_type_mismatch_handling():
  SETUP:
    CREATE parent config with list value for target field
    CREATE child config with scalar value for same field under merge-inherit path
    DEFINE expected mismatch policy outcome per CFG-005
  ACT:
    LOAD child configuration
  ASSERT:
    VERIFY mismatch triggers explicit type-conflict handling
    VERIFY resolver either rejects merge or applies configured fallback exactly as specified
    VERIFY no silent coercion changes the value type unexpectedly

## TESTNULLVALUEHANDLING

SPEC-ID: IMPL-TEST_CFG_005_P1::TESTNULLVALUEHANDLING
STEP T001: explicit null YAML values merge without panics and preserve intended field nil semantics

- [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: explicit null YAML values merge without panics and preserve intended field nil semantics.

PROCEDURE test_null_value_handling():
  SETUP:
    CREATE parent config with non-null values for nullable fields
    CREATE child config that sets selected inherited fields to null
    DEFINE expected nil semantics for merged result
  ACT:
    LOAD child configuration and execute merge
  ASSERT:
    VERIFY merge completes without runtime panic
    VERIFY null-designated fields remain null in final config
    VERIFY non-null unrelated fields continue to merge normally

## TESTWHITESPACESTRINGHANDLING

SPEC-ID: IMPL-TEST_CFG_005_P1::TESTWHITESPACESTRINGHANDLING
STEP T001: whitespace-only and padded string config values load and merge without unintended trimming unless specified

- [IMPL-TEST_CFG_005_P1] [ARCH-CFG_005] [REQ-CFG_001] [REQ-CFG_005] [REQ-CONFIGURATION] — How: whitespace-only and padded string config values load and merge without unintended trimming unless specified.

PROCEDURE test_whitespace_string_handling():
  SETUP:
    CREATE parent and child configs with whitespace-only strings and padded strings
    DEFINE fields with explicit trim policy and fields with preserve policy
    PREPARE expected merged string values for each policy case
  ACT:
    LOAD child configuration with inheritance and merge behavior
  ASSERT:
    VERIFY preserve-policy fields retain leading/trailing whitespace
    VERIFY trim-policy fields are normalized only where policy requires
    VERIFY whitespace-only strings are retained or normalized exactly per field policy
