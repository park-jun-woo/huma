//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Returns the allowed source extensions for a backend lang, with full-set fallback when unknown
package scanner

import "strings"

// allowedExts returns the set of extensions to scan for the given lang. When
// lang is empty or unrecognized it returns the full sourceExts fallback and
// false (unknown), so callers can flag low-confidence linking per §2.1.
func allowedExts(lang string) (map[string]bool, bool) {
	exts, ok := langExts[strings.ToLower(lang)]
	if !ok {
		return sourceExts, false
	}
	set := make(map[string]bool, len(exts))
	for _, e := range exts {
		set[e] = true
	}
	return set, true
}
