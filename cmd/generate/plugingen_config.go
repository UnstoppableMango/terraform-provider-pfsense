package main

import (
	"os"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/spf13/cobra"
	"github.com/unstoppablemango/terraform-provider-pfsense/pkg/plugingen"
	"github.com/unstoppablemango/terraform-provider-pfsense/pkg/read"
)

var plugingenConfig = &cobra.Command{
	Use:   "plugingen-config",
	Short: "Generates the terraform generator config",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oas, err := readOpenAPI(args[0])
		if err != nil {
			return err
		}

		cfg, err := plugingen.ConfigFor(oas)
		if err != nil {
			return err
		}

		dest, err := os.OpenFile(args[1], os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer dest.Close()

		return plugingen.WriteConfig(cfg, dest)
	},
}

func readOpenAPI(name string) (*v3.Document, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return read.OpenAPIV3(file)
}
