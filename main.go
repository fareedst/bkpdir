package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	dryRun bool
	note   string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "bkpdir",
		Short: "Directory archiving CLI for macOS and Linux",
	}

	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Show what would be done, but don't create archives")
	rootCmd.PersistentFlags().StringVar(&note, "note", "", "Optional note for the archive name")

	rootCmd.AddCommand(fullCmd())
	rootCmd.AddCommand(incCmd())
	rootCmd.AddCommand(listCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func fullCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "full",
		Short: "Create a full archive of the current directory",
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
			if err := CreateFullArchive(cfg, note, dryRun); err != nil {
				fmt.Println("Error creating full archive:", err)
				os.Exit(1)
			}
		},
	}
}

func incCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inc",
		Short: "Create an incremental archive of the current directory",
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
			if err := CreateIncrementalArchive(cfg, note, dryRun); err != nil {
				fmt.Println("Error creating incremental archive:", err)
				os.Exit(1)
			}
		},
	}
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all archives for the current directory",
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
				fmt.Println(a.Name)
			}
		},
	}
}
