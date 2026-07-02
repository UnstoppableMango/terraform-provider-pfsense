package integration_test

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testExternalProviders = map[string]resource.ExternalProvider{
	"pfsense": {Source: "registry.terraform.io/unstoppablemango/pfsense"},
}

func providerHCL(host string) string {
	return fmt.Sprintf(`
provider "pfsense" {
  host         = %q
  client_id    = "test"
  client_token = "test"
}
`, host)
}
