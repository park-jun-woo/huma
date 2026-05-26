//ff:func feature=analyzer type=parser control=iteration dimension=1
//ff:what Analyzes Go source using AST to extract gin response status codes from handler functions
package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// Analyze parses a Go file and extracts response branches from the specified handler.
func (g *GoAnalyzer) Analyze(file string, handlerName string, startLine, endLine int) ([]ResponseBranch, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", file, err)
	}

	funcMap := make(map[string]*ast.FuncDecl)
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcMap[fd.Name.Name] = fd
	}

	target, ok := funcMap[handlerName]
	if !ok {
		return nil, nil
	}
	if target.Body == nil {
		return nil, nil
	}

	var branches []ResponseBranch
	ast.Inspect(target.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				branches = append(branches, traceHelperCall(fset, file, funcMap, ident.Name, call)...)
			}
			return true
		}

		name := sel.Sel.Name
		switch name {
		case "JSON", "AbortWithStatusJSON", "Status", "AbortWithStatus":
			if len(call.Args) < 1 {
				return true
			}
			status := extractStatus(call.Args[0])
			if status <= 0 {
				return true
			}
			pos := fset.Position(call.Pos())
			line := strings.TrimSpace(nodeString(fset, call))
			branches = append(branches, ResponseBranch{
				Status: status,
				File:   file,
				Line:   pos.Line,
				Code:   line,
			})
		}

		return true
	})

	return branches, nil
}
