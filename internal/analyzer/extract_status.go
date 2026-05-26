//ff:func feature=analyzer type=helper control=sequence
//ff:what Resolves a status code from an AST expression (integer literal or http.StatusXxx constant)
package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
)

// extractStatus resolves the status code from an AST expression.
func extractStatus(expr ast.Expr) int {
	if lit, ok := expr.(*ast.BasicLit); ok && lit.Kind == token.INT {
		v, err := strconv.Atoi(lit.Value)
		if err != nil {
			return 0
		}
		return v
	}

	return extractSelectorStatus(expr)
}
