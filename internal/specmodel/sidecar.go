// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

import (
	"path/filepath"
	"strings"

	"bkpdir/internal/specparse"
)

// SidecarForImpl loads the pseudocode sidecar for an IMPL token.
func SidecarForImpl(repoRoot, impl string) (*specparse.Document, error) {
	path := filepath.Join(repoRoot, "tied", "implementation-decisions", impl+"-pseudocode.md")
	return specparse.ParseFile(path)
}

// BlockHasFormalMeta reports whether block has SPEC-ID and at least one STEP or PROCEDURE.
func BlockHasFormalMeta(b specparse.Block) bool {
	if b.SpecID == "" {
		return false
	}
	return len(b.Steps) > 0 || len(b.Procedures) > 0
}

// BlocksMissingLead returns runtime block names without block leads.
func BlocksMissingLead(doc *specparse.Document) []string {
	var out []string
	for _, b := range doc.Blocks {
		if b.Lead == "" && !isExemptBlock(b.Name) {
			out = append(out, b.Name)
		}
	}
	return out
}

func isExemptBlock(name string) bool {
	if name == "Summary contract" {
		return true
	}
	return strings.HasPrefix(name, "EMBEDDED")
}
