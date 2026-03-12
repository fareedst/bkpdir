// This file is part of bkpdir
//
// Package main provides Git integration for archive naming and metadata in BkpDir.
// It handles Git branch and commit hash detection for archive naming.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [REQ-GIT_INTEGRATION] Git integration for archive naming and metadata
// [ARCH-GIT_INTEGRATION] Git command-line integration architecture
// [IMPL-GIT_CLI] Git command-line interface implementation
package main

import (
	"bkpdir/pkg/git"
)

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// This file serves as an adapter layer to maintain backward compatibility
// while delegating to the extracted Git integration package.

// GitError is now provided by the pkg/git package
// [IMPL-GIT_CLI] Type alias for backward compatibility
type GitError = git.GitError

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// Adapter functions maintaining the original API while delegating to pkg/git

// IsGitRepository checks if the given directory is a Git repository.
// It uses git rev-parse to check if the directory is inside a git work tree.
// [IMPL-GIT_CLI] Delegates to git.IsGitRepository
func IsGitRepository(dir string) bool {
	return git.IsGitRepository(dir)
}

// GetGitBranch returns the current Git branch name.
// It returns an empty string if not in a Git repository or on error.
// [IMPL-GIT_CLI] Delegates to git.GetGitBranch
func GetGitBranch(dir string) string {
	return git.GetGitBranch(dir)
}

// GetGitShortHash returns the short commit hash of the current HEAD.
// It returns an empty string if not in a Git repository or on error.
// [IMPL-GIT_CLI] Delegates to git.GetGitShortHash
func GetGitShortHash(dir string) string {
	return git.GetGitShortHash(dir)
}

// GetGitInfo returns both the branch name and commit hash.
// It returns empty strings if not in a Git repository or on error.
// [IMPL-GIT_CLI] Delegates to git.GetGitInfo
func GetGitInfo(dir string) (branch, hash string) {
	return git.GetGitInfo(dir)
}

// IsGitWorkingDirectoryClean checks if the Git working directory is clean (no uncommitted changes).
// It returns true if the working directory is clean, false if there are changes or not in a Git repository.
// [IMPL-GIT_CLI] Delegates to git.IsGitWorkingDirectoryClean
func IsGitWorkingDirectoryClean(dir string) bool {
	return git.IsGitWorkingDirectoryClean(dir)
}

// GetGitInfoWithStatus returns branch name, commit hash, and working directory status.
// It returns empty strings and false for status if not in a Git repository or on error.
// [IMPL-GIT_CLI] Delegates to git.GetGitInfoWithStatus
func GetGitInfoWithStatus(dir string) (branch, hash string, isClean bool) {
	return git.GetGitInfoWithStatus(dir)
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// Submodule adapter functions

// IsGitSubmodule checks if the given directory is a Git submodule.
// It returns true if the directory is a submodule, false otherwise.
// [IMPL-GIT_CLI] Delegates to git.IsGitSubmodule
func IsGitSubmodule(dir string) bool {
	return git.IsGitSubmodule(dir)
}

// GetGitSubmodules returns information about all submodules in the given directory.
// It returns an empty slice if there are no submodules or if not in a Git repository.
// [IMPL-GIT_CLI] Delegates to git.GetGitSubmodules
func GetGitSubmodules(dir string) []git.SubmoduleInfo {
	return git.GetGitSubmodules(dir)
}

// GetGitSubmoduleStatus returns the status of a specific submodule.
// It returns "unknown" if the submodule doesn't exist or if not in a Git repository.
// [IMPL-GIT_CLI] Delegates to git.GetGitSubmoduleStatus
func GetGitSubmoduleStatus(dir, path string) string {
	return git.GetGitSubmoduleStatus(dir, path)
}
