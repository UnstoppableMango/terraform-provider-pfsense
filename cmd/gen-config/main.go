package main

import (
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
	"github.com/unstoppablemango/terraform-provider-pfsense/pkg"
)

var root = &cobra.Command{
	Use:   "gen-config",
	Short: "Generates the terraform generator config",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		src, dest := args[0], args[1]

		if err := pkg.GenerateConfig(ctx, src, dest); err != nil {
			cli.Fail(err)
		}
	},
}

func main() {
	if err := root.Execute(); err != nil {
		cli.Fail(err)
	}
}
