// ⭐ EXTRACT-009: Testing utility extraction - 🔧
package testutil

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// DefaultCliTestHelper provides standard CLI testing functionality.
type DefaultCliTestHelper struct{}

// NewCliTestHelper creates a new CLI test helper.
func NewCliTestHelper() CliTestHelper {
	return &DefaultCliTestHelper{}
}

// CreateTestCommand creates a test cobra command.
// Extracted from main_test.go createTestRootCmd pattern.
//
// ⭐ EXTRACT-009: Test command creation - 🔧
func (h *DefaultCliTestHelper) CreateTestCommand(name string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: "Test command for " + name,
		Run:   runFunc,
	}

	// Disable usage on error for cleaner test output
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	return cmd
}

// ExecuteCommand executes a command with arguments and returns output and error.
// This captures both stdout and stderr for testing.
//
// ⭐ EXTRACT-009: Command execution with output capture - 🔧
func (h *DefaultCliTestHelper) ExecuteCommand(t *testing.T, cmd *cobra.Command, args []string) (string, error) {
	t.Helper()

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	// Set arguments
	cmd.SetArgs(args)

	// Execute command
	err := cmd.Execute()

	// Combine stdout and stderr for complete output
	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderr.String()
	}

	return output, err
}

// CaptureOutput captures stdout/stderr during function execution.
//
// ⭐ EXTRACT-009: Output capture utility - 🔧
func (h *DefaultCliTestHelper) CaptureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	t.Helper()

	// Save original stdout/stderr
	origStdout := os.Stdout
	origStderr := os.Stderr

	// Create pipes for capturing output
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}

	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}

	// Replace stdout/stderr
	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter

	// Channels to collect output
	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	// Start goroutines to read output
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, stdoutReader)
		stdoutChan <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, stderrReader)
		stderrChan <- buf.String()
	}()

	// Execute function
	fn()

	// Close writers to signal end of output
	stdoutWriter.Close()
	stderrWriter.Close()

	// Restore original stdout/stderr
	os.Stdout = origStdout
	os.Stderr = origStderr

	// Collect output
	stdout = <-stdoutChan
	stderr = <-stderrChan

	return stdout, stderr
}

// CreateTestRootCommand creates a test root command with common setup.
// Extracted from main_test.go createTestRootCmd function.
//
// ⭐ EXTRACT-009: Test root command creation - 🔧
func CreateTestRootCommand(appName string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   appName,
		Short: "Test application for " + appName,
		Long:  "Test application for " + appName + " with common CLI patterns",
	}

	// Common CLI setup
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	return rootCmd
}

// ExecuteCommandWithInput executes a command with simulated stdin input.
//
// ⭐ EXTRACT-009: Command execution with input - 🔧
func ExecuteCommandWithInput(t *testing.T, cmd *cobra.Command, args []string, input string) (string, error) {
	t.Helper()

	// Save original stdin
	origStdin := os.Stdin

	// Create pipe for input
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create input pipe: %v", err)
	}

	// Replace stdin
	os.Stdin = reader

	// Write input and close writer
	go func() {
		defer writer.Close()
		writer.Write([]byte(input))
	}()

	// Execute command
	helper := NewCliTestHelper()
	output, cmdErr := helper.ExecuteCommand(t, cmd, args)

	// Restore stdin
	os.Stdin = origStdin
	reader.Close()

	return output, cmdErr
}

// AssertCommandSuccess asserts that a command executed successfully.
//
// ⭐ EXTRACT-009: Command success assertion - 🔧
func AssertCommandSuccess(t *testing.T, cmd *cobra.Command, args []string) string {
	t.Helper()

	helper := NewCliTestHelper()
	output, err := helper.ExecuteCommand(t, cmd, args)

	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	return output
}

// AssertCommandError asserts that a command failed with an error.
//
// ⭐ EXTRACT-009: Command error assertion - 🔧
func AssertCommandError(t *testing.T, cmd *cobra.Command, args []string) (string, error) {
	t.Helper()

	helper := NewCliTestHelper()
	output, err := helper.ExecuteCommand(t, cmd, args)

	if err == nil {
		t.Fatalf("Expected command to fail, but it succeeded\nOutput: %s", output)
	}

	return output, err
}

// AssertCommandOutput asserts that a command produces expected output.
//
// ⭐ EXTRACT-009: Command output assertion - 🔧
func AssertCommandOutput(t *testing.T, cmd *cobra.Command, args []string, expectedOutput string) {
	t.Helper()

	output := AssertCommandSuccess(t, cmd, args)

	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain %q, got %q", expectedOutput, output)
	}
}

// WithTestCommand executes a function with a test command and cleans up automatically.
//
// ⭐ EXTRACT-009: Test command with cleanup - 🔧
func WithTestCommand(t *testing.T, name string, runFunc func(cmd *cobra.Command, args []string), fn func(cmd *cobra.Command)) {
	t.Helper()

	helper := NewCliTestHelper()
	cmd := helper.CreateTestCommand(name, runFunc)
	fn(cmd)
}

// CreateSubcommand creates a subcommand and adds it to a parent command.
//
// ⭐ EXTRACT-009: Subcommand creation - 🔧
func CreateSubcommand(parent *cobra.Command, name, short string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	subCmd := &cobra.Command{
		Use:   name,
		Short: short,
		Run:   runFunc,
	}

	subCmd.SilenceUsage = true
	subCmd.SilenceErrors = true

	parent.AddCommand(subCmd)
	return subCmd
}

// AddStringFlag adds a string flag to a command.
//
// ⭐ EXTRACT-009: String flag addition - 🔧
func AddStringFlag(cmd *cobra.Command, name, shorthand, defaultValue, usage string) {
	cmd.Flags().StringP(name, shorthand, defaultValue, usage)
}

// AddBoolFlag adds a boolean flag to a command.
//
// ⭐ EXTRACT-009: Boolean flag addition - 🔧
func AddBoolFlag(cmd *cobra.Command, name, shorthand string, defaultValue bool, usage string) {
	cmd.Flags().BoolP(name, shorthand, defaultValue, usage)
}

// AddIntFlag adds an integer flag to a command.
//
// ⭐ EXTRACT-009: Integer flag addition - 🔧
func AddIntFlag(cmd *cobra.Command, name, shorthand string, defaultValue int, usage string) {
	cmd.Flags().IntP(name, shorthand, defaultValue, usage)
}

// Package-level convenience functions

var defaultCliHelper = NewCliTestHelper()

// CreateTestCommand is a package-level convenience function for creating test commands.
//
// ⭐ EXTRACT-009: Package-level test command - 🔧
func CreateTestCommand(name string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	return defaultCliHelper.CreateTestCommand(name, runFunc)
}

// ExecuteCommand is a package-level convenience function for executing commands.
//
// ⭐ EXTRACT-009: Package-level command execution - 🔧
func ExecuteCommand(t *testing.T, cmd *cobra.Command, args []string) (string, error) {
	return defaultCliHelper.ExecuteCommand(t, cmd, args)
}

// CaptureOutput is a package-level convenience function for capturing output.
//
// ⭐ EXTRACT-009: Package-level output capture - 🔧
func CaptureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	return defaultCliHelper.CaptureOutput(t, fn)
}
