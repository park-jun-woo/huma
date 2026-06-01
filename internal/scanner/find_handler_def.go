//ff:func feature=scan type=engine control=iteration dimension=1
//ff:what Scans all candidate files for a handler definition, returning every matching file:line
package scanner

// findHandlerDef scans EVERY candidate file (no early return) for a definition
// whose normalized name matches handlerName, returning all matches. The caller
// links only when exactly one file matches and rejects ambiguous (multi-file)
// matches per §2.5, replacing the old first-match-wins behavior that produced
// BUG-002's arbitrary alphabetical pick.
func findHandlerDef(files []string, handlerName string) []handlerMatch {
	target := normalizeSymbol(handlerName)
	if target == "" {
		return nil
	}
	var matches []handlerMatch
	for _, f := range files {
		if line := scanForHandler(f, target); line > 0 {
			matches = append(matches, handlerMatch{File: f, Line: line})
		}
	}
	return matches
}
