//ff:func feature=scan type=reader control=sequence
//ff:what Walks a root directory collecting source files with recognized code extensions
package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// sourceExts is the set of source file extensions link-source will scan.
var sourceExts = map[string]bool{
	".go": true, ".py": true, ".js": true, ".ts": true,
	".java": true, ".cs": true, ".php": true, ".rs": true,
}

// collectSourceFiles walks root and returns paths of recognized source files.
func collectSourceFiles(root string) []string {
	var files []string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if sourceExts[strings.ToLower(filepath.Ext(path))] {
			files = append(files, path)
		}
		return nil
	})
	return files
}
