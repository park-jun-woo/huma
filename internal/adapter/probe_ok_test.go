package adapter

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbeOK_Returns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	if !probeOK(http.DefaultClient, ts.URL) {
		t.Fatal("expected true for 200 response")
	}
}

func TestProbeOK_Returns500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	if probeOK(http.DefaultClient, ts.URL) {
		t.Fatal("expected false for 500 response")
	}
}

func TestProbeOK_ConnectionRefused(t *testing.T) {
	if probeOK(http.DefaultClient, "http://localhost:19999") {
		t.Fatal("expected false for connection refused")
	}
}
