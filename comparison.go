// This file is part of bkpdir
//
// Package main provides directory comparison functionality for the BkpDir application.
// It handles comparing directories to detect changes and identical states.

// Directory comparison for incremental archive detection
// [ARCH-DIRECTORY_COMPARISON] Snapshot-based directory comparison system
// [ARCH-PACKAGE_EXTRACTION] Uses extracted fileops package for comparison

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"bkpdir/pkg/fileops"
)


// - [IMPL-PACKAGE_EXTRACTION] [ARCH-PACKAGE_EXTRACTION] [REQ-MAINTAINABILITY] — How: for each package in a phase move code to pkg/, define interfaces first, update imports, and require go build plus all tests green.
// - [IMPL-PACKAGE_EXTRACTION] [ARCH-PACKAGE_EXTRACTION] [REQ-MAINTAINABILITY] — How: expose type aliases and delegating wrappers at root with deprecation comments pointing to new pkg paths.

// Legacy type aliases for backward compatibility
type (
	FileInfo          = fileops.FileInfo
	DirectorySnapshot = fileops.DirectorySnapshot
)

// Legacy function wrappers for backward compatibility

// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: delegate to fileops.CreateDirectorySnapshot with exclusion patterns applied during walk.
func CreateDirectorySnapshot(rootPath string, excludePatterns []string) (*DirectorySnapshot, error) {
	return fileops.CreateDirectorySnapshot(rootPath, excludePatterns)
}

// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: delegate to fileops.CreateArchiveSnapshot to enumerate zip member files with metadata.
func CreateArchiveSnapshot(archivePath string) (*DirectorySnapshot, error) {
	return fileops.CreateArchiveSnapshot(archivePath)
}

// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: delegate equality check of two DirectorySnapshot values to fileops.CompareSnapshots.
func CompareSnapshots(snapshot1, snapshot2 *DirectorySnapshot) bool {
	return fileops.CompareSnapshots(snapshot1, snapshot2)
}

// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] [REQ-PERFORMANCE] — How: build dir and archive snapshots with exclusions and compare for full structural equality.
func IsDirectoryIdenticalToArchive(dirPath, archivePath string, excludePatterns []string) (bool, error) {
	return fileops.IsDirectoryIdenticalToArchive(dirPath, archivePath, excludePatterns)
}

// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: list archives, keep full archives only, sort by name, return path of last entry.
func FindMostRecentArchive(archiveDir string) (string, error) {
	archives, err := ListArchives(archiveDir)
	if err != nil {
		return "", err
	}

	if len(archives) == 0 {
		return "", nil
	}

	// Filter full archives only (skip incremental archives)
	var fullArchives []Archive
	for _, archive := range archives {
		if !archive.IsIncremental {
			fullArchives = append(fullArchives, archive)
		}
	}

	if len(fullArchives) == 0 {
		return "", nil
	}

	// Sort by name (archive names include timestamps that are alphabetically sortable)
	// Most recent archive will be last when sorted ascending
	sort.Slice(fullArchives, func(i, j int) bool {
		return fullArchives[i].Name < fullArchives[j].Name
	})

	// Return the last (most recent) archive
	mostRecent := fullArchives[len(fullArchives)-1]
	return mostRecent.Path, nil
}

// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: FindMostRecentArchive then IsDirectoryIdenticalToArchive against cwd.
func CheckForIdenticalArchive(dirPath, archiveDir string, excludePatterns []string) (bool, string, error) {
	mostRecentArchive, err := FindMostRecentArchive(archiveDir)
	if err != nil {
		return false, "", err
	}

	if mostRecentArchive == "" {
		// No archives exist
		return false, "", nil
	}

	// Check if directory is identical to most recent archive
	identical, err := IsDirectoryIdenticalToArchive(dirPath, mostRecentArchive, excludePatterns)
	if err != nil {
		return false, "", err
	}

	return identical, mostRecentArchive, nil
}

// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: build human-readable directory listing from CreateDirectorySnapshot for diagnostics.
func GetDirectoryTreeSummary(dirPath string, excludePatterns []string) (string, error) {
	snapshot, err := CreateDirectorySnapshot(dirPath, excludePatterns)
	if err != nil {
		return "", err
	}

	summary := fmt.Sprintf("Directory: %s\nFiles: %d\n", dirPath, len(snapshot.Files))
	for _, file := range snapshot.Files {
		if file.IsDir {
			summary += fmt.Sprintf("  [DIR]  %s\n", file.RelativePath)
		} else {
			summary += fmt.Sprintf("  [FILE] %s (%d bytes)\n", file.RelativePath, file.Size)
		}
	}

	return summary, nil
}

// - [IMPL-DIRECTORY_COMPARISON] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] [REQ-FILE_BACKUP] — How: build human-readable file listing from CreateArchiveSnapshot for diagnostics and tests.
func GetArchiveTreeSummary(archivePath string) (string, error) {
	snapshot, err := CreateArchiveSnapshot(archivePath)
	if err != nil {
		return "", err
	}

	summary := fmt.Sprintf("Archive: %s\nFiles: %d\n", archivePath, len(snapshot.Files))
	for _, file := range snapshot.Files {
		summary += fmt.Sprintf("  [FILE] %s (%d bytes)\n", file.RelativePath, file.Size)
	}

	return summary, nil
}

// DiffResult holds the results of a diff comparison
type DiffResult struct {
	Added    []string
	Modified []string
	Deleted  []string
}

// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: select incrementals whose name prefix matches base full archive _update= pattern; return latest by name.
func findLatestIncrementalArchive(archiveDir string, baseFullArchive *Archive) (*Archive, error) {
	archives, err := ListArchives(archiveDir)
	if err != nil {
		return nil, err
	}
	if len(archives) == 0 {
		return nil, nil
	}

	// Extract base name from full archive (without .zip extension)
	baseName := strings.TrimSuffix(baseFullArchive.Name, ".zip")
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: findLatestIncrementalArchive - Looking for incrementals based on: %s\n", baseName)
	}

	// Filter incremental archives that are based on the given full archive
	// Incremental archives are named: BASENAME_update=...
	var matchingIncrementals []Archive
	expectedPrefix := baseName + "_update="

	for i := range archives {
		if archives[i].IsIncremental {
			incrementalBaseName := strings.TrimSuffix(archives[i].Name, ".zip")
			if !strings.HasPrefix(incrementalBaseName, expectedPrefix) {
				if debug {
					fmt.Fprintf(os.Stderr, "DEBUG: findLatestIncrementalArchive - Skipping %s (doesn't match prefix %s)\n", archives[i].Name, expectedPrefix)
				}
				continue // Skip incrementals not based on this full archive
			}
			if debug {
				fmt.Fprintf(os.Stderr, "DEBUG: findLatestIncrementalArchive - Considering incremental: %s\n", archives[i].Name)
			}
			matchingIncrementals = append(matchingIncrementals, archives[i])
		}
	}

	if len(matchingIncrementals) == 0 {
		return nil, nil
	}

	// Sort by name (archive names include timestamps that are alphabetically sortable)
	// Most recent archive will be last when sorted ascending
	sort.Slice(matchingIncrementals, func(i, j int) bool {
		return matchingIncrementals[i].Name < matchingIncrementals[j].Name
	})

	// Return the last (most recent) incremental archive
	latestIncremental := matchingIncrementals[len(matchingIncrementals)-1]
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: findLatestIncrementalArchive - Selected latest: %s\n", latestIncremental.Name)
	}
	return &latestIncremental, nil
}

// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: load full zip snapshot, overlay incremental zip files by RelativePath, return merged DirectorySnapshot.
func ReconstructArchiveState(archiveDir string) (*DirectorySnapshot, error) {
	// Find most recent full archive
	latestFullArchive, err := findLatestFullArchive(archiveDir)
	if err != nil {
		return nil, fmt.Errorf("failed to find full archive: %w", err)
	}

	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - Found latest full archive: %s (mod time: %v)\n", latestFullArchive.Name, latestFullArchive.CreationTime)
	}

	// Load full archive snapshot
	fullSnapshot, err := CreateArchiveSnapshot(latestFullArchive.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to load full archive snapshot: %w", err)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - Loaded full archive snapshot: %s, %d files\n", latestFullArchive.Path, len(fullSnapshot.Files))
	}

	// Find most recent incremental archive that is based on the latest full archive
	latestIncremental, err := findLatestIncrementalArchive(archiveDir, latestFullArchive)
	if err != nil {
		return nil, fmt.Errorf("failed to find incremental archive: %w", err)
	}

	// If no incremental archive exists, return full archive snapshot
	if latestIncremental == nil {
		if debug {
			fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - No incremental archive found for base %s\n", latestFullArchive.Name)
		}
		return fullSnapshot, nil
	}

	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - Found latest incremental archive: %s\n", latestIncremental.Name)
	}

	// Load incremental archive snapshot
	incrementalSnapshot, err := CreateArchiveSnapshot(latestIncremental.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to load incremental archive snapshot: %w", err)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - Loaded incremental archive snapshot: %s, %d files\n", latestIncremental.Path, len(incrementalSnapshot.Files))
		fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - Full archive has %d files, incremental has %d files\n", len(fullSnapshot.Files), len(incrementalSnapshot.Files))
	}

	// Merge incremental changes on top of full archive
	// Create a map for efficient lookup
	fullMap := make(map[string]FileInfo)
	for _, file := range fullSnapshot.Files {
		fullMap[file.RelativePath] = file
	}

	// Apply incremental changes (incremental files override/add to full archive)
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - Applying %d incremental files to full archive\n", len(incrementalSnapshot.Files))
	}
	for _, file := range incrementalSnapshot.Files {
		if debug {
			if _, exists := fullMap[file.RelativePath]; exists {
				fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - Incremental overriding file: %s (size=%d, hash=%s)\n", file.RelativePath, file.Size, file.Hash[:16])
			} else {
				fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - Incremental adding new file: %s (size=%d, hash=%s)\n", file.RelativePath, file.Size, file.Hash[:16])
			}
		}
		fullMap[file.RelativePath] = file
	}

	// Convert map back to slice and sort
	reconstructedFiles := make([]FileInfo, 0, len(fullMap))
	for _, file := range fullMap {
		reconstructedFiles = append(reconstructedFiles, file)
	}

	// Sort by relative path for consistency
	sort.Slice(reconstructedFiles, func(i, j int) bool {
		return reconstructedFiles[i].RelativePath < reconstructedFiles[j].RelativePath
	})

	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: ReconstructArchiveState - Reconstructed state has %d files\n", len(reconstructedFiles))
	}

	return &DirectorySnapshot{Files: reconstructedFiles}, nil
}

// - [IMPL-DIFF_COMMAND] [ARCH-CLI_COMMANDS] [ARCH-DIFF_COMMAND] [ARCH-DIRECTORY_COMPARISON] [REQ-DIFF_COMMAND] — How: snapshot cwd, compare file maps by path/size/hash; classify added, modified, deleted (files only).
func CalculateDiff(cwd string, reconstructedState *DirectorySnapshot, excludePatterns []string) (*DiffResult, error) {
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: CalculateDiff - Reconstructed state has %d files\n", len(reconstructedState.Files))
	}
	// Create snapshot of current directory
	currentSnapshot, err := CreateDirectorySnapshot(cwd, excludePatterns)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory snapshot: %w", err)
	}
	if debug {
		fmt.Fprintf(os.Stderr, "DEBUG: CalculateDiff - Current directory has %d files\n", len(currentSnapshot.Files))
	}

	// Create maps for efficient lookup
	// Filter out directories from current snapshot since archives only contain files
	currentMap := make(map[string]FileInfo)
	for _, file := range currentSnapshot.Files {
		// Archives only contain files (directories are implicit), so we should only
		// compare files to avoid false positives for directories
		if !file.IsDir {
			currentMap[file.RelativePath] = file
		}
	}

	reconstructedMap := make(map[string]FileInfo)
	for _, file := range reconstructedState.Files {
		// Archive snapshots should only contain files, but filter just in case
		if !file.IsDir {
			reconstructedMap[file.RelativePath] = file
		}
	}

	// Find added and modified files
	var added []string
	var modified []string
	for path, currentFile := range currentMap {
		reconstructedFile, exists := reconstructedMap[path]
		if !exists {
			// File exists in current directory but not in reconstructed state
			if debug {
				fmt.Fprintf(os.Stderr, "DEBUG: CalculateDiff - File added: %s\n", path)
			}
			added = append(added, path)
		} else {
			// File exists in both - check if modified
			if currentFile.Size != reconstructedFile.Size ||
				currentFile.Hash != reconstructedFile.Hash {
				if debug {
					fmt.Fprintf(os.Stderr, "DEBUG: CalculateDiff - File modified: %s (current size=%d hash=%s mtime=%v, reconstructed size=%d hash=%s mtime=%v)\n",
						path, currentFile.Size, currentFile.Hash[:16], currentFile.ModTime, reconstructedFile.Size, reconstructedFile.Hash[:16], reconstructedFile.ModTime)
				}
				modified = append(modified, path)
			}
		}
	}

	// Find deleted files
	var deleted []string
	for path := range reconstructedMap {
		if _, exists := currentMap[path]; !exists {
			deleted = append(deleted, path)
		}
	}

	// Sort results for consistent output
	sort.Strings(added)
	sort.Strings(modified)
	sort.Strings(deleted)

	return &DiffResult{
		Added:    added,
		Modified: modified,
		Deleted:  deleted,
	}, nil
}
