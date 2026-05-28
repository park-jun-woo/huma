//ff:func feature=scan type=parser control=iteration dimension=1
//ff:what Extracts HTTP methods from Edge Function source content using regex patterns
package scanner

import "strings"

func extractMethods(content string) []string {
	seen := map[string]bool{}
	var methods []string

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		collectLineMethods(line, edgeFuncMethodPositive, seen, &methods)
		collectLineMethods(line, edgeFuncMethodCase, seen, &methods)
	}

	if len(methods) == 0 {
		return []string{"POST"}
	}
	return methods
}
