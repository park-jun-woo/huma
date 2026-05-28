package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestNewRustAdapter_Fields(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
		Server: config.ServerConfig{
			Build: "cargo build",
			Start: "./target/debug/my-app",
			Ready: "/health",
			Env:   map[string]string{"RUST_LOG": "info"},
		},
	}

	a := NewRustAdapter(cfg)
	if a == nil {
		t.Fatal("expected non-nil adapter")
	}
	if a.baseURL != "http://localhost:8080" {
		t.Fatalf("expected baseURL http://localhost:8080, got %s", a.baseURL)
	}
	if a.cfg != &cfg.Server {
		t.Fatal("expected cfg to point to cfg.Server")
	}
}
