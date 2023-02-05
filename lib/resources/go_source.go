package resources

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				Optional: true,
				Computed: true,
			},
			"package_name": schema.StringAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"import": schema.ListNestedBlock{
				NestedObject: ImportSpec,
			},
			"func": schema.ListNestedBlock{
				NestedObject: FuncDecl,
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

	contents := renderGoSource(ctx, &plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Contents = types.StringValue(contents)

	if err := os.WriteFile(path, []byte(contents), os.ModePerm); err != nil {
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
	var plan goSourceResourceModel
	{
		diags := req.Plan.Get(ctx, &plan)
		resp.Diagnostics.Append(diags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	path := filepath.Join(r.baseDir, plan.Filename.ValueString())

	contents := renderGoSource(ctx, &plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Contents = types.StringValue(contents)

	if err := os.WriteFile(path, []byte(contents), os.ModePerm); err != nil {
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

func (r *goSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state goSourceResourceModel
	{
		diags := req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	// Remove the file
	if err := os.Remove(state.Filename.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error removing file",
			"Unable to delete file from disk: "+err.Error(),
		)
		return
	}

	// Then remove any empty directories above it in the filesystem.
	dir := filepath.Dir(state.Filename.ValueString())
	for dir != "/" {
		ctx = tflog.SetField(ctx, "dir", dir)
		tflog.Debug(ctx, "empty-dir removal loop")
		entries, err := os.ReadDir(dir)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error searching for empty directories to remove",
				"Unable to find empty directories to remove after file deletion: "+err.Error(),
			)
			return
		}

		if len(entries) > 0 {
			break
		}

		if err := os.Remove(dir); err != nil {
			resp.Diagnostics.AddError(
				"Error removing empty directory",
				"Unable to remove empty directory after file deletion: "+err.Error(),
			)
			return
		}

		dir = filepath.Dir(dir)
	}
}