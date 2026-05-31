//ff:func feature=hurlcheck type=helper control=sequence
//ff:what Records an entry shape and appends a copy-paste or missing-error-body advisory if warranted
package hurlcheck

import "fmt"

// appendShapeAdvisory records the entry's shape and appends an advisory when a
// reused shape (different status) or an error status lacking a body assertion
// is detected.
func appendShapeAdvisory(advisories []string, seen map[entryShape]int, e HurlEntry) []string {
	s := entryShape{method: e.Method, url: e.URL, asserts: e.Asserts}
	if prev, ok := seen[s]; ok && prev != e.Status {
		return append(advisories, fmt.Sprintf(
			"%s %s: status %d reuses the assertion shape of status %d (possible copy-paste — vary the body assertions)",
			e.Method, e.URL, e.Status, prev))
	}
	seen[s] = e.Status
	if e.Status >= 400 && e.Asserts == 0 {
		return append(advisories, fmt.Sprintf(
			"%s %s: error status %d has no body assertion (assert the error schema)",
			e.Method, e.URL, e.Status))
	}
	return advisories
}
