// This file is part of bkpdir
//
// Package testutil provides testing infrastructure for complex scenarios,
// specifically disk space simulation utilities for controlled testing.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// TEST-INFRA-001-B: Disk space simulation framework
// DECISION-REF: DEC-004 (Error handling)
// IMPLEMENTATION-NOTES: Uses filesystem interface wrapper to simulate space constraints without requiring large files

package testutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"syscall"
)

// DiskSpaceSimulator provides controlled disk space simulation for testing
type DiskSpaceSimulator struct {
	mu                  sync.RWMutex
	availableSpace      int64
	totalSpace          int64
	usedSpace           int64
	writeOperations     int
	failurePoints       map[int]error
	spaceExhaustionRate float64 // 0.0-1.0, how much space to reduce per operation
	exhaustionMode      ExhaustionMode
	recoveryPoints      map[int]int64    // operation number -> space to add back
	injectionPoints     map[string]error // file path -> error to inject
	writeCounter        int
	isEnabled           bool
	originalFS          FileSystemInterface
	stats               SimulationStats
}

// ExhaustionMode determines how disk space is reduced during operations
type ExhaustionMode int

const (
	// LinearExhaustion reduces space by a fixed amount each operation
	LinearExhaustion ExhaustionMode = iota
	// ProgressiveExhaustion reduces space by increasing amounts each operation
	ProgressiveExhaustion
	// RandomExhaustion reduces space by random amounts
	RandomExhaustion
	// ImmediateExhaustion triggers disk full immediately
	ImmediateExhaustion
)

// SimulationStats tracks statistics about disk space simulation
type SimulationStats struct {
	TotalWrites         int
	FailedWrites        int
	SpaceExhausted      int
	RecoveryOperations  int
	InjectedErrors      int
	AverageSpaceReduced float64
	TotalSpaceReduced   int64
	LargestFileWritten  int64
	SmallestFileWritten int64
}

// DiskSpaceConfig configures disk space simulation parameters
type DiskSpaceConfig struct {
	InitialSpace      int64 // Initial available space in bytes
	TotalSpace        int64 // Total disk space in bytes
	ExhaustionMode    ExhaustionMode
	ExhaustionRate    float64          // 0.0-1.0, rate of space reduction
	FailurePoints     map[int]error    // operation number -> error to inject
	RecoveryPoints    map[int]int64    // operation number -> space to recover
	InjectionPoints   map[string]error // file path -> error to inject
	MinSpaceThreshold int64            // Minimum space before triggering ENOSPC
}

// FileSystemInterface abstracts file system operations for simulation
type FileSystemInterface interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Create(name string) (*os.File, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	MkdirAll(path string, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	Remove(name string) error
	RemoveAll(path string) error
}

// RealFileSystem implements FileSystemInterface for actual file operations
type RealFileSystem struct{}

func (rfs *RealFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (rfs *RealFileSystem) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (rfs *RealFileSystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func (rfs *RealFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (rfs *RealFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (rfs *RealFileSystem) Remove(name string) error {
	return os.Remove(name)
}

func (rfs *RealFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// SimulatedFileSystem wraps a real filesystem with disk space simulation
type SimulatedFileSystem struct {
	simulator  *DiskSpaceSimulator
	underlying FileSystemInterface
}

func (sfs *SimulatedFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if err := sfs.simulator.checkSpaceAndInject(filename, int64(len(data))); err != nil {
		return err
	}
	return sfs.underlying.WriteFile(filename, data, perm)
}

func (sfs *SimulatedFileSystem) Create(name string) (*os.File, error) {
	if err := sfs.simulator.checkSpaceAndInject(name, 0); err != nil {
		return nil, err
	}
	return sfs.underlying.Create(name)
}

func (sfs *SimulatedFileSystem) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	// Only check space for write operations
	if flag&(os.O_WRONLY|os.O_RDWR|os.O_CREATE) != 0 {
		if err := sfs.simulator.checkSpaceAndInject(name, 0); err != nil {
			return nil, err
		}
	}
	return sfs.underlying.OpenFile(name, flag, perm)
}

func (sfs *SimulatedFileSystem) MkdirAll(path string, perm os.FileMode) error {
	if err := sfs.simulator.checkSpaceAndInject(path, 1024); err != nil { // Assume 1KB for directory
		return err
	}
	return sfs.underlying.MkdirAll(path, perm)
}

func (sfs *SimulatedFileSystem) Stat(name string) (os.FileInfo, error) {
	return sfs.underlying.Stat(name)
}

func (sfs *SimulatedFileSystem) Remove(name string) error {
	// Try to get file size before removal for space recovery
	var fileSize int64
	if info, err := sfs.underlying.Stat(name); err == nil && !info.IsDir() {
		fileSize = info.Size()
	}

	err := sfs.underlying.Remove(name)
	if err == nil && fileSize > 0 {
		// Only recover space if the file actually existed and was removed successfully
		sfs.simulator.recoverSpace(fileSize)
	}
	return err
}

func (sfs *SimulatedFileSystem) RemoveAll(path string) error {
	// Calculate space to recover before removal
	var totalSize int64
	if info, err := sfs.underlying.Stat(path); err == nil {
		if info.IsDir() {
			// Walk the directory to calculate total size
			filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() {
					totalSize += info.Size()
				}
				return nil
			})
		} else {
			// It's a single file
			totalSize = info.Size()
		}
	}

	err := sfs.underlying.RemoveAll(path)
	if err == nil && totalSize > 0 {
		// Only recover space if removal was successful
		sfs.simulator.recoverSpace(totalSize)
	}
	return err
}

// NewDiskSpaceSimulator creates a new disk space simulator with the given configuration
func NewDiskSpaceSimulator(config DiskSpaceConfig) *DiskSpaceSimulator {
	if config.TotalSpace == 0 {
		config.TotalSpace = 100 * 1024 * 1024 // Default 100MB total space
	}

	// Only set default if both InitialSpace and TotalSpace were zero (completely unspecified)
	initialSpace := config.InitialSpace
	if config.InitialSpace == 0 && config.TotalSpace == 100*1024*1024 {
		// This means both were unspecified, use default 50%
		initialSpace = config.TotalSpace / 2
	}

	if config.MinSpaceThreshold == 0 {
		config.MinSpaceThreshold = 1024 // Default 1KB threshold
	}

	return &DiskSpaceSimulator{
		availableSpace:      initialSpace,
		totalSpace:          config.TotalSpace,
		usedSpace:           config.TotalSpace - initialSpace,
		failurePoints:       config.FailurePoints,
		spaceExhaustionRate: config.ExhaustionRate,
		exhaustionMode:      config.ExhaustionMode,
		recoveryPoints:      config.RecoveryPoints,
		injectionPoints:     config.InjectionPoints,
		isEnabled:           true,
		originalFS:          &RealFileSystem{},
		stats: SimulationStats{
			SmallestFileWritten: int64(^uint64(0) >> 1), // Max int64
		},
	}
}

// GetFileSystem returns a FileSystemInterface that simulates disk space constraints
func (ds *DiskSpaceSimulator) GetFileSystem() FileSystemInterface {
	return &SimulatedFileSystem{
		simulator:  ds,
		underlying: ds.originalFS,
	}
}

// checkSpaceAndInject checks available space and injects errors based on configuration
func (ds *DiskSpaceSimulator) checkSpaceAndInject(filename string, requiredSpace int64) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if !ds.isEnabled {
		return nil
	}

	ds.writeCounter++
	ds.stats.TotalWrites++

	// Check for specific error injection points
	if err, exists := ds.injectionPoints[filename]; exists {
		ds.stats.InjectedErrors++
		return err
	}

	// Check for operation-specific failure points
	if err, exists := ds.failurePoints[ds.writeCounter]; exists {
		ds.stats.FailedWrites++
		return err
	}

	// Check for space recovery points
	if recoverySpace, exists := ds.recoveryPoints[ds.writeCounter]; exists {
		ds.recoverSpaceUnlocked(recoverySpace)
		ds.stats.RecoveryOperations++
	}

	// Simulate space exhaustion based on mode
	spaceToReduce := ds.calculateSpaceReduction(requiredSpace)
	ds.reduceAvailableSpace(spaceToReduce)

	// Check if we have enough space for the operation
	if ds.availableSpace < requiredSpace || ds.availableSpace < 0 {
		ds.stats.FailedWrites++
		ds.stats.SpaceExhausted++
		return &os.PathError{
			Op:   "write",
			Path: filename,
			Err:  syscall.ENOSPC,
		}
	}

	// Track file size statistics
	if requiredSpace > 0 {
		if requiredSpace > ds.stats.LargestFileWritten {
			ds.stats.LargestFileWritten = requiredSpace
		}
		if requiredSpace < ds.stats.SmallestFileWritten {
			ds.stats.SmallestFileWritten = requiredSpace
		}
	}

	// Consume the space
	ds.usedSpace += requiredSpace
	ds.availableSpace -= requiredSpace

	return nil
}

// calculateSpaceReduction determines how much space to reduce based on exhaustion mode
func (ds *DiskSpaceSimulator) calculateSpaceReduction(operationSize int64) int64 {
	switch ds.exhaustionMode {
	case LinearExhaustion:
		reduction := int64(float64(ds.totalSpace) * ds.spaceExhaustionRate)
		ds.stats.TotalSpaceReduced += reduction
		return reduction

	case ProgressiveExhaustion:
		// Increase reduction amount with each operation
		reduction := int64(float64(ds.totalSpace) * ds.spaceExhaustionRate * float64(ds.writeCounter) / 10.0)
		ds.stats.TotalSpaceReduced += reduction
		return reduction

	case RandomExhaustion:
		// Random reduction between 0 and max rate
		maxReduction := int64(float64(ds.totalSpace) * ds.spaceExhaustionRate)
		reduction := int64(float64(maxReduction) * (float64(ds.writeCounter%7) / 6.0)) // Pseudo-random
		ds.stats.TotalSpaceReduced += reduction
		return reduction

	case ImmediateExhaustion:
		// Exhaust all available space immediately
		reduction := ds.availableSpace
		ds.stats.TotalSpaceReduced += reduction
		return reduction

	default:
		return 0
	}
}

// reduceAvailableSpace reduces available space by the specified amount
func (ds *DiskSpaceSimulator) reduceAvailableSpace(amount int64) {
	ds.availableSpace -= amount
	if ds.availableSpace < 0 {
		ds.availableSpace = 0
	}

	// Update average space reduced
	if ds.stats.TotalWrites > 0 {
		ds.stats.AverageSpaceReduced = float64(ds.stats.TotalSpaceReduced) / float64(ds.stats.TotalWrites)
	}
}

// recoverSpace adds space back to the available pool (unlocked version)
func (ds *DiskSpaceSimulator) recoverSpaceUnlocked(amount int64) {
	ds.availableSpace += amount
	ds.usedSpace -= amount

	// Ensure we don't exceed total space or go below zero
	if ds.availableSpace > ds.totalSpace {
		ds.availableSpace = ds.totalSpace
		ds.usedSpace = 0
	}
	if ds.usedSpace < 0 {
		ds.usedSpace = 0
		ds.availableSpace = ds.totalSpace
	}
}

// recoverSpace adds space back to the available pool
func (ds *DiskSpaceSimulator) recoverSpace(amount int64) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.recoverSpaceUnlocked(amount)
}

// GetAvailableSpace returns the current available space
func (ds *DiskSpaceSimulator) GetAvailableSpace() int64 {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.availableSpace
}

// GetUsedSpace returns the current used space
func (ds *DiskSpaceSimulator) GetUsedSpace() int64 {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.usedSpace
}

// GetTotalSpace returns the total space
func (ds *DiskSpaceSimulator) GetTotalSpace() int64 {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.totalSpace
}

// GetStats returns current simulation statistics
func (ds *DiskSpaceSimulator) GetStats() SimulationStats {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	stats := ds.stats
	if stats.SmallestFileWritten == int64(^uint64(0)>>1) {
		stats.SmallestFileWritten = 0
	}
	return stats
}

// Reset resets the simulator to initial state
func (ds *DiskSpaceSimulator) Reset() {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.availableSpace = ds.totalSpace / 2 // Reset to 50% available
	ds.usedSpace = ds.totalSpace / 2
	ds.writeCounter = 0
	ds.stats = SimulationStats{
		SmallestFileWritten: int64(^uint64(0) >> 1),
	}
}

// SetEnabled enables or disables the simulation
func (ds *DiskSpaceSimulator) SetEnabled(enabled bool) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.isEnabled = enabled
}

// IsEnabled returns whether simulation is currently enabled
func (ds *DiskSpaceSimulator) IsEnabled() bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.isEnabled
}

// AddFailurePoint adds a failure point at the specified operation number
func (ds *DiskSpaceSimulator) AddFailurePoint(operation int, err error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if ds.failurePoints == nil {
		ds.failurePoints = make(map[int]error)
	}
	ds.failurePoints[operation] = err
}

// AddRecoveryPoint adds a space recovery point at the specified operation number
func (ds *DiskSpaceSimulator) AddRecoveryPoint(operation int, spaceToRecover int64) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if ds.recoveryPoints == nil {
		ds.recoveryPoints = make(map[int]int64)
	}
	ds.recoveryPoints[operation] = spaceToRecover
}

// AddInjectionPoint adds an error injection point for a specific file path
func (ds *DiskSpaceSimulator) AddInjectionPoint(filePath string, err error) {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if ds.injectionPoints == nil {
		ds.injectionPoints = make(map[string]error)
	}
	ds.injectionPoints[filePath] = err
}

// RemoveFailurePoint removes a failure point
func (ds *DiskSpaceSimulator) RemoveFailurePoint(operation int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	delete(ds.failurePoints, operation)
}

// RemoveRecoveryPoint removes a recovery point
func (ds *DiskSpaceSimulator) RemoveRecoveryPoint(operation int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	delete(ds.recoveryPoints, operation)
}

// RemoveInjectionPoint removes an injection point
func (ds *DiskSpaceSimulator) RemoveInjectionPoint(filePath string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	delete(ds.injectionPoints, filePath)
}

// GetOperationCount returns the current operation count
func (ds *DiskSpaceSimulator) GetOperationCount() int {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.writeCounter
}

// DiskSpaceError represents errors related to disk space simulation
type DiskSpaceError struct {
	Message   string
	Operation string
	Path      string
	Available int64
	Required  int64
	Err       error
}

func (e *DiskSpaceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("disk space error during %s on %s: %s (available: %d, required: %d): %v",
			e.Operation, e.Path, e.Message, e.Available, e.Required, e.Err)
	}
	return fmt.Sprintf("disk space error during %s on %s: %s (available: %d, required: %d)",
		e.Operation, e.Path, e.Message, e.Available, e.Required)
}

func (e *DiskSpaceError) Unwrap() error {
	return e.Err
}

// NewDiskSpaceError creates a new DiskSpaceError
func NewDiskSpaceError(operation, path, message string, available, required int64) *DiskSpaceError {
	return &DiskSpaceError{
		Message:   message,
		Operation: operation,
		Path:      path,
		Available: available,
		Required:  required,
	}
}

// CreateDiskFullError creates a standard "no space left on device" error
func CreateDiskFullError(path string) error {
	return &os.PathError{
		Op:   "write",
		Path: path,
		Err:  syscall.ENOSPC,
	}
}

// CreateQuotaExceededError creates a quota exceeded error
func CreateQuotaExceededError(path string) error {
	return &os.PathError{
		Op:   "write",
		Path: path,
		Err:  errors.New("disk quota exceeded"),
	}
}

// CreateDeviceFullError creates a device full error
func CreateDeviceFullError(path string) error {
	return &os.PathError{
		Op:   "write",
		Path: path,
		Err:  errors.New("device full"),
	}
}

// DiskSpaceScenario represents a complete test scenario with disk space constraints
type DiskSpaceScenario struct {
	Name            string
	Description     string
	InitialSpace    int64
	TotalSpace      int64
	ExhaustionMode  ExhaustionMode
	ExhaustionRate  float64
	FailurePoints   map[int]error
	RecoveryPoints  map[int]int64
	InjectionPoints map[string]error
	ExpectedResults ScenarioResults
}

// ScenarioResults defines expected results for a disk space scenario
type ScenarioResults struct {
	ShouldFailAtOperation  int
	ExpectedAvailableSpace int64
	ExpectedFailureCount   int
	ExpectedRecoveryCount  int
	ShouldTriggerCleanup   bool
}

// PredefinedScenarios provides common disk space testing scenarios
var PredefinedScenarios = map[string]DiskSpaceScenario{
	"GradualExhaustion": {
		Name:           "Gradual Space Exhaustion",
		Description:    "Gradually reduce available space over multiple operations",
		InitialSpace:   10 * 1024 * 1024, // 10MB
		TotalSpace:     20 * 1024 * 1024, // 20MB
		ExhaustionMode: LinearExhaustion,
		ExhaustionRate: 0.1, // 10% per operation
		ExpectedResults: ScenarioResults{
			ShouldFailAtOperation: 10,
			ExpectedFailureCount:  1,
		},
	},
	"ImmediateFailure": {
		Name:           "Immediate Disk Full",
		Description:    "Trigger disk full error immediately",
		InitialSpace:   1024,              // 1KB
		TotalSpace:     100 * 1024 * 1024, // 100MB
		ExhaustionMode: ImmediateExhaustion,
		ExpectedResults: ScenarioResults{
			ShouldFailAtOperation: 1,
			ExpectedFailureCount:  1,
		},
	},
	"SpaceRecovery": {
		Name:           "Space Recovery Test",
		Description:    "Test behavior when space becomes available again",
		InitialSpace:   1024,             // 1KB
		TotalSpace:     10 * 1024 * 1024, // 10MB
		ExhaustionMode: LinearExhaustion,
		ExhaustionRate: 0.5,
		RecoveryPoints: map[int]int64{
			3: 5 * 1024 * 1024, // Recover 5MB at operation 3
		},
		ExpectedResults: ScenarioResults{
			ShouldFailAtOperation: 1,
			ExpectedFailureCount:  1,
			ExpectedRecoveryCount: 1,
		},
	},
	"ProgressiveExhaustion": {
		Name:           "Progressive Space Exhaustion",
		Description:    "Progressively increase space reduction with each operation",
		InitialSpace:   50 * 1024 * 1024,  // 50MB
		TotalSpace:     100 * 1024 * 1024, // 100MB
		ExhaustionMode: ProgressiveExhaustion,
		ExhaustionRate: 0.05, // 5% base rate
		ExpectedResults: ScenarioResults{
			ShouldFailAtOperation: 15,
			ExpectedFailureCount:  1,
		},
	},
}

// GetScenario returns a predefined scenario by name
func GetScenario(name string) (DiskSpaceScenario, bool) {
	scenario, exists := PredefinedScenarios[name]
	return scenario, exists
}

// GetAllScenarioNames returns names of all predefined scenarios
func GetAllScenarioNames() []string {
	names := make([]string, 0, len(PredefinedScenarios))
	for name := range PredefinedScenarios {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// RunScenario creates and configures a simulator for the given scenario
func RunScenario(scenario DiskSpaceScenario) *DiskSpaceSimulator {
	config := DiskSpaceConfig{
		InitialSpace:    scenario.InitialSpace,
		TotalSpace:      scenario.TotalSpace,
		ExhaustionMode:  scenario.ExhaustionMode,
		ExhaustionRate:  scenario.ExhaustionRate,
		FailurePoints:   scenario.FailurePoints,
		RecoveryPoints:  scenario.RecoveryPoints,
		InjectionPoints: scenario.InjectionPoints,
	}

	return NewDiskSpaceSimulator(config)
}

// SimulateArchiveCreation simulates archive creation under disk space constraints
func SimulateArchiveCreation(simulator *DiskSpaceSimulator, archiveSize int64, numFiles int) error {
	fs := simulator.GetFileSystem()

	// Create archive directory
	if err := fs.MkdirAll("/tmp/test_archive", 0755); err != nil {
		return fmt.Errorf("failed to create archive directory: %w", err)
	}

	// Simulate writing archive files
	avgFileSize := archiveSize / int64(numFiles)
	for i := 0; i < numFiles; i++ {
		filename := fmt.Sprintf("/tmp/test_archive/file_%d.txt", i)
		data := make([]byte, avgFileSize)

		if err := fs.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filename, err)
		}
	}

	return nil
}

// SimulateBackupOperation simulates file backup under disk space constraints
func SimulateBackupOperation(simulator *DiskSpaceSimulator, sourceSize int64, backupPath string) error {
	fs := simulator.GetFileSystem()

	// Create backup directory
	backupDir := filepath.Dir(backupPath)
	if err := fs.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Simulate writing backup file
	data := make([]byte, sourceSize)
	if err := fs.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}
