package analyzer

import (
	"encoding/json"
	"testing"
)

func TestParseResponses_Valid(t *testing.T) {
	data := json.RawMessage(`[
		{"status": 201, "line": 70, "code": "c.JSON(http.StatusCreated, body)"},
		{"status": 400, "line": 55},
		{"status": 409, "line": 62}
	]`)

	branches := ParseResponses(data, "handler.go")
	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}
	if branches[0].Status != 201 {
		t.Fatalf("expected 201, got %d", branches[0].Status)
	}
	if branches[0].File != "handler.go" {
		t.Fatalf("expected handler.go, got %s", branches[0].File)
	}
	if branches[0].Line != 70 {
		t.Fatalf("expected line 70, got %d", branches[0].Line)
	}
	if branches[0].Code != "c.JSON(http.StatusCreated, body)" {
		t.Fatalf("unexpected code: %s", branches[0].Code)
	}
}

func TestParseResponses_Empty(t *testing.T) {
	branches := ParseResponses(nil, "handler.go")
	if branches != nil {
		t.Fatalf("expected nil, got %v", branches)
	}
}

func TestParseResponses_InvalidJSON(t *testing.T) {
	data := json.RawMessage(`not json`)
	branches := ParseResponses(data, "handler.go")
	if branches != nil {
		t.Fatalf("expected nil, got %v", branches)
	}
}

func TestParseResponses_SkipZeroStatus(t *testing.T) {
	data := json.RawMessage(`[{"status": 0, "line": 10}, {"status": 200, "line": 20}]`)
	branches := ParseResponses(data, "handler.go")
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 200 {
		t.Fatalf("expected 200, got %d", branches[0].Status)
	}
}
