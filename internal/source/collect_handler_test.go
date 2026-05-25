package source

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestCollectHandler_Found(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "handler.go")
	content := `package main

func GetUser(c interface{}) {
	x := 1
	_ = x
}

func Other() {}
`
	os.WriteFile(file, []byte(content), 0o644)

	f, _ := os.Open(file)
	defer f.Close()

	pattern := regexp.MustCompile(`^func\s+GetUser\b`)
	lines, startLine, err := collectHandler(f, pattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startLine == 0 {
		t.Fatal("expected non-zero start line")
	}
	if len(lines) == 0 {
		t.Fatal("expected non-empty lines")
	}
}

func TestCollectHandler_SeekError(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "handler.go")
	// Create a file with a line longer than 64KB to trigger scanner error
	longLine := make([]byte, 1024*1024)
	for i := range longLine {
		longLine[i] = 'x'
	}
	os.WriteFile(file, longLine, 0o644)

	f, _ := os.Open(file)
	defer f.Close()

	pattern := regexp.MustCompile(`^func\s+GetUser\b`)
	_, _, err := collectHandler(f, pattern)
	if err == nil {
		t.Fatal("expected error from scanner")
	}
}

func TestCollectHandler_ReadUntilError(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "handler.go")
	// Handler line first, then a very long line to trigger scanner error during readUntilNextFunc
	longLine := make([]byte, 1024*1024)
	for i := range longLine {
		longLine[i] = 'x'
	}
	content := append([]byte("func GetUser(c interface{}) {\n"), longLine...)
	os.WriteFile(file, content, 0o644)

	f, _ := os.Open(file)
	defer f.Close()

	pattern := regexp.MustCompile(`^func\s+GetUser\b`)
	_, _, err := collectHandler(f, pattern)
	if err == nil {
		t.Fatal("expected error from readUntilNextFunc scanner")
	}
}

func TestCollectHandler_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "handler.go")
	content := `package main

func Other() {}
`
	os.WriteFile(file, []byte(content), 0o644)

	f, _ := os.Open(file)
	defer f.Close()

	pattern := regexp.MustCompile(`^func\s+GetUser\b`)
	lines, startLine, err := collectHandler(f, pattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startLine != 0 {
		t.Fatal("expected 0 start line for not found")
	}
	if lines != nil {
		t.Fatal("expected nil lines")
	}
}
