// Package specconformance compares pkg/* behavior to internal/specmodel oracles.
// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-ATOMIC_OPS]
package specconformance_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"bkpdir/internal/specmodel"
	"bkpdir/pkg/fileops"
)

// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: embed semantic tokens inline and replicate critical REQ/ARCH/IMPL context so each document layer stays self-contained with registry links.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: verify tokens in code and tests exist in the registry with bidirectional links and no orphans.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: document input/output and invariants per feature and bind each contract to requirement tokens.
// - [IMPL-DOC_ENHANCEMENT] [ARCH-DOCUMENTATION_ARCHITECTURE] [REQ-DOC_001] — How: traverse dependency graph from registry cross-refs to list impacted docs, code, and tests for a proposed token change.
func TestAtomicWriter_StateMachine_REQ_PSEUDOCODE_FORMAL_VERIFICATION(t *testing.T) {
	properties := gopter.NewProperties(gopter.DefaultTestParameters())
	properties.Property("commit then write fails like spec", prop.ForAll(
		func(name string) bool {
			if len(name) == 0 {
				return true
			}
			dir, err := os.MkdirTemp("", "spec-pbt-*")
			if err != nil {
				return false
			}
			defer os.RemoveAll(dir)
			target := filepath.Join(dir, name+".txt")
			w, err := fileops.NewAtomicWriter(target)
			if err != nil {
				return false
			}
			if _, err := w.Write([]byte("x")); err != nil {
				_ = w.Close()
				return false
			}
			if err := w.Commit(); err != nil {
				_ = w.Close()
				return false
			}
			_, err = w.Write([]byte("y"))
			return err != nil
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) > 0 && !strings.Contains(s, "/")
		}),
	))
	properties.Property("model rollback when not committed on close", prop.ForAll(
		func(name string) bool {
			if len(name) == 0 {
				return true
			}
			dir, err := os.MkdirTemp("", "spec-pbt-*")
			if err != nil {
				return false
			}
			defer os.RemoveAll(dir)
			target := filepath.Join(dir, name+".dat")
			model := specmodel.WriterState{TargetPath: target}
			model, err = model.Apply(specmodel.EvNew, 0)
			if err != nil {
				return false
			}
			model, err = model.Apply(specmodel.EvClose, 0)
			if err != nil {
				return false
			}
			w, err := fileops.NewAtomicWriter(target)
			if err != nil {
				return false
			}
			err = w.Close()
			if err != nil {
				return false
			}
			_, statErr := os.Stat(target)
			return os.IsNotExist(statErr) && model.Closed && !model.Committed
		},
		gen.AlphaString().SuchThat(func(s string) bool {
			return len(s) > 0 && !strings.Contains(s, "/")
		}),
	))
	properties.TestingRun(t)
}

func TestAtomicWriter_ContractPost_REQ_PSEUDOCODE_FORMAL_VERIFICATION(t *testing.T) {
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, "IMPL-ATOMIC_OPS")
	if err != nil {
		t.Fatal(err)
	}
	// NEWATOMICWRITER POST after EvNew
	st := specmodel.WriterState{TargetPath: "/tmp/safe/out.txt"}
	st, err = st.Apply(specmodel.EvNew, 0)
	if err != nil {
		t.Fatal(err)
	}
	block, _ := specmodel.BlockByName(doc, "NEWATOMICWRITER")
	if err := specmodel.EvalBlockPost(block, specmodel.WriterEnvFromState(st)); err != nil {
		t.Fatalf("NEWATOMICWRITER POST: %v", err)
	}

	// ATOMICWRITER_COMMIT POST
	st, err = st.Apply(specmodel.EvCommit, 0)
	if err != nil {
		t.Fatal(err)
	}
	AssertBlockContracts(t, "IMPL-ATOMIC_OPS", "ATOMICWRITER_COMMIT", specmodel.WriterEnvFromState(st))

	// ATOMICWRITER_CLOSE POST on rolled-back close
	st2 := specmodel.WriterState{TargetPath: "/tmp/safe/x.txt"}
	st2, _ = st2.Apply(specmodel.EvNew, 0)
	st2, _ = st2.Apply(specmodel.EvClose, 0)
	AssertBlockContracts(t, "IMPL-ATOMIC_OPS", "ATOMICWRITER_CLOSE", specmodel.WriterEnvFromState(st2))

	// BRANCH B001 on commit when already committed
	st3 := specmodel.WriterState{TargetPath: "/tmp/y.txt", Committed: true, Closed: true}
	blockCommit, _ := specmodel.BlockByName(doc, "ATOMICWRITER_COMMIT")
	ok, err := specmodel.EvalBlockBranch(blockCommit, "B001", specmodel.WriterEnvFromState(st3))
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected B001 committed guard true")
	}
}
