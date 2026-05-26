//ff:func feature=analyzer type=helper control=sequence
//ff:what Resolves an http.StatusXxx selector expression to its integer status code value
package analyzer

import "go/ast"

// extractSelectorStatus resolves http.StatusXxx from a selector expression.
func extractSelectorStatus(expr ast.Expr) int {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return 0
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok || ident.Name != "http" {
		return 0
	}
	v, ok := resolveHTTPStatus(sel.Sel.Name)
	if !ok {
		return 0
	}
	return v
}
