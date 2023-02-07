package resources

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"terraform-provider-caiac/lib/astutil"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type goSourceResourceModel struct {
	Filename    types.String `tfsdk:"filename"`
	Contents    types.String `tfsdk:"contents"`
	PackageName types.String `tfsdk:"package_name"`
	Imports     []TImport    `tfsdk:"import"`
	Funcs       []TFunc      `tfsdk:"func"`
}

var ImportSpec = &schema.NestedBlockObject{
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
	tflog.Debug(context.Background(), "OH HEY", nil)
	return &ast.ImportSpec{
		Name: astutil.MaybeNewIdent(i.Name),
		Path: astutil.NewStringLiteral(i.Path),
	}
}

var FuncDecl = &schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Required: true,
		},
	},
	Blocks: map[string]schema.Block{
		"signature": schema.SingleNestedBlock{
			Blocks: map[string]schema.Block{
				"param":  Params,
				"result": Results,
			},
		},
		"body": Body,
	},
}

var Params = schema.ListNestedBlock{
	NestedObject: *Field,
}
var Results = schema.ListNestedBlock{
	NestedObject: *Field,
}
var Field = &schema.NestedBlockObject{
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

var Body = schema.SingleNestedBlock{
	Blocks: map[string]schema.Block{
		"statement": schema.ListNestedBlock{
			NestedObject: Statement,
		},
	},
}

type TBody struct {
	Statements []TStatement `tfsdk:"statement"`
}

func (b *TBody) toAst() *ast.BlockStmt {
	stmts := []ast.Stmt{}
	for _, stmt := range b.Statements {
		stmts = append(stmts, stmt.toAst())
	}
	return &ast.BlockStmt{
		List: stmts,
	}
}

type stmtKind = string

// Define just enough to get through the Tour of Go
const (
	KExpr   stmtKind = "expression"
	KReturn stmtKind = "return"
)

var Statement = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"kind": schema.StringAttribute{
			Required: true,
		},
	},
	Blocks: map[string]schema.Block{
		"expression": Expression,
	},
}

type TStatement struct {
	Kind stmtKind     `tfsdk:"kind"`
	Expr *TExpression `tfsdk:"expression"`
}

func (s *TStatement) toAst() ast.Stmt {
	switch s.Kind {
	case KExpr:
		return &ast.ExprStmt{X: s.Expr.toAst()}
	default:
		return nil
	}
}

type exprKind = string

const (
	KCall       exprKind = "call"
	KSelector   exprKind = "selector"
	KLiteral    exprKind = "literal"
	KIdentifier exprKind = "identifier"
)

var Expression = &schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
		"kind": schema.StringAttribute{
			Required: true,
		},
	},
	Blocks: map[string]schema.Block{
		"literal": schema.SingleNestedBlock{
			Attributes: Literal.Attributes,
		},
		"selector":   Selector,
		"identifier": Identifier,
		"call":       Call,
	},
}

// lazily initialize some Expression.Blocks entries to avoid initialization loops
func init() {
	Expression.Blocks["call"] = Call
}

type TExpression struct {
	Kind       exprKind     `tfsdk:"kind"`
	Selector   *TSelector   `tfsdk:"selector"`
	Call       *TCall       `tfsdk:"call"`
	Literal    *TLiteral    `tfsdk:"literal"`
	Identifier *TIdentifier `tfsdk:"identifier"`
}

func (e *TExpression) toAst() ast.Expr {
	fmt.Println("[TEspression.toAst]")
	switch e.Kind {
	case KCall:
		return e.Call.toAst()
	case KSelector:
		return e.Selector.toAst()
	case KLiteral:
		return e.Literal.toAst()
	case KIdentifier:
		return e.Identifier.toAst()
	default:
		return nil
	}
}

type TIdentifier struct {
	Name string `tfsdk:"name"`
}

func (i *TIdentifier) toAst() *ast.Ident {
	return ast.NewIdent(i.Name)
}

var Identifier = &schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Optional: true,
		},
	},
}

var Call = &schema.SingleNestedBlock{
	Blocks: map[string]schema.Block{
		"func": Selector,
		"arg": schema.ListNestedBlock{
			NestedObject: Literal,
		},
	},
}

type TCall struct {
	Func *TSelector `tfsdk:"func"`
	Args []TLiteral `tfsdk:"arg"`
}

func (c *TCall) toAst() *ast.CallExpr {
	fmt.Println("[TCall.toAst]")
	args := []ast.Expr{}
	for _, arg := range c.Args {
		args = append(args, arg.toAst())
	}
	return &ast.CallExpr{
		Fun:  c.Func.toAst(),
		Args: args,
	}
}

var Selector = schema.SingleNestedBlock{
	Attributes: map[string]schema.Attribute{
		"from": schema.StringAttribute{
			Optional: true,
		},
		"prop": schema.StringAttribute{
			Optional: true,
		},
	},
}

type TSelector struct {
	From *string `tfsdk:"from"`
	Prop string  `tfsdk:"prop"`
}

func (s *TSelector) toAst() *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   astutil.MaybeNewIdent(s.From),
		Sel: ast.NewIdent(s.Prop),
	}
}

type litKind = string

const (
	LitIdent  litKind = "identifier"
	LitString litKind = "string"
	LitInt    litKind = "int"
)

type TLiteral struct {
	Kind  litKind `tfsdk:"kind"`
	Value string  `tfsdk:"value"`
}

func (l *TLiteral) toAst() *ast.BasicLit {
	switch l.Kind {
	case LitIdent:
		return &ast.BasicLit{Kind: token.IDENT, Value: l.Value}
	case LitString:
		return &ast.BasicLit{Kind: token.STRING, Value: `"` + l.Value + `"`}
	case LitInt:
		return &ast.BasicLit{Kind: token.INT, Value: l.Value}
	default:
		return nil
	}
}

var Literal = schema.NestedBlockObject{
	Attributes: map[string]schema.Attribute{
		"kind": schema.StringAttribute{
			Optional: true,
		},
		"value": schema.StringAttribute{
			Optional: true,
		},
	},
}

type TFunc struct {
	Name      string      `tfsdk:"name"`
	Signature *TSignature `tfsdk:"signature"`
	Body      *TBody      `tfsdk:"body"`
}

func (f *TFunc) toAst() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(f.Name),
		Type: f.Signature.toAst(),
		Body: f.Body.toAst(),
	}
}
