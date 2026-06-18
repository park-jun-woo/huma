//ff:func feature=gate type=helper control=sequence
//ff:what Marshals a value to JSON and base64url-encodes it (no padding) as a JWT segment
package humaquest

import (
	"encoding/base64"
	"encoding/json"
)

// jsonSegment marshals v to JSON and base64url-encodes it (no padding) as a JWT
// header or payload segment.
func jsonSegment(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
