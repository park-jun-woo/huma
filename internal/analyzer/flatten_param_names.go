//ff:func feature=analyzer type=helper control=iteration dimension=1
//ff:what Flattens AST function parameter fields into a flat list of parameter names
package analyzer

import "go/ast"

// flattenParamNames extracts all parameter names from a field list in order.
func flattenParamNames(params *ast.FieldList) []string {
	var names []string
	for _, field := range params.List {
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
	}
	return names
}
