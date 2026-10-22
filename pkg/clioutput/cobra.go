package clioutput

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// RegisterPersistentFlags documents global operator-output controls on the root command.
func RegisterPersistentFlags(cmd *cobra.Command) {
	// No -p shorthand on persistent flags — subcommands may use -p (e.g. server --port).
	cmd.PersistentFlags().Bool("plain", false, "Plain help text (no ANSI decoration)")
	cmd.PersistentFlags().Bool("no-color", false, "Plain help (no-color alias)")
	cmd.PersistentFlags().Bool("stdout", false, "Write help to stdout instead of stderr")
}

// ApplyToRoot registers operator-output persistent flags and routes help to the resolved policy.
func ApplyToRoot(cmd *cobra.Command, policy Policy) {
	RegisterPersistentFlags(cmd)
	cmd.SetOut(policy.DataWriter())
	cmd.SetErr(policy.DiagWriter())
	if policy.Plain {
		cmd.SetHelpTemplate(plainHelpTemplate())
	}
	installHelpFunc(cmd, policy)
	visitCommands(cmd, func(c *cobra.Command) {
		c.SetOut(policy.DataWriter())
		c.SetErr(policy.DiagWriter())
	})
}

func installHelpFunc(cmd *cobra.Command, policy Policy) {
	cmd.InitDefaultHelpFlag()
	defaultHelp := cmd.HelpFunc()
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		w := policy.HelpWriter()
		oldOut := c.OutOrStdout()
		c.SetOut(w)
		defer c.SetOut(oldOut)
		defaultHelp(c, args)
	})
	for _, sub := range cmd.Commands() {
		installHelpFunc(sub, policy)
	}
}

func visitCommands(cmd *cobra.Command, fn func(*cobra.Command)) {
	fn(cmd)
	for _, sub := range cmd.Commands() {
		visitCommands(sub, fn)
	}
}

func plainHelpTemplate() string {
	return `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
}

// PrintExecError writes execution failures to stderr per PRINT_EXEC_ERROR.
func PrintExecError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// CaptureHelp runs help for cmd into separate stdout/stderr buffers using policy (composition tests).
func CaptureHelp(cmd *cobra.Command, policy Policy, args []string) (stdout, stderr string, err error) {
	var outBuf, errBuf bytes.Buffer
	p := ResolvePolicy(PeeledFlags{}, WithWriters(&outBuf, &errBuf), WithTTYFunc(func(io.Writer) bool { return true }))
	p.HelpDest = policy.HelpDest
	p.Plain = policy.Plain
	ApplyToRoot(cmd, p)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return outBuf.String(), errBuf.String(), err
}
