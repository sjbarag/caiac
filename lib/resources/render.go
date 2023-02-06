package resources

import (
	"context"
	"go/ast"
	"go/format"
	"go/token"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func renderGoSource(ctx context.Context, model *goSourceResourceModel, diags *diag.Diagnostics) string {
	imports, err := makeImportSpecAstNodes(ctx, model.Imports)
	if err != nil {
		diags.AddError(
			"Error converting HCL to AST",
			"Unable to convert HCL imports to AST imports: "+err.Error(),
		)
		return ""
	}

	functions, err := makeFuncDecl(ctx, model.Funcs)
	if err != nil {
		diags.AddError(
			"Error converting HCL to AST",
			"Unable to convert HCL function declarations to AST function declarations: "+err.Error(),
		)
		return ""
	}

	decls := []ast.Decl{imports}
	decls = append(decls, functions...)

	f := &ast.File{
		Name:  ast.NewIdent(model.PackageName.ValueString()),
		Decls: decls,
	}
	fset := token.NewFileSet()

	contents := new(strings.Builder)
	if err := format.Node(contents, fset, f); err != nil {
		diags.AddError(
			"Error printing AST",
			"Unable to serialize AST to string: "+err.Error(),
		)
		return ""
	}

	return contents.String()
}

func makeImportSpecAstNodes(ctx context.Context, imports []TImport) (ast.Decl, error) {
	specs := []ast.Spec{}

	for _, theImport := range imports {
		specs = append(specs, theImport.toAst())
	}

	return &ast.GenDecl{Tok: token.IMPORT, Specs: specs}, nil
}

func makeFuncDecl(ctx context.Context, funcs []TFunc) ([]ast.Decl, error) {
	decls := []ast.Decl{}

	for _, theFunc := range funcs {
		decls = append(decls, theFunc.toAst())
	}

	return decls, nil
}
