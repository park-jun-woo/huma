//ff:func feature=scan type=helper control=sequence
//ff:what Returns true if the path ends with .go and is not a test file
package scanner

import "strings"

func isGoSource(path string) bool {
	return strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go")
}
