// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// OracleEntry describes one IMPL reference oracle.
type OracleEntry struct {
	Impl             string   `yaml:"impl"`
	SpecmodelFile    string   `yaml:"specmodel_file"`
	ConformanceTest  string   `yaml:"conformance_test"`
	PkgPaths         []string `yaml:"pkg_paths"`
	Status           string   `yaml:"status"`
}

// OracleRegistry is tied/spec/oracle-registry.yaml.
type OracleRegistry struct {
	OracleTargets int           `yaml:"oracle_targets"`
	Oracles       []OracleEntry `yaml:"oracles"`
}

// LoadOracleRegistry reads tied/spec/oracle-registry.yaml from repoRoot.
func LoadOracleRegistry(repoRoot string) (*OracleRegistry, error) {
	path := filepath.Join(repoRoot, "tied", "spec", "oracle-registry.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var reg OracleRegistry
	if err := yaml.Unmarshal(data, &reg); err != nil {
		return nil, err
	}
	return &reg, nil
}

// VerifiedCount returns oracles with status verified.
func (r *OracleRegistry) VerifiedCount() int {
	n := 0
	for _, o := range r.Oracles {
		if o.Status == "verified" {
			n++
		}
	}
	return n
}

// EntryForImpl returns the registry row for impl token.
func (r *OracleRegistry) EntryForImpl(impl string) (OracleEntry, bool) {
	for _, o := range r.Oracles {
		if o.Impl == impl {
			return o, true
		}
	}
	return OracleEntry{}, false
}
