// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-SPEC_CTL]
package specparse_test

import (
	"strings"
	"testing"

	"bkpdir/internal/specparse"
)

func TestSpecctlValidate_REQ_PSEUDOCODE_FORMAL_VERIFICATION(t *testing.T) {
	const body = `# [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]

## VALIDATE_COMMAND
SPEC-ID: IMPL-SPEC_CTL::VALIDATE_COMMAND
STEP T001: PARSE paths
- [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION] — How: parse files.

PROCEDURE VALIDATE(paths):
  RETURN success
`
	sc := specparse.ParseReader(strings.NewReader(body))
	doc, err := specparse.Parse("test.md", sc)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Blocks) != 1 {
		t.Fatalf("blocks: %d", len(doc.Blocks))
	}
	if doc.Blocks[0].SpecID != "IMPL-SPEC_CTL::VALIDATE_COMMAND" {
		t.Fatalf("spec id: %q", doc.Blocks[0].SpecID)
	}
	diags := specparse.ValidateDocument(doc, &specparse.Registry{FormalSpec: map[string]bool{"IMPL-SPEC_CTL": true}}, true)
	if specparse.HasErrors(diags) {
		t.Fatalf("unexpected errors: %+v", diags)
	}
}
