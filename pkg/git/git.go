// [REQ-GIT_INTEGRATION] Git integration for repository detection and metadata extraction
// [ARCH-GIT_INTEGRATION] Git command-line integration architecture
// [IMPL-GIT_CLI] Git command-line interface implementation
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

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// Config represents Git integration configuration options
type Config struct {
	// Basic Git integration settings
	Enabled         bool // Enable/disable Git integration (default: true)
	IncludeInfo     bool // Include Git info in operations (default: false)
	ShowDirtyStatus bool // Show dirty status indicator (default: false)

	// Git command configuration
	Command          string // Git command path (default: "git")
	WorkingDirectory string // Working directory for Git operations (default: ".")

	// Git behavior settings
	RequireCleanRepo  bool // Fail operations if repository is dirty (default: false)
	AutoDetectRepo    bool // Automatically detect Git repositories (default: true)
	IncludeSubmodules bool // Include submodule information (default: false)

	// Git information inclusion
	IncludeBranch bool // Include branch name in operations (default: true)
	IncludeHash   bool // Include commit hash in operations (default: true)
	IncludeStatus bool // Include working directory status (default: true)

	// Git command timeouts and limits
	CommandTimeout    string // Timeout for Git commands (default: "30s")
	MaxSubmoduleDepth int    // Maximum submodule recursion depth (default: 3)

	// Legacy compatibility fields (to be deprecated)
	IncludeDirtyStatus bool   // Legacy: use ShowDirtyStatus instead
	GitCommand         string // Legacy: use Command instead
}

// DefaultConfig returns a Config with sensible defaults
// [IMPL-GIT_CLI] Provides default configuration for git operations
func DefaultConfig() *Config {
	return &Config{
		Enabled:           true,
		IncludeInfo:       false,
		ShowDirtyStatus:   false,
		Command:           "git",
		WorkingDirectory:  ".",
		RequireCleanRepo:  false,
		AutoDetectRepo:    true,
		IncludeSubmodules: false,
		IncludeBranch:     true,
		IncludeHash:       true,
		IncludeStatus:     true,
		CommandTimeout:    "30s",
		MaxSubmoduleDepth: 3,
		// Legacy compatibility
		IncludeDirtyStatus: true,  // Legacy default
		GitCommand:         "git", // Legacy default
	}
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// GitError represents an error that occurred during Git operations.
// It includes the operation that failed and the underlying error.
type GitError struct {
	Operation string
	Err       error
}

func (e *GitError) Error() string {
	return fmt.Sprintf("git %s failed: %v", e.Operation, e.Err)
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// Info represents Git repository information
type Info struct {
	Branch      string
	Hash        string
	IsClean     bool
	IsRepo      bool
	IsSubmodule bool
	Submodules  []SubmoduleInfo
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// SubmoduleInfo represents information about a Git submodule
type SubmoduleInfo struct {
	Name   string // Submodule name
	Path   string // Submodule path relative to repository root
	URL    string // Submodule remote URL
	Hash   string // Current commit hash of the submodule
	Status string // Submodule status (e.g., "clean", "dirty", "uninitialized")
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// Repository defines the interface for Git operations
type Repository interface {
	IsRepository() bool
	GetBranch() (string, error)
	GetShortHash() (string, error)
	IsWorkingDirectoryClean() (bool, error)
	GetInfo() (*Info, error)
	GetInfoWithStatus() (*Info, error)
	// [IMPL-GIT_CLI] Submodule interface methods
	IsSubmodule() (bool, error)
	GetSubmodules() ([]SubmoduleInfo, error)
	GetSubmoduleStatus(path string) (string, error)
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
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

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// executeGitCommand runs a Git command with the configured parameters
func (r *Repo) executeGitCommand(args ...string) (string, error) {
	// [IMPL-GIT_CLI] Resolves git command path with legacy fallback
	gitCmd := r.config.Command
	if gitCmd == "" {
		gitCmd = r.config.GitCommand // Legacy fallback
	}
	if gitCmd == "" {
		gitCmd = "git" // Ultimate fallback
	}

	cmd := exec.Command(gitCmd, args...)
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

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// IsRepository checks if the configured directory is a Git repository
func (r *Repo) IsRepository() bool {
	out, err := r.executeGitCommand("rev-parse", "--is-inside-work-tree")
	return err == nil && out == "true"
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// GetBranch returns the current Git branch name
func (r *Repo) GetBranch() (string, error) {
	if !r.IsRepository() {
		return "", &GitError{Operation: "branch detection", Err: fmt.Errorf("not a git repository")}
	}
	return r.executeGitCommand("rev-parse", "--abbrev-ref", "HEAD")
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// GetShortHash returns the short commit hash of the current HEAD
func (r *Repo) GetShortHash() (string, error) {
	if !r.IsRepository() {
		return "", &GitError{Operation: "hash extraction", Err: fmt.Errorf("not a git repository")}
	}
	return r.executeGitCommand("rev-parse", "--short", "HEAD")
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
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

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
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

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// GetInfoWithStatus returns complete Git information including working directory status
func (r *Repo) GetInfoWithStatus() (*Info, error) {
	info, err := r.GetInfo()
	if err != nil || !info.IsRepo {
		return info, err
	}

	// [IMPL-GIT_CLI] Check dirty status based on config flags
	includeDirtyStatus := r.config.ShowDirtyStatus || r.config.IncludeDirtyStatus
	if includeDirtyStatus {
		info.IsClean, err = r.IsWorkingDirectoryClean()
		if err != nil {
			return info, err
		}
	}

	// [IMPL-GIT_CLI] Include submodule info if configured
	if r.config.IncludeSubmodules {
		info.IsSubmodule, err = r.IsSubmodule()
		if err != nil {
			return info, err
		}

		info.Submodules, err = r.GetSubmodules()
		if err != nil {
			return info, err
		}
	}

	return info, nil
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// IsSubmodule checks if the current directory is a Git submodule
func (r *Repo) IsSubmodule() (bool, error) {
	if !r.IsRepository() {
		return false, nil
	}

	out, err := r.executeGitCommand("rev-parse", "--show-superproject-working-tree")
	if err != nil {
		return false, nil
	}

	return strings.TrimSpace(out) != "", nil
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// GetSubmodules returns information about all submodules in the repository
func (r *Repo) GetSubmodules() ([]SubmoduleInfo, error) {
	if !r.IsRepository() {
		return nil, &GitError{Operation: "submodule listing", Err: fmt.Errorf("not a git repository")}
	}

	out, err := r.executeGitCommand("submodule", "status", "--recursive")
	if err != nil {
		// If submodule command fails, there might be no submodules
		return []SubmoduleInfo{}, nil
	}

	var submodules []SubmoduleInfo
	lines := strings.Split(strings.TrimSpace(out), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		submodule, err := r.parseSubmoduleStatusLine(line)
		if err != nil {
			continue // Skip malformed lines
		}

		submodules = append(submodules, submodule)
	}

	return submodules, nil
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// GetSubmoduleStatus returns the status of a specific submodule
func (r *Repo) GetSubmoduleStatus(path string) (string, error) {
	if !r.IsRepository() {
		return "", &GitError{Operation: "submodule status", Err: fmt.Errorf("not a git repository")}
	}

	out, err := r.executeGitCommand("submodule", "status", path)
	if err != nil {
		return "", &GitError{Operation: "submodule status", Err: err}
	}

	// Parse the status from the first character of the output
	if len(out) == 0 {
		return "unknown", nil
	}

	switch out[0] {
	case ' ':
		return "clean", nil
	case '+':
		return "dirty", nil
	case '-':
		return "uninitialized", nil
	case 'U':
		return "conflict", nil
	default:
		return "unknown", nil
	}
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// parseSubmoduleStatusLine parses a line from git submodule status output
func (r *Repo) parseSubmoduleStatusLine(line string) (SubmoduleInfo, error) {
	// Git submodule status format: [status][hash] [path] [(description)]
	// Status characters: ' ' (clean), '+' (dirty), '-' (uninitialized), 'U' (conflict)

	if len(line) < 2 {
		return SubmoduleInfo{}, fmt.Errorf("invalid submodule status line: %s", line)
	}

	status := "unknown"
	switch line[0] {
	case ' ':
		status = "clean"
	case '+':
		status = "dirty"
	case '-':
		status = "uninitialized"
	case 'U':
		status = "conflict"
	}

	// Remove status character and parse the rest
	line = line[1:]
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return SubmoduleInfo{}, fmt.Errorf("invalid submodule status format: %s", line)
	}

	hash := parts[0]
	path := parts[1]

	// Extract submodule name from path (last component)
	name := path
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		name = path[idx+1:]
	}

	url, err := r.getSubmoduleURL(path)
	if err != nil {
		url = "" // URL is optional
	}

	return SubmoduleInfo{
		Name:   name,
		Path:   path,
		URL:    url,
		Hash:   hash,
		Status: status,
	}, nil
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// getSubmoduleURL gets the remote URL for a submodule
func (r *Repo) getSubmoduleURL(path string) (string, error) {
	out, err := r.executeGitCommand("config", "--file", ".gitmodules", "--get", "submodule."+path+".url")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// Package-level convenience functions that maintain the original API

// IsGitRepository checks if the given directory is a Git repository
// [IMPL-GIT_CLI] Creates ephemeral Repo, delegates to IsRepository
func IsGitRepository(dir string) bool {
	config := &Config{WorkingDirectory: dir, GitCommand: "git"}
	repo := &Repo{config: config}
	return repo.IsRepository()
}

// GetGitBranch returns the current Git branch name for the given directory
// [IMPL-GIT_CLI] Creates ephemeral Repo, delegates to GetBranch
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
// [IMPL-GIT_CLI] Creates ephemeral Repo, delegates to GetShortHash
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
// [IMPL-GIT_CLI] Creates ephemeral Repo, delegates to GetInfo
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
// [IMPL-GIT_CLI] Creates ephemeral Repo, delegates to IsWorkingDirectoryClean
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
// [IMPL-GIT_CLI] Creates ephemeral Repo with dirty-status enabled, delegates to GetInfoWithStatus
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

// [IMPL-GIT_CLI] [ARCH-GIT_INTEGRATION] [REQ-GIT_INTEGRATION]
// Submodule convenience functions

// IsGitSubmodule checks if the given directory is a Git submodule
// [IMPL-GIT_CLI] Creates ephemeral Repo, delegates to IsSubmodule
func IsGitSubmodule(dir string) bool {
	config := &Config{WorkingDirectory: dir, GitCommand: "git"}
	repo := &Repo{config: config}
	isSubmodule, err := repo.IsSubmodule()
	if err != nil {
		return false
	}
	return isSubmodule
}

// GetGitSubmodules returns information about all submodules in the given directory
// [IMPL-GIT_CLI] Creates ephemeral Repo, delegates to GetSubmodules
func GetGitSubmodules(dir string) []SubmoduleInfo {
	config := &Config{WorkingDirectory: dir, GitCommand: "git"}
	repo := &Repo{config: config}
	submodules, err := repo.GetSubmodules()
	if err != nil {
		return []SubmoduleInfo{}
	}
	return submodules
}

// GetGitSubmoduleStatus returns the status of a specific submodule
// [IMPL-GIT_CLI] Creates ephemeral Repo, delegates to GetSubmoduleStatus
func GetGitSubmoduleStatus(dir, path string) string {
	config := &Config{WorkingDirectory: dir, GitCommand: "git"}
	repo := &Repo{config: config}
	status, err := repo.GetSubmoduleStatus(path)
	if err != nil {
		return "unknown"
	}
	return status
}
