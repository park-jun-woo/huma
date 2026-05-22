//ff:func feature=scan type=helper control=selection
//ff:what Returns SkipDir for vendor, node_modules, and .git directories
package scanner

import "path/filepath"

func skipDir(path string) error {
	base := filepath.Base(path)
	switch base {
	case "vendor", "node_modules", ".git":
		return filepath.SkipDir
	}
	return nil
}
