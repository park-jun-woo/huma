//ff:func feature=hurlcheck type=engine control=iteration dimension=1
//ff:what Returns the set of status codes tested by at least one non-vacuous (graded) hurl entry
package hurlcheck

// NonVacuousStatuses returns the set of status codes tested by at least one
// non-vacuous entry (A-grade >= 1, a real HTTP status assertion that is not
// skipped). Vacuous entries never count toward coverage (§3.5, CV-10).
func NonVacuousStatuses(entries []HurlEntry) map[int]bool {
	out := make(map[int]bool)
	for _, e := range entries {
		if e.Skip || e.Status == 0 || e.Grade == 0 {
			continue
		}
		out[e.Status] = true
	}
	return out
}
