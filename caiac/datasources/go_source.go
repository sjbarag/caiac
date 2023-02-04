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
	_ datasource.DataSource              = &goSourceDataSource{}
	_ datasource.DataSourceWithConfigure = &goSourceDataSource{}
)

func NewGoSourceDataSource() datasource.DataSource {
	return &goSourceDataSource{}
}

type goSourceDataSourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Contents types.String `tfsdk:"contents"`
}

type goSourceDataSource struct {
	baseDir string
}

func (d *goSourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	dsd, ok := req.ProviderData.(*DataSourceData)
	if !ok {
		return
	}
	d.baseDir = dsd.BaseDir
}

func (d *goSourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_go_source"
}

func (d *goSourceDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

func (d *goSourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state goSourceDataSourceModel

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
