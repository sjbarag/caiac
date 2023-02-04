package caiac

import (
	"context"
	"os"

	"terraform-provider-caiac/caiac/datasources"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &caiacProvider{}
)

func New() provider.Provider {
	return &caiacProvider{}
}

type caiacProvider struct {
	BaseDir string
}

type caiacProviderModel struct {
	BaseDir types.String `tfsdk:"base_dir"`
}

// Metadata returns the provider type name.
func (p *caiacProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "caiac"
}

// Schema defines the provider-level schema for configuration data.
func (p *caiacProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"base_dir": schema.StringAttribute{
				Optional:    true,
				Description: "The base directory, to which all other paths are relative.",
			},
		},
	}
}

// Configure would normally perform some client-specific configuration, but there's nothing to do.
func (p *caiacProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config caiacProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.BaseDir.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_dir"),
			"Unknown CaIaC base directory",
			"The CaIaC provider requires a base directory to know where files are, "+
				" but it received an unknown value. Please set it statically, or set CAIAC_BASE_DIR.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	baseDir := os.Getenv("CAIAC_BASE_DIR")

	if !config.BaseDir.IsNull() {
		baseDir = config.BaseDir.ValueString()
	}

	baseDir, err := os.Getwd()

	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("base_dir"),
			"Missing CaIaC base directory",
			"The CaIaC provider requires a base directory to know where files are, "+
				" but one wasn't provided and the current working directory couldn't be found: "+
				err.Error(),
		)
		return
	}

	resp.DataSourceData = &datasources.DataSourceData{
		BaseDir: baseDir,
	}
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