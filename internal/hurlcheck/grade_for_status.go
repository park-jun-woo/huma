//ff:func feature=hurlcheck type=helper control=iteration dimension=1
//ff:what Returns the best assertion-depth grade among entries asserting a given status
package hurlcheck

// gradeForStatus returns the maximum A-grade among non-skipped entries that
// assert the given status. Returns 0 if no such entry exists.
func gradeForStatus(entries []HurlEntry, status int) int {
	best := 0
	for _, e := range entries {
		if e.Status != status || e.Skip {
			continue
		}
		if e.Grade > best {
			best = e.Grade
		}
	}
	return best
}
