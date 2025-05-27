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

// GitError represents an error that occurred during Git operations
type GitError struct {
	Operation string
	Err       error
}

func (e *GitError) Error() string {
	return fmt.Sprintf("git %s failed: %v", e.Operation, e.Err)
}

// IsGitRepository returns true if the directory is a git repository.
// It uses git rev-parse to check if the directory is inside a git work tree.
func IsGitRepository(dir string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = dir
	out, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// GetGitBranch returns the current git branch name.
// Returns empty string if not in a git repository or on error.
func GetGitBranch(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GetGitShortHash returns the short commit hash of the current HEAD.
// Returns empty string if not in a git repository or on error.
func GetGitShortHash(dir string) string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// GetGitInfo returns both branch name and commit hash.
// Returns empty strings if not in a git repository or on error.
func GetGitInfo(dir string) (branch, hash string) {
	if !IsGitRepository(dir) {
		return "", ""
	}
	return GetGitBranch(dir), GetGitShortHash(dir)
}
