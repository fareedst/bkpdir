# [IMPL-CONFIG_SCHEMA_FLEX] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION]

## Summary contract

Schema-agnostic configuration loading in pkg/config clones defaults, merges discovered YAML files, applies environment overrides, validates, and adapts results back to the legacy Config type.

INPUT: root path, defaultConfig interface{}
OUTPUT: merged configuration value or map of ConfigValue with sources
DATA: GenericConfigLoader, search paths, env provider, validator

## GENERIC_CONFIG_LOAD

SPEC-ID: IMPL-CONFIG_SCHEMA_FLEX::GENERIC_CONFIG_LOAD
STEP T001: CLONE default and merge each existing search-path YAML file
STEP T002: APPLY environment overrides and VALIDATE schema

- [IMPL-CONFIG_SCHEMA_FLEX] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: clone defaultConfig, merge each existing search-path YAML file, apply environment overrides, then validate schema.

PROCEDURE GENERIC_CONFIG_LOAD(root, defaultConfig):
  config = CLONE(defaultConfig)
  FOR EACH path IN GetConfigSearchPaths():
    expanded = ExpandPath(path)
    IF FileExists(expanded) THEN loadConfigFromFile(config, expanded)  # non-fatal per-file errors
  APPLY environment overrides to config
  VALIDATE config WHEN validator present
  RETURN config

## CONFIG_ADAPTER_LOAD

SPEC-ID: IMPL-CONFIG_SCHEMA_FLEX::CONFIG_ADAPTER_LOAD
STEP T001: DELEGATE to GenericConfigLoader and type-assert to Config

- [IMPL-CONFIG_SCHEMA_FLEX] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: delegate to GenericConfigLoader with DefaultConfig() then type-assert to *Config or fall back to defaults.

PROCEDURE CONFIG_ADAPTER_LOAD(root):
  genericCfg = loader.LoadConfig(root, DefaultConfig())
  IF genericCfg is *Config THEN RETURN genericCfg
  RETURN DefaultConfig()
