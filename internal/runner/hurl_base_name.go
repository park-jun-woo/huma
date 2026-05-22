//ff:func feature=runner type=helper control=sequence
//ff:what Generates the base hurl file name without directory prefix
package runner

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/hurlfill/internal/scanner"
)

func hurlFileName(ep *scanner.Endpoint) string {
	method := strings.ToLower(ep.Method)
	path := strings.ReplaceAll(ep.Path, "/", "_")
	path = strings.ReplaceAll(path, ":", "")
	path = strings.TrimPrefix(path, "_")
	return fmt.Sprintf("%s_%s.hurl", method, path)
}
