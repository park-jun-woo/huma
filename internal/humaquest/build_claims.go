//ff:func feature=gate type=helper control=iteration dimension=1
//ff:what Builds the JWT payload claim set from AuthConfig.Claims (emitted as given) plus a generated exp one hour out
package humaquest

import "time"

// buildClaims copies the configured string claims into a fresh payload map and adds
// an exp one hour from now. Claims are emitted as their given (string) values — see
// mintToken's documented type-matching limitation.
func buildClaims(configured map[string]string) map[string]any {
	claims := make(map[string]any, len(configured)+1)
	for k, v := range configured {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	return claims
}
