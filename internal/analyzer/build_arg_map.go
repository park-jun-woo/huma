//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Maps function parameter names to caller argument status code values for propagation
package analyzer

import "go/ast"

// buildArgMap maps parameter names to the caller's argument status code values.
func buildArgMap(fd *ast.FuncDecl, call *ast.CallExpr) map[string]int {
	m := make(map[string]int)
	if fd.Type.Params == nil {
		return m
	}

	names := flattenParamNames(fd.Type.Params)
	for i, name := range names {
		if i >= len(call.Args) {
			break
		}
		status := extractStatus(call.Args[i])
		if status > 0 {
			m[name] = status
		}
	}

	return m
}
