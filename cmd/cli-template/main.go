// [REQ-DOC_016] Semantic token coverage audit marker
// [ARCH-TOKEN_SYSTEM] Token system architecture
// [IMPL-TOKEN_COVERAGE_AUDIT] Applied audit annotation

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
package main

import (
	"context"
	"os"

	"cli-template/cmd"
)

func main() {
	ctx := context.Background()

	// Execute the root command
	if err := cmd.Execute(ctx); err != nil {
		os.Exit(1)
	}
}
