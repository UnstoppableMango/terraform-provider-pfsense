package main

import (
	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
)

var rootCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate various files",
}

func init() {
	rootCmd.AddCommand(plugingenConfig)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		cli.Fail(err)
	}
}
