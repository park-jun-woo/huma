package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestNewDenoAdapter_Fields(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:54321",
		Server: config.ServerConfig{
			Start: "supabase functions serve",
			Ready: "/functions/v1/health",
			Env:   map[string]string{"DENO_ENV": "test"},
		},
	}

	a := NewDenoAdapter(cfg)
	if a == nil {
		t.Fatal("expected non-nil adapter")
	}
	if a.baseURL != "http://localhost:54321" {
		t.Fatalf("expected baseURL http://localhost:54321, got %s", a.baseURL)
	}
	if a.cfg != &cfg.Server {
		t.Fatal("expected cfg to point to cfg.Server")
	}
}
