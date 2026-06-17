package humaquest

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
	"github.com/park-jun-woo/reins/pkg/quest"
)

// sampleEndpoint is a small, fully-populated endpoint used to drive Render.
func sampleEndpoint() scanner.Endpoint {
	return scanner.Endpoint{
		ID:      "GET_/api/v1/users",
		Method:  "GET",
		Path:    "/api/v1/users",
		Handler: "ListUsers",
		Source:  "handlers/users.go",
		Line:    42,
	}
}

// renderItem builds an Item with the sample endpoint payload plus the given
// attempt log, ready to feed to Render.
func renderItem(t *testing.T, tries int, attempts ...quest.Attempt) *quest.Item {
	t.Helper()
	it := &quest.Item{Key: "GET_/api/v1/users", State: quest.TODO, Tries: tries, Log: attempts}
	ep := sampleEndpoint()
	if err := it.SetPayload(&ep); err != nil {
		t.Fatalf("SetPayload: %v", err)
	}
	return it
}

// chdir switches into dir for the duration of the test, restoring cwd after.
func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
}

// liveManifest writes a manifest.yaml with a server.start command (live mode).
const liveManifest = `apiVersion: yongol/v1
kind: Project
metadata:
  name: render-test
backend:
  lang: go
  framework: gin
  module: github.com/test/test
testing:
  base_url: http://localhost:8080
  hurl_dir: hurl
  server:
    build: go build -cover -o app
    start: ./app
    ready: /health
`

func TestRender_PhaseHeaders(t *testing.T) {
	tests := []struct {
		name       string
		manifest   string // "" → no manifest.yaml (static fallback via ErrNoManifest)
		attempts   []quest.Attempt
		wantHeader string
	}{
		{
			name:       "fresh item, no manifest → static TODO",
			manifest:   "",
			attempts:   nil,
			wantHeader: "# TODO  GET /api/v1/users",
		},
		{
			name:       "fresh item, live manifest → live TODO",
			manifest:   liveManifest,
			attempts:   nil,
			wantHeader: "# TODO  GET /api/v1/users",
		},
		{
			name:     "last attempt FAIL → IMPROVE",
			manifest: "",
			attempts: []quest.Attempt{
				{Try: 1, Outcome: string(quest.OutFail), Reason: "R1: status 404 uncovered"},
			},
			wantHeader: "# IMPROVE  GET /api/v1/users",
		},
		{
			name:     "last attempt REVIEW → UNVERIFIED",
			manifest: "",
			attempts: []quest.Attempt{
				{Try: 1, Outcome: string(quest.OutReview), Reason: "no oracle"},
			},
			wantHeader: "# UNVERIFIED  GET /api/v1/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			if tt.manifest != "" {
				if err := os.WriteFile(filepath.Join(dir, "manifest.yaml"), []byte(tt.manifest), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			chdir(t, dir)

			tries := 0
			if len(tt.attempts) > 0 {
				tries = len(tt.attempts)
			}
			it := renderItem(t, tries, tt.attempts...)

			out, err := (humaDef{}).Render(nil, it)
			if err != nil {
				t.Fatalf("Render: %v", err)
			}
			if !strings.HasPrefix(out, tt.wantHeader) {
				t.Errorf("prompt header = %q..., want prefix %q", firstLine(out), tt.wantHeader)
			}
		})
	}
}

// TestRender_ImproveCarriesReason proves the IMPROVE prompt feeds back the last
// attempt's Reason (recovered via lastReason).
func TestRender_ImproveCarriesReason(t *testing.T) {
	chdir(t, t.TempDir())
	it := renderItem(t, 1, quest.Attempt{
		Try: 1, Outcome: string(quest.OutFail), Reason: "R7: missing 422 branch",
	})

	out, err := (humaDef{}).Render(nil, it)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(out, "R7: missing 422 branch") {
		t.Errorf("IMPROVE prompt does not carry the last reason; got:\n%s", out)
	}
}

// TestRender_ReadOnly locks in the invariant that Render never mutates the Item.
func TestRender_ReadOnly(t *testing.T) {
	chdir(t, t.TempDir())
	it := renderItem(t, 2,
		quest.Attempt{Try: 1, Outcome: string(quest.OutFail), Reason: "first"},
		quest.Attempt{Try: 2, Outcome: string(quest.OutReview), Reason: "second"},
	)

	beforeTries := it.Tries
	beforeLog := append([]quest.Attempt(nil), it.Log...)
	beforePayload := append(json.RawMessage(nil), it.Payload...)

	if _, err := (humaDef{}).Render(nil, it); err != nil {
		t.Fatalf("Render: %v", err)
	}

	if it.Tries != beforeTries {
		t.Errorf("Render mutated Tries: %d → %d", beforeTries, it.Tries)
	}
	if len(it.Log) != len(beforeLog) {
		t.Fatalf("Render mutated Log length: %d → %d", len(beforeLog), len(it.Log))
	}
	for i := range beforeLog {
		if it.Log[i] != beforeLog[i] {
			t.Errorf("Render mutated Log[%d]: %+v → %+v", i, beforeLog[i], it.Log[i])
		}
	}
	if string(it.Payload) != string(beforePayload) {
		t.Errorf("Render mutated Payload:\n before: %s\n after:  %s", beforePayload, it.Payload)
	}
}

// TestRender_DecodePayloadError exercises the error branch where the Item's
// payload is not a valid Endpoint JSON.
func TestRender_DecodePayloadError(t *testing.T) {
	chdir(t, t.TempDir())
	it := &quest.Item{
		Key:     "bad",
		State:   quest.TODO,
		Payload: json.RawMessage(`{"id": 12345`), // truncated → invalid JSON
	}

	if _, err := (humaDef{}).Render(nil, it); err == nil {
		t.Fatal("expected DecodePayload error, got nil")
	}
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
