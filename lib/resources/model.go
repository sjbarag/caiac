package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type goSourceResourceModel struct {
	Filename    types.String `tfsdk:"filename"`
	Contents    types.String `tfsdk:"contents"`
	PackageName types.String `tfsdk:"package_name"`
	Imports     types.List   `tfsdk:"import"`
	Funcs       types.List   `tfsdk:"func"`
}

var ImportSpec = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional: true,
		},
		"path": schema.StringAttribute{
			Required: true,
		},
	},
}

var FuncDecl = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required: true,
		},
	},
	Blocks: map[string]schema.Block{
		"signature": schema.SingleNestedBlock{
			Blocks: map[string]schema.Block{
				"param": schema.ListNestedBlock{
					NestedObject: Field,
				},
			},
		},
	},
}

var Field = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required: true,
		},
		"type": schema.StringAttribute{
			Required: true,
		},
	},
}
