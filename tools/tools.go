//go:build tools

package tools

import (
	_ "github.com/hashicorp/terraform-plugin-framework/providerserver"
	_ "github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)
