// This file is part of bkpdir
//
// Package main provides file statistics functionality for BkpDir.
// It handles gathering and formatting file statistics for output display.
//
// Copyright (c) 2024 BkpDir Contributors
// Licensed under the MIT License

// [REQ-OUTPUT_FORMATTING] File statistics and information gathering
// [ARCH-FILE_STATISTICS] File statistics gathering for enhanced output formatting
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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

// - [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: stat path and populate FileStatInfo with name, size, human size, mtime, mode, and type.
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

// - [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: scale bytes to TB/GB/MB/KB/B with one decimal for large units.
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

// - [IMPL-FILE_STATISTICS] [ARCH-FILE_STATISTICS] [REQ-OUT_002] [REQ-OUTPUT_FORMATTING] — How: classify regular, directory, symlink, device, pipe, socket, or other from mode bits.
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
