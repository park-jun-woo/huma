//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Searches an interface-keyed response map for a matching status code and returns the response object
package scanner

func lookupResponseIface(responses map[interface{}]interface{}, code int) (map[string]interface{}, bool) {
	for k, v := range responses {
		if toStatusCode(k) == code {
			obj, ok := v.(map[string]interface{})
			return obj, ok
		}
	}
	return nil, false
}
