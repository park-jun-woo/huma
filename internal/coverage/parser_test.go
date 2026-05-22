package coverage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCoverageFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage.out")
	content := `mode: atomic
github.com/example/handler.go:10.2,15.4 2 1
github.com/example/handler.go:16.5,20.3 1 0
`
	os.WriteFile(f, []byte(content), 0o644)

	blocks, err := ParseCoverageFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}
	if blocks[0].StartLine != 10 || blocks[0].EndLine != 15 || blocks[0].Count != 1 {
		t.Fatalf("unexpected block 0: %+v", blocks[0])
	}
	if blocks[1].StartLine != 16 || blocks[1].EndLine != 20 || blocks[1].Count != 0 {
		t.Fatalf("unexpected block 1: %+v", blocks[1])
	}
}

func TestParseCoverageFile_FileNotFound(t *testing.T) {
	_, err := ParseCoverageFile("/nonexistent/coverage.out")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseCoverageFile_MalformedLines(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage.out")
	content := `mode: set
this-is-not-valid
github.com/example/handler.go:10.2,15.4 2 1
also-bad-line
`
	os.WriteFile(f, []byte(content), 0o644)

	blocks, err := ParseCoverageFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only the valid line should be parsed
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
}

func TestParseCoverageFile_ScannerError(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage.out")
	// Create a file with a line that exceeds bufio.Scanner's default buffer size (64KB)
	longLine := make([]byte, 1024*1024) // 1MB line
	for i := range longLine {
		longLine[i] = 'x'
	}
	os.WriteFile(f, longLine, 0o644)

	_, err := ParseCoverageFile(f)
	if err == nil {
		t.Fatal("expected scanner error for extremely long line")
	}
}

func TestParseCoverageFile_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "coverage.out")
	os.WriteFile(f, []byte(""), 0o644)

	blocks, err := ParseCoverageFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blocks) != 0 {
		t.Fatalf("expected 0 blocks, got %d", len(blocks))
	}
}

func TestFilterBlocks(t *testing.T) {
	blocks := []Block{
		{File: "github.com/example/handler.go", StartLine: 10, EndLine: 15, Count: 1},
		{File: "github.com/example/handler.go", StartLine: 20, EndLine: 25, Count: 0},
		{File: "github.com/example/other.go", StartLine: 10, EndLine: 15, Count: 1},
		{File: "github.com/example/handler.go", StartLine: 5, EndLine: 8, Count: 1},
	}

	filtered := FilterBlocks(blocks, "handler.go", 10, 20)
	// Should include block 0 (10-15 overlaps 10-20) and block 1 (20-25 overlaps 10-20)
	// Should exclude block 2 (wrong file) and block 3 (5-8 < 10)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 filtered blocks, got %d", len(filtered))
	}
}

func TestFilterBlocks_NoOverlap(t *testing.T) {
	blocks := []Block{
		{File: "github.com/example/handler.go", StartLine: 1, EndLine: 5, Count: 1},
	}
	filtered := FilterBlocks(blocks, "handler.go", 10, 20)
	if len(filtered) != 0 {
		t.Fatalf("expected 0, got %d", len(filtered))
	}
}

func TestParseLine_Valid(t *testing.T) {
	b, err := parseLine("github.com/example/handler.go:10.2,15.4 2 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.File != "github.com/example/handler.go" {
		t.Fatalf("unexpected file: %s", b.File)
	}
	if b.StartLine != 10 || b.EndLine != 15 || b.Count != 1 {
		t.Fatalf("unexpected block: %+v", b)
	}
}

func TestParseLine_NoSpace(t *testing.T) {
	_, err := parseLine("nospacehere")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLine_NoSecondSpace(t *testing.T) {
	_, err := parseLine("onespace only")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseLine_InvalidCount(t *testing.T) {
	_, err := parseLine("file:1.2,3.4 2 notanumber")
	if err == nil {
		t.Fatal("expected error for invalid count")
	}
}

func TestParseLine_NoColon(t *testing.T) {
	_, err := parseLine("noposition 2 1")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseLine_NoComma(t *testing.T) {
	_, err := parseLine("file:1.2 2 1")
	if err == nil {
		t.Fatal("expected error for missing comma")
	}
}

func TestParseLine_InvalidStartLine(t *testing.T) {
	_, err := parseLine("file:abc.2,3.4 2 1")
	if err == nil {
		t.Fatal("expected error for invalid start line")
	}
}

func TestParseLine_InvalidEndLine(t *testing.T) {
	_, err := parseLine("file:1.2,abc.4 2 1")
	if err == nil {
		t.Fatal("expected error for invalid end line")
	}
}

func TestParseLineNum_WithDot(t *testing.T) {
	n, err := parseLineNum("42.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 42 {
		t.Fatalf("expected 42, got %d", n)
	}
}

func TestParseLineNum_WithoutDot(t *testing.T) {
	n, err := parseLineNum("42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 42 {
		t.Fatalf("expected 42, got %d", n)
	}
}

func TestParseLineNum_Invalid(t *testing.T) {
	_, err := parseLineNum("abc")
	if err == nil {
		t.Fatal("expected error for invalid line number")
	}
}

func TestMatchFile_ExactMatch(t *testing.T) {
	if !matchFile("handler.go", "handler.go") {
		t.Fatal("expected match")
	}
}

func TestMatchFile_SuffixMatch(t *testing.T) {
	if !matchFile("github.com/example/handler.go", "handler.go") {
		t.Fatal("expected match")
	}
}

func TestMatchFile_NoMatch(t *testing.T) {
	if matchFile("github.com/example/other.go", "handler.go") {
		t.Fatal("expected no match")
	}
}

func TestMatchFile_BackslashNormalization(t *testing.T) {
	if !matchFile("github.com\\example\\handler.go", "example/handler.go") {
		t.Fatal("expected match with backslash normalization")
	}
}
