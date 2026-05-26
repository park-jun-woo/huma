//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Counts the number of overlapping strings between two keyword slices
package cmd

func countKeywordOverlap(a, b []string) int {
	set := make(map[string]bool, len(b))
	for _, k := range b {
		set[k] = true
	}
	count := 0
	for _, k := range a {
		if set[k] {
			count++
		}
	}
	return count
}
