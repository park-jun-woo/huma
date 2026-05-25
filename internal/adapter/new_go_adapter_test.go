package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestNewGoAdapter_Fields(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:3000",
		Server: config.ServerConfig{
			Build: "go build -cover -o app ./cmd/server",
			Start: "./app",
			Ready: "http://localhost:3000/health",
			Env:   map[string]string{"GIN_MODE": "release"},
		},
	}

	a := NewGoAdapter(cfg)
	if a == nil {
		t.Fatal("expected non-nil adapter")
	}
	if a.baseURL != "http://localhost:3000" {
		t.Fatalf("expected baseURL http://localhost:3000, got %s", a.baseURL)
	}
	if a.cfg != &cfg.Server {
		t.Fatal("expected cfg to point to cfg.Server")
	}
	if a.coverDir != coverDir {
		t.Fatalf("expected coverDir %s, got %s", coverDir, a.coverDir)
	}
}
