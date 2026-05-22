//ff:func feature=scan type=helper control=sequence
//ff:what Generates a deterministic SHA-256-based ID from method and path
package scanner

import (
	"crypto/sha256"
	"fmt"
)

func makeID(method, path string) string {
	h := sha256.Sum256([]byte(method + " " + path))
	return fmt.Sprintf("%x", h[:8])
}
