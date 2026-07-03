package integration_test

// providerHCL returns required_providers + provider config.
// dev_overrides in TF_CLI_CONFIG_FILE redirects to the local binary.
// host/credentials will be added once the provider schema is implemented.
func providerHCL() string {
	return `
terraform {
  required_providers {
    pfsense = {
      source = "registry.terraform.io/unstoppablemango/pfsense"
    }
  }
}

provider "pfsense" {}
`
}
