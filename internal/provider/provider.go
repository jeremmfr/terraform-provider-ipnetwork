package provider

import (
	"context"

	"github.com/jeremmfr/terraform-provider-ipnetwork/version"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider              = &ipnetworkProvider{}
	_ provider.ProviderWithFunctions = &ipnetworkProvider{}
)

type ipnetworkProvider struct{}

func New() provider.Provider {
	return &ipnetworkProvider{}
}

const (
	providerName = "ipnetwork"
)

func (p *ipnetworkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = providerName
	resp.Version = version.Get()
}

func (p *ipnetworkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *ipnetworkProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *ipnetworkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *ipnetworkProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		newAddressFunction,
		newCidrFunction,
		newExpand6Function,
		newPrefixFunction,
		newPtrFunction,
	}
}

func (p *ipnetworkProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	// no-op
}
