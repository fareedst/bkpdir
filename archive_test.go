package main

import "testing"

func TestGenerateArchiveName(t *testing.T) {
	tests := []struct {
		prefix        string
		timestamp     string
		gitBranch     string
		gitHash       string
		note          string
		isGit         bool
		isIncremental bool
		baseName      string
		expect        string
		name          string
	}{
		{"HOME", "2024-06-01-12-00", "main", "abc123", "", true, false, "", "HOME-2024-06-01-12-00=main=abc123.zip", "full with git"},
		{"", "2024-06-01-12-00", "", "", "note", false, false, "", "2024-06-01-12-00=note.zip", "full with note only"},
		{"", "2024-06-01-12-00", "dev", "def456", "backup", true, true, "HOME-2024-06-01-12-00=main=abc123", "HOME-2024-06-01-12-00=main=abc123_update=2024-06-01-12-00=dev=def456=backup.zip", "incremental with git and note"},
		{"", "2024-06-01-12-00", "", "", "", false, true, "ARCHIVE-2024-06-01-12-00", "ARCHIVE-2024-06-01-12-00_update=2024-06-01-12-00.zip", "incremental no git, no note"},
		{"", "2024-06-01-12-00", "", "", "", false, false, "", "2024-06-01-12-00.zip", "full no git, no note"},
		{"", "2024-06-01-12-00", "", "", "", true, false, "", "2024-06-01-12-00.zip", "full git, no branch/hash"},
		{"", "2024-06-01-12-00", "", "", "test note", false, false, "", "2024-06-01-12-00=test note.zip", "full with note containing spaces"},
		{"", "2024-06-01-12-00", "", "", "test-note", false, true, "ARCHIVE-2024-06-01-12-00", "ARCHIVE-2024-06-01-12-00_update=2024-06-01-12-00=test-note.zip", "incremental with note only"},
		{"", "2024-06-01-12-00", "main", "abc123", "test note", true, false, "", "2024-06-01-12-00=main=abc123=test note.zip", "full with git and note with spaces"},
	}
	for _, tt := range tests {
		got := GenerateArchiveName(tt.prefix, tt.timestamp, tt.gitBranch, tt.gitHash, tt.note, tt.isGit, tt.isIncremental, tt.baseName)
		if got != tt.expect {
			t.Errorf("%s: GenerateArchiveName() = %q, want %q", tt.name, got, tt.expect)
		}
	}
}
