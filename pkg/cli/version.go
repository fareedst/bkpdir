package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// DefaultVersionManager provides standard version handling functionality
type DefaultVersionManager struct{}

// NewVersionManager creates a new version manager
// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: construct DefaultCommandBuilder with a non-nil FlagManager (defaulting when nil).
func NewVersionManager() VersionManager {
	return &DefaultVersionManager{}
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: format version, build date, and platform into display string.
func (vm *DefaultVersionManager) FormatVersion(info BuildInfo) string {
	return fmt.Sprintf("%s (compiled %s) [%s]", info.Version, info.Date, info.Platform)
}

// - [IMPL-CLI_FRAMEWORK] [ARCH-CLI_FRAMEWORK] [REQ-USABILITY] [REQ-IMMUTABLE_CLI_COMMANDS] — How: create version subcommand that prints FormatVersion on run.
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
