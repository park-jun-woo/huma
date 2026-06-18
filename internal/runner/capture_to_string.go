//ff:func feature=runner type=helper control=sequence
//ff:what Coerces a hurl capture value (JSON string/number/bool) to its hurl-variable string form: strings are unquoted, everything else keeps its raw JSON text
package runner

import (
	"encoding/json"
	"strconv"
)

// captureToString coerces a hurl capture value (string, number, or bool) to its
// hurl-variable string form. A JSON string is unquoted; everything else (number,
// bool, null, object) keeps its raw JSON text.
func captureToString(raw json.RawMessage) string {
	s := string(raw)
	if unq, err := strconv.Unquote(s); err == nil {
		return unq
	}
	return s
}
