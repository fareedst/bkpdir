package clioutput

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var cobraOrigHelpFuncs sync.Map // *cobra.Command -> func(*cobra.Command, []string)

// RegisterPersistentFlags documents global operator-output controls on the root command.
func RegisterPersistentFlags(cmd *cobra.Command) {
	if cmd.Annotations != nil && cmd.Annotations["clioutput.flagsRegistered"] == "1" {
		return
	}
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	// No -p shorthand on persistent flags — subcommands may use -p (e.g. server --port).
	cmd.PersistentFlags().Bool("plain", false, "Plain help text (no ANSI decoration)")
	cmd.PersistentFlags().Bool("no-color", false, "Plain help (no-color alias)")
	cmd.PersistentFlags().Bool("stdout", false, "Write help to stdout instead of stderr")
	cmd.Annotations["clioutput.flagsRegistered"] = "1"
}

// ApplyToRoot registers operator-output persistent flags and routes help to the resolved policy.
func ApplyToRoot(cmd *cobra.Command, policy Policy) {
	RegisterPersistentFlags(cmd)
	cmd.SetOut(policy.DataWriter())
	cmd.SetErr(policy.DiagWriter())
	installHelpFunc(cmd, policy)
	visitCommands(cmd, func(c *cobra.Command) {
		c.SetOut(policy.DataWriter())
		c.SetErr(policy.DiagWriter())
	})
}

func installHelpFunc(cmd *cobra.Command, policy Policy) {
	cmd.InitDefaultHelpFlag()
	var defaultHelp func(*cobra.Command, []string)
	if stored, ok := cobraOrigHelpFuncs.Load(cmd); ok {
		defaultHelp = stored.(func(*cobra.Command, []string))
	} else {
		defaultHelp = cmd.HelpFunc()
		cobraOrigHelpFuncs.Store(cmd, defaultHelp)
	}
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		var buf bytes.Buffer
		w := policy.HelpWriter()
		oldOut := c.OutOrStdout()
		c.SetOut(&buf)
		defer c.SetOut(oldOut)
		defaultHelp(c, args)
		rendered := RenderHelp(buf.String(), RenderOptions{Plain: policy.Plain})
		_, _ = io.WriteString(w, rendered)
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

// PrintExecError writes execution failures to stderr per PRINT_EXEC_ERROR.
func PrintExecError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

// CaptureHelp runs help for cmd into separate stdout/stderr buffers using policy (composition tests).
func CaptureHelp(cmd *cobra.Command, policy Policy, args []string) (stdout, stderr string, err error) {
	var outBuf, errBuf bytes.Buffer
	p := ResolvePolicy(PeeledFlags{}, WithWriters(&outBuf, &errBuf), WithEnv(map[string]string{}), WithTTYFunc(func(io.Writer) bool { return true }))
	p.HelpDest = policy.HelpDest
	p.Plain = policy.Plain
	ApplyToRoot(cmd, p)
	cmd.SetArgs(args)
	err = cmd.Execute()
	return outBuf.String(), errBuf.String(), err
}
