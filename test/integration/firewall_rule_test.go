package integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFirewallRule_schema(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ExternalProviders: testExternalProviders,
		Steps: []resource.TestStep{
			{
				Config: providerHCL(mockServerURL) + `
resource "pfsense_firewall_rule" "test" {
  type        = "pass"
  interface   = ["wan"]
  ipprotocol  = "inet"
  source      = "any"
  destination = "any"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("pfsense_firewall_rule.test", "type", "pass"),
					resource.TestCheckResourceAttr("pfsense_firewall_rule.test", "ipprotocol", "inet"),
					resource.TestCheckResourceAttr("pfsense_firewall_rule.test", "source", "any"),
					resource.TestCheckResourceAttr("pfsense_firewall_rule.test", "destination", "any"),
					resource.TestCheckResourceAttrSet("pfsense_firewall_rule.test", "id"),
				),
			},
		},
	})
}
