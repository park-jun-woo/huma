//ff:func feature=coverage type=helper control=sequence
//ff:what Checks if a module-qualified coverage path ends with the given local file path
package coverage

import "strings"

// matchFile checks if a module-qualified coverage path ends with the given file path.
// Coverage files use paths like "github.com/park-jun-woo/hurlfill/testdata/server/main.go"
// and we need to match against local paths like "testdata/server/main.go".
func matchFile(coveragePath, localPath string) bool {
	coveragePath = strings.ReplaceAll(coveragePath, "\\", "/")
	localPath = strings.ReplaceAll(localPath, "\\", "/")

	if coveragePath == localPath {
		return true
	}

	return strings.HasSuffix(coveragePath, "/"+localPath)
}
