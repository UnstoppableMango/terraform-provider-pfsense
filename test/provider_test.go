package integration_test

import (
	"fmt"
	"os"
)

// providerHCL returns required_providers + provider config.
// dev_overrides in TF_CLI_CONFIG_FILE redirects to the local binary.
// PFSENSE_MOCK_URL is set by TestMain to the mock server address.
func providerHCL() string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    pfsense = {
      source = "registry.terraform.io/unstoppablemango/pfsense"
    }
  }
}

provider "pfsense" {
  host     = %q
  username = "admin"
  password = "pfsense"
}
`, os.Getenv("PFSENSE_MOCK_URL"))
}
