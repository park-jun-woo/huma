//ff:func feature=hurlcheck type=engine control=iteration dimension=1
//ff:what Returns the non-vacuous tested status codes as a slice
package hurlcheck

// NonVacuousStatusList returns the non-vacuous tested statuses as a slice.
func NonVacuousStatusList(entries []HurlEntry) []int {
	set := NonVacuousStatuses(entries)
	out := make([]int, 0, len(set))
	for s := range set {
		out = append(out, s)
	}
	return out
}
