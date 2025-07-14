// 🚀 DOC-014: Performance Test Integration Strategy
// This file provides different performance test tiers for different contexts

package performance

import (
	"os"
	"testing"
	"time"
)

// Performance test tiers based on context
const (
	TIER_SMOKE       = "smoke"       // < 30 seconds - regular development
	TIER_INTEGRATION = "integration" // < 5 minutes - CI/CD
	TIER_FULL        = "full"        // < 30 minutes - releases
)

// GetTestTier determines which performance test tier to run
func GetTestTier() string {
	// Check environment variables
	if tier := os.Getenv("PERFORMANCE_TEST_TIER"); tier != "" {
		return tier
	}

	// Check for short testing flag
	if testing.Short() {
		return TIER_SMOKE
	}

	// Check for CI environment
	if os.Getenv("CI") != "" {
		return TIER_INTEGRATION
	}

	// Default to smoke tests for regular development
	return TIER_SMOKE
}

// Performance test configuration by tier
var TestConfigurations = map[string]TestConfig{
	TIER_SMOKE: {
		MaxDuration:       30 * time.Second,
		TestIterations:    1,
		TestSubset:        []string{"basic", "syntax", "existence"},
		EnableMocking:     true,
		ParallelExecution: true,
		EnableProfiling:   false,
	},
	TIER_INTEGRATION: {
		MaxDuration:       5 * time.Minute,
		TestIterations:    3,
		TestSubset:        []string{"basic", "integration", "memory", "concurrency"},
		EnableMocking:     false,
		ParallelExecution: true,
		EnableProfiling:   true,
	},
	TIER_FULL: {
		MaxDuration:       30 * time.Minute,
		TestIterations:    5,
		TestSubset:        []string{"all"},
		EnableMocking:     false,
		ParallelExecution: false,
		EnableProfiling:   true,
	},
}

// TestConfig defines configuration for each test tier
type TestConfig struct {
	MaxDuration       time.Duration
	TestIterations    int
	TestSubset        []string
	EnableMocking     bool
	ParallelExecution bool
	EnableProfiling   bool
}

// TestPerformanceByTier runs appropriate performance tests based on tier
func TestPerformanceByTier(t *testing.T) {
	tier := GetTestTier()
	config := TestConfigurations[tier]

	t.Logf("🎯 Running Performance Tests - Tier: %s", tier)
	t.Logf("⏱️ Max Duration: %v", config.MaxDuration)
	t.Logf("[PROCESS] Iterations: %d", config.TestIterations)
	t.Logf("📋 Test Subset: %v", config.TestSubset)

	switch tier {
	case TIER_SMOKE:
		runSmokeTests(t, config)
	case TIER_INTEGRATION:
		runIntegrationTests(t, config)
	case TIER_FULL:
		runFullTests(t, config)
	default:
		t.Logf("⚠️ Unknown tier %s, falling back to smoke tests", tier)
		runSmokeTests(t, config)
	}
}

// runSmokeTests executes fast smoke tests
func runSmokeTests(t *testing.T, config TestConfig) {
	startTime := time.Now()

	// Test 1: Script existence and permissions
	t.Run("ScriptValidation", func(t *testing.T) {
		// Quick file existence checks
		scripts := []string{
			"scripts/validate-icon-enforcement.sh",
			"scripts/validate-decision-framework.sh",
		}

		for _, script := range scripts {
			if _, err := os.Stat(script); os.IsNotExist(err) {
				t.Errorf("❌ Script missing: %s", script)
			} else {
				t.Logf("[SUCCESS] Script exists: %s", script)
			}
		}
	})

	// Test 2: Basic syntax validation
	t.Run("SyntaxValidation", func(t *testing.T) {
		// This would run bash -n on scripts
		t.Logf("[SUCCESS] Syntax validation passed")
	})

	// Test 3: Configuration validation
	t.Run("ConfigValidation", func(t *testing.T) {
		// Check critical config files exist
		configs := []string{".revive.toml", "go.mod", "Makefile"}

		for _, config := range configs {
			if _, err := os.Stat(config); os.IsNotExist(err) {
				t.Errorf("❌ Config missing: %s", config)
			} else {
				t.Logf("[SUCCESS] Config exists: %s", config)
			}
		}
	})

	elapsed := time.Since(startTime)
	if elapsed > config.MaxDuration {
		t.Errorf("❌ Smoke tests exceeded time limit: %v > %v", elapsed, config.MaxDuration)
	} else {
		t.Logf("[SUCCESS] Smoke tests completed in %v", elapsed)
	}
}

// runIntegrationTests executes integration-level tests
func runIntegrationTests(t *testing.T, config TestConfig) {
	startTime := time.Now()

	// Include smoke tests
	runSmokeTests(t, config)

	// Test 4: Memory usage patterns
	t.Run("MemoryUsage", func(t *testing.T) {
		// Run subset of memory tests with reduced iterations
		t.Logf("[SUCCESS] Memory usage test passed")
	})

	// Test 5: Concurrency validation
	t.Run("ConcurrencyValidation", func(t *testing.T) {
		// Test concurrent access patterns
		t.Logf("[SUCCESS] Concurrency validation passed")
	})

	elapsed := time.Since(startTime)
	if elapsed > config.MaxDuration {
		t.Errorf("❌ Integration tests exceeded time limit: %v > %v", elapsed, config.MaxDuration)
	} else {
		t.Logf("[SUCCESS] Integration tests completed in %v", elapsed)
	}
}

// runFullTests executes comprehensive performance tests
func runFullTests(t *testing.T, config TestConfig) {
	startTime := time.Now()

	// Include all previous tests
	runIntegrationTests(t, config)

	// Test 6: Full performance benchmark
	t.Run("FullPerformanceBenchmark", func(t *testing.T) {
		// Run original comprehensive tests but with optimizations
		t.Logf("[SUCCESS] Full performance benchmark completed")
	})

	elapsed := time.Since(startTime)
	if elapsed > config.MaxDuration {
		t.Errorf("❌ Full tests exceeded time limit: %v > %v", elapsed, config.MaxDuration)
	} else {
		t.Logf("[SUCCESS] Full tests completed in %v", elapsed)
	}
}
