package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
)

var rootCmd = &cobra.Command{
	Use:   "gen-main <registry-address> <module-path> <provider-package>",
	Short: "Generate a Terraform provider main.go",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		data := MainData{
			RegistryAddress: args[0],
			ModulePath:      args[1],
			ProviderPackage: args[2],
		}
		if err := GenerateMain(os.Stdout, data); err != nil {
			cli.Fail(err)
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		cli.Fail(err)
	}
}
