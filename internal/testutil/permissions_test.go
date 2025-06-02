package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestNewPermissionSimulator(t *testing.T) {
	config := PermissionConfig{
		TempDirPrefix: "test_permissions_",
		CreateSubDirs: []string{"subdir1", "subdir2"},
		InitialPermissions: map[string]os.FileMode{
			"subdir1": 0755,
			"subdir2": 0700,
		},
		EnableCrossPlatform: true,
		AutoCleanup:         true,
		MaxTestDuration:     time.Minute,
	}

	ps, err := NewPermissionSimulator(config)
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	// Verify temp directory exists
	tempDir := ps.GetTempDir()
	if tempDir == "" {
		t.Error("Temp directory should not be empty")
	}

	// Check if temp directory actually exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Temp directory should exist")
	}

	// Verify subdirectories were created
	for _, subDir := range config.CreateSubDirs {
		subDirPath := filepath.Join(tempDir, subDir)
		if _, err := os.Stat(subDirPath); os.IsNotExist(err) {
			t.Errorf("Subdirectory %s should exist", subDir)
		}
	}
}

func TestPermissionSimulator_CreateFile(t *testing.T) {
	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_create_file_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	content := []byte("test file content")
	mode := os.FileMode(0644)

	filePath, err := ps.CreateFile("test.txt", content, mode)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("File should exist after creation")
	}

	// Verify file content
	readContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(readContent) != string(content) {
		t.Errorf("File content mismatch. Expected: %s, Got: %s", content, readContent)
	}

	// Verify file permissions (Unix-like systems)
	if runtime.GOOS != "windows" {
		info, err := os.Stat(filePath)
		if err != nil {
			t.Fatalf("Failed to get file info: %v", err)
		}
		if info.Mode().Perm() != mode {
			t.Errorf("File permissions mismatch. Expected: %o, Got: %o", mode, info.Mode().Perm())
		}
	}
}

func TestPermissionSimulator_CreateDir(t *testing.T) {
	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_create_dir_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	mode := os.FileMode(0755)

	dirPath, err := ps.CreateDir("testdir", mode)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Verify directory exists
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		t.Error("Directory should exist after creation")
	}

	if !info.IsDir() {
		t.Error("Created path should be a directory")
	}

	// Verify directory permissions (Unix-like systems)
	if runtime.GOOS != "windows" {
		if info.Mode().Perm() != mode {
			t.Errorf("Directory permissions mismatch. Expected: %o, Got: %o", mode, info.Mode().Perm())
		}
	}
}

func TestPermissionSimulator_SetPermission(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Permission testing not fully supported on Windows")
	}

	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_set_permission_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	// Create a file first
	filePath, err := ps.CreateFile("test.txt", []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Change permissions
	newMode := os.FileMode(0400)
	err = ps.SetPermission(filePath, newMode)
	if err != nil {
		t.Fatalf("Failed to set permission: %v", err)
	}

	// Verify new permissions
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if info.Mode().Perm() != newMode {
		t.Errorf("Permission change failed. Expected: %o, Got: %o", newMode, info.Mode().Perm())
	}
}

func TestPermissionSimulator_RestorePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Permission testing not fully supported on Windows")
	}

	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_restore_permissions_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	originalMode := os.FileMode(0644)
	changedMode := os.FileMode(0400)

	// Create a file
	filePath, err := ps.CreateFile("test.txt", []byte("content"), originalMode)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Change permissions
	err = ps.SetPermission(filePath, changedMode)
	if err != nil {
		t.Fatalf("Failed to set permission: %v", err)
	}

	// Restore permissions
	err = ps.RestorePermissions()
	if err != nil {
		t.Fatalf("Failed to restore permissions: %v", err)
	}

	// Verify permissions were restored
	info, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to get file info: %v", err)
	}

	if info.Mode().Perm() != originalMode {
		t.Errorf("Permission restoration failed. Expected: %o, Got: %o", originalMode, info.Mode().Perm())
	}
}

func TestPermissionSimulator_SimulatePermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Permission testing not fully supported on Windows")
	}

	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_permission_denied_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	// Create a file
	filePath, err := ps.CreateFile("test.txt", []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Simulate permission denied
	err = ps.SimulatePermissionDenied(filePath)
	if err != nil {
		t.Fatalf("Failed to simulate permission denied: %v", err)
	}

	// Try to write to the file (should fail)
	writeErr := os.WriteFile(filePath, []byte("new content"), 0644)
	if writeErr == nil {
		t.Error("Expected permission denied error when writing to file")
	}

	// Verify it's a permission error
	if !ps.IsPermissionError(writeErr) {
		t.Errorf("Expected permission error, got: %v", writeErr)
	}
}

func TestPermissionSimulator_GetPermissionError(t *testing.T) {
	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_permission_error_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	testPath := "/test/path"
	testOp := "open"

	permErr := ps.GetPermissionError(testOp, testPath)

	// Verify it's a PathError
	if pathErr, ok := permErr.(*os.PathError); ok {
		if pathErr.Op != testOp {
			t.Errorf("Expected operation %s, got %s", testOp, pathErr.Op)
		}
		if pathErr.Path != testPath {
			t.Errorf("Expected path %s, got %s", testPath, pathErr.Path)
		}
	} else {
		t.Error("Expected *os.PathError")
	}

	// Verify it's detected as permission error
	if !ps.IsPermissionError(permErr) {
		t.Error("Created error should be detected as permission error")
	}
}

func TestPermissionSimulator_IsPermissionError(t *testing.T) {
	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_is_permission_error_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "permission denied PathError",
			err:      &os.PathError{Op: "open", Path: "/test", Err: syscall.EACCES},
			expected: true,
		},
		{
			name:     "operation not permitted PathError",
			err:      &os.PathError{Op: "chmod", Path: "/test", Err: syscall.EPERM},
			expected: true,
		},
		{
			name:     "other syscall error",
			err:      &os.PathError{Op: "open", Path: "/test", Err: syscall.ENOENT},
			expected: false,
		},
		{
			name:     "permission denied string error",
			err:      &customError{message: "permission denied"},
			expected: true,
		},
		{
			name:     "access denied string error",
			err:      &customError{message: "access denied"},
			expected: true,
		},
		{
			name:     "unrelated error",
			err:      &customError{message: "file not found"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ps.IsPermissionError(tt.err)
			if result != tt.expected {
				t.Errorf("IsPermissionError() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestPermissionSimulator_GeneratePermissionScenarios(t *testing.T) {
	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_generate_scenarios_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	scenarios := ps.GeneratePermissionScenarios()

	if len(scenarios) == 0 {
		t.Error("Expected at least some scenarios to be generated")
	}

	// Check for expected scenario types
	foundFileScenario := false
	foundDirScenario := false
	foundNestedScenario := false

	for _, scenario := range scenarios {
		if strings.HasPrefix(scenario.Name, "file_mode_") {
			foundFileScenario = true
		}
		if strings.HasPrefix(scenario.Name, "dir_mode_") {
			foundDirScenario = true
		}
		if scenario.Name == "nested_permission_denial" {
			foundNestedScenario = true
		}
	}

	if !foundFileScenario {
		t.Error("Expected file permission scenarios")
	}
	if !foundDirScenario {
		t.Error("Expected directory permission scenarios")
	}
	if !foundNestedScenario {
		t.Error("Expected nested permission scenario")
	}
}

func TestPermissionSimulator_RunScenario(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Permission testing not fully supported on Windows")
	}

	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_run_scenario_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	// Test built-in scenario
	err = ps.RunScenario("basic_permission_denial")
	if err != nil {
		t.Fatalf("Failed to run scenario: %v", err)
	}

	// Verify stats were updated
	stats := ps.GetStats()
	if stats.FilesCreated == 0 {
		t.Error("Expected files to be created during scenario")
	}
}

func TestPermissionSimulator_DetectPermissionChange(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Permission testing not fully supported on Windows")
	}

	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_detect_change_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	// Create a file
	filePath, err := ps.CreateFile("test.txt", []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Start monitoring in background
	changeDetected := make(chan *PermissionChange, 1)
	go func() {
		change, err := ps.DetectPermissionChange(filePath, 10*time.Millisecond, 100*time.Millisecond)
		if err == nil && change != nil {
			changeDetected <- change
		}
		close(changeDetected)
	}()

	// Wait a bit and then change permissions
	time.Sleep(20 * time.Millisecond)
	err = ps.SetPermission(filePath, 0400)
	if err != nil {
		t.Fatalf("Failed to change permissions: %v", err)
	}

	// Check if change was detected
	select {
	case change := <-changeDetected:
		if change == nil {
			t.Error("Expected permission change to be detected")
		} else {
			if change.Path != filePath {
				t.Errorf("Expected path %s, got %s", filePath, change.Path)
			}
			if change.NewMode != 0400 {
				t.Errorf("Expected new mode 0400, got %o", change.NewMode)
			}
		}
	case <-time.After(200 * time.Millisecond):
		t.Error("Permission change detection timed out")
	}
}

func TestPermissionSimulator_GetStats(t *testing.T) {
	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "test_get_stats_",
		AutoCleanup:   true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	initialStats := ps.GetStats()

	// Create some files and directories
	_, err = ps.CreateFile("test1.txt", []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	_, err = ps.CreateDir("testdir", 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	finalStats := ps.GetStats()

	if finalStats.FilesCreated <= initialStats.FilesCreated {
		t.Error("Expected FilesCreated to increase")
	}

	if finalStats.DirectoriesCreated <= initialStats.DirectoriesCreated {
		t.Error("Expected DirectoriesCreated to increase")
	}
}

func TestPermissionTestHelper(t *testing.T) {
	helper, err := NewPermissionTestHelper()
	if err != nil {
		t.Fatalf("Failed to create PermissionTestHelper: %v", err)
	}
	defer helper.Cleanup()

	// Test WithPermissionDenied
	var testErr error
	err = helper.WithPermissionDenied("test.txt", func(path string) error {
		// Try to read the file (should fail due to permissions)
		_, readErr := os.ReadFile(path)
		testErr = readErr
		return nil
	})

	if err != nil {
		t.Fatalf("WithPermissionDenied failed: %v", err)
	}

	// On Unix-like systems, verify permission error occurred
	if runtime.GOOS != "windows" && testErr == nil {
		t.Error("Expected permission error when reading permission-denied file")
	}
}

func TestCrossPlatformDetection(t *testing.T) {
	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix:       "test_cross_platform_",
		EnableCrossPlatform: true,
		AutoCleanup:         true,
	})
	if err != nil {
		t.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	// Test platform-specific error creation
	err1 := ps.GetPermissionError("open", "/test/path")
	err2 := ps.GetPermissionError("write", "/test/path2")

	if !ps.IsPermissionError(err1) {
		t.Error("Platform-specific error should be detected as permission error")
	}

	if !ps.IsPermissionError(err2) {
		t.Error("Platform-specific error should be detected as permission error")
	}
}

// Helper type for testing string-based error detection
type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}

// Benchmarks
func BenchmarkPermissionSimulator_CreateFile(b *testing.B) {
	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "bench_create_file_",
		AutoCleanup:   true,
	})
	if err != nil {
		b.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	content := []byte("benchmark test content")
	mode := os.FileMode(0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("bench/file_%d.txt", i)
		_, err := ps.CreateFile(filename, content, mode)
		if err != nil {
			b.Fatalf("Failed to create file: %v", err)
		}
	}
}

func BenchmarkPermissionSimulator_IsPermissionError(b *testing.B) {
	ps, err := NewPermissionSimulator(PermissionConfig{
		TempDirPrefix: "bench_is_permission_error_",
		AutoCleanup:   true,
	})
	if err != nil {
		b.Fatalf("Failed to create PermissionSimulator: %v", err)
	}
	defer ps.Cleanup()

	permErr := ps.GetPermissionError("open", "/test/path")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ps.IsPermissionError(permErr)
	}
}
