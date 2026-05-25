package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestNewNodeAdapter_Fields(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:4000",
		Server: config.ServerConfig{
			Build: "npm run build",
			Start: "node server.js",
			Ready: "http://localhost:4000/health",
		},
	}

	a := NewNodeAdapter(cfg)
	if a == nil {
		t.Fatal("expected non-nil adapter")
	}
	if a.baseURL != "http://localhost:4000" {
		t.Fatalf("expected baseURL http://localhost:4000, got %s", a.baseURL)
	}
	if a.cfg != &cfg.Server {
		t.Fatal("expected cfg to point to cfg.Server")
	}
	if a.coverDir != nodeCoverDir {
		t.Fatalf("expected coverDir %s, got %s", nodeCoverDir, a.coverDir)
	}
}
