package main

import "testing"

func TestShouldExcludeFile(t *testing.T) {
	patterns := []string{".git/", "vendor/", "*.tmp", "foo/bar/*"}
	tests := []struct {
		file     string
		excluded bool
	}{
		{".git/config", true},
		{"vendor/package/file.go", true},
		{"main.go", false},
		{"temp.tmp", true},
		{"foo/bar/baz.txt", true},
		{"foo/baz.txt", false},
	}
	for _, tt := range tests {
		if got := ShouldExcludeFile(tt.file, patterns); got != tt.excluded {
			t.Errorf("ShouldExcludeFile(%q) = %v, want %v", tt.file, got, tt.excluded)
		}
	}
}
