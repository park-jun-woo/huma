package scanner

import (
	"encoding/json"
	"testing"
)

func TestExtractOpenAPIResponses(t *testing.T) {
	op := map[string]interface{}{
		"responses": map[string]interface{}{
			"200": map[string]interface{}{"description": "OK"},
			"400": map[string]interface{}{"description": "Bad Request"},
			"401": map[string]interface{}{"description": "Unauthorized"},
		},
	}

	got := extractOpenAPIResponses(op)
	if got == nil {
		t.Fatal("expected non-nil responses")
	}

	type entry struct {
		Status int `json:"status"`
	}
	var entries []entry
	if err := json.Unmarshal(got, &entries); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Status != 200 || entries[1].Status != 400 || entries[2].Status != 401 {
		t.Fatalf("unexpected status codes: %+v", entries)
	}
}

func TestExtractOpenAPIResponses_SkipsDefault(t *testing.T) {
	op := map[string]interface{}{
		"responses": map[string]interface{}{
			"200":     map[string]interface{}{"description": "OK"},
			"default": map[string]interface{}{"description": "Error"},
		},
	}

	got := extractOpenAPIResponses(op)
	if got == nil {
		t.Fatal("expected non-nil responses")
	}

	type entry struct {
		Status int `json:"status"`
	}
	var entries []entry
	if err := json.Unmarshal(got, &entries); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != 200 {
		t.Fatalf("expected 200, got %d", entries[0].Status)
	}
}

func TestExtractOpenAPIResponses_NoResponses(t *testing.T) {
	op := map[string]interface{}{
		"operationId": "Test",
	}

	got := extractOpenAPIResponses(op)
	if got != nil {
		t.Fatalf("expected nil, got %s", string(got))
	}
}

func TestExtractOpenAPIResponses_AllNonNumeric(t *testing.T) {
	op := map[string]interface{}{
		"responses": map[string]interface{}{
			"default": map[string]interface{}{"description": "Error"},
			"2XX":     map[string]interface{}{"description": "Success"},
		},
	}

	got := extractOpenAPIResponses(op)
	if got != nil {
		t.Fatalf("expected nil for all non-numeric, got %s", string(got))
	}
}
