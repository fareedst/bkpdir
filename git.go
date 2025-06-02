// This file is part of bkpdir
//
// Package main provides Git integration for archive naming and metadata in BkpDir.
// It handles Git branch and commit hash detection for archive naming.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// REFACTOR-001: Git provider interface contracts defined
// REFACTOR-001: Dependency analysis - clean boundary confirmed

// GIT-001: Git error handling structure
// IMMUTABLE-REF: Error Handling Requirements
// TEST-REF: TestGitIntegration
// DECISION-REF: DEC-005
// GitError represents an error that occurred during Git operations.
// It includes the operation that failed and the underlying error.
type GitError struct {
	Operation string
	Err       error
}

func (e *GitError) Error() string {
	return fmt.Sprintf("git %s failed: %v", e.Operation, e.Err)
}

// GIT-001: Git repository detection implementation
// IMMUTABLE-REF: Git Integration Requirements
// TEST-REF: TestGitIntegration
// DECISION-REF: DEC-005
// IsGitRepository checks if the given directory is a Git repository.
// It uses git rev-parse to check if the directory is inside a git work tree.
func IsGitRepository(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = dir
	out, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// GIT-002: Git branch extraction implementation
// IMMUTABLE-REF: Git Integration Requirements, Archive Naming Convention
// TEST-REF: TestGitNaming
// DECISION-REF: DEC-005
// GetGitBranch returns the current Git branch name.
// It returns an empty string if not in a Git repository or on error.
func GetGitBranch(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GIT-002: Git commit hash extraction implementation
// IMMUTABLE-REF: Git Integration Requirements, Archive Naming Convention
// TEST-REF: TestGitNaming
// DECISION-REF: DEC-005
// GetGitShortHash returns the short commit hash of the current HEAD.
// It returns an empty string if not in a Git repository or on error.
func GetGitShortHash(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GIT-001: Combined Git information extraction
// IMMUTABLE-REF: Git Integration Requirements
// TEST-REF: TestGitIntegration
// DECISION-REF: DEC-005
// GetGitInfo returns both the branch name and commit hash.
// It returns empty strings if not in a Git repository or on error.
func GetGitInfo(dir string) (branch, hash string) {
	if !IsGitRepository(dir) {
		return "", ""
	}
	return GetGitBranch(dir), GetGitShortHash(dir)
}

// GIT-003: Git working directory state detection
// IMMUTABLE-REF: Git Integration Requirements
// TEST-REF: TestGitStatus
// DECISION-REF: DEC-005
// IsGitWorkingDirectoryClean checks if the Git working directory is clean (no uncommitted changes).
// It returns true if the working directory is clean, false if there are changes or not in a Git repository.
func IsGitWorkingDirectoryClean(dir string) bool {
	if !IsGitRepository(dir) {
		return false
	}

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) == 0
}

// GIT-003: Combined Git information extraction with status
// IMMUTABLE-REF: Git Integration Requirements
// TEST-REF: TestGitStatus
// DECISION-REF: DEC-005
// GetGitInfoWithStatus returns branch name, commit hash, and working directory status.
// It returns empty strings and false for status if not in a Git repository or on error.
func GetGitInfoWithStatus(dir string) (branch, hash string, isClean bool) {
	if !IsGitRepository(dir) {
		return "", "", false
	}
	return GetGitBranch(dir), GetGitShortHash(dir), IsGitWorkingDirectoryClean(dir)
}
