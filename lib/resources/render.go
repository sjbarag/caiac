package resources

import (
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
	"terraform-provider-caiac/lib/astutil"
	"terraform-provider-caiac/lib/tfutil"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func renderGoSource(ctx context.Context, model *goSourceResourceModel, diags diag.Diagnostics) string {
	imports, err := makeImportSpecAstNodes(ctx, model.Imports)
	if err != nil {
		diags.AddError(
			"Error converting HCL to AST",
			"Unable to convert HCL imports to AST imports: "+err.Error(),
		)
		return ""
	}

	f := &ast.File{
		Name:  ast.NewIdent(model.PackageName.ValueString()),
		Decls: []ast.Decl{imports},
	}

	contents := new(strings.Builder)
	if err := format.Node(contents, nil, f); err != nil {
		diags.AddError(
			"Error printing AST",
			"Unable to serialize AST to string: "+err.Error(),
		)
		return ""
	}

	return contents.String()
}

func makeImportSpecAstNodes(ctx context.Context, in types.List) (ast.Decl, error) {
	specs := []ast.Spec{}

	elements := []types.Object{}
	diags := in.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		return nil, fmt.Errorf("unable to cast ImportSpecs list to []types.Object: %+v", diags)
	}

	for _, val := range elements {
		spec, err := importSpecToAst(ctx, val)
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}

	return &ast.GenDecl{Tok: token.IMPORT, Specs: specs}, nil
}

func importSpecToAst(ctx context.Context, in types.Object) (*ast.ImportSpec, error) {
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
