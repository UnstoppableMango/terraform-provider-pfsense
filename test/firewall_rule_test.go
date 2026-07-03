package integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccFirewallRule_scaffold verifies the provider binary loads, registers the
// firewall_rule resource type, and can create a resource. This test targets the
// current scaffold (schema: id only, Create hardcodes "example-id").
func TestAccFirewallRule_scaffold(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				Config: providerHCL() + `
resource "pfsense_firewall_rule" "test" {}
`,
				Check: resource.TestCheckResourceAttr(
					"pfsense_firewall_rule.test", "id", "example-id",
				),
			},
		},
	})
}
