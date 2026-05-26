package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindOpenAPIFile(t *testing.T) {
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	t.Run("finds openapi.yaml in current dir", func(t *testing.T) {
		dir := t.TempDir()
		os.Chdir(dir)

		os.WriteFile("openapi.yaml", []byte("openapi: '3.0.0'"), 0644)

		got := FindOpenAPIFile()
		if got != "openapi.yaml" {
			t.Fatalf("expected openapi.yaml, got %q", got)
		}
	})

	t.Run("finds api/openapi.yaml", func(t *testing.T) {
		dir := t.TempDir()
		os.Chdir(dir)

		os.MkdirAll("api", 0755)
		os.WriteFile(filepath.Join("api", "openapi.yaml"), []byte("openapi: '3.0.0'"), 0644)

		got := FindOpenAPIFile()
		if got != "api/openapi.yaml" {
			t.Fatalf("expected api/openapi.yaml, got %q", got)
		}
	})

	t.Run("returns empty when no file found", func(t *testing.T) {
		dir := t.TempDir()
		os.Chdir(dir)

		got := FindOpenAPIFile()
		if got != "" {
			t.Fatalf("expected empty string, got %q", got)
		}
	})

	t.Run("prefers openapi.yaml over api/openapi.yaml", func(t *testing.T) {
		dir := t.TempDir()
		os.Chdir(dir)

		os.WriteFile("openapi.yaml", []byte("openapi: '3.0.0'"), 0644)
		os.MkdirAll("api", 0755)
		os.WriteFile(filepath.Join("api", "openapi.yaml"), []byte("openapi: '3.0.0'"), 0644)

		got := FindOpenAPIFile()
		if got != "openapi.yaml" {
			t.Fatalf("expected openapi.yaml (first priority), got %q", got)
		}
	})
}
