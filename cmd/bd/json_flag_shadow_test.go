package main

import (
	"testing"

	"github.com/spf13/cobra"
)

// TestNoSubcommandShadowsRootJSONFlag guards the root --json contract.
//
// --json is a persistent flag on rootCmd. When a subcommand registers its own
// --json (local or persistent), pflag keeps that one and drops the inherited
// one, so rootCmd.PersistentFlags().Changed("json") stays false for that
// command. The PersistentPreRun then treats the flag as unset and overwrites
// jsonOutput with the configured default, and `bd <cmd> --json` prints text
// unless the config also says json. Every command must inherit the root flag
// instead.
func TestNoSubcommandShadowsRootJSONFlag(t *testing.T) {
	var offenders []string
	var walk func(cmd *cobra.Command)
	walk = func(cmd *cobra.Command) {
		for _, sub := range cmd.Commands() {
			if sub.LocalNonPersistentFlags().Lookup("json") != nil || sub.PersistentFlags().Lookup("json") != nil {
				offenders = append(offenders, sub.CommandPath())
			}
			walk(sub)
		}
	}
	walk(rootCmd)

	for _, path := range offenders {
		t.Errorf("%q registers a local --json flag that shadows the root persistent flag", path)
	}
}
