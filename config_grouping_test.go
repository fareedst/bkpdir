package main

import (
	"testing"
)

func TestGetFieldImportance(t *testing.T) {
	tests := []struct {
		name     string
		category string
		want     int
	}{
		{"archive_dir_path", "basic_settings", ImportanceCritical},
		{"backup_dir_path", "basic_settings", ImportanceCritical},
		{"use_current_dir_name", "basic_settings", ImportanceHigh},
		{"git.enabled", "git_settings", ImportanceHigh},
		{"format_created_archive", "output_format", ImportanceLow},
		{"status_disk_full", "status_codes", ImportanceLow},
		{"unknown_field", "basic_settings", ImportanceMedium},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldImportance(tt.name, tt.category); got != tt.want {
				t.Errorf("getFieldImportance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCategoryPriority(t *testing.T) {
	// Verify critical categories have lower priority values (appear first)
	if CategoryPriority["basic_settings"] >= CategoryPriority["advanced"] {
		t.Error("basic_settings should have higher priority (lower value) than advanced")
	}
	if CategoryPriority["archive_settings"] >= CategoryPriority["status_codes"] {
		t.Error("archive_settings should have higher priority (lower value) than status_codes")
	}
}
