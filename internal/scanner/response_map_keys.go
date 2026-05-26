//ff:func feature=scan type=helper control=selection
//ff:what Extracts keys from a YAML map as interface slice handling both string-keyed and interface-keyed maps
package scanner

func responseMapKeys(responsesRaw interface{}) []interface{} {
	switch responses := responsesRaw.(type) {
	case map[string]interface{}:
		keys := make([]interface{}, 0, len(responses))
		for k := range responses {
			keys = append(keys, k)
		}
		return keys
	case map[interface{}]interface{}:
		keys := make([]interface{}, 0, len(responses))
		for k := range responses {
			keys = append(keys, k)
		}
		return keys
	default:
		return nil
	}
}
