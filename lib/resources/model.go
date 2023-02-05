package resources

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"

	"terraform-provider-caiac/lib/astutil"
	"terraform-provider-caiac/lib/tfutil"

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

func MakeImportSpecAstNodes(ctx context.Context, in types.List) (ast.Decl, error) {
	specs := []ast.Spec{}

	elements := []types.Object{}
	diags := in.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, fmt.Errorf("unable to cast ImportSpecs list to []types.Object: %+v", diags)
	}

	for _, val := range elements {
		spec, err := ImportSpecToAst(ctx, val)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}

	return &ast.GenDecl{Tok: token.IMPORT, Specs: specs}, nil
}

func ImportSpecToAst(ctx context.Context, in types.Object) (*ast.ImportSpec, error) {
	attr := in.Attributes()

	var name *ast.Ident
	nameStr, err := tfutil.AttrValueToString(ctx, attr["name"])
	if err != nil {
		return nil, err
	}
	if nameStr != "" {
		name = ast.NewIdent(nameStr)
	}

	importPath, err := tfutil.AttrValueToString(ctx, attr["path"])
	if err != nil {
		return nil, err
	}

	out := &ast.ImportSpec{
		Name: name,
		Path: astutil.NewStringLiteral(importPath),
	}

	return out, nil
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
