//ff:func feature=scan type=helper control=selection
//ff:what Converts a YAML map key (int, float64, or string) to an HTTP status code integer
package scanner

import "strconv"

func toStatusCode(key interface{}) int {
	switch v := key.(type) {
	case int:
		return v
	case float64:
		return int(v)
	case string:
		code, err := strconv.Atoi(v)
		if err != nil {
			return 0
		}
		return code
	default:
		return 0
	}
}
