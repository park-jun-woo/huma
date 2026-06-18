//ff:func feature=gate type=engine control=sequence level=error
//ff:what Phase 009 / 2-B mint: hand-rolls an HS256 JWT from testing.auth (secret_env + claims) with no external dep, returning {"token": <jwt>} for injection — the login-free, seed-free fast path
package humaquest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/park-jun-woo/huma/internal/config"
)

// mintToken signs an HS256 JWT directly from AuthConfig instead of logging in. It is
// the 2-B Phase 009 path: fastest when the secret is known, and it bypasses both the
// login endpoint and user seeding. The token header is {"alg":"HS256","typ":"JWT"},
// the payload is the configured claims (plus an exp one hour out), and it is signed
// with the env secret named by SecretEnv. Returns {"token": <jwt>}.
//
// Known limitation (documented): the claim value types and the algorithm must match
// what the app expects — claims here are emitted as strings, so an app demanding a
// numeric `sub` or `exp` may reject the token, and RS256/ES256 apps are out of scope.
// For those, prefer the 2-A capture path.
func mintToken(cfg *config.Config) (map[string]string, error) {
	secret := os.Getenv(cfg.Auth.SecretEnv)
	if cfg.Auth.SecretEnv == "" || secret == "" {
		return nil, fmt.Errorf("mint: secret env %q is unset or empty", cfg.Auth.SecretEnv)
	}

	headerSeg, err := jsonSegment(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return nil, fmt.Errorf("mint: encode header: %w", err)
	}
	payloadSeg, err := jsonSegment(buildClaims(cfg.Auth.Claims))
	if err != nil {
		return nil, fmt.Errorf("mint: encode claims: %w", err)
	}

	signingInput := headerSeg + "." + payloadSeg
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return map[string]string{"token": signingInput + "." + sig}, nil
}
