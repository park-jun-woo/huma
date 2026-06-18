package runner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCaptureToString(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"string unquoted", `"hello"`, "hello"},
		{"empty string", `""`, ""},
		{"string with spaces", `"a b"`, "a b"},
		{"number int", `42`, "42"},
		{"number float", `3.14`, "3.14"},
		{"bool true", `true`, "true"},
		{"bool false", `false`, "false"},
		{"null", `null`, "null"},
		{"object kept raw", `{"k":1}`, `{"k":1}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureToString(json.RawMessage(tt.raw))
			if got != tt.want {
				t.Fatalf("captureToString(%s) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func TestMergeReportCaptures(t *testing.T) {
	// Multiple entries, multiple captures, value coercion, later overrides earlier.
	report := hurlJSONReport{
		Success: true,
		Entries: []hurlJSONEntry{
			{Captures: []hurlJSONCapture{
				{Name: "token", Value: json.RawMessage(`"abc"`)},
				{Name: "count", Value: json.RawMessage(`7`)},
			}},
			{Captures: []hurlJSONCapture{
				{Name: "flag", Value: json.RawMessage(`true`)},
				{Name: "token", Value: json.RawMessage(`"xyz"`)}, // overrides
			}},
		},
	}
	dst := map[string]string{}
	mergeReportCaptures(dst, report)
	want := map[string]string{"token": "xyz", "count": "7", "flag": "true"}
	if !reflect.DeepEqual(dst, want) {
		t.Fatalf("got %v, want %v", dst, want)
	}

	// Report with no entries leaves dst untouched.
	dst2 := map[string]string{"x": "1"}
	mergeReportCaptures(dst2, hurlJSONReport{Success: true})
	if !reflect.DeepEqual(dst2, map[string]string{"x": "1"}) {
		t.Fatalf("empty report mutated dst: %v", dst2)
	}
}

func TestCollectCaptures_SingleReport(t *testing.T) {
	out := `{"success":true,"entries":[{"captures":[{"name":"token","value":"abc"}]}]}`
	got, err := collectCaptures("setup.hurl", out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["token"] != "abc" {
		t.Fatalf("expected token=abc, got %v", got)
	}
}

func TestCollectCaptures_NoCaptures(t *testing.T) {
	out := `{"success":true,"entries":[{"captures":[]}]}`
	got, err := collectCaptures("setup.hurl", out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestCollectCaptures_MultipleReportsJSONLines(t *testing.T) {
	// JSON-lines: two report objects back to back.
	out := `{"success":true,"entries":[{"captures":[{"name":"a","value":"1"}]}]}` + "\n" +
		`{"success":true,"entries":[{"captures":[{"name":"b","value":2}]}]}`
	got, err := collectCaptures("setup.hurl", out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := map[string]string{"a": "1", "b": "2"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCollectCaptures_SuccessFalse(t *testing.T) {
	out := `{"success":false,"entries":[]}`
	_, err := collectCaptures("setup.hurl", out)
	if err == nil {
		t.Fatal("expected error when report success==false")
	}
}

func TestCollectCaptures_BadJSON(t *testing.T) {
	out := `{"success":true, not json`
	_, err := collectCaptures("setup.hurl", out)
	if err == nil {
		t.Fatal("expected parse error for malformed JSON")
	}
}

func TestCollectCaptures_EmptyOutput(t *testing.T) {
	// No JSON objects -> empty map, no error (dec.More() false immediately).
	got, err := collectCaptures("setup.hurl", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestRunJSON_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token":"jwt-123","count":5}`))
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	hurlFile := filepath.Join(tmpDir, "setup.hurl")
	content := fmt.Sprintf(
		"GET %s/login\nHTTP 200\n[Captures]\ntoken: jsonpath \"$.token\"\ncount: jsonpath \"$.count\"\n",
		ts.URL,
	)
	if err := os.WriteFile(hurlFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := RunJSON(hurlFile, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["token"] != "jwt-123" {
		t.Fatalf("expected token=jwt-123, got %v", got)
	}
	if got["count"] != "5" {
		t.Fatalf("expected count coerced to \"5\", got %q", got["count"])
	}
}

func TestRunJSON_WithVariables(t *testing.T) {
	// Exercise the variable-sorting/append branch with a variable used as base URL.
	var gotPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"ok"}`))
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	hurlFile := filepath.Join(tmpDir, "setup.hurl")
	content := "GET {{base}}/ping\nHTTP 200\n[Captures]\nid: jsonpath \"$.id\"\n"
	if err := os.WriteFile(hurlFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := RunJSON(hurlFile, map[string]string{"base": ts.URL, "extra": "z"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["id"] != "ok" {
		t.Fatalf("expected id=ok, got %v", got)
	}
	if gotPath != "/ping" {
		t.Fatalf("variable not substituted, path=%s", gotPath)
	}
}

func TestRunJSON_HurlFailure(t *testing.T) {
	// Server returns 404 but the hurl asserts 200 -> hurl exits non-zero -> error.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	hurlFile := filepath.Join(tmpDir, "setup.hurl")
	content := fmt.Sprintf("GET %s/login\nHTTP 200\n", ts.URL)
	if err := os.WriteFile(hurlFile, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := RunJSON(hurlFile, nil)
	if err == nil {
		t.Fatal("expected error when hurl assertion fails")
	}
}

func TestRunJSON_MissingFile(t *testing.T) {
	_, err := RunJSON("/nonexistent/setup.hurl", nil)
	if err == nil {
		t.Fatal("expected error for missing hurl file")
	}
}
