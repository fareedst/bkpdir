// [REQ-INCREMENTAL_ANCHOR_MODE] Incremental anchor mode (full vs chained)
package main

import "fmt"

// IncrementalAnchorMode selects how incrementals relate to the latest full archive.
type IncrementalAnchorMode string

const (
	// IncrementalAnchorModeFull: each inc zip is cumulative vs full; effective state uses latest inc only.
	IncrementalAnchorModeFull IncrementalAnchorMode = "full"
	// IncrementalAnchorModeChained: each inc is a delta since prior effective state; merge all incs in order.
	IncrementalAnchorModeChained IncrementalAnchorMode = "chained"
)

// - [IMPL-INCREMENTAL_ANCHOR_MODE] [ARCH-INCREMENTAL_ANCHOR_MODE] [REQ-INCREMENTAL_ANCHOR_MODE] — How: parse and validate incremental_anchor_mode config/CLI values.
func ParseIncrementalAnchorMode(s string) (IncrementalAnchorMode, error) {
	if s == "" {
		return IncrementalAnchorModeFull, nil
	}
	switch IncrementalAnchorMode(s) {
	case IncrementalAnchorModeFull, IncrementalAnchorModeChained:
		return IncrementalAnchorMode(s), nil
	default:
		return "", fmt.Errorf("invalid incremental_anchor_mode %q: want full or chained", s)
	}
}

// - [IMPL-INCREMENTAL_ANCHOR_MODE] [ARCH-INCREMENTAL_ANCHOR_MODE] [REQ-INCREMENTAL_ANCHOR_MODE] [REQ-CONFIGURATION] — How: CLI override wins, then config, default full.
func (c *Config) EffectiveIncrementalAnchorMode(cliOverride string) (IncrementalAnchorMode, error) {
	if cliOverride != "" {
		return ParseIncrementalAnchorMode(cliOverride)
	}
	return ParseIncrementalAnchorMode(c.IncrementalAnchorMode)
}
