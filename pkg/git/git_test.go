// ⭐ EXTRACT-004: Git Integration System - Comprehensive test suite for extracted Git package - 🧪
package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGitIntegration tests the Git integration functionality
func TestGitIntegration(t *testing.T) {
	// ⭐ EXTRACT-004: Git integration validation - 🧪

	// Create temporary directory for testing
	tmpDir, err := ioutil.TempDir("", "git_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("NonGitRepository", func(t *testing.T) {
		repo := NewRepositoryWithConfig(&Config{
			WorkingDirectory: tmpDir,
			GitCommand:       "git",
		})

		// Test repository detection
		if repo.IsRepository() {
			t.Error("Expected non-git directory to return false")
		}

		// Test info extraction
		info, err := repo.GetInfo()
		if err != nil {
			t.Errorf("GetInfo should not return error for non-git directory: %v", err)
		}
		if info.IsRepo {
			t.Error("Expected IsRepo to be false for non-git directory")
		}
		if info.Branch != "" || info.Hash != "" {
			t.Errorf("Expected empty git info for non-git directory, got branch=%s, hash=%s", info.Branch, info.Hash)
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

		repo := NewRepositoryWithConfig(&Config{
			WorkingDirectory: gitDir,
			GitCommand:       "git",
		})

		// Test git repository detection
		if !repo.IsRepository() {
			t.Error("Expected git directory to return true")
		}

		// Test git info extraction
		info, err := repo.GetInfo()
		if err != nil {
			t.Errorf("GetInfo failed: %v", err)
		}
		if !info.IsRepo {
			t.Error("Expected IsRepo to be true for git repository")
		}
		if info.Branch == "" {
			t.Error("Expected non-empty branch name in git repository")
		}
		if info.Hash == "" {
			t.Error("Expected non-empty commit hash in git repository")
		}

		// Verify hash is reasonable length
		if len(info.Hash) < 4 || len(info.Hash) > 40 {
			t.Errorf("Expected commit hash length between 4-40 characters, got %d: %s", len(info.Hash), info.Hash)
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Test error handling with invalid directory
		invalidDir := filepath.Join(tmpDir, "nonexistent")
		repo := NewRepositoryWithConfig(&Config{
			WorkingDirectory: invalidDir,
			GitCommand:       "git",
		})

		if repo.IsRepository() {
			t.Error("Expected nonexistent directory to return false")
		}

		info, err := repo.GetInfo()
		if err != nil {
			t.Errorf("GetInfo should not return error for invalid directory: %v", err)
		}
		if info.IsRepo {
			t.Error("Expected IsRepo to be false for invalid directory")
		}
	})
}

// TestGitConfiguration tests the Git configuration functionality
func TestGitConfiguration(t *testing.T) {
	// ⭐ EXTRACT-004: Git configuration support testing - 🧪

	t.Run("DefaultConfig", func(t *testing.T) {
		config := DefaultConfig()
		if config.WorkingDirectory != "." {
			t.Errorf("Expected default working directory '.', got %s", config.WorkingDirectory)
		}
		if !config.IncludeDirtyStatus {
			t.Error("Expected default IncludeDirtyStatus to be true")
		}
		if config.GitCommand != "git" {
			t.Errorf("Expected default GitCommand 'git', got %s", config.GitCommand)
		}
	})

	t.Run("CustomConfig", func(t *testing.T) {
		tmpDir, err := ioutil.TempDir("", "git_config_test_")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		config := &Config{
			WorkingDirectory:   tmpDir,
			IncludeDirtyStatus: false,
			GitCommand:         "git",
		}

		repo := NewRepositoryWithConfig(config)
		if !repo.IsRepository() != true { // Should be false for non-git directory
			// This is expected behavior
		}

		// Test that config is properly used
		repoImpl := repo.(*Repo)
		if repoImpl.config.WorkingDirectory != tmpDir {
			t.Errorf("Expected working directory %s, got %s", tmpDir, repoImpl.config.WorkingDirectory)
		}
		if repoImpl.config.IncludeDirtyStatus != false {
			t.Error("Expected IncludeDirtyStatus to be false")
		}
	})
}

// TestGitBranchExtraction tests branch extraction functionality
func TestGitBranchExtraction(t *testing.T) {
	// ⭐ EXTRACT-004: Git branch and hash utilities testing - 🧪

	// Skip if git is not available
	if !isGitAvailable() {
		t.Skip("Git not available, skipping branch extraction tests")
	}

	tmpDir, err := ioutil.TempDir("", "git_branch_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

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

	repo := NewRepositoryWithConfig(&Config{
		WorkingDirectory: gitDir,
		GitCommand:       "git",
	})

	t.Run("BranchExtraction", func(t *testing.T) {
		branch, err := repo.GetBranch()
		if err != nil {
			t.Errorf("GetBranch failed: %v", err)
		}
		if branch == "" {
			t.Error("Expected non-empty branch name")
		}

		// Create and switch to a new branch
		runGitCommand(t, gitDir, "checkout", "-b", "feature-branch")

		// Test new branch name
		newBranch, err := repo.GetBranch()
		if err != nil {
			t.Errorf("GetBranch failed after branch switch: %v", err)
		}
		if newBranch != "feature-branch" {
			t.Errorf("Expected branch name 'feature-branch', got %s", newBranch)
		}
	})

	t.Run("HashExtraction", func(t *testing.T) {
		hash1, err := repo.GetShortHash()
		if err != nil {
			t.Errorf("GetShortHash failed: %v", err)
		}
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
		hash2, err := repo.GetShortHash()
		if err != nil {
			t.Errorf("GetShortHash failed after second commit: %v", err)
		}
		if hash2 == "" {
			t.Error("Expected non-empty commit hash after second commit")
		}
		if hash1 == hash2 {
			t.Error("Expected hash to change after new commit")
		}
	})
}

// TestGitStatus tests Git status detection functionality
func TestGitStatus(t *testing.T) {
	// ⭐ EXTRACT-004: Git status detection testing - 🧪

	// Skip if git is not available
	if !isGitAvailable() {
		t.Skip("Git not available, skipping status tests")
	}

	tmpDir, err := ioutil.TempDir("", "git_status_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	gitDir := filepath.Join(tmpDir, "status_test")
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

	repo := NewRepositoryWithConfig(&Config{
		WorkingDirectory:   gitDir,
		IncludeDirtyStatus: true,
		GitCommand:         "git",
	})

	t.Run("CleanWorkingDirectory", func(t *testing.T) {
		isClean, err := repo.IsWorkingDirectoryClean()
		if err != nil {
			t.Errorf("IsWorkingDirectoryClean failed: %v", err)
		}
		if !isClean {
			t.Error("Expected clean working directory")
		}

		// Test with status included
		info, err := repo.GetInfoWithStatus()
		if err != nil {
			t.Errorf("GetInfoWithStatus failed: %v", err)
		}
		if !info.IsClean {
			t.Error("Expected IsClean to be true for clean working directory")
		}
	})

	t.Run("DirtyWorkingDirectory", func(t *testing.T) {
		// Make working directory dirty
		if err := ioutil.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
			t.Fatal(err)
		}

		isClean, err := repo.IsWorkingDirectoryClean()
		if err != nil {
			t.Errorf("IsWorkingDirectoryClean failed: %v", err)
		}
		if isClean {
			t.Error("Expected dirty working directory")
		}

		// Test with status included
		info, err := repo.GetInfoWithStatus()
		if err != nil {
			t.Errorf("GetInfoWithStatus failed: %v", err)
		}
		if info.IsClean {
			t.Error("Expected IsClean to be false for dirty working directory")
		}
	})

	t.Run("StatusDisabled", func(t *testing.T) {
		repoNoStatus := NewRepositoryWithConfig(&Config{
			WorkingDirectory:   gitDir,
			IncludeDirtyStatus: false,
			GitCommand:         "git",
		})

		info, err := repoNoStatus.GetInfoWithStatus()
		if err != nil {
			t.Errorf("GetInfoWithStatus failed: %v", err)
		}
		// IsClean should be false (default) when status checking is disabled
		if info.IsClean {
			t.Error("Expected IsClean to be false when status checking is disabled")
		}
	})
}

// TestBackwardCompatibilityFunctions tests the package-level convenience functions
func TestBackwardCompatibilityFunctions(t *testing.T) {
	// ⭐ EXTRACT-004: Backward compatibility testing - 🧪

	tmpDir, err := ioutil.TempDir("", "git_compat_test_")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("NonGitDirectory", func(t *testing.T) {
		if IsGitRepository(tmpDir) {
			t.Error("Expected non-git directory to return false")
		}

		branch, hash := GetGitInfo(tmpDir)
		if branch != "" || hash != "" {
			t.Errorf("Expected empty git info for non-git directory, got branch=%s, hash=%s", branch, hash)
		}

		if IsGitWorkingDirectoryClean(tmpDir) {
			t.Error("Expected non-git directory to return false for clean status")
		}

		branch, hash, isClean := GetGitInfoWithStatus(tmpDir)
		if branch != "" || hash != "" || isClean {
			t.Errorf("Expected empty git info with status for non-git directory, got branch=%s, hash=%s, isClean=%t", branch, hash, isClean)
		}
	})

	if isGitAvailable() {
		t.Run("GitDirectory", func(t *testing.T) {
			gitDir := filepath.Join(tmpDir, "git_repo")
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

			// Test backward compatibility functions
			if !IsGitRepository(gitDir) {
				t.Error("Expected git directory to return true")
			}

			branch := GetGitBranch(gitDir)
			if branch == "" {
				t.Error("Expected non-empty branch name")
			}

			hash := GetGitShortHash(gitDir)
			if hash == "" {
				t.Error("Expected non-empty commit hash")
			}

			branchFromInfo, hashFromInfo := GetGitInfo(gitDir)
			if branchFromInfo != branch || hashFromInfo != hash {
				t.Errorf("GetGitInfo mismatch: expected branch=%s, hash=%s, got branch=%s, hash=%s", branch, hash, branchFromInfo, hashFromInfo)
			}

			if !IsGitWorkingDirectoryClean(gitDir) {
				t.Error("Expected clean working directory")
			}

			branchWithStatus, hashWithStatus, isClean := GetGitInfoWithStatus(gitDir)
			if branchWithStatus != branch || hashWithStatus != hash || !isClean {
				t.Errorf("GetGitInfoWithStatus mismatch: expected branch=%s, hash=%s, isClean=true, got branch=%s, hash=%s, isClean=%t",
					branch, hash, branchWithStatus, hashWithStatus, isClean)
			}
		})
	}
}

// TestGitError tests Git error handling
func TestGitError(t *testing.T) {
	// ⭐ EXTRACT-004: Git error handling testing - 🧪

	t.Run("GitErrorFormatting", func(t *testing.T) {
		gitErr := &GitError{
			Operation: "test operation",
			Err:       fmt.Errorf("test error"),
		}

		expected := "git test operation failed: test error"
		if gitErr.Error() != expected {
			t.Errorf("Expected error message '%s', got '%s'", expected, gitErr.Error())
		}
	})

	t.Run("NonGitRepositoryErrors", func(t *testing.T) {
		tmpDir, err := ioutil.TempDir("", "git_error_test_")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpDir)

		repo := NewRepositoryWithConfig(&Config{
			WorkingDirectory: tmpDir,
			GitCommand:       "git",
		})

		// Test that operations return appropriate errors for non-git repositories
		_, err = repo.GetBranch()
		if err == nil {
			t.Error("Expected error for GetBranch on non-git repository")
		}
		if gitErr, ok := err.(*GitError); ok {
			if gitErr.Operation != "branch detection" {
				t.Errorf("Expected operation 'branch detection', got '%s'", gitErr.Operation)
			}
		} else {
			t.Errorf("Expected GitError, got %T", err)
		}

		_, err = repo.GetShortHash()
		if err == nil {
			t.Error("Expected error for GetShortHash on non-git repository")
		}

		_, err = repo.IsWorkingDirectoryClean()
		if err == nil {
			t.Error("Expected error for IsWorkingDirectoryClean on non-git repository")
		}
	})
}

// Helper functions for testing

func isGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func runGitCommand(t *testing.T, dir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Git command failed: git %s in %s: %v", strings.Join(args, " "), dir, err)
	}
}
