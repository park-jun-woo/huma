//ff:func feature=prompt type=helper control=sequence
//ff:what Computes integer percentage, returning 0 when total is zero
package prompt

func percent(n, total int) int {
	if total == 0 {
		return 0
	}
	return n * 100 / total
}
