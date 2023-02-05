package resources

import (
	"context"
	"go/ast"
	"go/format"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func renderGoSource(ctx context.Context, model *goSourceResourceModel, diags diag.Diagnostics) string {
	imports, err := MakeImportSpecAstNodes(ctx, model.Imports)
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
