package caiac

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-caiac/caiac/datasources"
)

var (
	_ provider.Provider = &caiacProvider{}
)

func New() provider.Provider {
	return &caiacProvider{}
}

type caiacProvider struct{}

// Metadata returns the provider type name.
func (p *caiacProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "caiac"
}

// Schema defines the provider-level schema for configuration data.
func (p *caiacProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure would normally perform some client-specific configuration, but there's nothing to do.
func (p *caiacProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *caiacProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewSourceGoDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *caiacProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}