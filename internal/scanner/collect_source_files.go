//ff:func feature=scan type=reader control=sequence
//ff:what Walks a root directory collecting source files limited to the backend language extensions
package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// collectSourceFiles walks root and returns paths of source files whose
// extension is allowed for lang. Unknown/empty lang falls back to all
// recognized source extensions (no regression for manifest-less repos).
func collectSourceFiles(root, lang string) []string {
	exts, _ := allowedExts(lang)
	var files []string
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if exts[strings.ToLower(filepath.Ext(path))] {
			files = append(files, path)
		}
		return nil
	})
	return files
}
