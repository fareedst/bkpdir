// Package speccheck implements traceability matrix and block-lead checks.
// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package speccheck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"bkpdir/internal/specparse"

	"gopkg.in/yaml.v3"
)

// MatrixRow is one traceability row.
type MatrixRow struct {
	SpecID  string
	StepID  string
	Block   string
	Impl    string
	InTests bool
	InCode  bool
}

// BuildMatrix collects STEP/SPEC references and scans Go tree.
func BuildMatrix(repoRoot string, sidecarPaths []string) ([]MatrixRow, error) {
	reg, _ := specparse.LoadRegistry(repoRoot)
	goRefs := scanGoReferences(repoRoot)
	var rows []MatrixRow

	for _, p := range sidecarPaths {
		doc, err := specparse.ParseFile(p)
		if err != nil {
			return nil, err
		}
		impl := specparse.ImplTokenFromPath(p)
		if !reg.FormalSpec[impl] {
			continue
		}
		for _, b := range doc.Blocks {
			if b.SpecID != "" {
				rows = append(rows, MatrixRow{
					SpecID:  b.SpecID,
					Block:   b.Name,
					Impl:    impl,
					InTests: goRefs[b.SpecID],
					InCode:  goRefs[b.SpecID],
				})
			}
			covered := blockCovered(b, goRefs)
			for _, s := range b.Steps {
				rows = append(rows, MatrixRow{
					SpecID:  b.SpecID,
					StepID:  s.ID,
					Block:   b.Name,
					Impl:    impl,
					InTests: covered,
					InCode:  covered,
				})
			}
		}
	}
	return rows, nil
}

var reGoRef = regexp.MustCompile(`(T\d{3}|IMPL-[A-Z0-9_]+::[A-Z0-9_]+|[A-Z][A-Z0-9_]{2,})`)

func normalizeIdent(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToUpper(r))
		}
	}
	return b.String()
}

func scanGoReferences(repoRoot string) map[string]bool {
	found := make(map[string]bool)
	reFunc := regexp.MustCompile(`func\s+(?:\([^)]+\)\s*)?(\w+)\s*\(`)
	reType := regexp.MustCompile(`type\s+(\w+)\s`)
	reVar := regexp.MustCompile(`var\s+(\w+)\s`)
	_ = filepath.Walk(repoRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			if info != nil && info.IsDir() {
				base := filepath.Base(path)
				if base == "vendor" || base == ".git" {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.Contains(path, "internal/specparse") || strings.Contains(path, "cmd/specctl") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(data)
		for _, m := range reGoRef.FindAllString(content, -1) {
			found[m] = true
		}
		for _, m := range reToken.FindAllString(content, -1) {
			found[m] = true
		}
		for _, re := range []*regexp.Regexp{reFunc, reType, reVar} {
			for _, m := range re.FindAllStringSubmatch(content, -1) {
				if len(m) > 1 {
					found[m[1]] = true
					found[normalizeIdent(m[1])] = true
				}
			}
		}
		return nil
	})
	return found
}

func identMatchesGo(name string, goRefs map[string]bool) bool {
	if goRefs[name] {
		return true
	}
	n := normalizeIdent(name)
	if n != "" && goRefs[n] {
		return true
	}
	for k := range goRefs {
		nk := normalizeIdent(k)
		if len(nk) < 4 {
			continue
		}
		if nk == n || strings.HasSuffix(n, nk) || strings.HasSuffix(nk, n) {
			return true
		}
	}
	return false
}

func blockCovered(b specparse.Block, goRefs map[string]bool) bool {
	if identMatchesGo(b.Name, goRefs) {
		return true
	}
	if b.SpecID != "" && goRefs[b.SpecID] {
		return true
	}
	for _, proc := range b.Procedures {
		if identMatchesGo(proc, goRefs) {
			return true
		}
	}
	return false
}

var reToken = regexp.MustCompile(`\[(REQ|ARCH|IMPL)-[^\]]+\]`)

// LoadWaivers reads tied/spec/coverage-waivers.yaml.
func LoadWaivers(repoRoot string) map[string]bool {
	out := make(map[string]bool)
	path := filepath.Join(repoRoot, "tied", "spec", "coverage-waivers.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return out
	}
	var raw struct {
		Waivers []struct {
			Impl  string `yaml:"impl"`
			Block string `yaml:"block"`
		} `yaml:"waivers"`
	}
	if yaml.Unmarshal(data, &raw) != nil {
		return out
	}
	for _, w := range raw.Waivers {
		out[w.Impl+"::"+w.Block] = true
	}
	return out
}

// CoverageDiagnostics returns errors for uncovered formal blocks (COVER-001).
func CoverageDiagnostics(rows []MatrixRow, checkCoverage bool, waivers map[string]bool) []specparse.Diagnostic {
	if !checkCoverage {
		return nil
	}
	seen := map[string]bool{}
	var diags []specparse.Diagnostic
	for _, r := range rows {
		if r.StepID == "" || seen[r.Block] {
			continue
		}
		seen[r.Block] = true
		if waivers[r.Impl+"::"+r.Block] {
			continue
		}
		if !r.InTests && !r.InCode {
			diags = append(diags, specparse.Diagnostic{
				Code: "COVER-001", Severity: "error",
				Message: fmt.Sprintf("block %s has no Go test/code reference for SPEC/steps", r.Block),
				Block:   r.Block,
			})
		}
	}
	return diags
}
