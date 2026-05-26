//ff:func feature=scan type=helper control=selection
//ff:what Looks up a response object by status code from a YAML responses map handling both string and interface key types
package scanner

import "fmt"

func lookupResponse(responsesRaw interface{}, code int) (map[string]interface{}, bool) {
	switch responses := responsesRaw.(type) {
	case map[string]interface{}:
		obj, ok := responses[fmt.Sprintf("%d", code)].(map[string]interface{})
		return obj, ok
	case map[interface{}]interface{}:
		return lookupResponseIface(responses, code)
	default:
		return nil, false
	}
}
