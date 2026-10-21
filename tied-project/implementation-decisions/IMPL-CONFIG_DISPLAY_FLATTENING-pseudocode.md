# [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006]

## Summary contract

Presents configuration keys in flat YAML-tag form in `bkpdir config` output even when values live in nested structs, using reflection metadata (YAMLName) and legacy top-level Git field compatibility during merge.

INPUT: Config struct, optional config root path for source attribution
OUTPUT: sorted configFieldInfo list and ConfigValueWithMetadata rows for display
DATA: yaml tags, field cache, legacy include_git_info / show_git_dirty_status keys

## REFLECTCONFIGFIELDS_FLAT_DISPLAY

SPEC-ID: IMPL-CONFIG_DISPLAY_FLATTENING::REFLECTCONFIGFIELDS_FLAT_DISPLAY
STEP T001: EXECUTE PROCEDURE reflectConfigFields

- [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: recurse nested structs but emit leaf fields with YAML tag names so nested Git and similar fields show as top-level keys in config output.

PROCEDURE reflectConfigFields(structType, structValue, prefix, category):
  FOR EACH exported field:
    yamlName = first segment of yaml struct tag or lowercase field name
    IF field is nested struct (not yaml.Node) THEN
      APPEND reflectConfigFields(nestedType, nestedValue, fieldPath, category)
    ELSE
      APPEND configFieldInfo with YAMLName = yamlName for flat display key

## GETALLCONFIGVALUESWITHSOURCES

SPEC-ID: IMPL-CONFIG_DISPLAY_FLATTENING::GETALLCONFIGVALUESWITHSOURCES
STEP T001: EXECUTE PROCEDURE GetAllConfigValuesWithSources

- [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: build ConfigValueWithMetadata from reflected fields, skip struct container rows, attach default and file source metadata for each flat key.

PROCEDURE GetAllConfigValuesWithSources(cfg, root):
  fields = GetAllConfigFields(cfg)
  defaultValues = map from DefaultConfig field paths
  FOR EACH field IN fields:
    IF field is non-pointer struct container THEN CONTINUE
    EMIT ConfigValueWithMetadata using field.YAMLName, value, source, default comparison

## DETERMINEFIELDCATEGORY

SPEC-ID: IMPL-CONFIG_DISPLAY_FLATTENING::DETERMINEFIELDCATEGORY
STEP T001: EXECUTE PROCEDURE determineFieldCategory

- [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: assign display grouping (status_codes, format_strings, backup_settings, etc.) from field name patterns for sorted config command sections.

PROCEDURE determineFieldCategory(fieldName, parentCategory):
  IF parentCategory set THEN RETURN parentCategory
  MATCH fieldName prefixes and substrings to category string

## MERGEGITSETTINGS_LEGACY_FLAT_KEYS

SPEC-ID: IMPL-CONFIG_DISPLAY_FLATTENING::MERGEGITSETTINGS_LEGACY_FLAT_KEYS
STEP T001: EXECUTE PROCEDURE mergeGitSettings

- [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: sync legacy top-level include_git_info and show_git_dirty_status with nested Git struct so YAML using either flat or nested git keys merges consistently for display and runtime.

PROCEDURE mergeGitSettings(dst, src, inheritContext, ...):
  APPLY legacy top-level git booleans before mergeGitConfigStruct
  KEEP dst.Git and legacy fields in sync when either representation is set
