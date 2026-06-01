//ff:func feature=scan type=engine control=selection
//ff:what Classifies handler-definition matches into link / ambiguous / ext-mismatch / not-found
package scanner

import (
	"path/filepath"
	"strings"
)

// classifyMatches turns the candidate set into an outcome: link the single
// allowed match, or reject (ambiguous / ext-mismatch / not found). On a single
// allowed match it sets Source/Line on ep (§2.5).
func classifyMatches(ep *Endpoint, matches []handlerMatch, root, lang string, extSet map[string]bool) (linkOutcome, string) {
	switch {
	case len(matches) == 0:
		return outcomeNotFound, ""
	case len(matches) > 1:
		return outcomeAmbiguous, ambiguousMessage(ep.Handler, matches, root, lang)
	case !extSet[strings.ToLower(filepath.Ext(matches[0].File))]:
		return outcomeExtMismatch, extMismatchMessage(ep.Handler, matches[0].File, root, lang)
	default:
		ep.Source = matches[0].File
		ep.Line = matches[0].Line
		return outcomeLinked, ""
	}
}
