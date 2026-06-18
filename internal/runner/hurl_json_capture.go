//ff:type feature=runner type=model
//ff:what hurlJSONCapture is one [Captures] name/value pair from a hurl --json report; the value is kept raw for type coercion
package runner

import "encoding/json"

// hurlJSONCapture is one [Captures] name/value pair (value kept raw for coercion).
type hurlJSONCapture struct {
	Name  string          `json:"name"`
	Value json.RawMessage `json:"value"`
}
