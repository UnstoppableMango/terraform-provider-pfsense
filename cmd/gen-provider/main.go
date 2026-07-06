package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
)

var data TemplateData

var rootCmd = &cobra.Command{
	Use:   "gen-provider <output> <schema>",
	Short: "Generate Terraform provider glue code",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		root, err := os.OpenRoot(args[0])
		if err != nil {
			cli.Fail(err)
		}
		s, err := ParseSchema(args[1])
		if err != nil {
			cli.Fail(err)
		}
		if err := Generate(root, data, s); err != nil {
			cli.Fail(err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&data.RegistryAddress, "registry-address", "",
		"The terraform registry address",
	)
	rootCmd.Flags().StringVar(&data.ModulePath, "module-path", "",
		"The provider module path",
	)
	rootCmd.Flags().StringVar(&data.ProviderPackage, "provider-package", "",
		"The provider package name",
	)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		cli.Fail(err)
	}
}
