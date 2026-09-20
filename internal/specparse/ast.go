// Package specparse parses IMPL formal pseudocode sidecars.
// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specparse

// Document is a parsed sidecar file.
type Document struct {
	Path        string
	H1Tokens    []string
	Blocks      []Block
	SymbolTable map[string]Symbol
}

// Block is one H2 section.
type Block struct {
	Name       string
	Line       int
	Lead       string
	SpecID     string
	PRE        []string
	POST       []string
	INV        []string
	Steps      []Step
	Branches   []Branch
	Errors     []ErrorDecl
	Procedures []string
	Contracts  map[string]string
}

// Step is STEP T001: action.
type Step struct {
	ID     string
	Action string
	Line   int
}

// Branch is BRANCH B001: condition.
type Branch struct {
	ID        string
	Condition string
	Line      int
}

// ErrorDecl is ERROR E001: description.
type ErrorDecl struct {
	ID          string
	Description string
	Line        int
}

// Symbol describes a resolved name.
type Symbol struct {
	Kind  string
	Block string
}

// Diagnostic is a validation finding.
type Diagnostic struct {
	Code     string
	Severity string
	Message  string
	File     string
	Block    string
	Line     int
	Column   int
}
