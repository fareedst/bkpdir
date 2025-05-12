package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ArchiveDirPath    string   `yaml:"archive_dir_path"`
	UseCurrentDirName bool     `yaml:"use_current_dir_name"`
	ExcludePatterns   []string `yaml:"exclude_patterns"`
}

func DefaultConfig() *Config {
	return &Config{
		ArchiveDirPath:    "../.bkpdir",
		UseCurrentDirName: true,
		ExcludePatterns:   []string{".git/", "vendor/"},
	}
}

func LoadConfig(root string) (*Config, error) {
	cfg := DefaultConfig()
	cfgPath := filepath.Join(root, ".bkpdir.yaml")
	if _, err := os.Stat(cfgPath); err == nil {
		f, err := os.Open(cfgPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		d := yaml.NewDecoder(f)
		if err := d.Decode(cfg); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
