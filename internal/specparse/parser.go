// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package specparse

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	reH1     = regexp.MustCompile(`^#\s+(.+)$`)
	reH2     = regexp.MustCompile(`^##\s+(.+)$`)
	reSpecID = regexp.MustCompile(`^SPEC-ID:\s*(\S+)\s*$`)
	rePRE    = regexp.MustCompile(`^PRE:\s*(.+)$`)
	rePOST   = regexp.MustCompile(`^POST:\s*(.+)$`)
	reINV    = regexp.MustCompile(`^INV:\s*(.+)$`)
	reSTEP   = regexp.MustCompile(`^STEP\s+(T\d+[A-Z]?):\s*(.+)$`)
	reBRANCH = regexp.MustCompile(`^BRANCH\s+(B\d+[A-Z]?):\s*(.+)$`)
	reERROR  = regexp.MustCompile(`^ERROR\s+(E\d+[A-Z]?):\s*(.+)$`)
	reProc   = regexp.MustCompile(`^PROCEDURE\s+(\w+)\s*\(`)
	reInput  = regexp.MustCompile(`^(INPUT|OUTPUT|DATA|CONTROL):\s*(.+)$`)
	reToken  = regexp.MustCompile(`\[(REQ|ARCH|IMPL)-[^\]]+\]`)
)

// ParseReader parses from an io.Reader.
func ParseReader(r io.Reader) *bufio.Scanner {
	return bufio.NewScanner(r)
}

// ParseFile reads and parses a pseudocode sidecar.
func ParseFile(path string) (*Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(path, bufio.NewScanner(f))
}

// Parse parses sidecar content from a scanner.
func Parse(path string, sc *bufio.Scanner) (*Document, error) {
	doc := &Document{Path: path, SymbolTable: make(map[string]Symbol)}
	var block *Block
	lineNo := 0

	for sc.Scan() {
		lineNo++
		line := strings.TrimRight(sc.Text(), "\r")
		trim := strings.TrimSpace(line)

		if m := reH1.FindStringSubmatch(trim); m != nil {
			doc.H1Tokens = reToken.FindAllString(m[1], -1)
			continue
		}
		if m := reH2.FindStringSubmatch(trim); m != nil {
			if block != nil {
				doc.Blocks = append(doc.Blocks, *block)
			}
			block = &Block{
				Name:      strings.TrimSpace(m[1]),
				Line:      lineNo,
				Contracts: make(map[string]string),
			}
			continue
		}
		if block == nil {
			continue
		}

		switch {
		case strings.HasPrefix(trim, "- ["):
			block.Lead = trim
		case reSpecID.MatchString(trim):
			block.SpecID = reSpecID.FindStringSubmatch(trim)[1]
			doc.SymbolTable[block.SpecID] = Symbol{Kind: "spec_id", Block: block.Name}
		case rePRE.MatchString(trim):
			block.PRE = append(block.PRE, rePRE.FindStringSubmatch(trim)[1])
		case rePOST.MatchString(trim):
			block.POST = append(block.POST, rePOST.FindStringSubmatch(trim)[1])
		case reINV.MatchString(trim):
			block.INV = append(block.INV, reINV.FindStringSubmatch(trim)[1])
		case reSTEP.MatchString(trim):
			sm := reSTEP.FindStringSubmatch(trim)
			block.Steps = append(block.Steps, Step{ID: sm[1], Action: sm[2], Line: lineNo})
			doc.SymbolTable[fmt.Sprintf("%s::%s", block.Name, sm[1])] = Symbol{Kind: "step", Block: block.Name}
		case reBRANCH.MatchString(trim):
			sm := reBRANCH.FindStringSubmatch(trim)
			block.Branches = append(block.Branches, Branch{ID: sm[1], Condition: sm[2], Line: lineNo})
		case reERROR.MatchString(trim):
			sm := reERROR.FindStringSubmatch(trim)
			block.Errors = append(block.Errors, ErrorDecl{ID: sm[1], Description: sm[2], Line: lineNo})
		case reProc.MatchString(trim):
			name := reProc.FindStringSubmatch(trim)[1]
			block.Procedures = append(block.Procedures, name)
			doc.SymbolTable[name] = Symbol{Kind: "procedure", Block: block.Name}
		case reInput.MatchString(trim):
			sm := reInput.FindStringSubmatch(trim)
			block.Contracts[sm[1]] = sm[2]
		}
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if block != nil {
		doc.Blocks = append(doc.Blocks, *block)
	}
	return doc, nil
}
