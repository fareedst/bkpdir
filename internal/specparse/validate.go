// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specparse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Registry holds IMPL tokens requiring formal_spec rules.
type Registry struct {
	FormalSpec map[string]bool `yaml:"formal_spec"`
}

// LoadRegistry reads tied/spec/formal-spec-registry.yaml.
func LoadRegistry(repoRoot string) (*Registry, error) {
	path := filepath.Join(repoRoot, "tied", "spec", "formal-spec-registry.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return &Registry{FormalSpec: map[string]bool{}}, nil
	}
	var raw struct {
		FormalSpec []string `yaml:"formal_spec"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	r := &Registry{FormalSpec: make(map[string]bool)}
	for _, t := range raw.FormalSpec {
		r.FormalSpec[t] = true
	}
	return r, nil
}

// ImplTokenFromPath extracts IMPL-TOKEN from filename.
func ImplTokenFromPath(path string) string {
	base := filepath.Base(path)
	if !strings.HasSuffix(base, "-pseudocode.md") {
		return ""
	}
	return strings.TrimSuffix(base, "-pseudocode.md")
}

// ValidateDocument runs Layer B checks; returns diagnostics.
func ValidateDocument(doc *Document, reg *Registry, formal bool) []Diagnostic {
	var out []Diagnostic
	if doc == nil {
		out = append(out, Diagnostic{
			Code: "PARSE-001", Severity: "error", Message: "nil document",
		})
		return out
	}

	if len(doc.Blocks) == 0 {
		out = append(out, Diagnostic{
			Code: "PARSE-001", Severity: "error", Message: "no H2 blocks found",
			File: doc.Path, Line: 1,
		})
	}

	seenSpec := map[string]int{}
	seenStep := map[string]int{}

	for _, b := range doc.Blocks {
		if b.Lead == "" && !strings.HasPrefix(b.Name, "Summary") {
			out = append(out, Diagnostic{
				Code: "SHAPE-001", Severity: "warning", Message: "missing block lead comment",
				File: doc.Path, Block: b.Name, Line: b.Line,
			})
		}

		if formal && !strings.HasPrefix(b.Name, "Summary") && !strings.HasPrefix(b.Name, "EMBEDDED") {
			if b.SpecID == "" {
				out = append(out, Diagnostic{
					Code: "SHAPE-001", Severity: "error",
					Message: "formal_spec requires SPEC-ID",
					File:    doc.Path, Block: b.Name, Line: b.Line,
				})
			}
			if len(b.Steps) == 0 && len(b.Procedures) == 0 {
				out = append(out, Diagnostic{
					Code: "SHAPE-001", Severity: "error",
					Message: "formal_spec block needs STEP or PROCEDURE",
					File:    doc.Path, Block: b.Name, Line: b.Line,
				})
			}
		}

		if b.SpecID != "" {
			if prev, ok := seenSpec[b.SpecID]; ok {
				out = append(out, Diagnostic{
					Code: "RESOLVE-002", Severity: "error",
					Message: fmt.Sprintf("duplicate SPEC-ID (first at block line %d)", prev),
					File:    doc.Path, Block: b.Name, Line: b.Line,
				})
			}
			seenSpec[b.SpecID] = b.Line
		}

		for _, s := range b.Steps {
			key := b.Name + "::" + s.ID
			if prev, ok := seenStep[key]; ok {
				out = append(out, Diagnostic{
					Code: "RESOLVE-002", Severity: "error",
					Message: fmt.Sprintf("duplicate STEP %s (line %d)", s.ID, prev),
					File:    doc.Path, Block: b.Name, Line: s.Line,
				})
			}
			seenStep[key] = s.Line
		}
	}

	return out
}

// ValidatePaths parses and validates multiple sidecar files.
func ValidatePaths(paths []string, repoRoot string) ([]Diagnostic, error) {
	reg, err := LoadRegistry(repoRoot)
	if err != nil {
		return nil, err
	}
	var all []Diagnostic
	for _, p := range paths {
		doc, err := ParseFile(p)
		if err != nil {
			all = append(all, Diagnostic{
				Code: "PARSE-001", Severity: "error",
				Message: err.Error(), File: p, Line: 1,
			})
			continue
		}
		token := ImplTokenFromPath(p)
		formal := reg.FormalSpec[token]
		all = append(all, ValidateDocument(doc, reg, formal)...)
	}
	return all, nil
}

// HasErrors returns true if any error-severity diagnostic exists.
func HasErrors(diags []Diagnostic) bool {
	for _, d := range diags {
		if d.Severity == "error" {
			return true
		}
	}
	return false
}
