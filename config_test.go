package main

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	// No config file: should get defaults
	cfg, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.ArchiveDirPath != "../.bkpdir" {
		t.Errorf("default ArchiveDirPath = %q, want %q", cfg.ArchiveDirPath, "../.bkpdir")
	}
	// Write a config file
	configYAML := "archive_dir_path: /tmp/archives\nuse_current_dir_name: false\nexclude_patterns:\n  - node_modules/\n  - '*.log'\n"
	os.WriteFile(filepath.Join(dir, ".bkpdir.yml"), []byte(configYAML), 0644)
	cfg2, err := LoadConfig(dir)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg2.ArchiveDirPath != "/tmp/archives" {
		t.Errorf("YAML ArchiveDirPath = %q, want /tmp/archives", cfg2.ArchiveDirPath)
	}
	if cfg2.UseCurrentDirName {
		t.Errorf("YAML UseCurrentDirName = true, want false")
	}
	if len(cfg2.ExcludePatterns) != 2 || cfg2.ExcludePatterns[0] != "node_modules/" || cfg2.ExcludePatterns[1] != "*.log" {
		t.Errorf("YAML ExcludePatterns = %v, want [node_modules/ *.log]", cfg2.ExcludePatterns)
	}
}
