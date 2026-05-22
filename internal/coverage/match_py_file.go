//ff:func feature=coverage type=helper control=sequence
//ff:what Checks if a coverage.py file path matches the given local handler file path
package coverage

import "strings"

// matchPyFile checks if a coverage.py file path matches the given handler file path.
func matchPyFile(coveragePath, localPath string) bool {
	coveragePath = strings.ReplaceAll(coveragePath, "\\", "/")
	localPath = strings.ReplaceAll(localPath, "\\", "/")

	if coveragePath == localPath {
		return true
	}

	return strings.HasSuffix(coveragePath, "/"+localPath) ||
		strings.HasSuffix(localPath, "/"+coveragePath)
}
