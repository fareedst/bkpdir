package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Archive struct {
	Name          string
	Path          string
	CreationTime  time.Time
	IsIncremental bool
	GitBranch     string
	GitHash       string
	Note          string
	BaseArchive   string // for incremental
}

// GenerateArchiveName creates an archive name according to the spec.
func GenerateArchiveName(prefix, timestamp, gitBranch, gitHash, note string, isGit, isIncremental bool, baseName string) string {
	name := ""
	if prefix != "" {
		name += prefix + "-"
	}
	name += timestamp
	if isGit {
		name += "=" + gitBranch + "=" + gitHash
	}
	if isIncremental && baseName != "" {
		name = baseName + "_update=" + timestamp
		if isGit {
			name += "=" + gitBranch + "=" + gitHash
		}
	}
	if note != "" {
		name += "=" + note
	}
	name += ".zip"
	return name
}

// ListArchives lists all archives in the archive directory for the current source.
func ListArchives(archiveDir string) ([]Archive, error) {
	var archives []Archive
	dirEntries, err := os.ReadDir(archiveDir)
	if err != nil {
		return nil, err
	}
	for _, entry := range dirEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".zip") {
			continue
		}
		archives = append(archives, Archive{
			Name: entry.Name(),
			Path: filepath.Join(archiveDir, entry.Name()),
		})
	}
	return archives, nil
}

// Stub: CreateFullArchive creates a full archive (to be implemented)
func CreateFullArchive(cfg *Config, note string, dryRun bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Determine archive directory
	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}
	if !dryRun {
		if err := os.MkdirAll(archiveDir, 0o755); err != nil {
			return err
		}
	}

	// Walk directory and collect files
	var files []string
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if ShouldExcludeFile(rel, cfg.ExcludePatterns) {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return err
	}

	// Git info
	isGit := IsGitRepository(cwd)
	gitBranch, gitHash := "", ""
	if isGit {
		gitBranch = GetGitBranch(cwd)
		gitHash = GetGitShortHash(cwd)
	}

	// Archive name
	timestamp := time.Now().Format("2006-01-02-15-04")
	prefix := ""
	if !cfg.UseCurrentDirName {
		prefix = filepath.Base(cwd)
	}
	archiveName := GenerateArchiveName(prefix, timestamp, gitBranch, gitHash, note, isGit, false, "")
	archivePath := filepath.Join(archiveDir, archiveName)

	if dryRun {
		fmt.Println("[Dry Run] Files to include:")
		for _, f := range files {
			fmt.Println("  ", f)
		}
		fmt.Println("[Dry Run] Archive would be:", archivePath)
		return nil
	}

	// Create zip archive
	f, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	zipw := zip.NewWriter(f)
	defer zipw.Close()

	for _, rel := range files {
		abs := filepath.Join(cwd, rel)
		info, err := os.Lstat(abs)
		if err != nil {
			return err
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = rel
		hdr.Method = zip.Deflate
		w, err := zipw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rf, err := os.Open(abs)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, rf)
			rf.Close()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Created archive:", archivePath)
	return nil
}

// Stub: CreateIncrementalArchive creates an incremental archive (to be implemented)
func CreateIncrementalArchive(cfg *Config, note string, dryRun bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Determine archive directory
	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}
	if !dryRun {
		if err := os.MkdirAll(archiveDir, 0o755); err != nil {
			return err
		}
	}

	// Find latest full archive
	archives, err := ListArchives(archiveDir)
	if err != nil {
		return err
	}
	var lastFull Archive
	var lastFullTime time.Time
	for _, a := range archives {
		if !strings.Contains(a.Name, "_update=") {
			// Parse timestamp from name
			parts := strings.Split(a.Name, "=")
			base := parts[0]
			base = strings.TrimSuffix(base, ".zip")
			base = strings.TrimPrefix(base, filepath.Base(archiveDir)+"-")
			// Try to parse time
			ts := base
			if dash := strings.LastIndex(ts, "-"); dash != -1 {
				// Try to get YYYY-MM-DD-hh-mm
				if len(ts) >= 16 {
					ts = ts[len(ts)-16:]
				}
			}
			t, err := time.Parse("2006-01-02-15-04", ts)
			if err == nil && t.After(lastFullTime) {
				lastFull = a
				lastFullTime = t
			}
		}
	}
	if lastFull.Name == "" {
		return fmt.Errorf("no full archive found; create a full archive first")
	}

	// Get mod time of last full archive
	lastFullInfo, err := os.Stat(lastFull.Path)
	if err != nil {
		return err
	}
	lastFullMod := lastFullInfo.ModTime()

	// Walk directory and collect changed files
	var changed []string
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if ShouldExcludeFile(rel, cfg.ExcludePatterns) {
			return nil
		}
		if info.ModTime().After(lastFullMod) {
			changed = append(changed, rel)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Git info
	isGit := IsGitRepository(cwd)
	gitBranch, gitHash := "", ""
	if isGit {
		gitBranch = GetGitBranch(cwd)
		gitHash = GetGitShortHash(cwd)
	}

	// Archive name
	timestamp := time.Now().Format("2006-01-02-15-04")
	archiveName := GenerateArchiveName("", timestamp, gitBranch, gitHash, note, isGit, true, strings.TrimSuffix(lastFull.Name, ".zip"))
	archivePath := filepath.Join(archiveDir, archiveName)

	if dryRun {
		fmt.Println("[Dry Run] Files changed since last full archive:")
		for _, f := range changed {
			fmt.Println("  ", f)
		}
		fmt.Println("[Dry Run] Incremental archive would be:", archivePath)
		return nil
	}

	// Create zip archive
	f, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	zipw := zip.NewWriter(f)
	defer zipw.Close()

	for _, rel := range changed {
		abs := filepath.Join(cwd, rel)
		info, err := os.Lstat(abs)
		if err != nil {
			return err
		}
		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		hdr.Name = rel
		hdr.Method = zip.Deflate
		w, err := zipw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rf, err := os.Open(abs)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, rf)
			rf.Close()
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Created incremental archive:", archivePath)
	return nil
}
