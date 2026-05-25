package config

import "testing"

func TestURLVar_Found(t *testing.T) {
	cfg := &Config{
		BaseURL:       "http://localhost:8080",
		HurlVariables: map[string]string{"host": "http://localhost:8080"},
	}
	got := cfg.URLVar()
	if got != "host" {
		t.Fatalf("expected 'host', got %q", got)
	}
}

func TestURLVar_NotFound(t *testing.T) {
	cfg := &Config{
		BaseURL:       "http://localhost:8080",
		HurlVariables: map[string]string{"other": "http://example.com"},
	}
	got := cfg.URLVar()
	if got != "base_url" {
		t.Fatalf("expected 'base_url', got %q", got)
	}
}

func TestURLVar_EmptyVars(t *testing.T) {
	cfg := &Config{
		BaseURL:       "http://localhost:8080",
		HurlVariables: nil,
	}
	got := cfg.URLVar()
	if got != "base_url" {
		t.Fatalf("expected 'base_url', got %q", got)
	}
}
