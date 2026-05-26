//ff:func feature=scan type=helper control=sequence
//ff:what Generates an operationId from method and path when operationId is missing in OpenAPI
package scanner

import "strings"

func generateOperationID(method, path string) string {
	clean := strings.ReplaceAll(path, "/", "_")
	clean = strings.ReplaceAll(clean, "{", "")
	clean = strings.ReplaceAll(clean, "}", "")
	if strings.HasPrefix(clean, "_") {
		clean = clean[1:]
	}
	return method + "_" + clean
}
