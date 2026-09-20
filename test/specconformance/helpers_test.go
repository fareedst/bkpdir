// Package specconformance compares pkg/* behavior to internal/specmodel oracles.
package specconformance_test

import (
	"os"
	"path/filepath"
	"testing"

	"bkpdir/internal/specmodel"
	"bkpdir/internal/specmodel/contract"
	"bkpdir/internal/specparse"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for d := wd; ; {
		if _, err := os.Stat(filepath.Join(d, "go.mod")); err == nil {
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			t.Fatal("go.mod not found")
		}
		d = parent
	}
}

// AssertBlockContracts evaluates PRE/POST/BRANCH from sidecar against env.
func AssertBlockContracts(t *testing.T, implToken, blockName string, env contract.Env) {
	t.Helper()
	root := repoRoot(t)
	doc, err := specmodel.SidecarForImpl(root, implToken)
	if err != nil {
		t.Fatalf("load sidecar %s: %v", implToken, err)
	}
	block, ok := specmodel.BlockByName(doc, blockName)
	if !ok {
		t.Fatalf("block %s not in %s sidecar", blockName, implToken)
	}
	for _, pre := range block.PRE {
		ok, err := contract.EvalBool(pre, env)
		if err != nil {
			t.Fatalf("PRE %q: %v", pre, err)
		}
		if !ok {
			t.Fatalf("PRE failed: %q", pre)
		}
	}
	if err := specmodel.EvalBlockPost(block, env); err != nil {
		t.Fatal(err)
	}
}

func sidecarBlock(t *testing.T, impl, name string) specparse.Block {
	t.Helper()
	doc, err := specmodel.SidecarForImpl(repoRoot(t), impl)
	if err != nil {
		t.Fatal(err)
	}
	b, ok := specmodel.BlockByName(doc, name)
	if !ok {
		t.Fatalf("block %s missing", name)
	}
	return b
}
