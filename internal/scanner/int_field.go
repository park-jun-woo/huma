//ff:func feature=scan type=helper control=selection
//ff:what Extracts an int value from a map by key, returning 0 if missing or wrong type
package scanner

func intField(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case int:
		return n
	case float64:
		return int(n)
	default:
		return 0
	}
}
