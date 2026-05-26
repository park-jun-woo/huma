package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	_, err := Load()
	if !errors.Is(err, ErrNoManifest) {
		t.Fatalf("expected ErrNoManifest, got %v", err)
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	yaml := `apiVersion: yongol/v1
kind: Project
metadata:
  name: test-project
backend:
  lang: go
  framework: gin
  module: github.com/test/project
testing:
  base_url: http://localhost:9000
  hurl_dir: tests/hurl
  server:
    build: go build -cover -o app
    start: ./app
    ready: /health
`
	os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte(yaml), 0o644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseURL != "http://localhost:9000" {
		t.Fatalf("expected http://localhost:9000, got %s", cfg.BaseURL)
	}
	if cfg.HurlDir != "tests/hurl" {
		t.Fatalf("expected tests/hurl, got %s", cfg.HurlDir)
	}
	if cfg.Server.Build != "go build -cover -o app" {
		t.Fatalf("unexpected build: %s", cfg.Server.Build)
	}
	if cfg.Scan.Lang != "go" {
		t.Fatalf("expected go, got %s", cfg.Scan.Lang)
	}
}

func TestLoad_InvalidYAML2(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte("{{invalid"), 0o644)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_PermissionDenied2(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	f := filepath.Join(tmpDir, "manifest.yaml")
	os.WriteFile(f, []byte("apiVersion: yongol/v1"), 0o644)
	os.Chmod(f, 0o000)
	t.Cleanup(func() { os.Chmod(f, 0o644) })

	_, err := Load()
	if err == nil {
		t.Fatal("expected permission error")
	}
}
