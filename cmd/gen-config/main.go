package main

import (
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
)

var root = &cobra.Command{
	Use: "gen-config",
	Short: "Generates the terraform generator config",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO
	},
}

func main() {
	if err := root.Execute(); err != nil {
		cli.Fail(err)
	}
}

