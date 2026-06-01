//ff:func feature=scan type=helper control=sequence
//ff:what Formats the UNVERIFIED reason when the only candidate file's extension mismatches the lang
package scanner

import (
	"fmt"
	"path/filepath"
)

// extMismatchMessage formats the reason an endpoint was left UNVERIFIED because
// its only candidate file's extension does not belong to the backend language.
func extMismatchMessage(handler, file, root, lang string) string {
	return fmt.Sprintf(
		"Skipped link: %s → only %s candidate under %s (%s; left UNVERIFIED)\n  fix: point --link-source at the backend source root or correct the handler name",
		handler, filepath.Ext(file), root, langLabel(lang))
}
