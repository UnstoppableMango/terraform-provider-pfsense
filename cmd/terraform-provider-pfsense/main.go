package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/unstoppablemango/terraform-provider-pfsense/provider_pfsense"
)

var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/unstoppablemango/pfsense",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider_pfsense.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
