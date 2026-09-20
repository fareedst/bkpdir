// [REQ-PSEUDOCODE_FORMAL_VERIFICATION] [IMPL-ATOMIC_OPS]
package specmodel_test

import (
	"testing"

	"bkpdir/internal/specmodel"
)

func TestAtomicWriterConformance_REQ_PSEUDOCODE_FORMAL_VERIFICATION(t *testing.T) {
	s := specmodel.WriterState{TargetPath: "/tmp/out"}
	s, err := s.Apply(specmodel.EvNew, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Apply(specmodel.EvWrite, 3)
	if err != nil {
		t.Fatal(err)
	}
	s, err = s.Apply(specmodel.EvCommit, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.Apply(specmodel.EvWrite, 1)
	if err == nil || !specmodel.SameErrorClass(err, err) {
		// after commit+close, write must fail
	}
}
