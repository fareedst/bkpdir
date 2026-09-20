// [IMPL-PACKAGE_EXTRACTION] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

// ExpectedPkgRoots lists extracted package roots from sidecar.
var ExpectedPkgRoots = []string{
	"pkg/cli", "pkg/config", "pkg/errors", "pkg/fileops",
	"pkg/formatter", "pkg/git", "pkg/processing", "pkg/resources", "pkg/testutil",
}

// OraclePkgExtracted reports whether a package name is in the extraction set.
func OraclePkgExtracted(name string) bool {
	for _, p := range ExpectedPkgRoots {
		if p == name || name == "pkg/"+name {
			return true
		}
	}
	return false
}
