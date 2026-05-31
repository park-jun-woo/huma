//ff:func feature=rule type=helper control=sequence
//ff:what Returns the coverage-verdict validation rules (C-01 to C-04) for the cheese-resistant gate
package rule

var (
	C01 = Rule{ID: "C-01", Level: "ERROR", Description: "No-signal verdict cannot PASS — downgraded to UNVERIFIED"}
	C02 = Rule{ID: "C-02", Level: "ERROR", Description: "Denominator is monotonic — input spec cannot shrink ground-truth branches"}
	C03 = Rule{ID: "C-03", Level: "WARNING", Description: "Assertion depth below required level"}
	C04 = Rule{ID: "C-04", Level: "WARNING", Description: "DONE requires an unreachable.yaml reason artifact for uncovered branches"}
)

// CoverageRules returns all coverage-verdict validation rules.
func CoverageRules() []Rule {
	return []Rule{C01, C02, C03, C04}
}
