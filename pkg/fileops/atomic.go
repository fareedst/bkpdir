// [IMPL-ATOMIC_OPS] [ARCH-RESOURCE_MANAGEMENT] [REQ-RESOURCE_MANAGEMENT]
// Package fileops provides file operations and utilities for CLI applications.
//
// This file contains atomic file operation patterns for safe file writing.
package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// AtomicWriter provides atomic file writing capabilities
// [IMPL-ATOMIC_OPS] AtomicWriter — temp-file lifecycle for a single atomic write target
type AtomicWriter struct {
	targetPath  string
	tempPath    string
	tempFile    *os.File
	isCommitted bool
	isClosed    bool
}

// AtomicOp defines the interface for atomic file operations
// [IMPL-ATOMIC_OPS] AtomicOp — contract for atomic write/commit/rollback
type AtomicOp interface {
	Write(data []byte) (int, error)
	WriteString(s string) (int, error)
	Commit() error
	Rollback() error
	Close() error
}

// NewAtomicWriter creates a new atomic writer for the target path
// [IMPL-ATOMIC_OPS] NewAtomicWriter — creates temp file in same dir for same-filesystem rename
func NewAtomicWriter(targetPath string) (*AtomicWriter, error) {
	// Validate the target path
	if err := ValidatePath(targetPath); err != nil {
		return nil, fmt.Errorf("invalid target path: %v", err)
	}

	// Create temporary file in the same directory as target
	dir := filepath.Dir(targetPath)
	base := filepath.Base(targetPath)

	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("cannot create directory %s: %v", dir, err)
	}

	tempFile, err := os.CreateTemp(dir, base+".tmp.*")
	if err != nil {
		return nil, fmt.Errorf("cannot create temporary file: %v", err)
	}

	return &AtomicWriter{
		targetPath:  targetPath,
		tempPath:    tempFile.Name(),
		tempFile:    tempFile,
		isCommitted: false,
		isClosed:    false,
	}, nil
}

// Write writes data to the temporary file
// [IMPL-ATOMIC_OPS] AtomicWriter.Write — writes bytes to temp file, guards closed state
func (aw *AtomicWriter) Write(data []byte) (int, error) {
	if aw.isClosed {
		return 0, fmt.Errorf("writer is closed")
	}
	if aw.tempFile == nil {
		return 0, fmt.Errorf("no temporary file available")
	}

	return aw.tempFile.Write(data)
}

// WriteString writes a string to the temporary file
// [IMPL-ATOMIC_OPS] AtomicWriter.WriteString — converts string to bytes and delegates
func (aw *AtomicWriter) WriteString(s string) (int, error) {
	return aw.Write([]byte(s))
}

// Commit atomically moves the temporary file to the target location
// [IMPL-ATOMIC_OPS] AtomicWriter.Commit — closes temp file then renames to target
func (aw *AtomicWriter) Commit() error {
	if aw.isCommitted {
		return fmt.Errorf("already committed")
	}
	if aw.isClosed && aw.tempFile != nil {
		return fmt.Errorf("writer is closed but not properly cleaned up")
	}

	// Close the temporary file first
	if aw.tempFile != nil {
		if err := aw.tempFile.Close(); err != nil {
			aw.cleanup()
			return fmt.Errorf("cannot close temporary file: %v", err)
		}
		aw.tempFile = nil
	}

	// Atomically move the file
	if err := os.Rename(aw.tempPath, aw.targetPath); err != nil {
		aw.cleanup()
		return fmt.Errorf("cannot commit file: %v", err)
	}

	aw.isCommitted = true
	aw.isClosed = true
	return nil
}

// Rollback removes the temporary file without committing
// [IMPL-ATOMIC_OPS] AtomicWriter.Rollback — removes temp file, prevents committed rollback
func (aw *AtomicWriter) Rollback() error {
	if aw.isCommitted {
		return fmt.Errorf("cannot rollback after commit")
	}

	return aw.cleanup()
}

// Close closes the writer and cleans up if not committed
// [IMPL-ATOMIC_OPS] AtomicWriter.Close — idempotent close, rolls back if not committed
func (aw *AtomicWriter) Close() error {
	if aw.isClosed {
		return nil
	}

	if !aw.isCommitted {
		return aw.Rollback()
	}

	aw.isClosed = true
	return nil
}

// cleanup removes the temporary file
// [IMPL-ATOMIC_OPS] AtomicWriter.cleanup — closes handle and removes temp from disk
func (aw *AtomicWriter) cleanup() error {
	var err error

	if aw.tempFile != nil {
		aw.tempFile.Close()
		aw.tempFile = nil
	}

	if aw.tempPath != "" {
		if removeErr := os.Remove(aw.tempPath); removeErr != nil && !os.IsNotExist(removeErr) {
			err = removeErr
		}
	}

	aw.isClosed = true
	return err
}

// AtomicCopy copies a file atomically from source to destination
// [IMPL-ATOMIC_OPS] AtomicCopy — streams src to AtomicWriter, sets perms, commits
func AtomicCopy(src, dst string) error {
	// Validate paths
	if err := ValidateReadable(src); err != nil {
		return fmt.Errorf("source file validation failed: %v", err)
	}

	if err := ValidatePath(dst); err != nil {
		return fmt.Errorf("destination path validation failed: %v", err)
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open source file: %v", err)
	}
	defer srcFile.Close()

	// Get source file info
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("cannot get source file info: %v", err)
	}

	// Create atomic writer for destination
	writer, err := NewAtomicWriter(dst)
	if err != nil {
		return fmt.Errorf("cannot create atomic writer: %v", err)
	}
	defer writer.Close()

	// Copy data
	if _, err := io.Copy(writer, srcFile); err != nil {
		return fmt.Errorf("copy failed: %v", err)
	}

	// Set permissions to match source
	if err := os.Chmod(writer.tempPath, srcInfo.Mode()); err != nil {
		return fmt.Errorf("cannot set permissions: %v", err)
	}

	// Commit the copy
	if err := writer.Commit(); err != nil {
		return fmt.Errorf("cannot commit copy: %v", err)
	}

	return nil
}

// AtomicWriteFile writes data to a file atomically
// [IMPL-ATOMIC_OPS] AtomicWriteFile — writes byte slice atomically with given permissions
func AtomicWriteFile(filename string, data []byte, perm os.FileMode) error {
	writer, err := NewAtomicWriter(filename)
	if err != nil {
		return err
	}
	defer writer.Close()

	if _, err := writer.Write(data); err != nil {
		return err
	}

	// Set permissions on temporary file
	if err := os.Chmod(writer.tempPath, perm); err != nil {
		return fmt.Errorf("cannot set permissions: %v", err)
	}

	return writer.Commit()
}

// AtomicWriteString writes a string to a file atomically
// [IMPL-ATOMIC_OPS] AtomicWriteString — converts string and delegates to AtomicWriteFile
func AtomicWriteString(filename, data string, perm os.FileMode) error {
	return AtomicWriteFile(filename, []byte(data), perm)
}
