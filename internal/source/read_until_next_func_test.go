package source

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadUntilNextFunc_StopsAtFunc(t *testing.T) {
	input := "func GetUser() {\n\tx := 1\n}\n\nfunc Other() {}\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Scan() // Read first line "func GetUser() {"

	lines, err := readUntilNextFunc(scanner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should include "func GetUser() {", "	x := 1", "}", ""
	// but stop before "func Other() {}"
	if len(lines) == 0 {
		t.Fatal("expected non-empty lines")
	}
	for _, l := range lines {
		if strings.HasPrefix(l, "func Other") {
			t.Fatal("should not include Other func")
		}
	}
}

func TestReadUntilNextFunc_ScannerError(t *testing.T) {
	longLine := make([]byte, 1024*1024)
	for i := range longLine {
		longLine[i] = 'x'
	}
	input := "func GetUser() {\n" + string(longLine) + "\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Scan() // first line

	_, err := readUntilNextFunc(scanner)
	if err == nil {
		t.Fatal("expected scanner error")
	}
}

func TestReadUntilNextFunc_ReadsToEOF(t *testing.T) {
	input := "func GetUser() {\n\tx := 1\n}\n"
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Scan()

	lines, err := readUntilNextFunc(scanner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 lines, got %d", len(lines))
	}
}
