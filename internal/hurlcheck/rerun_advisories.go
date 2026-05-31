//ff:func feature=hurlcheck type=engine control=iteration dimension=1
//ff:what Detects copy-pasted hurl entries (same shape, different status) and emits IMPROVE advisories
package hurlcheck

// RerunAdvisories detects entries that look copy-pasted across branches: the
// same (method, url, assertion-count) carried to a different status label, or an
// error status with no body variation. These are advisory IMPROVE hints only —
// they never hard-exclude from coverage (§3.5).
func RerunAdvisories(entries []HurlEntry) []string {
	seen := make(map[entryShape]int) // shape -> status of first occurrence
	var advisories []string
	for _, e := range entries {
		if e.Skip || e.Status == 0 {
			continue
		}
		advisories = appendShapeAdvisory(advisories, seen, e)
	}
	return advisories
}
