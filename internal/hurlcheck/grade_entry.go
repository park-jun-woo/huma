//ff:func feature=hurlcheck type=helper control=selection
//ff:what Computes the assertion-depth grade (0..3) for a single hurl entry
package hurlcheck

// gradeEntry computes the A-grade (0..3) for a single entry (§3.3).
func gradeEntry(e HurlEntry) int {
	switch {
	case e.Skip || e.Status == 0:
		return 0 // no status assertion / skipped — vacuous (CV-10)
	case e.Asserts == 0:
		return 1 // status == N only (CV-4)
	case e.Asserts >= 2:
		return 3 // status + shape + invariant (e.g. error schema)
	default:
		return 2 // status + 1 body/header assertion
	}
}
