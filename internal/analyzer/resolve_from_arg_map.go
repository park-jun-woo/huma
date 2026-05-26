//ff:func feature=analyzer type=helper control=sequence
//ff:what Resolves a parameter identifier to a status code using the caller argument map
package analyzer

import "go/ast"

// resolveFromArgMap resolves a parameter reference to a status code via the arg map.
func resolveFromArgMap(expr ast.Expr, argMap map[string]int) int {
	if ident, ok := expr.(*ast.Ident); ok {
		if v, ok := argMap[ident.Name]; ok {
			return v
		}
	}
	return 0
}
