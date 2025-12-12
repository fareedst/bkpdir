// This file is part of bkpdir
//
// Package main provides Git integration for archive naming and metadata in BkpDir.
// It handles Git branch and commit hash detection for archive naming.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [REQ:GIT_INTEGRATION] Git integration for archive naming and metadata
// [ARCH:GIT_INTEGRATION] Git command-line integration architecture
// [IMPL:GIT_CLI] Git command-line interface implementation
// GIT-INTEGRATION-001: Git integration specification - Git integration and metadata [ACTION:core-functionality]
// Source: git.go - GIT-INTEGRATION-001
// Impact: Core functionality requirement for Git integration

// SERVICE-GIT-001: Git service architecture decision - Git service implementation [ACTION:core-functionality]
// Source: git.go - SERVICE-GIT-001
// Impact: Git service implementation decision
package main

import (
	"bkpdir/pkg/git"
)

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// This file now serves as an adapter layer to maintain backward compatibility
// while using the extracted Git integration system.

// REFACTOR-001: See architecture.md - Interface Contracts [DECISION:maintenance]
// REFACTOR-001: See architecture.md - Interface Contracts [DECISION:maintenance]

// GIT-001: See specification.md - Git Information Extraction [DECISION:maintenance]
// GitError is now provided by the pkg/git package
type GitError = git.GitError

// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
// These functions maintain the original API while delegating to the extracted package

// GIT-001: See specification.md - Git Information Extraction [DECISION:maintenance]
// IsGitRepository checks if the given directory is a Git repository.
// It uses git rev-parse to check if the directory is inside a git work tree.
func IsGitRepository(dir string) bool {
	return git.IsGitRepository(dir)
}

// GIT-002: See specification.md - Git Branch and Commit Extraction [DECISION:maintenance]
// GetGitBranch returns the current Git branch name.
// It returns an empty string if not in a Git repository or on error.
func GetGitBranch(dir string) string {
	return git.GetGitBranch(dir)
}

// GIT-002: See specification.md - Git Branch and Commit Extraction [DECISION:maintenance]
// GetGitShortHash returns the short commit hash of the current HEAD.
// It returns an empty string if not in a Git repository or on error.
func GetGitShortHash(dir string) string {
	return git.GetGitShortHash(dir)
}

// GIT-001: See specification.md - Git Information Extraction [DECISION:maintenance]
// GetGitInfo returns both the branch name and commit hash.
// It returns empty strings if not in a Git repository or on error.
func GetGitInfo(dir string) (branch, hash string) {
	return git.GetGitInfo(dir)
}

// GIT-003: See specification.md - Git Status Detection [DECISION:maintenance]
// IsGitWorkingDirectoryClean checks if the Git working directory is clean (no uncommitted changes).
// It returns true if the working directory is clean, false if there are changes or not in a Git repository.
func IsGitWorkingDirectoryClean(dir string) bool {
	return git.IsGitWorkingDirectoryClean(dir)
}

// GIT-003: See specification.md - Git Status Detection [DECISION:maintenance]
// GetGitInfoWithStatus returns branch name, commit hash, and working directory status.
// It returns empty strings and false for status if not in a Git repository or on error.
func GetGitInfoWithStatus(dir string) (branch, hash string, isClean bool) {
	return git.GetGitInfoWithStatus(dir)
}

// GIT-004: See architecture.md - Git Submodule Support [DECISION:maintenance]
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
