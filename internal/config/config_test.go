package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.BaseURL != "http://localhost:8080" {
		t.Fatalf("expected BaseURL http://localhost:8080, got %s", cfg.BaseURL)
	}
	if cfg.HurlDir != "hurl" {
		t.Fatalf("expected HurlDir hurl, got %s", cfg.HurlDir)
	}
	if cfg.Scan.Lang != "go" {
		t.Fatalf("expected Scan.Lang go, got %s", cfg.Scan.Lang)
	}
}

func TestLoad_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should return defaults
	if cfg.BaseURL != "http://localhost:8080" {
		t.Fatalf("expected default BaseURL, got %s", cfg.BaseURL)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	yamlContent := `base_url: http://localhost:3000
hurl_dir: tests/hurl
scan:
  lang: go
server:
  build: go build -cover -o app
  start: ./app
  ready: http://localhost:3000/health
`
	os.WriteFile(filepath.Join(tmpDir, "huma.yaml"), []byte(yamlContent), 0o644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseURL != "http://localhost:3000" {
		t.Fatalf("expected http://localhost:3000, got %s", cfg.BaseURL)
	}
	if cfg.HurlDir != "tests/hurl" {
		t.Fatalf("expected tests/hurl, got %s", cfg.HurlDir)
	}
	if cfg.Server.Build != "go build -cover -o app" {
		t.Fatalf("unexpected build: %s", cfg.Server.Build)
	}
}

func TestLoad_PermissionDenied(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	f := filepath.Join(tmpDir, "huma.yaml")
	os.WriteFile(f, []byte("base_url: test"), 0o644)
	os.Chmod(f, 0o000)
	t.Cleanup(func() { os.Chmod(f, 0o644) })

	_, err := Load()
	if err == nil {
		t.Fatal("expected permission error")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, "huma.yaml"), []byte(":\ninvalid:\n  - [broken"), 0o644)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}
