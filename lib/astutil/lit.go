package astutil

import (
	"go/ast"
	"go/token"
)

func NewStringLiteral(val string) *ast.BasicLit {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: `"` + val + `"`,
	}
}
