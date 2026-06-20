package main

import (
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
)

var rootCmd = &cobra.Command{
	Use:   "patch-provider",
	Short: "Patch a Terraform provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := Patch(args[0]); err != nil {
			cli.Fail(err)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	if err := Execute(); err != nil {
		cli.Fail(err)
	}
}
