// Package cli tests for CLI framework components.
package cli
import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: construct DefaultCommandBuilder with a non-nil FlagManager (defaulting when nil).
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: create a subcommand with use, short, and long description fields.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: attach RunE handler to command and return command for fluent chaining.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: register each child command on parent and return parent.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: construct DefaultRootCommandBuilder with default flag and version managers when nil.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: build root command from AppInfo with formatted version, version template, and help when invoked without subcommand.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: materialize command from CommandTemplate including handler, prerun, postrun, flags, and subcommands.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: add persistent help and verbose flags via provided or embedded FlagManager.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: wire all managers and builders then create root command from AppInfo.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: attach subcommand to application root.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: run root command without extra signal context.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: run root with signal-handled context on command; exit non-zero on error.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: register dry-run, note, config, incremental, and list flags when FlagSet pointers are non-nil.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: add persistent help and verbose flags to command.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: return cancellable child context defaulting parent to background.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: parse duration string; on parse failure return cancel-only context without timeout.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: listen for interrupt/terminate and invoke cancel when signal received.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: wrap function as operation that returns canceled when Cancel was called.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: create cancelable context that cancels on INT/TERM or parent done.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: log Describe output when DryRun; otherwise run operation Execute.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: format version, build date, and platform into display string.
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: create version subcommand that prints FormatVersion on run.
func TestBuildInfo(t *testing.T) {
	info := BuildInfo{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Commit:   "abc123",
		Platform: "linux/amd64",
	}

	if info.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", info.Version)
	}
}

func TestAppInfo(t *testing.T) {
	build := BuildInfo{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Commit:   "abc123",
		Platform: "linux/amd64",
	}

	app := AppInfo{
		Name:  "testapp",
		Short: "Test application",
		Long:  "A test application for testing",
		Build: build,
	}

	if app.Name != "testapp" {
		t.Errorf("Expected name testapp, got %s", app.Name)
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify FormatVersion, CreateVersionTemplate, and CreateVersionCommand outputs.
func TestVersionManager(t *testing.T) {
	vm := NewVersionManager()

	info := BuildInfo{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Commit:   "abc123",
		Platform: "linux/amd64",
	}

	formatted := vm.FormatVersion(info)
	expected := "1.0.0 (compiled 2024-01-01) [linux/amd64]"
	if formatted != expected {
		t.Errorf("Expected %s, got %s", expected, formatted)
	}

	template := vm.CreateVersionTemplate(info)
	expectedTemplate := "version 1.0.0 (compiled 2024-01-01) [linux/amd64]\n"
	if template != expectedTemplate {
		t.Errorf("Expected %s, got %s", expectedTemplate, template)
	}

	cmd := vm.CreateVersionCommand(info)
	if cmd.Use != "version" {
		t.Errorf("Expected version command name, got %s", cmd.Use)
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify dry-run skips Execute and logs [DRY-RUN] prefix; non-dry-run runs operation.
func TestDryRunManager(t *testing.T) {
	drm := NewDryRunManager()

	var output bytes.Buffer
	ctx := CommandContext{
		Context:     context.Background(),
		Output:      &output,
		ErrorOutput: &output,
		DryRun:      true,
	}

	executed := false
	op := NewSimpleDryRunOperation("Test operation", func(ctx CommandContext) error {
		executed = true
		return nil
	})

	err := drm.Execute(ctx, op)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if executed {
		t.Error("Operation should not have been executed in dry-run mode")
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "[DRY-RUN] Test operation") {
		t.Errorf("Expected dry-run output, got: %s", outputStr)
	}

	// Test actual execution
	output.Reset()
	ctx.DryRun = false
	executed = false

	err = drm.Execute(ctx, op)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !executed {
		t.Error("Operation should have been executed when not in dry-run mode")
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify Create yields active context until cancel; WithTimeout ends after duration.
func TestContextManager(t *testing.T) {
	cm := NewContextManager()

	parent := context.Background()
	ctx, cancel := cm.Create(parent)
	defer cancel()

	if ctx == nil {
		t.Error("Context should not be nil")
	}

	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	default:
		// Expected
	}

	cancel()

	select {
	case <-ctx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be done after cancel")
	}
}

func TestContextManagerWithTimeout(t *testing.T) {
	cm := NewContextManager()

	parent := context.Background()
	ctx, cancel := cm.WithTimeout(parent, "100ms")
	defer cancel()

	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	case <-time.After(50 * time.Millisecond):
		// Expected - context should still be active
	}

	select {
	case <-ctx.Done():
		// Expected - context should timeout
	case <-time.After(200 * time.Millisecond):
		t.Error("Context should have timed out")
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify AddDryRunFlag and AddNoteFlag register expected flags.
func TestFlagManager(t *testing.T) {
	fm := NewFlagManager()
	cmd := &cobra.Command{Use: "test"}

	var dryRun bool
	err := fm.AddDryRunFlag(cmd, &dryRun)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	flag := cmd.PersistentFlags().Lookup("dry-run")
	if flag == nil {
		t.Error("Dry-run flag should be added")
	}

	var note string
	err = fm.AddNoteFlag(cmd, &note)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	noteFlag := cmd.Flags().Lookup("note")
	if noteFlag == nil {
		t.Error("Note flag should be added")
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify NewCommand sets use/short and WithHandler invokes RunE.
func TestCommandBuilder(t *testing.T) {
	fm := NewFlagManager()
	cb := NewCommandBuilder(fm)

	cmd := cb.NewCommand("test", "Test command", "A test command")
	if cmd.Use != "test" {
		t.Errorf("Expected command name 'test', got %s", cmd.Use)
	}

	if cmd.Short != "Test command" {
		t.Errorf("Expected short description 'Test command', got %s", cmd.Short)
	}

	handlerCalled := false
	handler := func(cmd *cobra.Command, args []string) error {
		handlerCalled = true
		return nil
	}

	cb.WithHandler(cmd, handler)
	if cmd.RunE == nil {
		t.Error("Handler should be set")
	}

	// Test handler execution
	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !handlerCalled {
		t.Error("Handler should have been called")
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify NewRootCommand sets name, short, and non-empty version from AppInfo.
func TestRootCommandBuilder(t *testing.T) {
	fm := NewFlagManager()
	vm := NewVersionManager()
	rb := NewRootCommandBuilder(fm, vm)

	buildInfo := BuildInfo{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Commit:   "abc123",
		Platform: "linux/amd64",
	}

	appInfo := AppInfo{
		Name:  "testapp",
		Short: "Test application",
		Long:  "A test application for testing",
		Build: buildInfo,
	}

	cmd := rb.NewRootCommand(appInfo)
	if cmd.Use != "testapp" {
		t.Errorf("Expected command name 'testapp', got %s", cmd.Use)
	}

	if cmd.Short != "Test application" {
		t.Errorf("Expected short description 'Test application', got %s", cmd.Short)
	}

	if cmd.Version == "" {
		t.Error("Version should be set")
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify NewCLIApp preserves AppInfo and AddCommand registers subcommand on root.
func TestCLIApp(t *testing.T) {
	buildInfo := BuildInfo{
		Version:  "1.0.0",
		Date:     "2024-01-01",
		Commit:   "abc123",
		Platform: "linux/amd64",
	}

	appInfo := AppInfo{
		Name:  "testapp",
		Short: "Test application",
		Long:  "A test application for testing",
		Build: buildInfo,
	}

	app := NewCLIApp(appInfo)
	if app.Info.Name != "testapp" {
		t.Errorf("Expected app name 'testapp', got %s", app.Info.Name)
	}

	if app.RootCommand == nil {
		t.Error("Root command should not be nil")
	}

	// Add a test command
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run: func(cmd *cobra.Command, args []string) {
			// Test command
		},
	}

	app.AddCommand(testCmd)

	// Verify command was added
	commands := app.RootCommand.Commands()
	found := false
	for _, cmd := range commands {
		if cmd.Use == "test" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Test command should have been added to root command")
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify Execute runs until Cancel then returns canceled.
func TestCancellableOperation(t *testing.T) {
	executed := false
	op := NewCancellableOperation(func(ctx context.Context) error {
		executed = true
		return nil
	})

	// Test normal execution
	err := op.Execute(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !executed {
		t.Error("Operation should have been executed")
	}

	// Test cancellation
	executed = false
	err = op.Cancel()
	if err != nil {
		t.Errorf("Unexpected error during cancel: %v", err)
	}

	err = op.Execute(context.Background())
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}

	if executed {
		t.Error("Operation should not have been executed after cancellation")
	}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: verify WithSignalHandling context is active until manual cancel.
func TestWithSignalHandling(t *testing.T) {
	ctx, cancel := WithSignalHandling(context.Background())
	defer cancel()

	if ctx == nil {
		t.Error("Context should not be nil")
	}

	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	default:
		// Expected
	}

	// Test manual cancellation
	cancel()

	select {
	case <-ctx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be done after cancel")
	}
}
