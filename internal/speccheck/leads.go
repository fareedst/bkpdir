// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package speccheck

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bkpdir/internal/specparse"

	"gopkg.in/yaml.v3"
)

// LeadScanDirPaths returns absolute paths under repoRoot scanned for block-lead comments.
// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: pkg/, internal/, cmd/, test/ when present; top-level *.go handled separately via repoRoot in CheckLeads.
func LeadScanDirPaths(repoRoot string) []string {
	rels := []string{"pkg", "internal", "cmd", "test"}
	var out []string
	for _, rel := range rels {
		p := filepath.Join(repoRoot, rel)
		if st, err := os.Stat(p); err == nil && st.IsDir() {
			out = append(out, p)
		}
	}
	return out
}

// LoadLeadCheckSkip reads impl tokens exempt from literal block-lead Go sync.
func LoadLeadCheckSkip(repoRoot string) map[string]bool {
	out := make(map[string]bool)
	path := filepath.Join(repoRoot, "tied", "spec", "coverage-waivers.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return out
	}
	var raw struct {
		LeadCheckSkip []string `yaml:"lead_check_skip"`
	}
	if yaml.Unmarshal(data, &raw) != nil {
		return out
	}
	for _, t := range raw.LeadCheckSkip {
		out[t] = true
	}
	return out
}

// CheckLeads verifies block lead lines appear literally in Go comments under scanDirs and optional repoRoot top-level .go files.
func CheckLeads(sidecarPaths []string, scanDirs []string, repoRoot string) ([]specparse.Diagnostic, error) {
	skipImpl := LoadLeadCheckSkip(repoRoot)
	var diags []specparse.Diagnostic
	for _, p := range sidecarPaths {
		doc, err := specparse.ParseFile(p)
		if err != nil {
			return nil, err
		}
		if skipImpl[specparse.ImplTokenFromPath(p)] {
			continue
		}
		goComments, err := collectGoComments(scanDirs)
		if err != nil {
			return nil, err
		}
		if repoRoot != "" {
			top, err := collectTopLevelGoComments(repoRoot)
			if err != nil {
				return nil, err
			}
			goComments = append(goComments, top...)
		}
		for _, b := range doc.Blocks {
			if b.Lead == "" || strings.HasPrefix(b.Name, "Summary") {
				continue
			}
			needle := leadNeedle(b.Lead)
			if !commentContains(goComments, needle) {
				diags = append(diags, specparse.Diagnostic{
					Code: "TRACE-002", Severity: "error",
					Message: fmt.Sprintf("block lead not found in Go comments: %s", truncate(needle, 80)),
					File:    p, Block: b.Name, Line: b.Line,
				})
			}
		}
	}
	return diags, nil
}

func leadNeedle(lead string) string {
	s := strings.TrimSpace(lead)
	if strings.HasPrefix(s, "- ") {
		s = strings.TrimPrefix(s, "- ")
	}
	return s
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func commentContains(comments []string, needle string) bool {
	for _, c := range comments {
		if strings.Contains(c, needle) {
			return true
		}
	}
	return false
}

func collectGoComments(dirs []string) ([]string, error) {
	var out []string
	for _, root := range dirs {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".go") {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			for _, line := range strings.Split(string(data), "\n") {
				trim := strings.TrimSpace(line)
				if strings.HasPrefix(trim, "//") {
					out = append(out, strings.TrimSpace(strings.TrimPrefix(trim, "//")))
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func collectTopLevelGoComments(repoRoot string) ([]string, error) {
	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, e.Name()))
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			trim := strings.TrimSpace(line)
			if strings.HasPrefix(trim, "//") {
				out = append(out, strings.TrimSpace(strings.TrimPrefix(trim, "//")))
			}
		}
	}
	return out, nil
}
