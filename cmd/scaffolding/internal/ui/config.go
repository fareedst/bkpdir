// 🔺 EXTRACT-008: Component selection interface - 🔍 Package discovery and configuration
package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

// ProjectConfig holds the configuration for generating a new CLI project
type ProjectConfig struct {
	ProjectName       string
	OutputPath        string
	ModuleName        string
	Author            string
	Description       string
	Template          string
	SelectedPackages  []string
	IncludeTests      bool
	IncludeDockerfile bool
	IncludeGitIgnore  bool
}

// AvailablePackages defines the packages that can be included in generated projects
var AvailablePackages = []Package{
	{Name: "config", Description: "Configuration management - schema-agnostic YAML/JSON loading", Required: true},
	{Name: "errors", Description: "Structured error handling with status codes", Required: true},
	{Name: "resources", Description: "Resource management with cleanup coordination", Required: false},
	{Name: "formatter", Description: "Template-based and printf-style output formatting", Required: false},
	{Name: "git", Description: "Git repository detection and information extraction", Required: false},
	{Name: "cli", Description: "CLI framework and command patterns", Required: true},
	{Name: "fileops", Description: "Safe file operations and comparisons", Required: false},
	{Name: "processing", Description: "Concurrent processing with worker pools", Required: false},
}

// Package represents an available package for inclusion
type Package struct {
	Name        string
	Description string
	Required    bool
}

// Templates defines the available project templates
var Templates = []Template{
	{Name: "minimal", Description: "Basic CLI with config, errors, and CLI framework", Packages: []string{"config", "errors", "cli"}},
	{Name: "standard", Description: "Adds formatting and Git integration", Packages: []string{"config", "errors", "cli", "formatter", "git"}},
	{Name: "advanced", Description: "Full feature set with file operations and processing", Packages: []string{"config", "errors", "resources", "formatter", "git", "cli", "fileops", "processing"}},
	{Name: "custom", Description: "Select your own package combination", Packages: []string{}},
}

// Template represents a project template
type Template struct {
	Name        string
	Description string
	Packages    []string
}

// CollectProjectConfig interactively collects project configuration from the user
func CollectProjectConfig() (*ProjectConfig, error) {
	config := &ProjectConfig{}

	// 🔺 EXTRACT-008: Project name and basic configuration - 🔍 User input validation
	if err := collectBasicInfo(config); err != nil {
		return nil, fmt.Errorf("collecting basic info: %w", err)
	}

	// 🔺 EXTRACT-008: Template selection interface - 🔍 Template and package selection
	if err := collectTemplateAndPackages(config); err != nil {
		return nil, fmt.Errorf("collecting template and packages: %w", err)
	}

	// 🔺 EXTRACT-008: Optional features configuration - 🔍 Additional options
	if err := collectOptionalFeatures(config); err != nil {
		return nil, fmt.Errorf("collecting optional features: %w", err)
	}

	return config, nil
}

func collectBasicInfo(config *ProjectConfig) error {
	// Project name
	projectPrompt := promptui.Prompt{
		Label:    "Project name",
		Default:  "my-cli-app",
		Validate: validateProjectName,
	}

	projectName, err := projectPrompt.Run()
	if err != nil {
		return err
	}
	config.ProjectName = projectName

	// Output path
	pwd, _ := os.Getwd()
	outputPrompt := promptui.Prompt{
		Label:    "Output directory",
		Default:  filepath.Join(pwd, projectName),
		Validate: validateOutputPath,
	}

	outputPath, err := outputPrompt.Run()
	if err != nil {
		return err
	}
	config.OutputPath = outputPath

	// Module name
	modulePrompt := promptui.Prompt{
		Label:   "Go module name",
		Default: fmt.Sprintf("github.com/your-username/%s", projectName),
	}

	moduleName, err := modulePrompt.Run()
	if err != nil {
		return err
	}
	config.ModuleName = moduleName

	// Author
	authorPrompt := promptui.Prompt{
		Label:   "Author name",
		Default: "Your Name",
	}

	author, err := authorPrompt.Run()
	if err != nil {
		return err
	}
	config.Author = author

	// Description
	descPrompt := promptui.Prompt{
		Label:   "Project description",
		Default: fmt.Sprintf("A CLI application built with extracted packages"),
	}

	description, err := descPrompt.Run()
	if err != nil {
		return err
	}
	config.Description = description

	return nil
}

func collectTemplateAndPackages(config *ProjectConfig) error {
	// Template selection
	templatePrompt := promptui.Select{
		Label: "Select project template",
		Items: Templates,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "▶ {{ .Name | cyan }} - {{ .Description }}",
			Inactive: "  {{ .Name | cyan }} - {{ .Description }}",
			Selected: "✓ Template: {{ .Name | green }}",
		},
	}

	templateIdx, _, err := templatePrompt.Run()
	if err != nil {
		return err
	}

	selectedTemplate := Templates[templateIdx]
	config.Template = selectedTemplate.Name

	// Package selection
	if selectedTemplate.Name == "custom" {
		return collectCustomPackages(config)
	} else {
		config.SelectedPackages = selectedTemplate.Packages
		fmt.Printf("📦 Selected packages: %s\n", strings.Join(selectedTemplate.Packages, ", "))
		return nil
	}
}

func collectCustomPackages(config *ProjectConfig) error {
	fmt.Println("\n📦 Select packages to include in your project:")

	var selectedPackages []string

	// Add required packages automatically
	for _, pkg := range AvailablePackages {
		if pkg.Required {
			selectedPackages = append(selectedPackages, pkg.Name)
		}
	}

	// Allow selection of optional packages
	for _, pkg := range AvailablePackages {
		if !pkg.Required {
			confirmPrompt := promptui.Prompt{
				Label:     fmt.Sprintf("Include %s (%s)", pkg.Name, pkg.Description),
				IsConfirm: true,
				Default:   "n",
			}

			result, err := confirmPrompt.Run()
			if err != nil && err != promptui.ErrAbort {
				return err
			}

			if result == "y" || result == "Y" {
				selectedPackages = append(selectedPackages, pkg.Name)
			}
		}
	}

	config.SelectedPackages = selectedPackages
	fmt.Printf("✓ Selected packages: %s\n", strings.Join(selectedPackages, ", "))
	return nil
}

func collectOptionalFeatures(config *ProjectConfig) error {
	// Include tests
	testsPrompt := promptui.Prompt{
		Label:     "Include example tests",
		IsConfirm: true,
		Default:   "y",
	}

	testsResult, err := testsPrompt.Run()
	if err != nil && err != promptui.ErrAbort {
		return err
	}
	config.IncludeTests = (testsResult == "y" || testsResult == "Y")

	// Include Dockerfile
	dockerPrompt := promptui.Prompt{
		Label:     "Include Dockerfile",
		IsConfirm: true,
		Default:   "n",
	}

	dockerResult, err := dockerPrompt.Run()
	if err != nil && err != promptui.ErrAbort {
		return err
	}
	config.IncludeDockerfile = (dockerResult == "y" || dockerResult == "Y")

	// Include .gitignore
	gitignorePrompt := promptui.Prompt{
		Label:     "Include .gitignore",
		IsConfirm: true,
		Default:   "y",
	}

	gitignoreResult, err := gitignorePrompt.Run()
	if err != nil && err != promptui.ErrAbort {
		return err
	}
	config.IncludeGitIgnore = (gitignoreResult == "y" || gitignoreResult == "Y")

	return nil
}

func validateProjectName(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check for valid Go package/directory name
	if matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_-]*$`, input); !matched {
		return fmt.Errorf("project name must start with a letter and contain only letters, numbers, hyphens, and underscores")
	}

	return nil
}

func validateOutputPath(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("output path cannot be empty")
	}

	// Check if path already exists
	if _, err := os.Stat(input); err == nil {
		return fmt.Errorf("path already exists: %s", input)
	}

	// Check if parent directory is writable
	parentDir := filepath.Dir(input)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		return fmt.Errorf("parent directory does not exist: %s", parentDir)
	}

	return nil
}
