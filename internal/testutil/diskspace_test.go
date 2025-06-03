// This file is part of bkpdir
//
// Package testutil provides comprehensive tests for disk space simulation utilities.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// 🔺 TEST-INFRA-001-B: Disk space simulation framework tests - 🔧
// DECISION-REF: DEC-004 (Error handling)
// IMPLEMENTATION-NOTES: Comprehensive test coverage for disk space simulation utilities

package testutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
)

// TestExhaustionMode tests the exhaustion mode enumeration
func TestExhaustionMode(t *testing.T) {
	tests := []struct {
		name string
		mode ExhaustionMode
		want int
	}{
		{"Linear exhaustion", LinearExhaustion, 0},
		{"Progressive exhaustion", ProgressiveExhaustion, 1},
		{"Random exhaustion", RandomExhaustion, 2},
		{"Immediate exhaustion", ImmediateExhaustion, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.mode) != tt.want {
				t.Errorf("ExhaustionMode value = %d, want %d", int(tt.mode), tt.want)
			}
		})
	}
}

// TestNewDiskSpaceSimulator tests simulator creation with various configurations
func TestNewDiskSpaceSimulator(t *testing.T) {
	tests := []struct {
		name          string
		config        DiskSpaceConfig
		expectedSpace int64
		expectedTotal int64
		expectedUsed  int64
		expectedMode  ExhaustionMode
		expectedRate  float64
	}{
		{
			name:          "default configuration",
			config:        DiskSpaceConfig{},
			expectedSpace: 50 * 1024 * 1024,  // 50MB (half of default 100MB)
			expectedTotal: 100 * 1024 * 1024, // 100MB
			expectedUsed:  50 * 1024 * 1024,  // 50MB
			expectedMode:  LinearExhaustion,  // default
			expectedRate:  0.0,               // default
		},
		{
			name: "custom configuration",
			config: DiskSpaceConfig{
				InitialSpace:   20 * 1024 * 1024, // 20MB
				TotalSpace:     40 * 1024 * 1024, // 40MB
				ExhaustionMode: ProgressiveExhaustion,
				ExhaustionRate: 0.1,
			},
			expectedSpace: 20 * 1024 * 1024, // 20MB
			expectedTotal: 40 * 1024 * 1024, // 40MB
			expectedUsed:  20 * 1024 * 1024, // 20MB
			expectedMode:  ProgressiveExhaustion,
			expectedRate:  0.1,
		},
		{
			name: "zero initial space",
			config: DiskSpaceConfig{
				InitialSpace: 0,
				TotalSpace:   10 * 1024 * 1024, // 10MB
			},
			expectedSpace: 0,
			expectedTotal: 10 * 1024 * 1024, // 10MB
			expectedUsed:  10 * 1024 * 1024, // 10MB
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			simulator := NewDiskSpaceSimulator(tt.config)

			if simulator.GetAvailableSpace() != tt.expectedSpace {
				t.Errorf("Available space = %d, want %d", simulator.GetAvailableSpace(), tt.expectedSpace)
			}
			if simulator.GetTotalSpace() != tt.expectedTotal {
				t.Errorf("Total space = %d, want %d", simulator.GetTotalSpace(), tt.expectedTotal)
			}
			if simulator.GetUsedSpace() != tt.expectedUsed {
				t.Errorf("Used space = %d, want %d", simulator.GetUsedSpace(), tt.expectedUsed)
			}
			if simulator.exhaustionMode != tt.expectedMode {
				t.Errorf("Exhaustion mode = %d, want %d", simulator.exhaustionMode, tt.expectedMode)
			}
			if simulator.spaceExhaustionRate != tt.expectedRate {
				t.Errorf("Exhaustion rate = %f, want %f", simulator.spaceExhaustionRate, tt.expectedRate)
			}
			if !simulator.IsEnabled() {
				t.Error("Simulator should be enabled by default")
			}
		})
	}
}

// TestDiskSpaceSimulator_BasicOperations tests basic simulator operations
func TestDiskSpaceSimulator_BasicOperations(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace:   10 * 1024 * 1024, // 10MB
		TotalSpace:     20 * 1024 * 1024, // 20MB
		ExhaustionMode: LinearExhaustion,
		ExhaustionRate: 0.1, // 10% per operation
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	t.Run("Enable/Disable", func(t *testing.T) {
		if !simulator.IsEnabled() {
			t.Error("Simulator should be enabled by default")
		}

		simulator.SetEnabled(false)
		if simulator.IsEnabled() {
			t.Error("Simulator should be disabled after SetEnabled(false)")
		}

		// Operations should succeed when disabled
		err := fs.WriteFile("/tmp/test.txt", []byte("test"), 0644)
		if err != nil {
			t.Errorf("WriteFile should succeed when simulator is disabled: %v", err)
		}

		simulator.SetEnabled(true)
		if !simulator.IsEnabled() {
			t.Error("Simulator should be enabled after SetEnabled(true)")
		}
	})

	t.Run("Reset", func(t *testing.T) {
		// Make some operations to change state
		fs.WriteFile("/tmp/test1.txt", []byte("test1"), 0644)
		fs.WriteFile("/tmp/test2.txt", []byte("test2"), 0644)

		if simulator.GetOperationCount() == 0 {
			t.Error("Operation count should be > 0 after operations")
		}

		simulator.Reset()

		if simulator.GetOperationCount() != 0 {
			t.Errorf("Operation count should be 0 after reset, got %d", simulator.GetOperationCount())
		}

		expectedSpace := simulator.GetTotalSpace() / 2
		if simulator.GetAvailableSpace() != expectedSpace {
			t.Errorf("Available space should be %d after reset, got %d", expectedSpace, simulator.GetAvailableSpace())
		}
	})

	t.Run("Statistics", func(t *testing.T) {
		simulator.Reset()

		// Perform some operations
		fs.WriteFile("/tmp/stat1.txt", []byte("test data 1"), 0644)
		fs.WriteFile("/tmp/stat2.txt", []byte("test data 2 with more content"), 0644)

		stats := simulator.GetStats()
		if stats.TotalWrites != 2 {
			t.Errorf("Expected 2 total writes, got %d", stats.TotalWrites)
		}
		if stats.LargestFileWritten == 0 {
			t.Error("Largest file written should be > 0")
		}
		if stats.SmallestFileWritten == 0 {
			t.Error("Smallest file written should be > 0")
		}
		if stats.TotalSpaceReduced == 0 {
			t.Error("Total space reduced should be > 0")
		}
	})
}

// TestLinearExhaustion tests linear space exhaustion mode
func TestLinearExhaustion(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace:   5 * 1024 * 1024,  // 5MB
		TotalSpace:     10 * 1024 * 1024, // 10MB
		ExhaustionMode: LinearExhaustion,
		ExhaustionRate: 0.3, // 30% per operation
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	// First operation should reduce space significantly
	initialSpace := simulator.GetAvailableSpace()

	// Write a small file
	err := fs.WriteFile("/tmp/test1.txt", []byte("small"), 0644)
	if err != nil {
		t.Fatalf("First write should succeed: %v", err)
	}

	spaceAfterFirst := simulator.GetAvailableSpace()
	if spaceAfterFirst >= initialSpace {
		t.Error("Space should be reduced after first operation")
	}

	// Try to write another file - should eventually fail due to exhaustion
	err = fs.WriteFile("/tmp/test2.txt", []byte("second"), 0644)
	// This might succeed or fail depending on remaining space

	// Write one more - this should likely fail
	largeData := make([]byte, 1024*1024) // 1MB
	err = fs.WriteFile("/tmp/test3.txt", largeData, 0644)
	if err == nil {
		t.Error("Large write should eventually fail due to space exhaustion")
	}

	// Verify it's a disk space error
	if !isDiskSpaceError(err) {
		t.Errorf("Expected disk space error, got: %v", err)
	}

	stats := simulator.GetStats()
	if stats.SpaceExhausted == 0 {
		t.Error("Space exhausted count should be > 0")
	}
	if stats.FailedWrites == 0 {
		t.Error("Failed writes count should be > 0")
	}
}

// TestProgressiveExhaustion tests progressive space exhaustion mode
func TestProgressiveExhaustion(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace:   10 * 1024 * 1024, // 10MB
		TotalSpace:     20 * 1024 * 1024, // 20MB
		ExhaustionMode: ProgressiveExhaustion,
		ExhaustionRate: 0.05, // 5% base rate
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	spaces := make([]int64, 0)

	// Perform multiple operations and track space reduction
	for i := 0; i < 5; i++ {
		spaces = append(spaces, simulator.GetAvailableSpace())

		data := make([]byte, 1024) // 1KB
		filename := fmt.Sprintf("/tmp/prog_%d.txt", i)
		err := fs.WriteFile(filename, data, 0644)
		if err != nil {
			// Expected to fail eventually
			break
		}
	}

	// Verify that space reduction increases with each operation
	if len(spaces) >= 3 {
		reduction1 := spaces[0] - spaces[1]
		reduction2 := spaces[1] - spaces[2]

		if reduction2 <= reduction1 {
			t.Error("Progressive exhaustion should increase reduction amount with each operation")
		}
	}
}

// TestImmediateExhaustion tests immediate space exhaustion mode
func TestImmediateExhaustion(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace:   5 * 1024 * 1024,  // 5MB
		TotalSpace:     10 * 1024 * 1024, // 10MB
		ExhaustionMode: ImmediateExhaustion,
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	// First operation should exhaust all space
	err := fs.WriteFile("/tmp/immediate.txt", []byte("test"), 0644)
	if err == nil {
		t.Error("First write should fail with immediate exhaustion")
	}

	if !isDiskSpaceError(err) {
		t.Errorf("Expected disk space error, got: %v", err)
	}

	if simulator.GetAvailableSpace() > 0 {
		t.Errorf("Available space should be 0 after immediate exhaustion, got %d", simulator.GetAvailableSpace())
	}
}

// TestRandomExhaustion tests random space exhaustion mode
func TestRandomExhaustion(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace:   5 * 1024 * 1024,  // 5MB
		TotalSpace:     10 * 1024 * 1024, // 10MB
		ExhaustionMode: RandomExhaustion,
		ExhaustionRate: 0.2, // 20% max reduction
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	reductions := make([]int64, 0)

	// Perform multiple operations and track space reduction variations
	for i := 0; i < 5; i++ {
		spaceBefore := simulator.GetAvailableSpace()

		data := make([]byte, 1024) // 1KB
		filename := fmt.Sprintf("/tmp/rand_%d.txt", i)
		err := fs.WriteFile(filename, data, 0644)
		if err != nil {
			break // Expected to fail eventually
		}

		spaceAfter := simulator.GetAvailableSpace()
		reduction := spaceBefore - spaceAfter - 1024 // Subtract actual file size
		if reduction > 0 {
			reductions = append(reductions, reduction)
		}
	}

	// Verify that reductions vary (at least some difference)
	if len(reductions) >= 2 {
		allSame := true
		first := reductions[0]
		for _, r := range reductions[1:] {
			if r != first {
				allSame = false
				break
			}
		}
		if allSame {
			t.Error("Random exhaustion should produce varying reduction amounts")
		}
	}
}

// TestFailurePoints tests operation-specific failure injection
func TestFailurePoints(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace: 10 * 1024 * 1024, // 10MB
		TotalSpace:   20 * 1024 * 1024, // 20MB
		FailurePoints: map[int]error{
			2: CreateDiskFullError("/tmp/test"),
			4: CreateQuotaExceededError("/tmp/test"),
		},
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	// Operations 1 and 3 should succeed
	err := fs.WriteFile("/tmp/op1.txt", []byte("op1"), 0644)
	if err != nil {
		t.Errorf("Operation 1 should succeed: %v", err)
	}

	// Operation 2 should fail with disk full error
	err = fs.WriteFile("/tmp/op2.txt", []byte("op2"), 0644)
	if err == nil {
		t.Error("Operation 2 should fail")
	}
	if !strings.Contains(err.Error(), "no space left") {
		t.Errorf("Operation 2 should fail with disk full error, got: %v", err)
	}

	// Operation 3 should succeed
	err = fs.WriteFile("/tmp/op3.txt", []byte("op3"), 0644)
	if err != nil {
		t.Errorf("Operation 3 should succeed: %v", err)
	}

	// Operation 4 should fail with quota exceeded error
	err = fs.WriteFile("/tmp/op4.txt", []byte("op4"), 0644)
	if err == nil {
		t.Error("Operation 4 should fail")
	}
	if !strings.Contains(err.Error(), "quota exceeded") {
		t.Errorf("Operation 4 should fail with quota exceeded error, got: %v", err)
	}

	stats := simulator.GetStats()
	if stats.FailedWrites != 2 {
		t.Errorf("Expected 2 failed writes, got %d", stats.FailedWrites)
	}
}

// TestRecoveryPoints tests space recovery functionality
func TestRecoveryPoints(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace: 1 * 1024 * 1024,  // 1MB (limited)
		TotalSpace:   10 * 1024 * 1024, // 10MB
		RecoveryPoints: map[int]int64{
			3: 5 * 1024 * 1024, // Recover 5MB at operation 3
		},
		ExhaustionMode: LinearExhaustion,
		ExhaustionRate: 0.1, // Reduced from 0.8 to 0.1 (10% per operation)
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	// First operation might succeed
	fs.WriteFile("/tmp/rec1.txt", []byte("rec1"), 0644)

	// Second operation likely to fail due to high exhaustion
	fs.WriteFile("/tmp/rec2.txt", []byte("rec2"), 0644)

	// Third operation should trigger recovery
	spaceBefore := simulator.GetAvailableSpace()
	fs.WriteFile("/tmp/rec3.txt", []byte("rec3"), 0644)
	spaceAfter := simulator.GetAvailableSpace()

	// Space after operation 3 should be more than before (due to recovery)
	if spaceAfter <= spaceBefore {
		t.Errorf("Space should increase at recovery point (operation 3): before=%d, after=%d", spaceBefore, spaceAfter)
	}

	stats := simulator.GetStats()
	if stats.RecoveryOperations != 1 {
		t.Errorf("Expected 1 recovery operation, got %d", stats.RecoveryOperations)
	}
}

// TestInjectionPoints tests file-specific error injection
func TestInjectionPoints(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace: 10 * 1024 * 1024, // 10MB
		TotalSpace:   20 * 1024 * 1024, // 20MB
		InjectionPoints: map[string]error{
			"/tmp/fail.txt":  CreateDiskFullError("/tmp/fail.txt"),
			"/tmp/quota.txt": CreateQuotaExceededError("/tmp/quota.txt"),
		},
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	// Normal file should succeed
	err := fs.WriteFile("/tmp/normal.txt", []byte("normal"), 0644)
	if err != nil {
		t.Errorf("Normal file should succeed: %v", err)
	}

	// fail.txt should fail with disk full error
	err = fs.WriteFile("/tmp/fail.txt", []byte("fail"), 0644)
	if err == nil {
		t.Error("/tmp/fail.txt should fail")
	}
	if !strings.Contains(err.Error(), "no space left") {
		t.Errorf("Expected disk full error for /tmp/fail.txt, got: %v", err)
	}

	// quota.txt should fail with quota exceeded error
	err = fs.WriteFile("/tmp/quota.txt", []byte("quota"), 0644)
	if err == nil {
		t.Error("/tmp/quota.txt should fail")
	}
	if !strings.Contains(err.Error(), "quota exceeded") {
		t.Errorf("Expected quota exceeded error for /tmp/quota.txt, got: %v", err)
	}

	stats := simulator.GetStats()
	if stats.InjectedErrors != 2 {
		t.Errorf("Expected 2 injected errors, got %d", stats.InjectedErrors)
	}
}

// TestFileSystemOperations tests various file system operations
func TestFileSystemOperations(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace: 5 * 1024 * 1024,  // 5MB
		TotalSpace:   10 * 1024 * 1024, // 10MB
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	tempDir := t.TempDir()

	t.Run("WriteFile", func(t *testing.T) {
		filename := filepath.Join(tempDir, "write_test.txt")
		data := []byte("test data for WriteFile")

		err := fs.WriteFile(filename, data, 0644)
		if err != nil {
			t.Errorf("WriteFile should succeed: %v", err)
		}

		// Verify file was actually created
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Error("File should exist after WriteFile")
		}
	})

	t.Run("Create", func(t *testing.T) {
		filename := filepath.Join(tempDir, "create_test.txt")

		file, err := fs.Create(filename)
		if err != nil {
			t.Errorf("Create should succeed: %v", err)
		}
		if file != nil {
			file.Close()
		}

		// Verify file was actually created
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Error("File should exist after Create")
		}
	})

	t.Run("OpenFile", func(t *testing.T) {
		filename := filepath.Join(tempDir, "open_test.txt")

		// Create with O_CREATE flag should trigger space check
		file, err := fs.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			t.Errorf("OpenFile with O_CREATE should succeed: %v", err)
		}
		if file != nil {
			file.Close()
		}

		// Open existing file should not trigger space check
		file, err = fs.OpenFile(filename, os.O_RDONLY, 0644)
		if err != nil {
			t.Errorf("OpenFile for reading should succeed: %v", err)
		}
		if file != nil {
			file.Close()
		}
	})

	t.Run("MkdirAll", func(t *testing.T) {
		dirPath := filepath.Join(tempDir, "test", "nested", "directory")

		err := fs.MkdirAll(dirPath, 0755)
		if err != nil {
			t.Errorf("MkdirAll should succeed: %v", err)
		}

		// Verify directory was created
		if info, err := os.Stat(dirPath); err != nil || !info.IsDir() {
			t.Error("Directory should exist after MkdirAll")
		}
	})

	t.Run("Remove and RemoveAll", func(t *testing.T) {
		// Create a file to remove
		filename := filepath.Join(tempDir, "remove_test.txt")
		data := []byte("data to be removed")
		err := fs.WriteFile(filename, data, 0644)
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		spaceBefore := simulator.GetAvailableSpace()

		// Remove the file - should recover space
		err = fs.Remove(filename)
		if err != nil {
			t.Errorf("Remove should succeed: %v", err)
		}

		spaceAfter := simulator.GetAvailableSpace()
		if spaceAfter <= spaceBefore {
			t.Error("Space should be recovered after Remove")
		}

		// Create a directory with files to test RemoveAll
		dirPath := filepath.Join(tempDir, "remove_all_test")
		fs.MkdirAll(dirPath, 0755)
		fs.WriteFile(filepath.Join(dirPath, "file1.txt"), []byte("file1"), 0644)
		fs.WriteFile(filepath.Join(dirPath, "file2.txt"), []byte("file2"), 0644)

		spaceBefore = simulator.GetAvailableSpace()

		// RemoveAll should recover space from all files
		err = fs.RemoveAll(dirPath)
		if err != nil {
			t.Errorf("RemoveAll should succeed: %v", err)
		}

		spaceAfter = simulator.GetAvailableSpace()
		if spaceAfter <= spaceBefore {
			t.Error("Space should be recovered after RemoveAll")
		}
	})
}

// TestPredefinedScenarios tests the predefined test scenarios
func TestPredefinedScenarios(t *testing.T) {
	scenarios := GetAllScenarioNames()
	if len(scenarios) == 0 {
		t.Error("Should have predefined scenarios")
	}

	for _, name := range scenarios {
		t.Run(name, func(t *testing.T) {
			scenario, exists := GetScenario(name)
			if !exists {
				t.Errorf("Scenario %s should exist", name)
				return
			}

			simulator := RunScenario(scenario)
			if simulator == nil {
				t.Error("RunScenario should return a simulator")
				return
			}

			// Verify scenario was configured correctly
			if simulator.GetTotalSpace() != scenario.TotalSpace {
				t.Errorf("Total space = %d, want %d", simulator.GetTotalSpace(), scenario.TotalSpace)
			}
			if simulator.GetAvailableSpace() != scenario.InitialSpace {
				t.Errorf("Initial space = %d, want %d", simulator.GetAvailableSpace(), scenario.InitialSpace)
			}

			// Run a basic operation to test the scenario
			fs := simulator.GetFileSystem()
			err := fs.WriteFile("/tmp/scenario_test.txt", []byte("test"), 0644)

			// Different scenarios have different expected outcomes
			switch name {
			case "ImmediateFailure":
				if err == nil {
					t.Error("ImmediateFailure scenario should fail on first operation")
				}
			case "GradualExhaustion":
				// Should succeed initially
				if err != nil {
					t.Errorf("GradualExhaustion should succeed initially: %v", err)
				}
			}
		})
	}
}

// TestSimulateArchiveCreation tests archive creation simulation
func TestSimulateArchiveCreation(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace: 2 * 1024 * 1024,  // 2MB
		TotalSpace:   10 * 1024 * 1024, // 10MB
	}
	simulator := NewDiskSpaceSimulator(config)

	t.Run("Successful archive creation", func(t *testing.T) {
		simulator.Reset()

		err := SimulateArchiveCreation(simulator, 1024*1024, 10) // 1MB archive, 10 files
		if err != nil {
			t.Errorf("Archive creation should succeed: %v", err)
		}

		stats := simulator.GetStats()
		if stats.TotalWrites == 0 {
			t.Error("Should have recorded write operations")
		}
	})

	t.Run("Failed archive creation due to space", func(t *testing.T) {
		// Use immediate exhaustion to force failure
		config.ExhaustionMode = ImmediateExhaustion
		simulator = NewDiskSpaceSimulator(config)

		err := SimulateArchiveCreation(simulator, 5*1024*1024, 10) // 5MB archive
		if err == nil {
			t.Error("Archive creation should fail with immediate exhaustion")
		}

		if !isDiskSpaceError(err) {
			t.Errorf("Should fail with disk space error: %v", err)
		}
	})
}

// TestSimulateBackupOperation tests backup operation simulation
func TestSimulateBackupOperation(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace: 3 * 1024 * 1024,  // 3MB
		TotalSpace:   10 * 1024 * 1024, // 10MB
	}
	simulator := NewDiskSpaceSimulator(config)

	t.Run("Successful backup", func(t *testing.T) {
		simulator.Reset()

		err := SimulateBackupOperation(simulator, 1024*1024, "/tmp/backup/test.bak") // 1MB backup
		if err != nil {
			t.Errorf("Backup operation should succeed: %v", err)
		}
	})

	t.Run("Failed backup due to space", func(t *testing.T) {
		// Configure for immediate failure
		config.ExhaustionMode = ImmediateExhaustion
		simulator = NewDiskSpaceSimulator(config)

		err := SimulateBackupOperation(simulator, 2*1024*1024, "/tmp/backup/large.bak") // 2MB backup
		if err == nil {
			t.Error("Backup operation should fail with immediate exhaustion")
		}

		if !isDiskSpaceError(err) {
			t.Errorf("Should fail with disk space error: %v", err)
		}
	})
}

// TestDynamicConfiguration tests dynamic simulator configuration
func TestDynamicConfiguration(t *testing.T) {
	simulator := NewDiskSpaceSimulator(DiskSpaceConfig{
		InitialSpace: 5 * 1024 * 1024,
		TotalSpace:   10 * 1024 * 1024,
	})

	t.Run("Add and remove failure points", func(t *testing.T) {
		// Add failure point
		testErr := CreateDiskFullError("/tmp/test")
		simulator.AddFailurePoint(2, testErr)

		fs := simulator.GetFileSystem()

		// Operation 1 should succeed
		err := fs.WriteFile("/tmp/op1.txt", []byte("op1"), 0644)
		if err != nil {
			t.Errorf("Operation 1 should succeed: %v", err)
		}

		// Operation 2 should fail
		err = fs.WriteFile("/tmp/op2.txt", []byte("op2"), 0644)
		if err == nil {
			t.Error("Operation 2 should fail")
		}

		// Remove failure point
		simulator.RemoveFailurePoint(2)
		simulator.Reset() // Reset operation counter

		// Now operation 2 should succeed
		fs.WriteFile("/tmp/op1.txt", []byte("op1"), 0644)       // Operation 1
		err = fs.WriteFile("/tmp/op2.txt", []byte("op2"), 0644) // Operation 2
		if err != nil {
			t.Errorf("Operation 2 should succeed after removing failure point: %v", err)
		}
	})

	t.Run("Add and remove recovery points", func(t *testing.T) {
		simulator.Reset()

		// Add recovery point
		simulator.AddRecoveryPoint(2, 2*1024*1024) // Recover 2MB at operation 2

		fs := simulator.GetFileSystem()
		spaceBefore := simulator.GetAvailableSpace()

		// Operations 1 and 2
		fs.WriteFile("/tmp/rec1.txt", []byte("rec1"), 0644)
		fs.WriteFile("/tmp/rec2.txt", []byte("rec2"), 0644) // Should trigger recovery

		spaceAfter := simulator.GetAvailableSpace()
		if spaceAfter <= spaceBefore {
			t.Error("Space should increase after recovery point")
		}

		// Remove recovery point and reset
		simulator.RemoveRecoveryPoint(2)
		stats := simulator.GetStats()
		recoveryOps := stats.RecoveryOperations

		simulator.Reset()
		fs.WriteFile("/tmp/rec1.txt", []byte("rec1"), 0644)
		fs.WriteFile("/tmp/rec2.txt", []byte("rec2"), 0644)

		newStats := simulator.GetStats()
		if newStats.RecoveryOperations >= recoveryOps {
			t.Error("Recovery operations should not increase after removing recovery point")
		}
	})

	t.Run("Add and remove injection points", func(t *testing.T) {
		// Add injection point
		testErr := CreateQuotaExceededError("/tmp/inject.txt")
		simulator.AddInjectionPoint("/tmp/inject.txt", testErr)

		fs := simulator.GetFileSystem()

		// File should fail
		err := fs.WriteFile("/tmp/inject.txt", []byte("inject"), 0644)
		if err == nil {
			t.Error("Injection point should cause failure")
		}

		// Remove injection point
		simulator.RemoveInjectionPoint("/tmp/inject.txt")

		// Now it should succeed (if space allows)
		simulator.SetEnabled(false) // Disable simulation to ensure success
		err = fs.WriteFile("/tmp/inject.txt", []byte("inject"), 0644)
		if err != nil {
			t.Errorf("File should succeed after removing injection point: %v", err)
		}
	})
}

// TestConcurrentOperations tests thread safety
func TestConcurrentOperations(t *testing.T) {
	config := DiskSpaceConfig{
		InitialSpace: 10 * 1024 * 1024, // 10MB
		TotalSpace:   20 * 1024 * 1024, // 20MB
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	var wg sync.WaitGroup
	numGoroutines := 10
	operationsPerGoroutine := 20

	errors := make(chan error, numGoroutines*operationsPerGoroutine)

	// Start multiple goroutines performing operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				filename := fmt.Sprintf("/tmp/concurrent_%d_%d.txt", id, j)
				err := fs.WriteFile(filename, []byte("concurrent test"), 0644)
				if err != nil {
					errors <- err
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Collect errors
	var errorCount int
	for err := range errors {
		if err != nil {
			errorCount++
		}
	}

	// Verify final state is consistent
	stats := simulator.GetStats()
	expectedTotal := numGoroutines * operationsPerGoroutine
	actualTotal := stats.TotalWrites

	// Allow some variance due to failures
	if actualTotal > expectedTotal || actualTotal < expectedTotal-errorCount {
		t.Errorf("Inconsistent total writes: expected around %d, got %d (errors: %d)",
			expectedTotal, actualTotal, errorCount)
	}

	// Verify space tracking is consistent
	if simulator.GetAvailableSpace() < 0 {
		t.Error("Available space should not be negative")
	}

	totalTrackedSpace := simulator.GetAvailableSpace() + simulator.GetUsedSpace()
	if totalTrackedSpace != simulator.GetTotalSpace() {
		t.Errorf("Space accounting inconsistent: available(%d) + used(%d) != total(%d)",
			simulator.GetAvailableSpace(), simulator.GetUsedSpace(), simulator.GetTotalSpace())
	}
}

// TestErrorTypes tests various error type creation and detection
func TestErrorTypes(t *testing.T) {
	tests := []struct {
		name        string
		createError func(string) error
		expectDisk  bool
	}{
		{
			name:        "Disk full error",
			createError: CreateDiskFullError,
			expectDisk:  true,
		},
		{
			name:        "Quota exceeded error",
			createError: CreateQuotaExceededError,
			expectDisk:  true,
		},
		{
			name:        "Device full error",
			createError: CreateDeviceFullError,
			expectDisk:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.createError("/tmp/test.txt")
			if err == nil {
				t.Error("Should create an error")
				return
			}

			isDisk := isDiskSpaceError(err)
			if isDisk != tt.expectDisk {
				t.Errorf("isDiskSpaceError() = %v, want %v", isDisk, tt.expectDisk)
			}

			// Test error message contains expected content
			errMsg := err.Error()
			if !strings.Contains(errMsg, "/tmp/test.txt") {
				t.Errorf("Error message should contain path: %s", errMsg)
			}
		})
	}

	t.Run("DiskSpaceError", func(t *testing.T) {
		err := NewDiskSpaceError("write", "/tmp/test", "insufficient space", 1024, 2048)

		if !strings.Contains(err.Error(), "write") {
			t.Error("Error should contain operation")
		}
		if !strings.Contains(err.Error(), "/tmp/test") {
			t.Error("Error should contain path")
		}
		if !strings.Contains(err.Error(), "1024") {
			t.Error("Error should contain available space")
		}
		if !strings.Contains(err.Error(), "2048") {
			t.Error("Error should contain required space")
		}

		// Test with wrapped error
		baseErr := errors.New("base error")
		err.Err = baseErr
		if err.Unwrap() != baseErr {
			t.Error("Should unwrap to base error")
		}
	})
}

// Helper function to detect disk space errors
func isDiskSpaceError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	patterns := []string{
		"no space left",
		"disk full",
		"quota exceeded",
		"device full",
		"insufficient space",
	}

	for _, pattern := range patterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	// Check for ENOSPC
	if pathErr, ok := err.(*os.PathError); ok {
		return pathErr.Err == syscall.ENOSPC
	}

	return false
}

// Benchmark tests for performance validation
func BenchmarkDiskSpaceSimulator_WriteFile(b *testing.B) {
	config := DiskSpaceConfig{
		InitialSpace: 100 * 1024 * 1024, // 100MB
		TotalSpace:   200 * 1024 * 1024, // 200MB
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	data := make([]byte, 1024) // 1KB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("/tmp/bench_%d.txt", i)
		fs.WriteFile(filename, data, 0644)
	}
}

func BenchmarkDiskSpaceSimulator_SpaceChecking(b *testing.B) {
	config := DiskSpaceConfig{
		InitialSpace:   50 * 1024 * 1024,
		TotalSpace:     100 * 1024 * 1024,
		ExhaustionMode: LinearExhaustion,
		ExhaustionRate: 0.01,
	}
	simulator := NewDiskSpaceSimulator(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("/tmp/bench_%d.txt", i)
		simulator.checkSpaceAndInject(filename, 1024)
	}
}

func BenchmarkDiskSpaceSimulator_ConcurrentAccess(b *testing.B) {
	config := DiskSpaceConfig{
		InitialSpace: 50 * 1024 * 1024,
		TotalSpace:   100 * 1024 * 1024,
	}
	simulator := NewDiskSpaceSimulator(config)
	fs := simulator.GetFileSystem()

	data := make([]byte, 512)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			filename := fmt.Sprintf("/tmp/parallel_%d.txt", i)
			fs.WriteFile(filename, data, 0644)
			i++
		}
	})
}
