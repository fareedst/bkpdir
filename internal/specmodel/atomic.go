// Package specmodel is a reference oracle for bounded IMPL domains.
// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specmodel

import "fmt"

// WriterState models AtomicWriter flags from IMPL-ATOMIC_OPS.
type WriterState struct {
	TargetPath string
	TempPath   string
	HasHandle  bool
	Committed  bool
	Closed     bool
}

// Event is an operation on the writer.
type Event int

const (
	EvNew Event = iota
	EvWrite
	EvCommit
	EvRollback
	EvClose
)

// Apply returns next state and error per spec semantics.
func (s WriterState) Apply(ev Event, dataLen int) (WriterState, error) {
	next := s
	switch ev {
	case EvNew:
		if s.TargetPath == "" {
			return s, fmt.Errorf("invalid target path")
		}
		next.HasHandle = true
		next.Committed = false
		next.Closed = false
		next.TempPath = s.TargetPath + ".tmp.gen"
	case EvWrite:
		if next.Closed {
			return s, fmt.Errorf("writer is closed")
		}
		if !next.HasHandle {
			return s, fmt.Errorf("no temporary file available")
		}
	case EvCommit:
		if next.Committed {
			return s, fmt.Errorf("already committed")
		}
		if next.Closed && next.HasHandle {
			return s, fmt.Errorf("writer is closed but not properly cleaned up")
		}
		next.HasHandle = false
		next.Committed = true
		next.Closed = true
	case EvRollback:
		if next.Committed {
			return s, fmt.Errorf("cannot rollback after commit")
		}
		next.HasHandle = false
		next.Closed = true
	case EvClose:
		if next.Closed {
			return next, nil
		}
		if !next.Committed {
			return next.Apply(EvRollback, 0)
		}
		next.Closed = true
	}
	return next, nil
}

// SameErrorClass reports whether two errors match spec failure class.
func SameErrorClass(want, got error) bool {
	if want == nil && got == nil {
		return true
	}
	if want == nil || got == nil {
		return false
	}
	return want.Error() == got.Error()
}
