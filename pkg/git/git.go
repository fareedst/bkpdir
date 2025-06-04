// ⭐ EXTRACT-004: Git Integration System - Extracted Git integration package - 🔧
// This file is part of bkpdir
//
// Package git provides Git integration for repository detection, metadata extraction,
// and status management. It offers a flexible Git operation framework suitable for
// many CLI applications.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// ⭐ EXTRACT-004: Git operation context and configuration - 🔧
// Config represents Git integration configuration options
type Config struct {
	// Directory to operate in (defaults to current directory)
	WorkingDirectory string
	// Whether to include dirty status in operations
	IncludeDirtyStatus bool
	// Custom Git command path (defaults to "git")
	GitCommand string
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		WorkingDirectory:   ".",
		IncludeDirtyStatus: true,
		GitCommand:         "git",
	}
}

// ⭐ EXTRACT-004: Git error handling structure - 🔧
// GitError represents an error that occurred during Git operations.
// It includes the operation that failed and the underlying error.
type GitError struct {
	Operation string
	Err       error
}

func (e *GitError) Error() string {
	return fmt.Sprintf("git %s failed: %v", e.Operation, e.Err)
}

// ⭐ EXTRACT-004: Git information structure - 🔧
// Info represents Git repository information
type Info struct {
	Branch  string
	Hash    string
	IsClean bool
	IsRepo  bool
}

// ⭐ EXTRACT-004: Git repository interface definition - 🔧
// Repository defines the interface for Git operations
type Repository interface {
	// IsRepository checks if the directory is a Git repository
	IsRepository() bool
	// GetBranch returns the current branch name
	GetBranch() (string, error)
	// GetShortHash returns the short commit hash
	GetShortHash() (string, error)
	// IsWorkingDirectoryClean checks if working directory is clean
	IsWorkingDirectoryClean() (bool, error)
	// GetInfo returns complete Git information
	GetInfo() (*Info, error)
	// GetInfoWithStatus returns Git information including status
	GetInfoWithStatus() (*Info, error)
}

// ⭐ EXTRACT-004: Git repository implementation - 🔧
// Repo implements the Repository interface using command-line Git
type Repo struct {
	config *Config
}

// NewRepository creates a new Git repository instance with default configuration
func NewRepository() Repository {
	return &Repo{config: DefaultConfig()}
}

// NewRepositoryWithConfig creates a new Git repository instance with custom configuration
func NewRepositoryWithConfig(config *Config) Repository {
	return &Repo{config: config}
}

// ⭐ EXTRACT-004: Generalized Git command execution framework - 🔧
// executeGitCommand runs a Git command with the configured parameters
func (r *Repo) executeGitCommand(args ...string) (string, error) {
	cmd := exec.Command(r.config.GitCommand, args...)
	cmd.Dir = r.config.WorkingDirectory
	out, err := cmd.Output()
	if err != nil {
		return "", &GitError{
			Operation: strings.Join(args, " "),
			Err:       err,
		}
	}
	return strings.TrimSpace(string(out)), nil
}

// ⭐ EXTRACT-004: Git repository detection implementation - 🔍
// IsRepository checks if the configured directory is a Git repository
func (r *Repo) IsRepository() bool {
	out, err := r.executeGitCommand("rev-parse", "--is-inside-work-tree")
	return err == nil && out == "true"
}

// ⭐ EXTRACT-004: Git branch extraction implementation - 🔍
// GetBranch returns the current Git branch name
func (r *Repo) GetBranch() (string, error) {
	if !r.IsRepository() {
		return "", &GitError{Operation: "branch detection", Err: fmt.Errorf("not a git repository")}
	}
	return r.executeGitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

// ⭐ EXTRACT-004: Git commit hash extraction implementation - 🔍
// GetShortHash returns the short commit hash of the current HEAD
func (r *Repo) GetShortHash() (string, error) {
	if !r.IsRepository() {
		return "", &GitError{Operation: "hash extraction", Err: fmt.Errorf("not a git repository")}
	}
	return r.executeGitCommand("rev-parse", "--short", "HEAD")
}

// ⭐ EXTRACT-004: Git working directory state detection - 🔍
// IsWorkingDirectoryClean checks if the Git working directory is clean
func (r *Repo) IsWorkingDirectoryClean() (bool, error) {
	if !r.IsRepository() {
		return false, &GitError{Operation: "status check", Err: fmt.Errorf("not a git repository")}
	}

	out, err := r.executeGitCommand("status", "--porcelain")
	if err != nil {
		return false, err
	}

	return len(out) == 0, nil
}

// ⭐ EXTRACT-004: Combined Git information extraction - 🔍
// GetInfo returns complete Git repository information
func (r *Repo) GetInfo() (*Info, error) {
	info := &Info{IsRepo: r.IsRepository()}
	if !info.IsRepo {
		return info, nil
	}

	var err error
	info.Branch, err = r.GetBranch()
	if err != nil {
		return info, err
	}

	info.Hash, err = r.GetShortHash()
	if err != nil {
		return info, err
	}

	return info, nil
}

// ⭐ EXTRACT-004: Combined Git information extraction with status - 🔍
// GetInfoWithStatus returns complete Git information including working directory status
func (r *Repo) GetInfoWithStatus() (*Info, error) {
	info, err := r.GetInfo()
	if err != nil || !info.IsRepo {
		return info, err
	}

	if r.config.IncludeDirtyStatus {
		info.IsClean, err = r.IsWorkingDirectoryClean()
		if err != nil {
			return info, err
		}
	}

	return info, nil
}

// ⭐ EXTRACT-004: Convenience functions for backward compatibility - 🔧
// Package-level convenience functions that maintain the original API

// IsGitRepository checks if the given directory is a Git repository
func IsGitRepository(dir string) bool {
	config := &Config{WorkingDirectory: dir, GitCommand: "git"}
	repo := &Repo{config: config}
	return repo.IsRepository()
}

// GetGitBranch returns the current Git branch name for the given directory
func GetGitBranch(dir string) string {
	config := &Config{WorkingDirectory: dir, GitCommand: "git"}
	repo := &Repo{config: config}
	branch, err := repo.GetBranch()
	if err != nil {
		return ""
	}
	return branch
}

// GetGitShortHash returns the short commit hash for the given directory
func GetGitShortHash(dir string) string {
	config := &Config{WorkingDirectory: dir, GitCommand: "git"}
	repo := &Repo{config: config}
	hash, err := repo.GetShortHash()
	if err != nil {
		return ""
	}
	return hash
}

// GetGitInfo returns both branch name and commit hash for the given directory
func GetGitInfo(dir string) (branch, hash string) {
	config := &Config{WorkingDirectory: dir, GitCommand: "git"}
	repo := &Repo{config: config}
	info, err := repo.GetInfo()
	if err != nil || !info.IsRepo {
		return "", ""
	}
	return info.Branch, info.Hash
}

// IsGitWorkingDirectoryClean checks if the Git working directory is clean
func IsGitWorkingDirectoryClean(dir string) bool {
	config := &Config{WorkingDirectory: dir, GitCommand: "git"}
	repo := &Repo{config: config}
	isClean, err := repo.IsWorkingDirectoryClean()
	if err != nil {
		return false
	}
	return isClean
}

// GetGitInfoWithStatus returns branch name, commit hash, and working directory status
func GetGitInfoWithStatus(dir string) (branch, hash string, isClean bool) {
	config := &Config{
		WorkingDirectory:   dir,
		IncludeDirtyStatus: true, // Enable dirty status for backward compatibility
		GitCommand:         "git",
	}
	repo := &Repo{config: config}
	info, err := repo.GetInfoWithStatus()
	if err != nil || !info.IsRepo {
		return "", "", false
	}
	return info.Branch, info.Hash, info.IsClean
}
