//ff:func feature=rule type=helper control=sequence
//ff:what Returns all adapter and server validation rules (A-01 to A-06)
package rule

var (
	A01 = Rule{ID: "A-01", Level: "ERROR", Description: "Server healthcheck failed"}
	A02 = Rule{ID: "A-02", Level: "ERROR", Description: "Server build command failed"}
	A03 = Rule{ID: "A-03", Level: "ERROR", Description: "Server start command failed"}
	A04 = Rule{ID: "A-04", Level: "ERROR", Description: "Server ready timeout"}
	A05 = Rule{ID: "A-05", Level: "ERROR", Description: "Coverage data collection failed"}
	A06 = Rule{ID: "A-06", Level: "WARNING", Description: "deps.ready check failed"}
)

// AdapterRules returns all adapter/server validation rules.
func AdapterRules() []Rule {
	return []Rule{A01, A02, A03, A04, A05, A06}
}
