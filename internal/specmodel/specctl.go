// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

import (
	"path/filepath"

	"bkpdir/internal/specparse"
)

// OracleSpecctlValidatePath models VALIDATE_COMMAND POST: no error-severity findings.
func OracleSpecctlValidatePath(repoRoot, sidecarRel string) (int, int, error) {
	path := filepath.Join(repoRoot, sidecarRel)
	reg, err := specparse.LoadRegistry(repoRoot)
	if err != nil {
		return 0, 0, err
	}
	doc, err := specparse.ParseFile(path)
	if err != nil {
		return 1, 0, err
	}
	token := specparse.ImplTokenFromPath(path)
	formal := reg.FormalSpec[token]
	diags := specparse.ValidateDocument(doc, reg, formal)
	errs, warns := 0, 0
	for _, d := range diags {
		if d.Severity == "error" {
			errs++
		}
		if d.Severity == "warning" {
			warns++
		}
	}
	return errs, warns, nil
}
