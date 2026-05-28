//ff:func feature=analyzer type=parser control=sequence
//ff:what Traces 1-depth helper function calls to find Gin, Fiber and Echo response status codes via argument propagation
package analyzer

import (
	"go/ast"
	"go/token"
)

// traceHelperCall does 1-depth call tracing: look for c.JSON calls inside the helper.
func traceHelperCall(fset *token.FileSet, file string, funcMap map[string]*ast.FuncDecl, funcName string, call *ast.CallExpr) []ResponseBranch {
	fd, ok := funcMap[funcName]
	if !ok || fd.Body == nil {
		return nil
	}

	argMap := buildArgMap(fd, call)

	var branches []ResponseBranch
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		innerCall, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := innerCall.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		switch sel.Sel.Name {
		case "JSON", "AbortWithStatusJSON", "Status", "AbortWithStatus", "SendStatus",
			"NoContent", "String", "HTML", "XML", "JSONBlob", "Blob", "Redirect":
			if len(innerCall.Args) < 1 {
				return true
			}
			status := extractStatus(innerCall.Args[0])
			if status <= 0 {
				status = resolveFromArgMap(innerCall.Args[0], argMap)
			}
			if status <= 0 {
				return true
			}
			pos := fset.Position(call.Pos())
			line := nodeString(fset, call)
			branches = append(branches, ResponseBranch{
				Status: status,
				File:   file,
				Line:   pos.Line,
				Code:   line,
			})
		}
		return true
	})

	return branches
}
