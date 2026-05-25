package source

import (
	"bufio"
	"regexp"
	"strings"
	"testing"
)

func TestSeekToHandler_Found(t *testing.T) {
	input := "package main\n\nfunc Handler() {\n}\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	pattern := regexp.MustCompile(`^func\s+Handler\b`)

	line, err := seekToHandler(scanner, pattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line != 3 {
		t.Fatalf("expected line 3, got %d", line)
	}
}

func TestSeekToHandler_NotFound(t *testing.T) {
	input := "package main\n\nfunc Other() {\n}\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	pattern := regexp.MustCompile(`^func\s+Handler\b`)

	line, err := seekToHandler(scanner, pattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line != 0 {
		t.Fatalf("expected 0, got %d", line)
	}
}

func TestSeekToHandler_ScannerError(t *testing.T) {
	longLine := make([]byte, 1024*1024)
	for i := range longLine {
		longLine[i] = 'x'
	}
	scanner := bufio.NewScanner(strings.NewReader(string(longLine)))
	pattern := regexp.MustCompile(`^func\s+Handler\b`)

	_, err := seekToHandler(scanner, pattern)
	if err == nil {
		t.Fatal("expected scanner error")
	}
}
