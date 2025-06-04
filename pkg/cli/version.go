package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ⭐ EXTRACT-005: Version and build info handling extracted from main.go - 🔧

// DefaultVersionManager provides standard version handling functionality
type DefaultVersionManager struct{}

// NewVersionManager creates a new version manager
func NewVersionManager() VersionManager {
	return &DefaultVersionManager{}
}

// FormatVersion formats version information for display
func (vm *DefaultVersionManager) FormatVersion(info BuildInfo) string {
	return fmt.Sprintf("%s (compiled %s) [%s]", info.Version, info.Date, info.Platform)
}

// CreateVersionCommand creates a version subcommand
func (vm *DefaultVersionManager) CreateVersionCommand(info BuildInfo) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Long:  "Display detailed version information including build date and platform.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(vm.FormatVersion(info))
		},
	}
}

// CreateVersionTemplate creates a version template string for root command
func (vm *DefaultVersionManager) CreateVersionTemplate(info BuildInfo) string {
	return fmt.Sprintf("version %s (compiled %s) [%s]\n",
		info.Version, info.Date, info.Platform)
}

// FormatLongDescription creates a formatted long description with version info
func (vm *DefaultVersionManager) FormatLongDescription(info AppInfo, baseLongDesc string) string {
	return fmt.Sprintf("%s version %s (compiled %s) [%s]\n\n%s",
		info.Name, info.Build.Version, info.Build.Date, info.Build.Platform, baseLongDesc)
}

// SetVersionInfo sets version information on a command
func (vm *DefaultVersionManager) SetVersionInfo(cmd *cobra.Command, info BuildInfo) {
	cmd.Version = vm.FormatVersion(info)
	template := vm.CreateVersionTemplate(info)
	cmd.SetVersionTemplate(template)
}
