package main

import (
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
)

var (
	dryRun   bool
	verify   bool
	checksum bool
	note     string
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "bkpdir",
		Short:   "Directory archiving CLI for macOS and Linux",
		Long:    fmt.Sprintf("bkpdir version %s (compiled %s)\n\nBkpDir is a command-line tool for archiving directories on macOS and Linux. It supports full and incremental backups, customizable exclusion patterns, Git-aware archive naming, and archive verification.", version, compileDate),
		Version: fmt.Sprintf("%s (compiled %s)", version, compileDate),
		Example: `  # Create a full backup
  bkpdir full -n "Initial backup"

  # Create an incremental backup with verification
  bkpdir inc -n "Changes after feature X" -v

  # List all backups
  bkpdir list

  # Verify a specific backup with checksums
  bkpdir verify backup-2024-03-20.zip -c`,
	}

	// Set the version template to show version in help output
	rootCmd.SetVersionTemplate(fmt.Sprintf("bkpdir version %s (compiled %s)\n", version, compileDate))

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "Show what would be done without creating archives")

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

func fullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "full [NOTE]",
		Short: "Create a full archive of the current directory",
		Long: `Create a complete archive of the current directory. The archive will include all files
except those matching patterns specified in .bkpdir.yml. If the directory is a Git repository,
the archive name will include the current branch and commit hash.`,
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
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				os.Exit(1)
			}
			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Println("Error loading config:", err)
				os.Exit(1)
			}
			// Use note from flag if provided, otherwise use positional argument
			backupNote := note
			if backupNote == "" && len(args) > 0 {
				backupNote = args[0]
			}
			if err := CreateFullArchive(cfg, backupNote, dryRun, verify); err != nil {
				fmt.Println("Error creating full archive:", err)
				os.Exit(1)
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
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				os.Exit(1)
			}
			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Println("Error loading config:", err)
				os.Exit(1)
			}
			// Use note from flag if provided, otherwise use positional argument
			backupNote := note
			if backupNote == "" && len(args) > 0 {
				backupNote = args[0]
			}
			if err := CreateIncrementalArchive(cfg, backupNote, dryRun, verify); err != nil {
				fmt.Println("Error creating incremental archive:", err)
				os.Exit(1)
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
				fmt.Println("Error getting current directory:", err)
				os.Exit(1)
			}
			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Println("Error loading config:", err)
				os.Exit(1)
			}
			archiveDir := cfg.ArchiveDirPath
			if cfg.UseCurrentDirName {
				archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
			}
			archives, err := ListArchives(archiveDir)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if len(archives) == 0 {
				fmt.Println("No archives found in", archiveDir)
				return
			}
			for _, a := range archives {
				status := ""
				if a.VerificationStatus != nil {
					if a.VerificationStatus.IsVerified {
						status = "[VERIFIED]"
					} else {
						status = "[FAILED]"
					}
				} else {
					status = "[UNVERIFIED]"
				}
				fmt.Println(a.Name, status)
			}
		},
	}
}

func verifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify ARCHIVE_NAME",
		Short: "Verify the integrity of an archive",
		Long: `Verify the integrity of a specific archive. This command performs several checks:
1. Validates the ZIP file structure
2. Checks for corrupted or incomplete archives
3. Verifies file headers and compression
4. Optionally verifies file contents against stored checksums`,
		Example: `  # Verify an archive
  bkpdir verify backup-2024-03-20.zip

  # Verify an archive with checksums
  bkpdir verify backup-2024-03-20.zip -c`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				os.Exit(1)
			}
			cfg, err := LoadConfig(cwd)
			if err != nil {
				fmt.Println("Error loading config:", err)
				os.Exit(1)
			}

			archiveName := args[0]
			archiveDir := cfg.ArchiveDirPath
			if cfg.UseCurrentDirName {
				archiveDir = filepath.Join(archiveDir, filepath.Base(cwd))
			}

			archivePath := filepath.Join(archiveDir, archiveName)
			if _, err := os.Stat(archivePath); os.IsNotExist(err) {
				fmt.Println("Error: Archive not found:", archivePath)
				os.Exit(1)
			}

			archive := &Archive{
				Name: archiveName,
				Path: archivePath,
			}

			// Verify archive structure
			status, err := VerifyArchive(archivePath)
			if err != nil {
				fmt.Println("Error verifying archive:", err)
				os.Exit(1)
			}

			// Verify checksums if requested
			if checksum {
				checksumStatus, err := VerifyChecksums(archivePath)
				if err != nil {
					fmt.Println("Error verifying checksums:", err)
					os.Exit(1)
				}

				// Merge statuses
				status.HasChecksums = checksumStatus.HasChecksums
				if !checksumStatus.IsVerified {
					status.IsVerified = false
					status.Errors = append(status.Errors, checksumStatus.Errors...)
				}
			}

			// Store verification status
			if err := StoreVerificationStatus(archive, status); err != nil {
				fmt.Println("Error storing verification status:", err)
				os.Exit(1)
			}

			// Display results
			if status.IsVerified {
				fmt.Println("Archive verification successful")
				if status.HasChecksums {
					fmt.Println("Checksum verification successful")
				}
				os.Exit(0)
			} else {
				fmt.Println("Archive verification failed:")
				for _, err := range status.Errors {
					fmt.Println("  -", err)
				}
				os.Exit(1)
			}
		},
	}
	cmd.Flags().BoolVarP(&checksum, "checksum", "c", false, "Verify file contents against stored checksums")
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("bkpdir version %s (compiled %s)\n", version, compileDate)
		},
	}
}
