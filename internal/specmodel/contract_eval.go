// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

import (
	"fmt"
	"strings"

	"bkpdir/internal/specmodel/contract"
	"bkpdir/internal/specparse"
)

// BlockByName returns the runtime block with the given H2 name.
func BlockByName(doc *specparse.Document, name string) (specparse.Block, bool) {
	for _, b := range doc.Blocks {
		if strings.EqualFold(b.Name, name) {
			return b, true
		}
	}
	return specparse.Block{}, false
}

// EvalBlockPost evaluates all POST: lines on a block against env.
func EvalBlockPost(b specparse.Block, env contract.Env) error {
	for _, post := range b.POST {
		ok, err := contract.EvalBool(post, env)
		if err != nil {
			return fmt.Errorf("POST %q: %w", post, err)
		}
		if !ok {
			return fmt.Errorf("POST failed: %q", post)
		}
	}
	return nil
}

// EvalBlockBranch evaluates BRANCH conditions (expected true when branch applies).
func EvalBlockBranch(b specparse.Block, id string, env contract.Env) (bool, error) {
	for _, br := range b.Branches {
		if br.ID == id {
			return contract.EvalBool(br.Condition, env)
		}
	}
	return false, fmt.Errorf("branch %s not found on block %s", id, b.Name)
}

// WriterEnvFromState maps WriterState to contract.Env.
func WriterEnvFromState(s WriterState) contract.Env {
	return contract.WriterEnv(s.TargetPath, s.TempPath, s.HasHandle, s.Committed, s.Closed)
}
