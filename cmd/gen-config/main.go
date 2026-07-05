package main

import (
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
)

var rootCmd = &cobra.Command{
	Use:   "gen-main <registry-address> <module-path> <provider-package>",
	Short: "Generate a Terraform provider main.go",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		if err := Generate(cmd.OutOrStdout()); err != nil {
			cli.Fail(err)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		cli.Fail(err)
	}
}
