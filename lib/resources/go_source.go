package resources

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &goSourceResource{}
	_ resource.ResourceWithConfigure = &goSourceResource{}
)

func NewGoSourceResource() resource.Resource {
	return &goSourceResource{}
}

type goSourceResource struct {
	baseDir string
}

func (r *goSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	rd, ok := req.ProviderData.(*ResourceData)
	if !ok {
		return
	}

	r.baseDir = rd.BaseDir
}

func (r *goSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_go_source"
}

func (r *goSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Required: true,
			},
			"contents": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *goSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan goSourceResourceModel
	{
		diags := req.Plan.Get(ctx, &plan)
		resp.Diagnostics.Append(diags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	path := filepath.Join(r.baseDir, plan.Filename.ValueString())
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		resp.Diagnostics.AddError(
			"Error creating directories",
			"Unable to create directory to hold new file: "+err.Error(),
		)
		return
	}
	if err := os.WriteFile(path, []byte(plan.Contents.ValueString()), os.ModePerm); err != nil {
		resp.Diagnostics.AddError(
			"Error writing file",
			"Unable to write file to disk: "+err.Error(),
		)
		return
	}

	{
		diags := resp.State.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *goSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state goSourceResourceModel
	{
		diags := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	contents, err := os.ReadFile(state.Filename.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading file",
			"unable to read file from disk: "+err.Error(),
		)
		return
	}

	state.Contents = types.StringValue(string(contents))
	{
		diags := resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *goSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *goSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}