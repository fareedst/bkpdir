# [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [ARCH-CFG_006] [REQ-CONFIGURATION] [REQ-CFG_006]

## Summary contract

Central configuration model: YAML-serializable Config aggregate, default factory, and legacy display adapters. Reflection field discovery and flat display are owned by [IMPL-CFG_006] and [IMPL-CONFIG_DISPLAY_FLATTENING].

INPUT: Config instance, optional filesystem root for source detection
OUTPUT: ConfigValue slices for simple and enhanced display
DATA: ConfigValue, configFieldInfo (shared type with CFG_006)

## CONFIG

SPEC-ID: IMPL-CONFIG_STRUCT::CONFIG
STEP T001: DEFINE aggregate fields for archive, backup, status, format, git, inheritance

- [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: define aggregate holding archive, backup, status, format, template, pattern, git, and inheritance settings for serialization and reflection.

PROCEDURE CONFIG aggregate:
  HOLDS archive paths, exclusion patterns, git nested config, backup paths
  HOLDS status codes for directory and file operations
  HOLDS printf and template format strings for operations and errors
  HOLDS regex patterns for archive, backup, config lines, timestamps
  HOLDS inherit list for configuration inheritance

## CONFIGVALUE

SPEC-ID: IMPL-CONFIG_STRUCT::CONFIGVALUE
STEP T001: DEFINE displayed entry with name, value, and source label

- [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: represent one displayed configuration entry with name, string value, and source label.

DATA ConfigValue:
  Name string
  Value string
  Source string

## CONFIGFIELDINFO

SPEC-ID: IMPL-CONFIG_STRUCT::CONFIGFIELDINFO
STEP T001: DEFINE reflection metadata for one config field

- [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: capture reflection metadata for one config field including YAML name, type, path, category, and importance.

DATA configFieldInfo:
  Name, YAMLName, Type, Path, Category string
  Kind type descriptor
  Value current value
  IsPointer, IsSlice, IsStruct bool
  Importance int

## DEFAULTCONFIG

SPEC-ID: IMPL-CONFIG_STRUCT::DEFAULTCONFIG
STEP T001: RETURN Config with preset defaults and nested git config

- [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: return new Config populated with preset defaults for all fields including nested git config.

PROCEDURE DEFAULTCONFIG():
Contract:
PRE: true
POST: true
INPUT: context for DEFAULTCONFIG
OUTPUT: result of DEFAULTCONFIG
EFFECTS: State.Config, pure
  RETURN Config with default archive/backup paths, status codes, format strings, templates, patterns
  INITIALIZE Git from DEFAULTGITCONFIG

## GETCONFIGVALUES

SPEC-ID: IMPL-CONFIG_STRUCT::GETCONFIGVALUES
STEP T001: BUILD ConfigValue list for key fields with Source config

- [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: expose subset of key fields as ConfigValue rows for simple --config display.

PROCEDURE GETCONFIGVALUES(cfg):
Contract:
PRE: true
POST: true
INPUT: cfg
OUTPUT: result of GETCONFIGVALUES
EFFECTS: State.Config, pure
  RETURN list of ConfigValue for archive_dir_path, use_current_dir_name, include_git_info, backup paths
  MARK each Source as config

## GETCONFIGVALUESWITHSOURCES

SPEC-ID: IMPL-CONFIG_STRUCT::GETCONFIGVALUESWITHSOURCES
STEP T001: DELEGATE to GetAllConfigValuesWithSources and convert to legacy format

- [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] [REQ-CFG_006] — How: delegate to reflection-based GetAllConfigValuesWithSources and convert to legacy ConfigValue format sorted by name.

PROCEDURE GETCONFIGVALUESWITHSOURCES(cfg, root):
Contract:
PRE: true
POST: true
INPUT: cfg, root
OUTPUT: result of GETCONFIGVALUESWITHSOURCES
EFFECTS: State.Config, pure
  enhanced = GETALLCONFIGVALUESWITHSOURCES(cfg, root)  # see IMPL-CONFIG_DISPLAY_FLATTENING
  FOR EACH enhanced: append ConfigValue from enhanced.ConfigValue
  RETURN legacyValues sorted by name

## EMBEDDED_MINITEST_DEFAULT_CONFIG

- [IMPL-CONFIG_STRUCT] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: verify DefaultConfig preset values for archive path, flags, and exclude patterns.

PROCEDURE EMBEDDED_MINITEST_DEFAULT_CONFIG():
Contract:
PRE: true
POST: true
INPUT: context for EMBEDDED_MINITEST_DEFAULT_CONFIG
OUTPUT: result of EMBEDDED_MINITEST_DEFAULT_CONFIG
EFFECTS: State.Config, pure
  cfg = DefaultConfig()
  ASSERT archive_dir_path, use_current_dir_name, exclude_patterns match expected defaults
