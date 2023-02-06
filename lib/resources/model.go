package resources

import (
	"go/ast"
	"terraform-provider-caiac/lib/astutil"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type goSourceResourceModel struct {
	Filename    types.String `tfsdk:"filename"`
	Contents    types.String `tfsdk:"contents"`
	PackageName types.String `tfsdk:"package_name"`
	Imports     []TImport    `tfsdk:"import"`
	Funcs       []TFunc      `tfsdk:"func"`
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

type TImport struct {
	Name *string `tfsdk:"name"`
	Path string  `tfsdk:"path"`
}

func (i *TImport) toAst() *ast.ImportSpec {
	return &ast.ImportSpec{
		Name: astutil.MaybeNewIdent(i.Name),
		Path: astutil.NewStringLiteral(i.Path),
	}
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
				"result": schema.ListNestedBlock{
					NestedObject: Field,
				},
			},
		},
	},
}

var Field = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional: true,
		},
		"type": schema.StringAttribute{
			Optional: true,
		},
	},
}

type TField struct {
	Name *string `tfsdk:"name"`
	Type *string `tfsdk:"type"`
}

func (f *TField) toAst() *ast.Field {
	names := []*ast.Ident{}
	name := astutil.MaybeNewIdent(f.Name)
	if name != nil {
		names = append(names, name)
	}
	return &ast.Field{
		Names: names,
		Type:  astutil.MaybeNewIdent(f.Type),
	}
}

type TSignature struct {
	Params  []TField `tfsdk:"param"`
	Results []TField `tfsdk:"result"`
}

func (s *TSignature) toAst() *ast.FuncType {
	if s == nil {
		return nil
	}

	params := []*ast.Field{}
	for _, param := range s.Params {
		params = append(params, param.toAst())
	}

	results := []*ast.Field{}
	for _, res := range s.Results {
		results = append(results, res.toAst())
	}

	return &ast.FuncType{
		Params: &ast.FieldList{
			List: params,
		},
		Results: &ast.FieldList{
			List: results,
		},
	}
}

type TFunc struct {
	Name      string     `tfsdk:"name"`
	Signature TSignature `tfsdk:"signature"`
}

func (f *TFunc) toAst() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(f.Name),
		Type: f.Signature.toAst(),
		Body: &ast.BlockStmt{},
	}
}
