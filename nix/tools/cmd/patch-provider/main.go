package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/unmango/go/cli"
	"sigs.k8s.io/controller-tools/pkg/loader"
)

func patch(root string) error {
	pkgs, err := loader.LoadRoots(root)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		fmt.Println("Got here: " + pkg.Dir)
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:   "patch-provider",
	Short: "Patch a Terraform provider",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := patch(args[0]); err != nil {
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
