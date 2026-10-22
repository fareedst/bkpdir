// Package clioutput implements the bkpdir operator output standard stream contract.
// [IMPL-BKPDIR_CLI_OUTPUT_CONFORMANCE] [ARCH-BKPDIR_CLI_OUTPUT_STANDARD] [REQ-BKPDIR_CLI_OUTPUT_STANDARD]
package clioutput

import (
	"io"
	"os"
)

// HelpDestination selects where successful help text is written.
type HelpDestination int

const (
	HelpToStderr HelpDestination = iota
	HelpToStdout
)

// Policy is the resolved operator output policy for one invocation.
type Policy struct {
	HelpDest HelpDestination
	Plain    bool
	stdout   io.Writer
	stderr   io.Writer
	isTTY    func(io.Writer) bool
}

// Option configures policy resolution (tests inject writers and TTY probes).
type Option func(*policyConfig)

type policyConfig struct {
	stdout io.Writer
	stderr io.Writer
	isTTY  func(io.Writer) bool
	env    map[string]string
}

func defaultConfig() policyConfig {
	return policyConfig{
		stdout: os.Stdout,
		stderr: os.Stderr,
		isTTY:  isTerminalWriter,
		env:    environMap(os.Environ()),
	}
}

// WithWriters sets stdout/stderr writers for tests.
func WithWriters(stdout, stderr io.Writer) Option {
	return func(c *policyConfig) {
		c.stdout = stdout
		c.stderr = stderr
	}
}

// WithEnv replaces the environment map used for TIED_CLI_* and NO_COLOR.
func WithEnv(env map[string]string) Option {
	return func(c *policyConfig) {
		c.env = env
	}
}

// WithTTYFunc overrides TTY detection (tests).
func WithTTYFunc(fn func(io.Writer) bool) Option {
	return func(c *policyConfig) {
		c.isTTY = fn
	}
}

// PeeledFlags holds global texture flags removed from argv before dispatch.
type PeeledFlags struct {
	Plain   bool
	Stdout  bool
	NoColor bool
}

// ResolvePolicy merges peeled flags and environment into a Policy.
func ResolvePolicy(peeled PeeledFlags, opts ...Option) Policy {
	cfg := defaultConfig()
	for _, o := range opts {
		o(&cfg)
	}
	dest := HelpToStderr
	if peeled.Stdout || envTruthy(cfg.env, "TIED_CLI_HELP_STDOUT") {
		dest = HelpToStdout
	}
	plain := peeled.Plain || peeled.NoColor ||
		envTruthy(cfg.env, "TIED_CLI_PLAIN") ||
		envNonEmpty(cfg.env, "NO_COLOR")
	hw := helpWriterFor(dest, cfg.stdout, cfg.stderr)
	if !plain && cfg.isTTY != nil && !cfg.isTTY(hw) {
		plain = true
	}
	return Policy{
		HelpDest: dest,
		Plain:    plain,
		stdout:   cfg.stdout,
		stderr:   cfg.stderr,
		isTTY:    cfg.isTTY,
	}
}

func helpWriterFor(dest HelpDestination, stdout, stderr io.Writer) io.Writer {
	if dest == HelpToStdout {
		return stdout
	}
	return stderr
}

// HelpWriter returns the writer for successful help text.
func (p Policy) HelpWriter() io.Writer {
	return helpWriterFor(p.HelpDest, p.stdout, p.stderr)
}

// DiagWriter returns the writer for diagnostics (always stderr in the contract).
func (p Policy) DiagWriter() io.Writer {
	return p.stderr
}

// DataWriter returns the writer for machine output (stdout).
func (p Policy) DataWriter() io.Writer {
	return p.stdout
}

func envTruthy(env map[string]string, key string) bool {
	v, ok := env[key]
	if !ok {
		return false
	}
	switch v {
	case "", "0", "false", "FALSE", "no", "NO":
		return false
	default:
		return true
	}
}

func envNonEmpty(env map[string]string, key string) bool {
	v, ok := env[key]
	return ok && v != ""
}

func environMap(environ []string) map[string]string {
	m := make(map[string]string, len(environ))
	for _, e := range environ {
		k, v, ok := splitEnv(e)
		if ok {
			m[k] = v
		}
	}
	return m
}

func splitEnv(e string) (key, val string, ok bool) {
	for i := 0; i < len(e); i++ {
		if e[i] == '=' {
			return e[:i], e[i+1:], true
		}
	}
	return "", "", false
}
