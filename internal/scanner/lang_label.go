//ff:func feature=scan type=helper control=sequence
//ff:what Renders the backend lang for skip messages, marking the unknown/fallback case explicitly
package scanner

// langLabel renders the backend lang for messages, marking the unknown/fallback
// case explicitly so users understand why ext-mismatch could still trigger.
func langLabel(lang string) string {
	if lang == "" {
		return "backend.lang=unknown"
	}
	return "backend.lang=" + lang
}
