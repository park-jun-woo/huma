package config

import "testing"

func TestDefaultConfig(t *testing.T) {
	cfg := Default()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.BaseURL != "http://localhost:8080" {
		t.Fatalf("expected base_url http://localhost:8080, got %s", cfg.BaseURL)
	}
	if cfg.HurlDir != "hurl" {
		t.Fatalf("expected hurl_dir 'hurl', got %s", cfg.HurlDir)
	}
	if cfg.HurlVariables["base_url"] != "http://localhost:8080" {
		t.Fatalf("expected hurl variable base_url, got %v", cfg.HurlVariables)
	}
	if cfg.Scan.Lang != "go" {
		t.Fatalf("expected scan.lang 'go', got %s", cfg.Scan.Lang)
	}
}
