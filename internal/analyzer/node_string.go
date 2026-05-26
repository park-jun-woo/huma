//ff:func feature=analyzer type=helper control=sequence
//ff:what Returns a simple string representation of an AST call expression for display
package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// nodeString returns a simple string representation of a call expression.
func nodeString(fset *token.FileSet, call *ast.CallExpr) string {
	start := fset.Position(call.Pos())
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if x, ok := sel.X.(*ast.Ident); ok {
			return strings.TrimSpace(fmt.Sprintf("%s.%s(...)", x.Name, sel.Sel.Name))
		}
		return strings.TrimSpace(fmt.Sprintf("%s(...)", sel.Sel.Name))
	}
	if ident, ok := call.Fun.(*ast.Ident); ok {
		return strings.TrimSpace(fmt.Sprintf("%s(...)", ident.Name))
	}
	return strings.TrimSpace(fmt.Sprintf("call at line %d", start.Line))
}
