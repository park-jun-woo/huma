//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Formats the UNVERIFIED reason when a handler name matched definitions in multiple files
package scanner

import (
	"fmt"
	"strings"
)

// ambiguousMessage formats the reason an endpoint was left UNVERIFIED because
// its handler name matched definitions in more than one candidate file.
func ambiguousMessage(handler string, matches []handlerMatch, root, lang string) string {
	names := make([]string, 0, len(matches))
	for _, m := range matches {
		names = append(names, m.File)
	}
	return fmt.Sprintf(
		"Skipped link: %s → %d candidates under %s (%s; left UNVERIFIED)\n  candidates: %s\n  fix: narrow --link-source to the backend source root or correct the handler name",
		handler, len(matches), root, langLabel(lang), strings.Join(names, ", "))
}
