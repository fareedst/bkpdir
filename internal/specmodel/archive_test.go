// [IMPL-ZIP_FORMAT] [ARCH-ARCHIVE_FORMAT] [REQ-FILE_BACKUP]
package specmodel

import "testing"

func TestOracleZipEntryPath_REQ_FILE_BACKUP(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"foo/bar", "foo/bar"},
		{`foo\bar`, "foo/bar"},
		{"./archive.go", "archive.go"},
		{`.git\HEAD`, ".git/HEAD"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := OracleZipEntryPath(tt.in); got != tt.want {
			t.Errorf("OracleZipEntryPath(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
