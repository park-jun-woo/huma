//ff:func feature=runner type=engine control=iteration dimension=1
//ff:what Searches for the hurl file in the hurl directory and current directory
package runner

import (
	"os"
	"path/filepath"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func FindHurlFile(ep *scanner.Endpoint, hurlDir string) string {
	name := hurlFileName(ep)
	candidates := []string{
		filepath.Join(hurlDir, name),
		name,
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}
