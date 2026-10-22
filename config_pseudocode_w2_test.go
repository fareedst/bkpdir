package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// SPEC-ID: IMPL-CONFIGURABLE_STRINGS::DEFAULT_CONFIG_ERROR_FORMATS
// STEP T001: DECLARE Format and Template error fields on Config struct
// STEP T002: SEED DefaultConfig with backward-compatible printf defaults
// - [IMPL-CONFIGURABLE_STRINGS] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: declare Format* and Template* error message fields on Config and seed DefaultConfig with backward-compatible printf defaults.
func TestIMPL_CONFIGURABLE_STRINGS_DEFAULT_CONFIG_ERROR_FORMATS_COVER001(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.FormatDiskFullError != "Disk full error: %v\n" {
		t.Errorf("FormatDiskFullError = %q, want legacy disk-full default", cfg.FormatDiskFullError)
	}
	if cfg.FormatPermissionError == "" || cfg.TemplateDiskFullError == "" {
		t.Error("expected non-empty Format* and Template* error defaults on DefaultConfig")
	}
}

// SPEC-ID: IMPL-CONFIGURABLE_STRINGS::GET_ERROR_FORMAT_STRINGS
// STEP T001: BUILD name-to-format-string map from cfg Format fields
// - [IMPL-CONFIGURABLE_STRINGS] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: return a name-to-format-string map so error handlers and tests read configurable messages without hardcoded literals.
func TestIMPL_CONFIGURABLE_STRINGS_GET_ERROR_FORMAT_STRINGS_COVER001(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FormatDiskFullError = "custom disk: %v"
	m := cfg.GetErrorFormatStrings()
	if m["disk_full"] != "custom disk: %v" {
		t.Errorf("disk_full = %q, want custom pattern bound from cfg", m["disk_full"])
	}
	for _, key := range []string{"permission", "failed_write_temp", "failed_access_file"} {
		if m[key] == "" {
			t.Errorf("expected non-empty map entry for %q", key)
		}
	}
}

// SPEC-ID: IMPL-CONFIGURABLE_STRINGS::FORMAT_DISK_FULL_ERROR
// STEP T001: SPRINT configured FormatDiskFullError pattern with err
// - [IMPL-CONFIGURABLE_STRINGS] [ARCH-CONFIG_SYSTEM] [REQ-CONFIGURATION] — How: sprintf the configured FormatDiskFullError pattern with the underlying error value.
func TestIMPL_CONFIGURABLE_STRINGS_FORMAT_DISK_FULL_ERROR_COVER001(t *testing.T) {
	cfg := DefaultConfig()
	cfg.FormatDiskFullError = "disk: %v"
	formatter := NewOutputFormatter(cfg)
	err := errors.New("ENOSPC")
	got := formatter.FormatDiskFullError(err)
	want := "disk: ENOSPC"
	if got != want {
		t.Errorf("FormatDiskFullError() = %q, want %q", got, want)
	}
}

// SPEC-ID: IMPL-CONFIG_DISPLAY_FLATTENING::REFLECTCONFIGFIELDS_FLAT_DISPLAY
// STEP T001: EXECUTE PROCEDURE reflectConfigFields
// - [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: recurse nested structs but emit leaf fields with YAML tag names so nested Git and similar fields show as top-level keys in config output.
func TestIMPL_CONFIG_DISPLAY_FLATTENING_reflectConfigFields_COVER001(t *testing.T) {
	cfg := DefaultConfig()
	rt := reflect.TypeOf(*cfg)
	rv := reflect.ValueOf(*cfg)
	fields := reflectConfigFields(rt, rv, "", "")
	if len(fields) < 50 {
		t.Fatalf("expected many discovered fields, got %d", len(fields))
	}
	foundArchiveDir := false
	for _, f := range fields {
		if f.YAMLName == "archive_dir_path" {
			foundArchiveDir = true
		}
		if strings.Contains(f.Path, "Git.") && f.YAMLName == "" {
			t.Errorf("Git nested field %s missing YAMLName for flat display", f.Path)
		}
	}
	if !foundArchiveDir {
		t.Error("expected flat yaml key archive_dir_path from reflectConfigFields")
	}
}

// SPEC-ID: IMPL-CONFIG_DISPLAY_FLATTENING::GETALLCONFIGVALUESWITHSOURCES
// STEP T001: EXECUTE PROCEDURE GetAllConfigValuesWithSources
// - [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: build ConfigValueWithMetadata from reflected fields, skip struct container rows, attach default and file source metadata for each flat key.
func TestIMPL_CONFIG_DISPLAY_FLATTENING_GetAllConfigValuesWithSources_COVER001(t *testing.T) {
	cfg := DefaultConfig()
	values := GetAllConfigValuesWithSources(cfg, ".")
	if len(values) < 50 {
		t.Fatalf("expected at least 50 display rows, got %d", len(values))
	}
	for _, v := range values {
		if v.ConfigValue.Name == "" {
			t.Fatal("display row missing flat ConfigValue.Name (YAMLName)")
		}
	}
}

// SPEC-ID: IMPL-CONFIG_DISPLAY_FLATTENING::DETERMINEFIELDCATEGORY
// STEP T001: EXECUTE PROCEDURE determineFieldCategory
// - [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: assign display grouping (status_codes, format_strings, backup_settings, etc.) from field name patterns for sorted config command sections.
func TestIMPL_CONFIG_DISPLAY_FLATTENING_determineFieldCategory_COVER001(t *testing.T) {
	if got := determineFieldCategory("StatusDiskFull", ""); got != "status_codes" {
		t.Errorf("StatusDiskFull category = %q, want status_codes", got)
	}
	if got := determineFieldCategory("FormatDiskFullError", ""); got != "format_strings" {
		t.Errorf("FormatDiskFullError category = %q, want format_strings", got)
	}
	if got := determineFieldCategory("ChildField", "backup_settings"); got != "backup_settings" {
		t.Errorf("parent category override = %q, want backup_settings", got)
	}
}

// SPEC-ID: IMPL-CONFIG_DISPLAY_FLATTENING::MERGEGITSETTINGS_LEGACY_FLAT_KEYS
// STEP T001: EXECUTE PROCEDURE mergeGitSettings
// - [IMPL-CONFIG_DISPLAY_FLATTENING] [IMPL-CFG_006] [ARCH-CFG_006] [REQ-CFG_006] — How: sync legacy top-level include_git_info and show_git_dirty_status with nested Git struct so YAML using either flat or nested git keys merges consistently for display and runtime.
func TestIMPL_CONFIG_DISPLAY_FLATTENING_mergeGitSettings_COVER001(t *testing.T) {
	dst := DefaultConfig()
	src := DefaultConfig()
	src.IncludeGitInfo = true
	src.ShowGitDirtyStatus = true
	rawSrc := map[string]interface{}{
		"include_git_info":       true,
		"show_git_dirty_status":  true,
	}
	mergeGitSettings(dst, src, true, DefaultConfig(), rawSrc, DefaultConfig(), dst, map[string]bool{})
	if !dst.IncludeGitInfo || dst.Git == nil || !dst.Git.IncludeInfo {
		t.Error("expected legacy include_git_info synced to dst and nested Git.IncludeInfo")
	}
	if !dst.ShowGitDirtyStatus || !dst.Git.ShowDirtyStatus {
		t.Error("expected legacy show_git_dirty_status synced to nested Git.ShowDirtyStatus")
	}
}
