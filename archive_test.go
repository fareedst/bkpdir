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
	}{
		{"HOME", "2024-06-01-12-00", "main", "abc123", "", true, false, "", "HOME-2024-06-01-12-00=main=abc123.zip"},
		{"", "2024-06-01-12-00", "", "", "note", false, false, "", "2024-06-01-12-00=note.zip"},
		{"", "2024-06-01-12-00", "dev", "def456", "backup", true, true, "HOME-2024-06-01-12-00=main=abc123", "HOME-2024-06-01-12-00=main=abc123_update=2024-06-01-12-00=dev=def456=backup.zip"},
	}
	for _, tt := range tests {
		got := GenerateArchiveName(tt.prefix, tt.timestamp, tt.gitBranch, tt.gitHash, tt.note, tt.isGit, tt.isIncremental, tt.baseName)
		if got != tt.expect {
			t.Errorf("GenerateArchiveName() = %q, want %q", got, tt.expect)
		}
	}
}
