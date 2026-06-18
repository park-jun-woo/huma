package humaquest

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestBuildClaims(t *testing.T) {
	// Empty configured map -> only exp.
	got := buildClaims(nil)
	if _, ok := got["exp"]; !ok {
		t.Fatal("expected exp claim")
	}
	if len(got) != 1 {
		t.Fatalf("expected only exp for nil input, got %d keys", len(got))
	}
	if _, ok := got["exp"].(int64); !ok {
		t.Fatalf("expected exp to be int64, got %T", got["exp"])
	}

	// Configured claims are copied as-is (string values) plus exp.
	configured := map[string]string{"sub": "alice", "role": "admin"}
	got = buildClaims(configured)
	if got["sub"] != "alice" || got["role"] != "admin" {
		t.Fatalf("configured claims not copied: %v", got)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 keys (sub, role, exp), got %d: %v", len(got), got)
	}
	// Mutating output must not affect input.
	got["sub"] = "mallory"
	if configured["sub"] != "alice" {
		t.Fatal("buildClaims must not mutate the input map")
	}
}

func TestJSONSegment(t *testing.T) {
	seg, err := jsonSegment(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Must be base64url no-padding decodable back to the JSON.
	raw, err := base64.RawURLEncoding.DecodeString(seg)
	if err != nil {
		t.Fatalf("segment not base64url-no-pad: %v", err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("segment not JSON: %v", err)
	}
	if decoded["alg"] != "HS256" || decoded["typ"] != "JWT" {
		t.Fatalf("unexpected decoded segment: %v", decoded)
	}
	if strings.Contains(seg, "=") {
		t.Fatal("segment must not contain padding")
	}
}

func TestJSONSegment_MarshalError(t *testing.T) {
	// channels are not JSON-marshalable -> error branch.
	_, err := jsonSegment(make(chan int))
	if err == nil {
		t.Fatal("expected marshal error for unmarshalable value")
	}
}

func TestSortedKeys(t *testing.T) {
	got := sortedKeys(map[string]string{"c": "1", "a": "2", "b": "3"})
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	if !sort.StringsAreSorted(got) {
		t.Fatal("keys not sorted")
	}
	// Empty map -> empty (non-nil) slice.
	if k := sortedKeys(map[string]string{}); len(k) != 0 {
		t.Fatalf("expected empty, got %v", k)
	}
}

func TestMergeVars(t *testing.T) {
	tests := []struct {
		name string
		base map[string]string
		over map[string]string
		want map[string]string
	}{
		{"both nil", nil, nil, map[string]string{}},
		{"base only", map[string]string{"a": "1"}, nil, map[string]string{"a": "1"}},
		{"over only", nil, map[string]string{"b": "2"}, map[string]string{"b": "2"}},
		{
			"over wins on conflict",
			map[string]string{"token": "stale", "x": "keep"},
			map[string]string{"token": "fresh"},
			map[string]string{"token": "fresh", "x": "keep"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeVars(tt.base, tt.over)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}

	// Inputs are not mutated.
	base := map[string]string{"a": "1"}
	over := map[string]string{"a": "2"}
	out := mergeVars(base, over)
	out["a"] = "999"
	if base["a"] != "1" || over["a"] != "2" {
		t.Fatal("mergeVars must not mutate inputs")
	}
}

func TestMintToken_SecretUnset(t *testing.T) {
	// SecretEnv empty.
	cfg := &config.Config{}
	if _, err := mintToken(cfg); err == nil {
		t.Fatal("expected error when SecretEnv is empty")
	}

	// SecretEnv set but env var unset/empty.
	cfg = &config.Config{}
	cfg.Auth.SecretEnv = "HUMA_TEST_MINT_SECRET_UNSET"
	t.Setenv("HUMA_TEST_MINT_SECRET_UNSET", "")
	if _, err := mintToken(cfg); err == nil {
		t.Fatal("expected error when secret env value is empty")
	}
}

func TestMintToken_Success(t *testing.T) {
	const secret = "topsecret-hs256-key"
	cfg := &config.Config{}
	cfg.Auth.SecretEnv = "HUMA_TEST_MINT_SECRET"
	cfg.Auth.Claims = map[string]string{"sub": "alice"}
	t.Setenv("HUMA_TEST_MINT_SECRET", secret)

	out, err := mintToken(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tok, ok := out["token"]
	if !ok {
		t.Fatal("expected token key")
	}

	parts := strings.Split(tok, ".")
	if len(parts) != 3 {
		t.Fatalf("expected 3-segment JWT, got %d segments", len(parts))
	}

	// Decode and verify header.
	hRaw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		t.Fatalf("header not base64url: %v", err)
	}
	var header map[string]string
	if err := json.Unmarshal(hRaw, &header); err != nil {
		t.Fatalf("header not JSON: %v", err)
	}
	if header["alg"] != "HS256" || header["typ"] != "JWT" {
		t.Fatalf("unexpected header: %v", header)
	}

	// Decode payload, check claim + exp present.
	pRaw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("payload not base64url: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(pRaw, &payload); err != nil {
		t.Fatalf("payload not JSON: %v", err)
	}
	if payload["sub"] != "alice" {
		t.Fatalf("expected sub=alice, got %v", payload["sub"])
	}
	if _, ok := payload["exp"]; !ok {
		t.Fatal("expected exp in payload")
	}

	// Verify the HS256 signature against the secret.
	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signingInput))
	wantSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if parts[2] != wantSig {
		t.Fatalf("signature mismatch:\ngot  %s\nwant %s", parts[2], wantSig)
	}

	// A wrong secret must NOT verify (proves the signature is meaningful).
	badMac := hmac.New(sha256.New, []byte("wrong-secret"))
	badMac.Write([]byte(signingInput))
	badSig := base64.RawURLEncoding.EncodeToString(badMac.Sum(nil))
	if parts[2] == badSig {
		t.Fatal("signature verified against the wrong secret")
	}
}
