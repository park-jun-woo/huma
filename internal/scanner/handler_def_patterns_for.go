//ff:func feature=scan type=engine control=selection
//ff:what Selects the ordered handler-definition regexes to apply to a file by its extension
package scanner

import (
	"path/filepath"
	"regexp"
	"strings"
)

// handlerDefPatterns returns the ordered list of definition regexes to apply to
// a file, chosen by file extension. Each regex captures the defined identifier
// in group 1. An empty slice means the extension defines no recognized shape.
func handlerDefPatterns(file string) []*regexp.Regexp {
	switch strings.ToLower(filepath.Ext(file)) {
	case ".go":
		return []*regexp.Regexp{goDefRe}
	case ".py":
		return []*regexp.Regexp{pyDefRe}
	case ".rs":
		return []*regexp.Regexp{rsDefRe}
	case ".java", ".cs":
		return []*regexp.Regexp{javaDefRe}
	case ".php":
		return []*regexp.Regexp{phpDefRe}
	case ".js", ".ts":
		return jsDefRes
	default:
		return nil
	}
}
