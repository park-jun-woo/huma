//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what Returns a string map's keys sorted, for stable transparency logging of captured variable names (keys only, never values)
package humaquest

import "sort"

// sortedKeys returns the keys of m sorted lexically, used to log which variables
// were captured/minted without ever logging their (secret) values.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
