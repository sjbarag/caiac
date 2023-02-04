package datasources

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = &sourceGoDataSource{}
)

func NewSourceGoDataSource() datasource.DataSource {
	return &sourceGoDataSource{}
}

type sourceGoDataSource struct{}

func (d *sourceGoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_go"
}

func (d *sourceGoDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Required:    true,
				Description: "The absolute path to the file on-disk.",
			},
			"contents": schema.StringAttribute{
				Computed:    true,
				Description: "The rendered content as it exists on-disk.",
			},
		},
	}
}

func (d *sourceGoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sourceGoDataSourceModel

	{
		diags := req.Config.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
	}

	contents, err := os.ReadFile(state.Filename.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read file from disk",
			err.Error(),
		)
		return
	}

	state.Contents = types.StringValue(string(contents))

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

type sourceGoDataSourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Contents types.String `tfsdk:"contents"`
}