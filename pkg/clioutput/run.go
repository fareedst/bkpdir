package clioutput

import (
	"os"

	"github.com/spf13/cobra"
)

// ExecuteRoot peels texture flags, applies policy, and runs the Cobra root command.
func ExecuteRoot(root *cobra.Command, args []string) {
	peeled, flags := PeelTextureArgs(args)
	policy := ResolvePolicy(flags)
	ApplyToRoot(root, policy)
	root.SetArgs(peeled)
	if err := root.Execute(); err != nil {
		PrintExecError(err)
		os.Exit(1)
	}
}
