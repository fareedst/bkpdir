// EXTRACT-010: Git-Aware Backup Tool Example - Integration of config + cli + git + fileops packages - 🔺
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"bkpdir/pkg/cli"
	"bkpdir/pkg/config"
	"bkpdir/pkg/fileops"
	"bkpdir/pkg/formatter"
	"bkpdir/pkg/git"
	"bkpdir/pkg/processing"
)

// GitBackupTool demonstrates Git-aware backup functionality
type GitBackupTool struct {
	configLoader *config.ConfigLoader
	gitRepo      git.RepositoryInterface
	fileOps      fileops.FileOperationsInterface
	formatter    formatter.FormatterInterface
	naming       processing.NamingProviderInterface
	appInfo      *cli.AppInfo
}

// NewGitBackupTool creates a new Git-aware backup tool
func NewGitBackupTool() *GitBackupTool {
	return &GitBackupTool{
		configLoader: config.NewConfigLoader(),
		gitRepo:      git.NewRepository(),
		fileOps:      fileops.NewFileOperations(),
		formatter:    formatter.NewTemplateFormatter(),
		naming:       processing.NewNamingProvider(),
		appInfo: &cli.AppInfo{
			Name:        "git-backup",
			Version:     "1.0.0",
			Description: "Git-aware backup tool using extracted packages",
		},
	}
}

// main demonstrates the integration
func main() {
	tool := NewGitBackupTool()

	rootCmd := &cobra.Command{
		Use:   tool.appInfo.Name,
		Short: tool.appInfo.Description,
		Long: `Git-Aware Backup Tool

Demonstrates integration of:
- pkg/config: Configuration management
- pkg/cli: CLI framework  
- pkg/git: Git repository integration
- pkg/fileops: File operations
- pkg/processing: Naming conventions
- pkg/formatter: Output formatting

Example:
  git-backup create --source=./myproject --note="before-refactor"`,
	}

	// Add create command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Git-aware backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateBackup(tool, cmd.Context())
		},
	}

	// Add list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List backups with Git information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListBackups(tool, cmd.Context())
		},
	}

	rootCmd.AddCommand(createCmd, listCmd)
	cli.AddVersionCommand(rootCmd, tool.appInfo)

	ctx := context.Background()
	if err := cli.ExecuteWithContext(rootCmd, &cli.CommandContext{Context: ctx}); err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
}

// runCreateBackup demonstrates Git-aware backup creation
func runCreateBackup(tool *GitBackupTool, ctx context.Context) error {
	fmt.Println("Creating Git-aware backup...")

	// 1. Load configuration (pkg/config)
	cfg, err := loadBackupConfig(tool.configLoader)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	sourcePath := cfg.GetString("source_path", ".")
	outputDir := cfg.GetString("output_directory", "./backups")

	// 2. Check if source is a Git repository (pkg/git)
	isRepo, err := tool.gitRepo.IsRepository(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to check Git repository: %w", err)
	}

	var gitInfo *git.RepositoryInfo
	if isRepo {
		gitInfo, err = tool.gitRepo.GetRepositoryInfo(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to get Git info: %w", err)
		}

		fmt.Printf("Git Repository Found:\n")
		fmt.Printf("  Branch: %s\n", gitInfo.Branch)
		fmt.Printf("  Hash: %s\n", gitInfo.Hash[:8])
		fmt.Printf("  Clean: %v\n", gitInfo.IsClean)
	}

	// 3. Generate Git-aware backup name (pkg/processing)
	backupName, err := generateGitBackupName(tool.naming, sourcePath, gitInfo)
	if err != nil {
		return fmt.Errorf("failed to generate backup name: %w", err)
	}

	backupPath := filepath.Join(outputDir, backupName)
	fmt.Printf("Backup will be created: %s\n", backupPath)

	// 4. Create backup directory (pkg/fileops)
	if err := tool.fileOps.CreateDirectory(outputDir); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// 5. List files to backup (pkg/fileops)
	files, err := tool.fileOps.ListFiles(sourcePath, fileops.ListOptions{
		IgnorePatterns: []string{".git", "node_modules", "*.tmp"},
		Recursive:      true,
	})
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	fmt.Printf("Found %d files to backup\n", len(files))

	// 6. Simulate backup creation with progress
	fmt.Println("Creating backup...")
	for i, file := range files {
		progress := float64(i+1) / float64(len(files)) * 100
		fmt.Printf("\rProgress: %.1f%% (%d/%d) %s", progress, i+1, len(files), filepath.Base(file))
		time.Sleep(10 * time.Millisecond) // Simulate work
	}
	fmt.Printf("\n")

	// 7. Create backup metadata
	metadata := map[string]interface{}{
		"backup_name": backupName,
		"source_path": sourcePath,
		"file_count":  len(files),
		"created_at":  time.Now().Format(time.RFC3339),
	}

	if gitInfo != nil {
		metadata["git"] = map[string]interface{}{
			"branch":   gitInfo.Branch,
			"hash":     gitInfo.Hash,
			"is_clean": gitInfo.IsClean,
		}
	}

	// 8. Format and display results (pkg/formatter)
	format := cfg.GetString("output_format", "table")
	switch format {
	case "json":
		output, _ := tool.formatter.FormatJSON(metadata)
		fmt.Println(output)
	case "yaml":
		output, _ := tool.formatter.FormatYAML(metadata)
		fmt.Println(output)
	default:
		fmt.Printf("\n✅ Backup created successfully!\n")
		fmt.Printf("Name: %s\n", backupName)
		fmt.Printf("Files: %d\n", len(files))
		if gitInfo != nil {
			fmt.Printf("Git: %s@%s (clean: %v)\n", gitInfo.Branch, gitInfo.Hash[:8], gitInfo.IsClean)
		}
	}

	return nil
}

// runListBackups demonstrates listing backups with Git information
func runListBackups(tool *GitBackupTool, ctx context.Context) error {
	fmt.Println("Listing Git-aware backups...")

	// Load configuration
	cfg, err := loadBackupConfig(tool.configLoader)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	backupDir := cfg.GetString("output_directory", "./backups")

	// List backup files
	files, err := tool.fileOps.ListFiles(backupDir, fileops.ListOptions{
		IncludePatterns: []string{"*.zip", "*.tar.gz"},
		Recursive:       false,
	})
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	if len(files) == 0 {
		fmt.Printf("No backups found in %s\n", backupDir)
		return nil
	}

	// Parse backup information from names
	var backups []map[string]interface{}
	for _, file := range files {
		backup := map[string]interface{}{
			"name": filepath.Base(file),
			"path": file,
		}

		// Try to parse Git information from filename
		if components, err := tool.naming.ParseName(filepath.Base(file), "archive"); err == nil {
			backup["timestamp"] = components.Timestamp.Format("2006-01-02 15:04:05")
			backup["git_branch"] = components.GitBranch
			backup["git_hash"] = components.GitHash
			backup["note"] = components.Note
		}

		// Get file size
		if info, err := tool.fileOps.GetFileInfo(file); err == nil {
			backup["size"] = fmt.Sprintf("%.1f MB", float64(info.Size())/1024/1024)
			backup["modified"] = info.ModTime().Format("2006-01-02 15:04:05")
		}

		backups = append(backups, backup)
	}

	// Format output
	format := cfg.GetString("output_format", "table")
	switch format {
	case "json":
		output, _ := tool.formatter.FormatJSON(backups)
		fmt.Println(output)
	case "yaml":
		output, _ := tool.formatter.FormatYAML(backups)
		fmt.Println(output)
	default:
		// Table format
		headers := []string{"Name", "Branch", "Hash", "Size", "Created"}
		var rows [][]string
		for _, backup := range backups {
			rows = append(rows, []string{
				backup["name"].(string),
				fmt.Sprintf("%v", backup["git_branch"]),
				fmt.Sprintf("%v", backup["git_hash"]),
				fmt.Sprintf("%v", backup["size"]),
				fmt.Sprintf("%v", backup["timestamp"]),
			})
		}

		output, _ := tool.formatter.FormatTable(headers, rows)
		fmt.Println(output)
	}

	return nil
}

// generateGitBackupName creates a backup name with Git information
func generateGitBackupName(naming processing.NamingProviderInterface, sourcePath string, gitInfo *git.RepositoryInfo) (string, error) {
	template := &processing.NamingTemplate{
		Prefix:          filepath.Base(sourcePath),
		Timestamp:       time.Now(),
		TimestampFormat: "2006-01-02T150405",
	}

	if gitInfo != nil {
		template.GitBranch = gitInfo.Branch
		template.GitHash = gitInfo.Hash[:8] // Short hash
		template.GitIsClean = gitInfo.IsClean
		template.ShowGitDirtyStatus = true
	}

	name, err := naming.GenerateName(template)
	if err != nil {
		return "", err
	}

	return name + ".zip", nil
}

// loadBackupConfig loads configuration for the backup tool
func loadBackupConfig(loader *config.ConfigLoader) (config.ConfigInterface, error) {
	sources := []config.ConfigSource{
		// Default configuration
		config.NewMapConfigSource("defaults", map[string]interface{}{
			"source_path":         ".",
			"output_directory":    "./backups",
			"output_format":       "table",
			"git.include_hash":    true,
			"git.max_hash_length": 8,
		}),

		// Environment variables
		config.NewEnvConfigSource("GITBACKUP"),

		// Configuration file
		config.NewFileConfigSource("backup-config.yaml"),
	}

	cfg, err := loader.LoadConfig(sources...)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return cfg, nil
}
