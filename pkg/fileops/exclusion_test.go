// [REQ-CONFIGURATION] [ARCH-EXCLUSION_PATTERNS] [IMPL-EXCLUSION_PATTERNS]
package fileops

import "testing"

func TestShouldExcludeFile_GlobAndExact_REQ_CONFIGURATION(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		patterns []string
		want     bool
	}{
		{name: "glob_star_txt", path: "logs/debug.log", patterns: []string{"*.log"}, want: true},
		{name: "no_match", path: "src/main.go", patterns: []string{"*.log"}, want: false},
		{name: "exact", path: "tmp/cache", patterns: []string{"tmp/cache"}, want: true},
		{name: "empty_patterns", path: "any.txt", patterns: nil, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldExcludeFile(tt.path, tt.patterns); got != tt.want {
				t.Errorf("ShouldExcludeFile(%q, %v) = %v, want %v", tt.path, tt.patterns, got, tt.want)
			}
		})
	}
}

func TestPatternMatcher_ShouldExclude_REQ_CONFIGURATION(t *testing.T) {
	pm := NewPatternMatcher([]string{"vendor/", "*.tmp"})
	if !pm.ShouldExclude("vendor/foo/bar.go") {
		t.Error("expected vendor/ to exclude nested path")
	}
	if !pm.ShouldExclude("scratch.tmp") {
		t.Error("expected *.tmp to match")
	}
	if pm.ShouldExclude("cmd/main.go") {
		t.Error("did not expect cmd/main.go to be excluded")
	}
}
