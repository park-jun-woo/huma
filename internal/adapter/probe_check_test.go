package adapter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbeCheck(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	if !ProbeCheck(ts.URL) {
		t.Fatal("expected true for 200")
	}
	// unreachable URL → false
	if ProbeCheck("http://127.0.0.1:1") {
		t.Fatal("expected false for unreachable server")
	}
}

func TestIsLineCovered(t *testing.T) {
	var nilR *CoverageResult
	if nilR.IsLineCovered(5) {
		t.Error("nil receiver → false")
	}
	if (&CoverageResult{CoveredLines: nil}).IsLineCovered(5) {
		t.Error("nil map → false")
	}
	r := &CoverageResult{CoveredLines: map[int]bool{5: true}}
	if !r.IsLineCovered(5) {
		t.Error("covered line → true")
	}
	if r.IsLineCovered(6) {
		t.Error("uncovered line → false")
	}
}
