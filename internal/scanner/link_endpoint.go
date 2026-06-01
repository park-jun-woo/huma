//ff:func feature=scan type=engine control=selection
//ff:what Attempts to link one endpoint to its handler definition, classifying skip reasons
package scanner

// linkEndpoint tries to set Source/Line on ep using files restricted to lang.
// It links only when exactly one file defines the handler AND that file's
// extension is allowed for lang; otherwise it rejects (ambiguous/ext-mismatch)
// and leaves the endpoint UNVERIFIED, returning a user-facing reason. extSet is
// the allowed-extension set for lang (full sourceExts when lang is unknown);
// (lang fallback is handled by extSet itself: unknown lang -> full sourceExts).
func linkEndpoint(ep *Endpoint, files []string, root, lang string, extSet map[string]bool) (linkOutcome, string) {
	switch {
	case ep.Source != "" || ep.Handler == "":
		return outcomeNoop, ""
	default:
		matches := findHandlerDef(files, ep.Handler)
		return classifyMatches(ep, matches, root, lang, extSet)
	}
}
