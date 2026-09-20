// Command specctl validates formal IMPL pseudocode sidecars.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bkpdir/internal/speccheck"
	"bkpdir/internal/specparse"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "validate":
		os.Exit(runValidate(os.Args[2:]))
	case "matrix":
		os.Exit(runMatrix(os.Args[2:]))
	case "check-leads":
		os.Exit(runCheckLeads(os.Args[2:]))
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: specctl <validate|matrix|check-leads> [flags] [paths...]\n")
}

func repoRoot() string {
	wd, _ := os.Getwd()
	for d := wd; ; {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			return wd
		}
		d = parent
	}
}

func expandPaths(args []string) ([]string, error) {
	if len(args) == 0 {
		args = []string{"tied/implementation-decisions/*-pseudocode.md"}
	}
	var out []string
	for _, pattern := range args {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 && strings.HasSuffix(pattern, ".md") {
			if _, err := os.Stat(pattern); err == nil {
				out = append(out, pattern)
			}
			continue
		}
		out = append(out, matches...)
	}
	return out, nil
}

func runValidate(args []string) int {
	// - [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: parse each file and fail on PARSE/SHAPE/RESOLVE errors; optional CRIT-001 via run_impl_logic_audit.py --req-criteria-strict; L4 mutation pilot and contract eval documented in scripts (not global CI).
	fs := flag.NewFlagSet("validate", flag.ExitOnError)
	fs.Parse(args)
	paths, err := expandPaths(fs.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "specctl validate: %v\n", err)
		return 1
	}
	root := repoRoot()
	diags, err := specparse.ValidatePaths(paths, root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "specctl validate: %v\n", err)
		return 1
	}
	printDiags(diags)
	if specparse.HasErrors(diags) {
		return 1
	}
	return 0
}

func runMatrix(args []string) int {
	// - [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: build REQ/IMPL/STEP to test/code matrix; L4 mutation score on pilot packages via go-mutesting; property conformance tests in test/specconformance.
	fs := flag.NewFlagSet("matrix", flag.ExitOnError)
	checkCoverage := fs.Bool("check-coverage", false, "fail on uncovered STEP rows")
	fs.Parse(args)
	paths, err := expandPaths(fs.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "specctl matrix: %v\n", err)
		return 1
	}
	root := repoRoot()
	rows, err := speccheck.BuildMatrix(root, paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "specctl matrix: %v\n", err)
		return 1
	}
	for _, r := range rows {
		fmt.Printf("%s\t%s\t%s\ttest=%v\tcode=%v\n", r.Impl, r.SpecID, r.StepID, r.InTests, r.InCode)
	}
	waivers := speccheck.LoadWaivers(root)
	diags := speccheck.CoverageDiagnostics(rows, *checkCoverage, waivers)
	printDiags(diags)
	if specparse.HasErrors(diags) {
		return 1
	}
	return 0
}

func runCheckLeads(args []string) int {
	// - [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: verify [PROC-IMPL_PSEUDOCODE_TOKENS] literal copy in Go sources.
	fs := flag.NewFlagSet("check-leads", flag.ExitOnError)
	fs.Parse(args)
	paths, err := expandPaths(fs.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "specctl check-leads: %v\n", err)
		return 1
	}
	root := repoRoot()
	diags, err := speccheck.CheckLeads(paths, speccheck.LeadScanDirPaths(root), root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "specctl check-leads: %v\n", err)
		return 1
	}
	printDiags(diags)
	if specparse.HasErrors(diags) {
		return 1
	}
	return 0
}

func printDiags(diags []specparse.Diagnostic) {
	for _, d := range diags {
		loc := d.File
		if d.Line > 0 {
			loc = fmt.Sprintf("%s:%d", d.File, d.Line)
		}
		if d.Block != "" {
			loc += " block=" + d.Block
		}
		fmt.Printf("[%s] %s %s: %s\n", d.Severity, d.Code, loc, d.Message)
	}
}
