package main

import (
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
)

var root = &cobra.Command{
	Use:   "patch-openapi",
	Short: "Patches the pfrest OpenAPI spec for conversion into a terraform provider",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		src, dest := args[0], args[1]

		if err := PatchSpec(ctx, src, dest); err != nil {
			cli.Fail(err)
		}
	},
}

func main() {
	if err := root.Execute(); err != nil {
		cli.Fail(err)
	}
}
