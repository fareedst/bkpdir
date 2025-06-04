package cli

import (
	"github.com/spf13/cobra"
)

// ⭐ EXTRACT-005: Cobra command patterns and flag handling extracted from main.go - 🔧

// DefaultFlagManager provides standard flag management functionality
type DefaultFlagManager struct{}

// NewFlagManager creates a new flag manager
func NewFlagManager() FlagManager {
	return &DefaultFlagManager{}
}

// AddGlobalFlags adds common global flags to a command
func (fm *DefaultFlagManager) AddGlobalFlags(cmd *cobra.Command) error {
	// Add commonly used global flags
	cmd.PersistentFlags().BoolP("help", "h", false, "Help for this command")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
	return nil
}

// AddDryRunFlag adds dry-run flag with consistent naming
func (fm *DefaultFlagManager) AddDryRunFlag(cmd *cobra.Command, target *bool) error {
	cmd.PersistentFlags().BoolVarP(target, "dry-run", "d", false,
		"Show what would be done without executing")
	return nil
}

// AddNoteFlag adds note flag for operations
func (fm *DefaultFlagManager) AddNoteFlag(cmd *cobra.Command, target *string) error {
	cmd.Flags().StringVarP(target, "note", "n", "",
		"Optional note to include with the operation")
	return nil
}

// AddConfigFlag adds configuration flag
func (fm *DefaultFlagManager) AddConfigFlag(cmd *cobra.Command, target *bool) error {
	cmd.Flags().BoolVar(target, "config", false,
		"Display configuration values and exit")
	return nil
}

// AddVerifyFlag adds verification flag for operations
func (fm *DefaultFlagManager) AddVerifyFlag(cmd *cobra.Command, target *bool) error {
	cmd.Flags().BoolVarP(target, "verify", "v", false,
		"Verify operation results")
	return nil
}

// AddChecksumFlag adds checksum flag for verification
func (fm *DefaultFlagManager) AddChecksumFlag(cmd *cobra.Command, target *bool) error {
	cmd.Flags().BoolVarP(target, "checksum", "c", false,
		"Perform checksum verification")
	return nil
}

// AddIncrementalFlag adds incremental flag for operations
func (fm *DefaultFlagManager) AddIncrementalFlag(cmd *cobra.Command, target *bool) error {
	cmd.Flags().BoolVarP(target, "incremental", "i", false,
		"Perform incremental operation")
	return nil
}

// AddListFlag adds list flag for displaying items
func (fm *DefaultFlagManager) AddListFlag(cmd *cobra.Command, target *string) error {
	cmd.Flags().StringVar(target, "list", "",
		"List items for the specified path")
	return nil
}

// FlagSet represents a set of flags to add to a command
type FlagSet struct {
	DryRun      *bool
	Note        *string
	Config      *bool
	Verify      *bool
	Checksum    *bool
	Incremental *bool
	List        *string
	Verbose     *bool
}

// AddFlags adds all configured flags to the command
func (fm *DefaultFlagManager) AddFlags(cmd *cobra.Command, flagSet FlagSet) error {
	if flagSet.DryRun != nil {
		if err := fm.AddDryRunFlag(cmd, flagSet.DryRun); err != nil {
			return err
		}
	}

	if flagSet.Note != nil {
		if err := fm.AddNoteFlag(cmd, flagSet.Note); err != nil {
			return err
		}
	}

	if flagSet.Config != nil {
		if err := fm.AddConfigFlag(cmd, flagSet.Config); err != nil {
			return err
		}
	}

	if flagSet.Verify != nil {
		if err := fm.AddVerifyFlag(cmd, flagSet.Verify); err != nil {
			return err
		}
	}

	if flagSet.Checksum != nil {
		if err := fm.AddChecksumFlag(cmd, flagSet.Checksum); err != nil {
			return err
		}
	}

	if flagSet.Incremental != nil {
		if err := fm.AddIncrementalFlag(cmd, flagSet.Incremental); err != nil {
			return err
		}
	}

	if flagSet.List != nil {
		if err := fm.AddListFlag(cmd, flagSet.List); err != nil {
			return err
		}
	}

	return nil
}
