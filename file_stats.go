// This file is part of bkpdir
//
// Package main provides file statistics functionality for BkpDir.
// It handles gathering and formatting file statistics for output display.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [REQ-OUTPUT_FORMATTING] File statistics and information gathering
// [ARCH-FILE_STATISTICS] File statistics gathering for enhanced output formatting
// [IMPL-FILE_STATISTICS] FileStatInfo structure with human-readable formatting
// STATS-FEATURES-001: Statistics features specification - File statistics and information gathering [ACTION:format-processing]
// Source: file_stats.go - STATS-FEATURES-001
// Impact: Core functionality requirement for statistics features

// SERVICE-STATS-001: Statistics service architecture decision - Statistics service implementation [ACTION:core-functionality]
// Source: file_stats.go - SERVICE-STATS-001
// Impact: Statistics service implementation decision
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUTPUT_FORMATTING]
// FileStatInfo contains statistical information about a file for output formatting
type FileStatInfo struct {
	Path      string      // Full file path
	Name      string      // File name only
	Size      int64       // File size in bytes
	SizeHuman string      // Human-readable size (1.2MB, 455KB, etc.)
	MTime     time.Time   // Modification time
	MTimeUnix int64       // Modification time as unix timestamp
	Mode      os.FileMode // File permissions/mode
	Type      string      // File type (regular, directory, symlink)
}

// [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUTPUT_FORMATTING]
// GatherFileStatInfo gathers statistical information about a file for output formatting.
// It returns a FileStatInfo struct containing path, size, modification time, permissions, and type.
func GatherFileStatInfo(path string) (*FileStatInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	return &FileStatInfo{
		Path:      path,
		Name:      filepath.Base(path),
		Size:      info.Size(),
		SizeHuman: formatHumanSize(info.Size()),
		MTime:     info.ModTime(),
		MTimeUnix: info.ModTime().Unix(),
		Mode:      info.Mode(),
		Type:      getFileType(info),
	}, nil
}

// [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUTPUT_FORMATTING]
// formatHumanSize converts a size in bytes to a human-readable format.
// Examples: 1024 -> "1.0KB", 1536 -> "1.5KB", 1048576 -> "1.0MB"
func formatHumanSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.1fTB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.1fGB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.1fMB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.1fKB", float64(size)/KB)
	default:
		return fmt.Sprintf("%dB", size)
	}
}

// [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUTPUT_FORMATTING]
// getFileType determines the type of a file based on its FileInfo.
// Returns "regular", "directory", "symlink", "device", "pipe", "socket", or "other".
func getFileType(info os.FileInfo) string {
	mode := info.Mode()
	switch {
	case mode.IsRegular():
		return "regular"
	case mode.IsDir():
		return "directory"
	case mode&os.ModeSymlink != 0:
		return "symlink"
	case mode&os.ModeDevice != 0:
		return "device"
	case mode&os.ModeNamedPipe != 0:
		return "pipe"
	case mode&os.ModeSocket != 0:
		return "socket"
	default:
		return "other"
	}
}
