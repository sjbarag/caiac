package datasources

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &sourceGoDataSource{}
	_ datasource.DataSourceWithConfigure = &sourceGoDataSource{}
)

func NewSourceGoDataSource() datasource.DataSource {
	return &sourceGoDataSource{}
}

type sourceGoDataSource struct {
	baseDir string
}

func (d *sourceGoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	dsd, ok := req.ProviderData.(*DataSourceData)
	if !ok {
		return
	}
	d.baseDir = dsd.BaseDir
}

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

	path := filepath.Join(d.baseDir, state.Filename.ValueString())

	ctx = tflog.SetField(ctx, "filename", state.Filename.ValueString())
	ctx = tflog.SetField(ctx, "path", path)
	tflog.Debug(ctx, "Reading file")

	contents, err := os.ReadFile(path)
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