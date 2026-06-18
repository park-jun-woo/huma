//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what Merges the static hurl_variables with the dynamically captured/minted vars, letting the dynamic ones (token, fixtures) override — the precedence Phase 009 needs so a captured {{token}} wins over any stale static value
package humaquest

// mergeVars returns base overlaid with over: every key in base is copied first,
// then every key in over is written, so over (the captured/minted extraVars) wins
// on conflict. A fresh map is returned (neither input is mutated) and nil inputs
// are handled. This is the precedence the cover loop relies on — a freshly captured
// token must override any hand-set hurl_variables.token.
func mergeVars(base, over map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(over))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range over {
		out[k] = v
	}
	return out
}
