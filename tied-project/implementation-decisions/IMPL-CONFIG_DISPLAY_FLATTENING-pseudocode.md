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
  # [IMPL-CONFIG_DISPLAY_FLATTENING] — Layer B procedure token link
Contract:
INPUT: structType, structValue, prefix, category
OUTPUT: appended configFieldInfo slice with flat YAMLName keys
PRE: structType and structValue describe the same struct kind
POST: nested struct fields are flattened; leaf fields carry yaml tag or lowercase name as YAMLName
EFFECTS: pure
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
  # [IMPL-CONFIG_DISPLAY_FLATTENING] — Layer B procedure token link
Contract:
INPUT: cfg, root path for source attribution
OUTPUT: ConfigValueWithMetadata rows for display
PRE: cfg is non-nil
POST: each emitted row uses flat YAMLName; struct container rows omitted; default comparison populated when root allows
EFFECTS: pure
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
  # [IMPL-CONFIG_DISPLAY_FLATTENING] — Layer B procedure token link
Contract:
INPUT: fieldName, parentCategory
OUTPUT: category string for config command grouping
PRE: fieldName is non-empty
POST: parentCategory wins when set; otherwise category derived from fieldName prefix rules
EFFECTS: pure
  IF parentCategory set THEN RETURN parentCategory
  MATCH fieldName prefixes and substrings to category string

## MERGEGITSETTINGS_LEGACY_FLAT_KEYS

SPEC-ID: IMPL-CONFIG_DISPLAY_FLATTENING::MERGEGITSETTINGS_LEGACY_FLAT_KEYS
STEP T001: EXECUTE PROCEDURE mergeGitSettings

- [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: sync legacy top-level include_git_info and show_git_dirty_status with nested Git struct so YAML using either flat or nested git keys merges consistently for display and runtime.

PROCEDURE mergeGitSettings(dst, src, inheritContext, ...):
  # [IMPL-CONFIG_DISPLAY_FLATTENING] — Layer B procedure token link
Contract:
INPUT: dst Config, src Config, inheritContext and merge options
OUTPUT: dst with Git and legacy flat git booleans aligned
PRE: dst and src are non-nil
POST: legacy top-level git booleans applied before nested merge; dst.Git matches effective flat key values
EFFECTS: mutates dst
DATA_TRANSITION: dst.Git and legacy include_git_info / show_git_dirty_status reflect merged src values
  APPLY legacy top-level git booleans before mergeGitConfigStruct
  KEEP dst.Git and legacy fields in sync when either representation is set
