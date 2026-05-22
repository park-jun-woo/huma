//ff:func feature=coverage type=parser control=iteration dimension=1
//ff:what Filters missing lines from coverage.py data to those within a handler line range
package coverage

// filterMissingLines returns the subset of missingLines that fall within [startLine, endLine].
func filterMissingLines(missingLines []int, startLine, endLine int) []int {
	var filtered []int
	for _, line := range missingLines {
		if line >= startLine && line <= endLine {
			filtered = append(filtered, line)
		}
	}
	return filtered
}
