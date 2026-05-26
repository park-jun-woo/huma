package config

import (
	"errors"
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

	_, err := Load()
	if !errors.Is(err, ErrNoManifest) {
		t.Fatalf("expected ErrNoManifest, got %v", err)
	}
}

func TestLoad_ValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	yamlContent := `apiVersion: yongol/v1
kind: Project
metadata:
  name: test
backend:
  lang: go
  framework: gin
  module: github.com/test/test
testing:
  base_url: http://localhost:3000
  hurl_dir: tests/hurl
  server:
    build: go build -cover -o app
    start: ./app
    ready: /health
`
	os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte(yamlContent), 0o644)

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

	f := filepath.Join(tmpDir, "manifest.yaml")
	os.WriteFile(f, []byte("apiVersion: test"), 0o644)
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

	os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte(":\ninvalid:\n  - [broken"), 0o644)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	yamlContent := `apiVersion: yongol/v1
kind: Project
metadata:
  name: test
backend: {}
testing: {}
`
	os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte(yamlContent), 0o644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseURL != "http://localhost:8080" {
		t.Fatalf("expected default BaseURL, got %s", cfg.BaseURL)
	}
	if cfg.HurlDir != "hurl" {
		t.Fatalf("expected default HurlDir, got %s", cfg.HurlDir)
	}
	if cfg.Scan.Lang != "go" {
		t.Fatalf("expected default lang go, got %s", cfg.Scan.Lang)
	}
}

func TestLoad_DepsConfig(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	yamlContent := `apiVersion: yongol/v1
kind: Project
metadata:
  name: test
backend:
  lang: go
testing:
  deps:
    up: "docker compose up -d"
    down: "docker compose down"
    ready: "localhost:5432"
  server:
    start: "./app"
    ready: "/health"
`
	os.WriteFile(filepath.Join(tmpDir, "manifest.yaml"), []byte(yamlContent), 0o644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Deps.Up != "docker compose up -d" {
		t.Fatalf("expected deps.up, got %s", cfg.Deps.Up)
	}
	if cfg.Deps.Down != "docker compose down" {
		t.Fatalf("expected deps.down, got %s", cfg.Deps.Down)
	}
}
