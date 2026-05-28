//ff:func feature=adapter type=helper control=sequence
//ff:what Checks if a JaCoCo source file name matches the handler file by base name or package path
package adapter

import (
	"path/filepath"
	"strings"
)

// matchSourceFile checks if a JaCoCo source file matches the handler file.
func matchSourceFile(sfName, baseName, handlerFile, pkgName string) bool {
	if sfName == baseName {
		return true
	}
	pkgPath := strings.ReplaceAll(pkgName, "/", string(filepath.Separator))
	fullPath := filepath.Join(pkgPath, sfName)
	return strings.HasSuffix(handlerFile, fullPath)
}
