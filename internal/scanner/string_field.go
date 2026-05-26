//ff:func feature=scan type=helper control=sequence
//ff:what Extracts a string value from a map by key, returning empty string if missing or wrong type
package scanner

func stringField(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
