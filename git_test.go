// This file is part of bkpdir
// LINT-001: Lint compliance

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGitIntegration tests the Git integration functionality for GIT-001 feature
func TestGitIntegration(t *testing.T) {
	// GIT-001: Git integration validation
	// TEST-REF: Feature tracking matrix GIT-001
	// IMMUTABLE-REF: Git Integration Requirements

	// Create temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "git_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("NonGitRepository", func(t *testing.T) {
		// Test with a non-git directory
		if IsGitRepository(tmpDir) {
			t.Error("Expected non-git directory to return false")
		}

		branch, hash := GetGitInfo(tmpDir)
		if branch != "" || hash != "" {
			t.Errorf("Expected empty git info for non-git directory, got branch=%s, hash=%s", branch, hash)
		}
	})

	t.Run("GitRepository", func(t *testing.T) {
		// Skip if git is not available
		if !isGitAvailable() {
			t.Skip("Git not available, skipping git repository tests")
		}

		// Initialize git repository
		gitDir := filepath.Join(tmpDir, "git_repo")
		if err := os.MkdirAll(gitDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Initialize git repo
		runGitCommand(t, gitDir, "init")
		runGitCommand(t, gitDir, "config", "user.email", "test@example.com")
		runGitCommand(t, gitDir, "config", "user.name", "Test User")

		// Create a test file and make initial commit
		testFile := filepath.Join(gitDir, "test.txt")
		if err := ioutil.WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}

		runGitCommand(t, gitDir, "add", "test.txt")
		runGitCommand(t, gitDir, "commit", "-m", "Initial commit")

		// Test git repository detection
		if !IsGitRepository(gitDir) {
			t.Error("Expected git directory to return true")
		}

		// Test git info extraction
		branch, hash := GetGitInfo(gitDir)
		if branch == "" {
			t.Error("Expected non-empty branch name in git repository")
		}
		if hash == "" {
			t.Error("Expected non-empty commit hash in git repository")
		}

		// Verify branch is typically "main" or "master"
		if branch != "main" && branch != "master" {
			t.Logf("Branch name is %s (expected main or master, but this may vary)", branch)
		}

		// Verify hash is reasonable length (git short hash is typically 7 characters)
		if len(hash) < 4 || len(hash) > 40 {
			t.Errorf("Expected commit hash length between 4-40 characters, got %d: %s", len(hash), hash)
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Test error handling with invalid directory
		invalidDir := filepath.Join(tmpDir, "nonexistent")

		if IsGitRepository(invalidDir) {
			t.Error("Expected nonexistent directory to return false")
		}

		branch, hash := GetGitInfo(invalidDir)
		if branch != "" || hash != "" {
			t.Errorf("Expected empty git info for invalid directory, got branch=%s, hash=%s", branch, hash)
		}
	})
}

// TestGitNaming tests the Git naming functionality for GIT-002 feature
func TestGitNaming(t *testing.T) {
	// GIT-002: Git naming validation
	// TEST-REF: Feature tracking matrix GIT-002
	// IMMUTABLE-REF: Git Integration Requirements

	// Create temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "git_naming_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("BranchExtraction", func(t *testing.T) {
		// Skip if git is not available
		if !isGitAvailable() {
			t.Skip("Git not available, skipping branch extraction tests")
		}

		gitDir := filepath.Join(tmpDir, "branch_test")
		if err := os.MkdirAll(gitDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Initialize git repo
		runGitCommand(t, gitDir, "init")
		runGitCommand(t, gitDir, "config", "user.email", "test@example.com")
		runGitCommand(t, gitDir, "config", "user.name", "Test User")

		// Create initial commit
		testFile := filepath.Join(gitDir, "test.txt")
		if err := ioutil.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
		runGitCommand(t, gitDir, "add", "test.txt")
		runGitCommand(t, gitDir, "commit", "-m", "Initial commit")

		// Test branch extraction
		branch := GetGitBranch(gitDir)
		if branch == "" {
			t.Error("Expected non-empty branch name")
		}

		// Create and switch to a new branch
		runGitCommand(t, gitDir, "checkout", "-b", "feature-branch")

		// Test new branch name
		newBranch := GetGitBranch(gitDir)
		if newBranch != "feature-branch" {
			t.Errorf("Expected branch name 'feature-branch', got %s", newBranch)
		}
	})

	t.Run("HashExtraction", func(t *testing.T) {
		// Skip if git is not available
		if !isGitAvailable() {
			t.Skip("Git not available, skipping hash extraction tests")
		}

		gitDir := filepath.Join(tmpDir, "hash_test")
		if err := os.MkdirAll(gitDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Initialize git repo
		runGitCommand(t, gitDir, "init")
		runGitCommand(t, gitDir, "config", "user.email", "test@example.com")
		runGitCommand(t, gitDir, "config", "user.name", "Test User")

		// Create initial commit
		testFile := filepath.Join(gitDir, "test.txt")
		if err := ioutil.WriteFile(testFile, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
		runGitCommand(t, gitDir, "add", "test.txt")
		runGitCommand(t, gitDir, "commit", "-m", "Initial commit")

		// Test hash extraction
		hash1 := GetGitShortHash(gitDir)
		if hash1 == "" {
			t.Error("Expected non-empty commit hash")
		}

		// Make another commit
		if err := ioutil.WriteFile(testFile, []byte("updated test"), 0644); err != nil {
			t.Fatal(err)
		}
		runGitCommand(t, gitDir, "add", "test.txt")
		runGitCommand(t, gitDir, "commit", "-m", "Second commit")

		// Test that hash changed
		hash2 := GetGitShortHash(gitDir)
		if hash2 == "" {
			t.Error("Expected non-empty commit hash after second commit")
		}
		if hash1 == hash2 {
			t.Error("Expected commit hash to change after new commit")
		}
	})

	t.Run("NonGitDirectory", func(t *testing.T) {
		// Test with non-git directory
		nonGitDir := filepath.Join(tmpDir, "non_git")
		if err := os.MkdirAll(nonGitDir, 0755); err != nil {
			t.Fatal(err)
		}

		branch := GetGitBranch(nonGitDir)
		if branch != "" {
			t.Errorf("Expected empty branch for non-git directory, got %s", branch)
		}

		hash := GetGitShortHash(nonGitDir)
		if hash != "" {
			t.Errorf("Expected empty hash for non-git directory, got %s", hash)
		}
	})
}

// TestGitStatus tests the Git status detection functionality for GIT-003 feature
func TestGitStatus(t *testing.T) {
	// GIT-003: Git status validation
	// TEST-REF: Feature tracking matrix GIT-003
	// IMMUTABLE-REF: Git Integration Requirements

	// Create temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "git_status_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("NonGitRepository", func(t *testing.T) {
		// Test with a non-git directory
		if IsGitWorkingDirectoryClean(tmpDir) {
			t.Error("Expected non-git directory to return false")
		}

		branch, hash, isClean := GetGitInfoWithStatus(tmpDir)
		if branch != "" || hash != "" || isClean {
			t.Errorf("Expected empty git info for non-git directory, got branch=%s, hash=%s, isClean=%v", branch, hash, isClean)
		}
	})

	t.Run("CleanGitRepository", func(t *testing.T) {
		// Skip if git is not available
		if !isGitAvailable() {
			t.Skip("Git not available, skipping git status tests")
		}

		// Initialize git repository
		gitDir := filepath.Join(tmpDir, "clean_repo")
		if err := os.MkdirAll(gitDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Initialize git repo
		runGitCommand(t, gitDir, "init")
		runGitCommand(t, gitDir, "config", "user.email", "test@example.com")
		runGitCommand(t, gitDir, "config", "user.name", "Test User")

		// Create a test file and make initial commit
		testFile := filepath.Join(gitDir, "test.txt")
		if err := ioutil.WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}

		runGitCommand(t, gitDir, "add", "test.txt")
		runGitCommand(t, gitDir, "commit", "-m", "Initial commit")

		// Test clean working directory
		if !IsGitWorkingDirectoryClean(gitDir) {
			t.Error("Expected clean working directory to return true")
		}

		// Test git info with status
		branch, hash, isClean := GetGitInfoWithStatus(gitDir)
		if branch == "" {
			t.Error("Expected non-empty branch name")
		}
		if hash == "" {
			t.Error("Expected non-empty commit hash")
		}
		if !isClean {
			t.Error("Expected clean working directory")
		}
	})

	t.Run("DirtyGitRepository", func(t *testing.T) {
		// Skip if git is not available
		if !isGitAvailable() {
			t.Skip("Git not available, skipping git status tests")
		}

		// Initialize git repository
		gitDir := filepath.Join(tmpDir, "dirty_repo")
		if err := os.MkdirAll(gitDir, 0755); err != nil {
			t.Fatal(err)
		}

		// Initialize git repo
		runGitCommand(t, gitDir, "init")
		runGitCommand(t, gitDir, "config", "user.email", "test@example.com")
		runGitCommand(t, gitDir, "config", "user.name", "Test User")

		// Create a test file and make initial commit
		testFile := filepath.Join(gitDir, "test.txt")
		if err := ioutil.WriteFile(testFile, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}

		runGitCommand(t, gitDir, "add", "test.txt")
		runGitCommand(t, gitDir, "commit", "-m", "Initial commit")

		// Modify the file to create dirty state
		if err := ioutil.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
			t.Fatal(err)
		}

		// Test dirty working directory
		if IsGitWorkingDirectoryClean(gitDir) {
			t.Error("Expected dirty working directory to return false")
		}

		// Test git info with status
		branch, hash, isClean := GetGitInfoWithStatus(gitDir)
		if branch == "" {
			t.Error("Expected non-empty branch name")
		}
		if hash == "" {
			t.Error("Expected non-empty commit hash")
		}
		if isClean {
			t.Error("Expected dirty working directory")
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Test error handling with invalid directory
		invalidDir := filepath.Join(tmpDir, "nonexistent")

		if IsGitWorkingDirectoryClean(invalidDir) {
			t.Error("Expected nonexistent directory to return false")
		}

		branch, hash, isClean := GetGitInfoWithStatus(invalidDir)
		if branch != "" || hash != "" || isClean {
			t.Errorf("Expected empty git info for invalid directory, got branch=%s, hash=%s, isClean=%v", branch, hash, isClean)
		}
	})
}

// Helper functions

// isGitAvailable checks if git command is available in the system
func isGitAvailable() bool {
	cmd := exec.Command("git", "--version")
	err := cmd.Run()
	return err == nil
}

// runGitCommand runs a git command and fails the test if it returns an error
func runGitCommand(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Git command failed: %s\nOutput: %s\nError: %v", strings.Join(args, " "), string(out), err)
	}
}
