package provider_pfsense

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	// "github.com/unstoppablemango/terraform-provider-pfsense/resource_firewall_rule"
)

var _ provider.Provider = (*pfsenseProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider { return &pfsenseProvider{version: version} }
}

type pfsenseProvider struct{ version string }

func (p *pfsenseProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pfsense"
	resp.Version = p.version
}

func (p *pfsenseProvider) Schema(ctx context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	// resp.Schema = PfsenseProviderSchema(ctx)
}

func (p *pfsenseProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p *pfsenseProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p *pfsenseProvider) Resources(_ context.Context) []func() resource.Resource {
	// return []func() resource.Resource{resource_firewall_rule.NewFirewallRuleResource}
	return nil
}
