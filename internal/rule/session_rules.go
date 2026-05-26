//ff:func feature=rule type=helper control=sequence
//ff:what Returns all session validation rules (S-01 to S-03)
package rule

var (
	S01 = Rule{ID: "S-01", Level: "ERROR", Description: "No session found"}
	S02 = Rule{ID: "S-02", Level: "ERROR", Description: "Session file corrupt"}
	S03 = Rule{ID: "S-03", Level: "WARNING", Description: "Session has stale entries"}
)

// SessionRules returns all session validation rules.
func SessionRules() []Rule {
	return []Rule{S01, S02, S03}
}
