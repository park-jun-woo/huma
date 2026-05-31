//ff:func feature=hurlcheck type=engine control=iteration dimension=1
//ff:what Computes the minimum assertion-depth grade across a set of expected client statuses
package hurlcheck

// MinAGrade returns the minimum A-grade across the supplied client statuses,
// using only entries that assert each status. Statuses with no asserting entry
// contribute 0. Returns 0 for an empty status set.
func MinAGrade(entries []HurlEntry, statuses []int) int {
	if len(statuses) == 0 {
		return 0
	}
	min := 3
	for _, s := range statuses {
		g := gradeForStatus(entries, s)
		if g < min {
			min = g
		}
	}
	return min
}
