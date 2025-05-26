package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	version = "1.0.0"
)

var (
	compileDate = "2024-03-20" // This is a placeholder - actual value is set during build via -ldflags
	platform    = "unknown"    // This is a placeholder - actual value is set during build via -ldflags
)

var (
	dryRun     bool
	verify     bool
	checksum   bool
	note       string
	showConfig bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "bkpdir",
		Short:   "Directory archiving CLI for macOS and Linux",
		Long:    fmt.Sprintf("bkpdir version %s (compiled %s) [%s]\n\nBkpDir is a command-line tool for archiving directories on macOS and Linux. It supports full and incremental backups, customizable exclusion patterns, Git-aware archive naming, and archive verification.", version, compileDate, platform),
		Version: fmt.Sprintf("%s (compiled %s) [%s]", version, compileDate, platform),
		Example: `  # Create a full backup
  bkpdir full -n "Initial backup"

  # Create an incremental backup with verification
  bkpdir inc -n "Changes after feature X" -v

  # List all backups
  bkpdir list

  # Verify a specific backup with checksums
  bkpdir verify backup-2024-03-20.zip -c

  # Show configuration
  bkpdir --config`,
		Run: func(cmd *cobra.Command, args []string) {
			// Handle --config flag when no subcommand is provided
			if showConfig {
				handleConfigCommand()
				return
			}
			// If no config flag and no subcommand, show help
			cmd.Help()
		},
	}

	// Set the version template to show version in help output
	rootCmd.SetVersionTemplate(fmt.Sprintf("bkpdir version %s (compiled %s) [%s]\n", version, compileDate, platform))

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be done without creating archives")
	rootCmd.PersistentFlags().BoolVar(&showConfig, "config", false, "Display configuration values and exit")

	rootCmd.AddCommand(fullCmd())
	rootCmd.AddCommand(incCmd())
	rootCmd.AddCommand(listCmd())
	rootCmd.AddCommand(verifyCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleConfigCommand() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	cfg, err := LoadConfig(cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(cfg.StatusConfigError)
	}

	formatter := NewOutputFormatter(cfg)
	configValues := GetConfigValues(cfg)

	for _, cv := range configValues {
		formatter.PrintConfigValue(cv.Name, cv.Value, cv.Source)
	}
}

func fullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "full [NOTE]",
		Short: "Create a full archive of the current directory",
		Long: `Create a complete archive of the current directory. The archive will include all files
except those matching patterns specified in .bkpdir.yml. If the directory is a Git repository,
the archive name will include the current branch and commit hash.

Before creating an archive, the command compares the directory with its most recent archive.
If the directory is identical to the most recent archive, no new archive is created.`,
		Example: `  # Create a full backup
  bkpdir full

  # Create a full backup with a note
  bkpdir full -n "Initial backup"

  # Create a full backup with verification
  bkpdir full -v

  # Show what would be backed up without creating archive
  bkpdir full -d`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}

			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(cfg.StatusConfigError)
			}

			formatter := NewOutputFormatter(cfg)

			// Use note from flag if provided, otherwise use positional argument
			backupNote := note
			if backupNote == "" && len(args) > 0 {
				backupNote = args[0]
			}

			if err := CreateFullArchiveEnhanced(ctx, cfg, formatter, backupNote, dryRun, verify); err != nil {
				exitCode := HandleArchiveError(err, cfg, formatter)
				os.Exit(exitCode)
			}
		},
	}
	cmd.Flags().StringVarP(&note, "note", "n", "", "Add a note to the archive name")
	cmd.Flags().BoolVarP(&verify, "verify", "v", false, "Verify the archive after creation")
	return cmd
}

func incCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inc [NOTE]",
		Short: "Create an incremental archive of the current directory",
		Long: `Create an incremental archive containing only files changed since the last full archive.
This command requires at least one full archive to exist. The archive will include only files
that have been modified since the last full backup.`,
		Example: `  # Create an incremental backup
  bkpdir inc

  # Create an incremental backup with a note
  bkpdir inc -n "Changes after feature X"

  # Create an incremental backup with verification
  bkpdir inc -v

  # Show what would be backed up without creating archive
  bkpdir inc -d`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}

			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(cfg.StatusConfigError)
			}

			formatter := NewOutputFormatter(cfg)

			// Use note from flag if provided, otherwise use positional argument
			backupNote := note
			if backupNote == "" && len(args) > 0 {
				backupNote = args[0]
			}

			if err := CreateIncrementalArchiveEnhanced(ctx, cfg, formatter, backupNote, dryRun, verify); err != nil {
				exitCode := HandleArchiveError(err, cfg, formatter)
				os.Exit(exitCode)
			}
		},
	}
	cmd.Flags().StringVarP(&note, "note", "n", "", "Add a note to the archive name")
	cmd.Flags().BoolVarP(&verify, "verify", "v", false, "Verify the archive after creation")
	return cmd
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all archives for the current directory",
		Long: `Display all archives associated with the current directory. The output includes
the archive name and its verification status (VERIFIED, UNVERIFIED, or FAILED).`,
		Example: `  # List all backups
  bkpdir list`,
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}

			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(cfg.StatusConfigError)
			}

			formatter := NewOutputFormatter(cfg)

			if err := ListArchivesEnhanced(cfg, formatter); err != nil {
				exitCode := HandleArchiveError(err, cfg, formatter)
				os.Exit(exitCode)
			}
		},
	}
}

func verifyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "verify [ARCHIVE_NAME]",
		Short: "Verify archive integrity",
		Long: `Verify the integrity of archives. If no archive name is provided, verifies all archives.
Use the --checksum flag to include content verification.`,
		Example: `  # Verify all archives
  bkpdir verify

  # Verify a specific archive
  bkpdir verify myproject-2024-03-20-15-30.zip

  # Verify with checksum validation
  bkpdir verify --checksum myproject-2024-03-20-15-30.zip`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
				os.Exit(1)
			}

			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				os.Exit(cfg.StatusConfigError)
			}

			formatter := NewOutputFormatter(cfg)

			var archiveName string
			if len(args) > 0 {
				archiveName = args[0]
			}

			if err := VerifyArchiveEnhanced(cfg, formatter, archiveName, checksum); err != nil {
				exitCode := HandleArchiveError(err, cfg, formatter)
				os.Exit(exitCode)
			}
		},
	}
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("bkpdir version %s (compiled %s) [%s]\n", version, compileDate, platform)
		},
	}
}

// Enhanced archive creation with directory comparison
func CreateFullArchiveEnhanced(ctx context.Context, cfg *Config, formatter *OutputFormatter, note string, dryRun bool, verify bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return NewArchiveErrorWithCause("Failed to get current directory", cfg.StatusDirectoryNotFound, err)
	}

	// Validate directory
	if err := ValidateDirectoryPath(cwd, cfg); err != nil {
		return err
	}

	// Determine archive directory
	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}

	// Check for identical archive before creating
	identical, existingArchive, err := CheckForIdenticalArchive(cwd, archiveDir, cfg.ExcludePatterns)
	if err != nil {
		// If we can't check, continue with creation (might be first archive)
		fmt.Printf("Debug: Error checking for identical archive: %v\n", err)
	} else if identical {
		// Directory is identical to existing archive
		fmt.Printf("Debug: Directory is identical to existing archive: %s\n", existingArchive)
		formatter.PrintIdenticalArchive(existingArchive)
		os.Exit(cfg.StatusDirectoryIsIdenticalToExistingArchive)
	} else {
		fmt.Printf("Debug: Directory is not identical to existing archive: %s\n", existingArchive)
	}

	// Create resource manager for cleanup
	rm := NewResourceManager()
	defer rm.CleanupWithPanicRecovery()

	// Create archive directory
	if !dryRun {
		if err := SafeMkdirAll(archiveDir, 0755, cfg); err != nil {
			return err
		}
	}

	// Use the original CreateFullArchive function but with enhanced error handling
	if err := CreateFullArchive(cfg, note, dryRun, verify); err != nil {
		return err
	}

	return nil
}

// Enhanced incremental archive creation
func CreateIncrementalArchiveEnhanced(ctx context.Context, cfg *Config, formatter *OutputFormatter, note string, dryRun bool, verify bool) error {
	// Use the original CreateIncrementalArchive function but with enhanced error handling
	if err := CreateIncrementalArchive(cfg, note, dryRun, verify); err != nil {
		return err
	}
	return nil
}

// Enhanced archive listing with formatting
func ListArchivesEnhanced(cfg *Config, formatter *OutputFormatter) error {
	cwd, err := os.Getwd()
	if err != nil {
		return NewArchiveErrorWithCause("Failed to get current directory", cfg.StatusDirectoryNotFound, err)
	}

	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}

	archives, err := ListArchives(archiveDir)
	if err != nil {
		return NewArchiveErrorWithCause("Failed to list archives", 1, err)
	}

	if len(archives) == 0 {
		fmt.Printf("No archives found in %s\n", archiveDir)
		return nil
	}

	for _, a := range archives {
		status := ""
		if a.VerificationStatus != nil {
			if a.VerificationStatus.IsVerified {
				status = " [VERIFIED]"
			} else {
				status = " [FAILED]"
			}
		} else {
			status = " [UNVERIFIED]"
		}

		// Use enhanced formatting with extraction if possible
		creationTime := a.CreationTime.Format("2006-01-02 15:04:05")
		output := formatter.FormatListArchiveWithExtraction(a.Name, creationTime)
		fmt.Print(output + status + "\n")
	}

	return nil
}

// Enhanced archive verification
func VerifyArchiveEnhanced(cfg *Config, formatter *OutputFormatter, archiveName string, withChecksum bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return NewArchiveErrorWithCause("Failed to get current directory", cfg.StatusDirectoryNotFound, err)
	}

	archiveDir := cfg.ArchiveDirPath
	if cfg.UseCurrentDirName {
		archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
	}

	if archiveName != "" {
		// Verify specific archive
		archivePath := filepath.Join(archiveDir, archiveName)
		archive := &Archive{
			Name: archiveName,
			Path: archivePath,
		}

		var status *VerificationStatus
		if withChecksum {
			status, err = VerifyChecksums(archivePath)
			if err != nil {
				return NewArchiveErrorWithCause("Archive checksum verification failed", 1, err)
			}
		} else {
			status, err = VerifyArchive(archivePath)
			if err != nil {
				return NewArchiveErrorWithCause("Archive verification failed", 1, err)
			}
		}

		// Store verification status
		if err := StoreVerificationStatus(archive, status); err != nil {
			// Don't fail if we can't store status, just warn
			fmt.Printf("Warning: Could not store verification status: %v\n", err)
		}

		if status.IsVerified {
			fmt.Printf("Archive %s verified successfully\n", archiveName)
		} else {
			fmt.Printf("Archive %s verification failed:\n", archiveName)
			for _, errMsg := range status.Errors {
				fmt.Printf("  - %s\n", errMsg)
			}
			return NewArchiveError("Archive verification failed", 1)
		}
	} else {
		// Verify all archives
		archives, err := ListArchives(archiveDir)
		if err != nil {
			return NewArchiveErrorWithCause("Failed to list archives", 1, err)
		}

		allPassed := true
		for _, archive := range archives {
			var status *VerificationStatus
			if withChecksum {
				status, err = VerifyChecksums(archive.Path)
				if err != nil {
					fmt.Printf("Archive %s checksum verification failed: %v\n", archive.Name, err)
					allPassed = false
					continue
				}
			} else {
				status, err = VerifyArchive(archive.Path)
				if err != nil {
					fmt.Printf("Archive %s verification failed: %v\n", archive.Name, err)
					allPassed = false
					continue
				}
			}

			// Store verification status
			if err := StoreVerificationStatus(&archive, status); err != nil {
				// Don't fail if we can't store status, just warn
				fmt.Printf("Warning: Could not store verification status for %s: %v\n", archive.Name, err)
			}

			if status.IsVerified {
				fmt.Printf("Archive %s verified successfully\n", archive.Name)
			} else {
				fmt.Printf("Archive %s verification failed:\n", archive.Name)
				for _, errMsg := range status.Errors {
					fmt.Printf("  - %s\n", errMsg)
				}
				allPassed = false
			}
		}

		if !allPassed {
			return NewArchiveError("Some archives failed verification", 1)
		}
	}

	return nil
}
