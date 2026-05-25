package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BaseURL != "http://localhost:8080" {
		t.Fatalf("expected default base_url, got %s", cfg.BaseURL)
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	yaml := `base_url: http://localhost:9000
hurl_dir: tests/hurl
`
	os.WriteFile(filepath.Join(tmpDir, "huma.yaml"), []byte(yaml), 0o644)

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
}

func TestLoad_InvalidYAML2(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	os.WriteFile(filepath.Join(tmpDir, "huma.yaml"), []byte("{{invalid"), 0o644)

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

	f := filepath.Join(tmpDir, "huma.yaml")
	os.WriteFile(f, []byte("base_url: test"), 0o644)
	os.Chmod(f, 0o000)
	t.Cleanup(func() { os.Chmod(f, 0o644) })

	_, err := Load()
	if err == nil {
		t.Fatal("expected permission error")
	}
}
