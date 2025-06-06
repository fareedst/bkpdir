// This file is part of bkpdir
//
// Package main provides Git integration for archive naming and metadata in BkpDir.
// It handles Git branch and commit hash detection for archive naming.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License
package main

import (
	"bkpdir/pkg/git"
)

// ⭐ EXTRACT-004: Git integration adapter using extracted pkg/git package - 🔧
// This file now serves as an adapter layer to maintain backward compatibility
// while using the extracted Git integration system.

// 🔶 REFACTOR-001: Git provider interface contracts defined - 🔧
// 🔶 REFACTOR-001: Dependency analysis - clean boundary confirmed - 🔧

// 🔺 GIT-001: Git error handling structure - 🔧
// GitError is now provided by the pkg/git package
type GitError = git.GitError

// ⭐ EXTRACT-004: Backward compatibility functions using pkg/git - 🔧
// These functions maintain the original API while delegating to the extracted package

// 🔺 GIT-001: Git repository detection implementation - 🔍
// IsGitRepository checks if the given directory is a Git repository.
// It uses git rev-parse to check if the directory is inside a git work tree.
func IsGitRepository(dir string) bool {
	return git.IsGitRepository(dir)
}

// 🔺 GIT-002: Git branch extraction implementation - 🔍
// GetGitBranch returns the current Git branch name.
// It returns an empty string if not in a Git repository or on error.
func GetGitBranch(dir string) string {
	return git.GetGitBranch(dir)
}

// 🔺 GIT-002: Git commit hash extraction implementation - 🔍
// GetGitShortHash returns the short commit hash of the current HEAD.
// It returns an empty string if not in a Git repository or on error.
func GetGitShortHash(dir string) string {
	return git.GetGitShortHash(dir)
}

// 🔺 GIT-001: Combined Git information extraction - 🔍
// GetGitInfo returns both the branch name and commit hash.
// It returns empty strings if not in a Git repository or on error.
func GetGitInfo(dir string) (branch, hash string) {
	return git.GetGitInfo(dir)
}

// 🔺 GIT-003: Git working directory state detection - 🔍
// IsGitWorkingDirectoryClean checks if the Git working directory is clean (no uncommitted changes).
// It returns true if the working directory is clean, false if there are changes or not in a Git repository.
func IsGitWorkingDirectoryClean(dir string) bool {
	return git.IsGitWorkingDirectoryClean(dir)
}

// 🔺 GIT-003: Combined Git information extraction with status - 🔍
// GetGitInfoWithStatus returns branch name, commit hash, and working directory status.
// It returns empty strings and false for status if not in a Git repository or on error.
func GetGitInfoWithStatus(dir string) (branch, hash string, isClean bool) {
	return git.GetGitInfoWithStatus(dir)
}

// 🔶 GIT-004: Git submodule support functions - 🔧
// These functions provide Git submodule detection and management capabilities.

// IsGitSubmodule checks if the given directory is a Git submodule.
// It returns true if the directory is a submodule, false otherwise.
func IsGitSubmodule(dir string) bool {
	return git.IsGitSubmodule(dir)
}

// GetGitSubmodules returns information about all submodules in the given directory.
// It returns an empty slice if there are no submodules or if not in a Git repository.
func GetGitSubmodules(dir string) []git.SubmoduleInfo {
	return git.GetGitSubmodules(dir)
}

// GetGitSubmoduleStatus returns the status of a specific submodule.
// It returns "unknown" if the submodule doesn't exist or if not in a Git repository.
func GetGitSubmoduleStatus(dir, path string) string {
	return git.GetGitSubmoduleStatus(dir, path)
}
