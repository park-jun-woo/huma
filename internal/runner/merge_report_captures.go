//ff:func feature=runner type=helper control=iteration dimension=2
//ff:what Copies one hurl report's per-entry [Captures] name/value pairs into the running captures map, coercing each value to its hurl-variable string form
package runner

// mergeReportCaptures writes every capture in report's entries into dst, coercing
// each raw JSON value to its string form. Later entries override earlier ones on a
// name clash (hurl evaluates entries in order).
func mergeReportCaptures(dst map[string]string, report hurlJSONReport) {
	for _, e := range report.Entries {
		for _, c := range e.Captures {
			dst[c.Name] = captureToString(c.Value)
		}
	}
}
