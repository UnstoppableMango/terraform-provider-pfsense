package main

import (
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
)

var root = &cobra.Command{
	Use:   "patch-openapi",
	Short: "Patches the pfrest OpenAPI spec for conversion into a terraform provider",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func main() {
	if err := root.Execute(); err != nil {
		cli.Fail(err)
	}
}
