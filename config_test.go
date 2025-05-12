package main

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.ArchiveDirPath != "../.bkpdir" {
		t.Errorf("ArchiveDirPath = %q, want %q", cfg.ArchiveDirPath, "../.bkpdir")
	}
	if !cfg.UseCurrentDirName {
		t.Errorf("UseCurrentDirName = false, want true")
	}
	if len(cfg.ExcludePatterns) != 2 || cfg.ExcludePatterns[0] != ".git/" || cfg.ExcludePatterns[1] != "vendor/" {
		t.Errorf("ExcludePatterns = %v, want [.git/ vendor/]", cfg.ExcludePatterns)
	}
}
