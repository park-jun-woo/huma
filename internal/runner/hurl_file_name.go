//ff:func feature=runner type=helper control=sequence
//ff:what Generates the expected hurl file path from an endpoint's method, path, and hurl directory
package runner

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/huma/internal/scanner"
)

// HurlFileName returns the expected .hurl file name for an endpoint.
func HurlFileName(ep *scanner.Endpoint, hurlDir string) string {
	method := strings.ToLower(ep.Method)
	path := strings.ReplaceAll(ep.Path, "/", "_")
	path = strings.ReplaceAll(path, ":", "")
	path = strings.TrimPrefix(path, "_")
	return filepath.Join(hurlDir, fmt.Sprintf("%s_%s.hurl", method, path))
}
