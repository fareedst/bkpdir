// This file is part of bkpdir
//
// Package testutil provides testing infrastructure for complex scenarios,
// specifically permission simulation utilities for controlled testing.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// 🔺 TEST-INFRA-001-C: Permission testing framework - 🔧
// DECISION-REF: DEC-004 (Error handling)
// IMPLEMENTATION-NOTES: Uses temporary directories with controlled permissions for safe testing

package testutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// PermissionSimulator provides controlled permission simulation for testing
type PermissionSimulator struct {
	mu                sync.RWMutex
	tempDir           string
	originalPerms     map[string]os.FileMode
	permissionChanges map[string]PermissionChange
	scenarios         map[string]PermissionScenario
	isEnabled         bool
	stats             PermissionStats
	osType            OSType
	cleanup           []func() error
}

// OSType represents the operating system type for cross-platform permission handling
type OSType int

const (
	OSUnix OSType = iota
	OSWindows
)

// PermissionChange represents a permission change to apply during testing
type PermissionChange struct {
	Path         string
	NewMode      os.FileMode
	OriginalMode os.FileMode
	ChangeTime   time.Time
	Description  string
}

// PermissionScenario defines a complete permission testing scenario
type PermissionScenario struct {
	Name             string
	Description      string
	SetupFiles       []FileSetup
	PermissionSteps  []PermissionStep
	ExpectedErrors   []string
	RestoreOnCleanup bool
}

// FileSetup defines how to create files/directories for testing
type FileSetup struct {
	Path          string
	IsDir         bool
	Content       []byte
	Mode          os.FileMode
	CreateParents bool
}

// PermissionStep defines a permission change during a scenario
type PermissionStep struct {
	Path        string
	NewMode     os.FileMode
	Description string
	DelayAfter  time.Duration
	ExpectError bool
}

// PermissionStats tracks statistics about permission simulation
type PermissionStats struct {
	TotalChanges       int
	SuccessfulChanges  int
	FailedChanges      int
	FilesCreated       int
	DirectoriesCreated int
	PermissionErrors   int
	RestorationErrors  int
	CrossPlatformTests int
	TempDirsCreated    int
}

// PermissionConfig configures permission simulation parameters
type PermissionConfig struct {
	TempDirPrefix       string
	CreateSubDirs       []string
	InitialPermissions  map[string]os.FileMode
	EnableCrossPlatform bool
	AutoCleanup         bool
	MaxTestDuration     time.Duration
}

// NewPermissionSimulator creates a new permission simulator
func NewPermissionSimulator(config PermissionConfig) (*PermissionSimulator, error) {
	// Detect OS type
	osType := OSUnix
	if runtime.GOOS == "windows" {
		osType = OSWindows
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", config.TempDirPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	ps := &PermissionSimulator{
		tempDir:           tempDir,
		originalPerms:     make(map[string]os.FileMode),
		permissionChanges: make(map[string]PermissionChange),
		scenarios:         getBuiltinScenarios(),
		isEnabled:         true,
		osType:            osType,
		cleanup:           make([]func() error, 0),
		stats: PermissionStats{
			TempDirsCreated: 1,
		},
	}

	// Create subdirectories if specified
	for _, subdir := range config.CreateSubDirs {
		path := filepath.Join(tempDir, subdir)
		if err := os.MkdirAll(path, 0755); err != nil {
			ps.Cleanup()
			return nil, fmt.Errorf("failed to create subdirectory %s: %w", path, err)
		}
		ps.stats.DirectoriesCreated++
	}

	// Apply initial permissions
	for path, mode := range config.InitialPermissions {
		fullPath := filepath.Join(tempDir, path)
		if err := ps.SetPermission(fullPath, mode); err != nil {
			ps.Cleanup()
			return nil, fmt.Errorf("failed to set initial permission for %s: %w", path, err)
		}
	}

	return ps, nil
}

// GetTempDir returns the temporary directory path
func (ps *PermissionSimulator) GetTempDir() string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.tempDir
}

// CreateFile creates a file with specified permissions in the temp directory
func (ps *PermissionSimulator) CreateFile(relativePath string, content []byte, mode os.FileMode) (string, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	fullPath := filepath.Join(ps.tempDir, relativePath)

	// Create parent directories if needed
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directories: %w", err)
	}

	// Create file
	if err := os.WriteFile(fullPath, content, mode); err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	// Store original permissions for restoration
	ps.originalPerms[fullPath] = mode
	ps.stats.FilesCreated++

	return fullPath, nil
}

// CreateDir creates a directory with specified permissions in the temp directory
func (ps *PermissionSimulator) CreateDir(relativePath string, mode os.FileMode) (string, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	fullPath := filepath.Join(ps.tempDir, relativePath)

	if err := os.MkdirAll(fullPath, mode); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Store original permissions for restoration
	ps.originalPerms[fullPath] = mode
	ps.stats.DirectoriesCreated++

	return fullPath, nil
}

// SetPermission changes permissions on a file or directory
func (ps *PermissionSimulator) SetPermission(path string, mode os.FileMode) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.setPermissionLocked(path, mode)
}

// setPermissionLocked is the internal implementation that assumes the mutex is already held
func (ps *PermissionSimulator) setPermissionLocked(path string, mode os.FileMode) error {
	// Get current permissions if not already stored
	if _, exists := ps.originalPerms[path]; !exists {
		if info, err := os.Stat(path); err == nil {
			ps.originalPerms[path] = info.Mode()
		}
	}

	// Apply permission change
	if err := os.Chmod(path, mode); err != nil {
		ps.stats.FailedChanges++
		return fmt.Errorf("failed to change permissions: %w", err)
	}

	// Record the change
	ps.permissionChanges[path] = PermissionChange{
		Path:         path,
		NewMode:      mode,
		OriginalMode: ps.originalPerms[path],
		ChangeTime:   time.Now(),
		Description:  fmt.Sprintf("Changed from %v to %v", ps.originalPerms[path], mode),
	}

	ps.stats.TotalChanges++
	ps.stats.SuccessfulChanges++

	return nil
}

// RestorePermissions restores original permissions for all changed files
func (ps *PermissionSimulator) RestorePermissions() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.restorePermissionsLocked()
}

// restorePermissionsLocked is the internal implementation that assumes the mutex is already held
func (ps *PermissionSimulator) restorePermissionsLocked() error {
	var errs []error
	restored := 0

	for path, originalMode := range ps.originalPerms {
		if err := os.Chmod(path, originalMode); err != nil {
			errs = append(errs, fmt.Errorf("failed to restore %s: %w", path, err))
			ps.stats.RestorationErrors++
		} else {
			restored++
			delete(ps.permissionChanges, path)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("restoration errors: %v", errs)
	}

	return nil
}

// SimulatePermissionDenied creates a scenario where permission is denied
func (ps *PermissionSimulator) SimulatePermissionDenied(path string) error {
	var deniedMode os.FileMode

	// Cross-platform permission denial
	if ps.osType == OSWindows {
		// On Windows, remove write permission
		deniedMode = 0444 // Read-only
	} else {
		// On Unix, remove all permissions
		deniedMode = 0000
	}

	return ps.SetPermission(path, deniedMode)
}

// SimulateReadOnlyAccess creates a read-only access scenario
func (ps *PermissionSimulator) SimulateReadOnlyAccess(path string) error {
	return ps.SetPermission(path, 0444)
}

// SimulateWriteOnlyAccess creates a write-only access scenario (Unix only)
func (ps *PermissionSimulator) SimulateWriteOnlyAccess(path string) error {
	if ps.osType == OSWindows {
		// Windows doesn't support write-only, fallback to read-write
		return ps.SetPermission(path, 0666)
	}
	return ps.SetPermission(path, 0222)
}

// SimulateExecuteOnlyAccess creates an execute-only access scenario (Unix only)
func (ps *PermissionSimulator) SimulateExecuteOnlyAccess(path string) error {
	if ps.osType == OSWindows {
		// Windows doesn't support execute-only, fallback to full access
		return ps.SetPermission(path, 0777)
	}
	return ps.SetPermission(path, 0111)
}

// RunScenario executes a predefined permission scenario
func (ps *PermissionSimulator) RunScenario(scenarioName string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	scenario, exists := ps.scenarios[scenarioName]
	if !exists {
		return fmt.Errorf("scenario %s not found", scenarioName)
	}

	// Setup files and directories
	for _, setup := range scenario.SetupFiles {
		fullPath := filepath.Join(ps.tempDir, setup.Path)

		if setup.CreateParents {
			if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directories for %s: %w", setup.Path, err)
			}
		}

		if setup.IsDir {
			if err := os.MkdirAll(fullPath, setup.Mode); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", setup.Path, err)
			}
			ps.stats.DirectoriesCreated++
		} else {
			if err := os.WriteFile(fullPath, setup.Content, setup.Mode); err != nil {
				return fmt.Errorf("failed to create file %s: %w", setup.Path, err)
			}
			ps.stats.FilesCreated++
		}

		ps.originalPerms[fullPath] = setup.Mode
	}

	// Execute permission steps
	for _, step := range scenario.PermissionSteps {
		fullPath := filepath.Join(ps.tempDir, step.Path)
		err := ps.setPermissionLocked(fullPath, step.NewMode)

		if step.ExpectError && err == nil {
			return fmt.Errorf("expected error for step %s but none occurred", step.Description)
		} else if !step.ExpectError && err != nil {
			return fmt.Errorf("unexpected error for step %s: %w", step.Description, err)
		}

		if step.DelayAfter > 0 {
			time.Sleep(step.DelayAfter)
		}
	}

	return nil
}

// DetectPermissionChange monitors a path for permission changes
func (ps *PermissionSimulator) DetectPermissionChange(path string, interval time.Duration, duration time.Duration) (*PermissionChange, error) {
	// Get initial permissions
	initialInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial permissions: %w", err)
	}
	initialMode := initialInfo.Mode()

	// Monitor for changes
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	timeout := time.After(duration)

	for {
		select {
		case <-ticker.C:
			currentInfo, err := os.Stat(path)
			if err != nil {
				continue // File might be temporarily inaccessible
			}

			currentMode := currentInfo.Mode()
			if currentMode != initialMode {
				return &PermissionChange{
					Path:         path,
					NewMode:      currentMode,
					OriginalMode: initialMode,
					ChangeTime:   time.Now(),
					Description:  fmt.Sprintf("Detected change from %v to %v", initialMode, currentMode),
				}, nil
			}

		case <-timeout:
			return nil, fmt.Errorf("no permission change detected within %v", duration)
		}
	}
}

// GetPermissionError creates platform-specific permission errors
func (ps *PermissionSimulator) GetPermissionError(operation, path string) error {
	if ps.osType == OSWindows {
		// Use Windows-specific error number for access denied
		return &os.PathError{
			Op:   operation,
			Path: path,
			Err:  syscall.Errno(5), // ERROR_ACCESS_DENIED on Windows
		}
	} else {
		return &os.PathError{
			Op:   operation,
			Path: path,
			Err:  syscall.EACCES,
		}
	}
}

// IsPermissionError checks if an error is permission-related (cross-platform)
func (ps *PermissionSimulator) IsPermissionError(err error) bool {
	if err == nil {
		return false
	}

	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		if ps.osType == OSWindows {
			// Check for Windows access denied error (errno 5)
			return pathErr.Err == syscall.Errno(5)
		} else {
			return pathErr.Err == syscall.EACCES || pathErr.Err == syscall.EPERM
		}
	}

	// String-based detection as fallback
	errStr := err.Error()
	permissionPatterns := []string{
		"permission denied",
		"access is denied",
		"access denied",
		"operation not permitted",
		"not permitted",
		"EACCES",
		"EPERM",
	}

	for _, pattern := range permissionPatterns {
		if containsIgnoreCase(errStr, pattern) {
			return true
		}
	}

	return false
}

// GeneratePermissionScenarios creates systematic permission combinations
func (ps *PermissionSimulator) GeneratePermissionScenarios() []PermissionScenario {
	scenarios := make([]PermissionScenario, 0)

	// File permission scenarios
	fileModes := []os.FileMode{0000, 0400, 0200, 0100, 0600, 0644, 0755, 0777}
	for _, mode := range fileModes {
		scenarios = append(scenarios, PermissionScenario{
			Name:        fmt.Sprintf("file_mode_%04o", mode),
			Description: fmt.Sprintf("Test file operations with mode %04o", mode),
			SetupFiles: []FileSetup{
				{
					Path:    "test_file.txt",
					IsDir:   false,
					Content: []byte("test content"),
					Mode:    mode,
				},
			},
		})
	}

	// Directory permission scenarios
	dirModes := []os.FileMode{0000, 0400, 0500, 0600, 0700, 0755, 0777}
	for _, mode := range dirModes {
		scenarios = append(scenarios, PermissionScenario{
			Name:        fmt.Sprintf("dir_mode_%04o", mode),
			Description: fmt.Sprintf("Test directory operations with mode %04o", mode),
			SetupFiles: []FileSetup{
				{
					Path:  "test_dir",
					IsDir: true,
					Mode:  mode,
				},
			},
		})
	}

	// Nested permission scenarios
	scenarios = append(scenarios, PermissionScenario{
		Name:        "nested_permission_denial",
		Description: "Test nested directory with permission denial",
		SetupFiles: []FileSetup{
			{
				Path:  "parent",
				IsDir: true,
				Mode:  0755,
			},
			{
				Path:  "parent/restricted",
				IsDir: true,
				Mode:  0000,
			},
			{
				Path:    "parent/restricted/file.txt",
				IsDir:   false,
				Content: []byte("restricted content"),
				Mode:    0644,
			},
		},
	})

	return scenarios
}

// GetStats returns permission simulation statistics
func (ps *PermissionSimulator) GetStats() PermissionStats {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.stats
}

// GetPermissionChanges returns all recorded permission changes
func (ps *PermissionSimulator) GetPermissionChanges() map[string]PermissionChange {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	changes := make(map[string]PermissionChange)
	for k, v := range ps.permissionChanges {
		changes[k] = v
	}
	return changes
}

// Cleanup removes temporary files and restores permissions
func (ps *PermissionSimulator) Cleanup() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var errs []error

	// Restore all permissions first
	if err := ps.restorePermissionsLocked(); err != nil {
		errs = append(errs, err)
	}

	// Remove temporary directory
	if ps.tempDir != "" {
		if err := os.RemoveAll(ps.tempDir); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove temp directory: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %v", errs)
	}

	return nil
}

// Helper functions

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// getBuiltinScenarios returns predefined permission testing scenarios
func getBuiltinScenarios() map[string]PermissionScenario {
	scenarios := map[string]PermissionScenario{
		"basic_permission_denial": {
			Name:        "basic_permission_denial",
			Description: "Basic permission denial scenario",
			SetupFiles: []FileSetup{
				{
					Path:    "test.txt",
					IsDir:   false,
					Content: []byte("test content"),
					Mode:    0644,
				},
			},
			PermissionSteps: []PermissionStep{
				{
					Path:        "test.txt",
					NewMode:     0000,
					Description: "Remove all permissions",
					ExpectError: false,
				},
			},
			ExpectedErrors:   []string{"permission denied", "access denied"},
			RestoreOnCleanup: true,
		},
		"directory_access_denied": {
			Name:        "directory_access_denied",
			Description: "Directory access denial scenario",
			SetupFiles: []FileSetup{
				{
					Path:  "testdir",
					IsDir: true,
					Mode:  0755,
				},
				{
					Path:    "testdir/file.txt",
					IsDir:   false,
					Content: []byte("content"),
					Mode:    0644,
				},
			},
			PermissionSteps: []PermissionStep{
				{
					Path:        "testdir",
					NewMode:     0000,
					Description: "Remove directory access",
					ExpectError: false,
				},
			},
			ExpectedErrors:   []string{"permission denied"},
			RestoreOnCleanup: true,
		},
		"mixed_permissions": {
			Name:        "mixed_permissions",
			Description: "Mixed file and directory permissions",
			SetupFiles: []FileSetup{
				{
					Path:  "readable_dir",
					IsDir: true,
					Mode:  0755,
				},
				{
					Path:  "restricted_dir",
					IsDir: true,
					Mode:  0700,
				},
				{
					Path:    "readable_dir/public.txt",
					IsDir:   false,
					Content: []byte("public content"),
					Mode:    0644,
				},
				{
					Path:    "restricted_dir/private.txt",
					IsDir:   false,
					Content: []byte("private content"),
					Mode:    0600,
				},
			},
			PermissionSteps: []PermissionStep{
				{
					Path:        "restricted_dir",
					NewMode:     0000,
					Description: "Restrict directory access",
					ExpectError: false,
				},
				{
					Path:        "readable_dir/public.txt",
					NewMode:     0000,
					Description: "Remove file access",
					ExpectError: false,
				},
			},
			RestoreOnCleanup: true,
		},
	}

	return scenarios
}

// PermissionTestHelper provides high-level testing utilities
type PermissionTestHelper struct {
	simulator *PermissionSimulator
}

// NewPermissionTestHelper creates a new permission test helper
func NewPermissionTestHelper() (*PermissionTestHelper, error) {
	config := PermissionConfig{
		TempDirPrefix:       "bkpdir_perm_test_",
		EnableCrossPlatform: true,
		AutoCleanup:         true,
		MaxTestDuration:     time.Minute * 5,
	}

	simulator, err := NewPermissionSimulator(config)
	if err != nil {
		return nil, err
	}

	return &PermissionTestHelper{
		simulator: simulator,
	}, nil
}

// WithPermissionDenied executes a function with permission denied on the specified path
func (pth *PermissionTestHelper) WithPermissionDenied(relativePath string, fn func(path string) error) error {
	fullPath := filepath.Join(pth.simulator.GetTempDir(), relativePath)

	// Create the file/directory if it doesn't exist
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
			return fmt.Errorf("failed to create test file: %w", err)
		}
	}

	// Apply permission denial
	if err := pth.simulator.SimulatePermissionDenied(fullPath); err != nil {
		return fmt.Errorf("failed to simulate permission denial: %w", err)
	}

	// Execute function
	err := fn(fullPath)

	// Restore permissions
	if restoreErr := pth.simulator.RestorePermissions(); restoreErr != nil {
		if err == nil {
			err = restoreErr
		}
	}

	return err
}

// Cleanup cleans up the test helper
func (pth *PermissionTestHelper) Cleanup() error {
	return pth.simulator.Cleanup()
}

// GetSimulator returns the underlying simulator for advanced usage
func (pth *PermissionTestHelper) GetSimulator() *PermissionSimulator {
	return pth.simulator
}
