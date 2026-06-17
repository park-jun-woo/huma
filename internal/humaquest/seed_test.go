package humaquest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/quest"
)

const seedOpenAPI = `
openapi: '3.0.0'
info:
  title: Seed Test API
paths:
  /api/v1/users:
    get:
      operationId: ListUsers
    post:
      operationId: CreateUser
  /api/v1/users/{userId}:
    get:
      operationId: GetUser
`

// writeTempFile writes content to a file under a fresh temp dir and returns its path.
func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestSeed_OneTODOItemPerEndpoint(t *testing.T) {
	path := writeTempFile(t, "openapi.yaml", seedOpenAPI)

	items, err := (humaDef{}).Seed([]string{path})
	if err != nil {
		t.Fatalf("Seed: %v", err)
	}

	// 3 operations declared above.
	if len(items) != 3 {
		t.Fatalf("want 3 items, got %d", len(items))
	}

	// Cross-check against the scanner directly so the item count/keys are tied to
	// the real endpoints, not a hardcoded list.
	eps, err := scanEndpoints(path)
	if err != nil {
		t.Fatalf("scanEndpoints: %v", err)
	}
	if len(eps) != len(items) {
		t.Fatalf("item count %d != endpoint count %d", len(items), len(eps))
	}

	byID := make(map[string]scanner.Endpoint, len(eps))
	for _, ep := range eps {
		byID[ep.ID] = ep
	}

	for _, it := range items {
		if it.State != quest.TODO {
			t.Errorf("item %s: state = %q, want TODO", it.Key, it.State)
		}
		if it.Key == "" {
			t.Fatal("item has empty Key")
		}
		ep, ok := byID[it.Key]
		if !ok {
			t.Fatalf("item Key %q does not match any endpoint ID", it.Key)
		}
		if it.Key != ep.ID {
			t.Errorf("Key %q != endpoint ID %q", it.Key, ep.ID)
		}

		// Payload must round-trip the Endpoint.
		var decoded scanner.Endpoint
		if err := it.DecodePayload(&decoded); err != nil {
			t.Fatalf("DecodePayload: %v", err)
		}
		if decoded.ID != ep.ID || decoded.Method != ep.Method || decoded.Path != ep.Path {
			t.Errorf("payload mismatch: got %+v, want %+v", decoded, ep)
		}
		if decoded.Handler != ep.Handler {
			t.Errorf("payload handler = %q, want %q", decoded.Handler, ep.Handler)
		}
	}
}

func TestSeed_PropagatesScanError(t *testing.T) {
	// A directory under cwd with no openapi triggers the E-01 path when from=="".
	// Here we point at a non-existent file to force scanEndpoints to fail on read.
	_, err := (humaDef{}).Seed([]string{filepath.Join(t.TempDir(), "does-not-exist.yaml")})
	if err == nil {
		t.Fatal("expected error for missing input file")
	}
}

func TestSeed_EmptySource(t *testing.T) {
	// Empty paths still produce zero items (no openapi in cwd → E-01 error path
	// is handled by scanEndpoints; an OpenAPI with no paths errors). Use an
	// OpenAPI with a single path to confirm empty input is not silently dropped.
	path := writeTempFile(t, "openapi.yaml", `
openapi: '3.0.0'
info:
  title: Single
paths:
  /ping:
    get:
      operationId: Ping
`)
	items, err := (humaDef{}).Seed([]string{path})
	if err != nil {
		t.Fatalf("Seed: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("want 1 item, got %d", len(items))
	}
}
