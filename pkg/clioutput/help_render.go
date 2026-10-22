// [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD]
package clioutput

import (
	"strings"
)

const profileRuleWidth = 56

// RenderOptions configures appearance profile rendering.
type RenderOptions struct {
	Plain bool
}

// RenderHelp post-processes Cobra-generated help into the working appearance profile.
func RenderHelp(cobraText string, opts RenderOptions) string {
	if cobraText == "" {
		return ""
	}
	lines := strings.Split(cobraText, "\n")
	out := make([]string, 0, len(lines)+8)
	titleDone := false
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			out = append(out, line)
			continue
		}
		if isSectionHeader(trimmed) {
			out = append(out, renderSectionHeader(trimmed, opts.Plain))
			out = append(out, renderRule(opts.Plain))
			continue
		}
		if name, desc, ok := parseCommandRow(line); ok {
			out = append(out, renderCommandRow(line, name, desc, opts.Plain))
			continue
		}
		if !titleDone && !strings.HasPrefix(trimmed, "Usage:") {
			out = append(out, renderTitle(line, opts.Plain))
			titleDone = true
			continue
		}
		out = append(out, renderBody(line, opts.Plain))
	}
	return strings.Join(out, "\n")
}

func isSectionHeader(trimmed string) bool {
	if !strings.HasSuffix(trimmed, ":") {
		return false
	}
	switch trimmed {
	case "Usage:", "Aliases:", "Examples:", "Flags:", "Global Flags:", "Available Commands:", "Additional Help Topics:":
		return true
	default:
		// Cobra local flags sections: "Flags for backup:" etc.
		if strings.HasPrefix(trimmed, "Flags for ") && strings.HasSuffix(trimmed, ":") {
			return true
		}
		return false
	}
}

func parseCommandRow(line string) (name, desc string, ok bool) {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == line || trimmed == "" {
		return "", "", false
	}
	fields := strings.Fields(trimmed)
	if len(fields) < 2 {
		return "", "", false
	}
	name = fields[0]
	if strings.HasPrefix(name, "-") {
		return "", "", false
	}
	desc = strings.TrimSpace(strings.TrimPrefix(trimmed, name))
	if desc == "" {
		return "", "", false
	}
	return name, desc, true
}

func renderTitle(line string, plain bool) string {
	if plain {
		return line
	}
	return ansiBold + line + ansiReset
}

func renderBody(line string, plain bool) string {
	if plain {
		return line
	}
	return ansiDim + line + ansiReset
}

func renderSectionHeader(header string, plain bool) string {
	if plain {
		return header
	}
	return ansiBold + ansiCyan + "▶ " + header + ansiReset
}

func renderRule(plain bool) string {
	if plain {
		return strings.Repeat("-", profileRuleWidth)
	}
	return ansiDim + strings.Repeat("─", profileRuleWidth) + ansiReset
}

func renderCommandRow(original, name, desc string, plain bool) string {
	indent := original[:len(original)-len(strings.TrimLeft(original, " \t"))]
	if plain {
		pad := profileCommandNameWidth - len(name)
		if pad < 1 {
			pad = 1
		}
		return indent + name + strings.Repeat(" ", pad) + desc
	}
	pad := profileCommandNameWidth - len(name)
	if pad < 1 {
		pad = 1
	}
	return indent + ansiBold + name + ansiReset + strings.Repeat(" ", pad) + ansiDim + desc + ansiReset
}

const profileCommandNameWidth = 26

const (
	ansiReset = "\033[0m"
	ansiBold  = "\033[1m"
	ansiDim   = "\033[2m"
	ansiCyan  = "\033[36m"
)

