package adapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestRustWaitReady_NoReadyURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 2s sleep test")
	}
	a := &RustAdapter{cfg: &config.ServerConfig{Ready: ""}}
	err := a.WaitReady()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRustWaitReady_ServerReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := &RustAdapter{cfg: &config.ServerConfig{Ready: ts.URL}}
	err := a.WaitReady()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
