//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Finds the first matching entry file (index.ts, index.tsx, index.js) in a Supabase Edge Function directory
package scanner

import (
	"os"
	"path/filepath"
)

func findEntryFile(funcDir string) string {
	candidates := []string{"index.ts", "index.tsx", "index.js"}
	for _, name := range candidates {
		if _, err := os.Stat(filepath.Join(funcDir, name)); err == nil {
			return name
		}
	}
	return ""
}
