package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestNewPythonAdapter_Fields(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:5000",
		Server: config.ServerConfig{
			Build: "pip install -r requirements.txt",
			Start: "python app.py",
			Ready: "http://localhost:5000/health",
		},
	}

	a := NewPythonAdapter(cfg)
	if a == nil {
		t.Fatal("expected non-nil adapter")
	}
	if a.baseURL != "http://localhost:5000" {
		t.Fatalf("expected baseURL http://localhost:5000, got %s", a.baseURL)
	}
	if a.cfg != &cfg.Server {
		t.Fatal("expected cfg to point to cfg.Server")
	}
}
