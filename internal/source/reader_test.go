package source

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadHandler_Success(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	content := `package main

import "net/http"

func GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", 400)
		return
	}
	w.Write([]byte("ok"))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("created"))
}
`
	os.WriteFile(f, []byte(content), 0o644)

	src, startLine, endLine, err := ReadHandler(f, "GetUser")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startLine != 5 {
		t.Fatalf("expected start line 5, got %d", startLine)
	}
	if endLine != 12 {
		t.Fatalf("expected end line 12, got %d", endLine)
	}
	if !strings.Contains(src, "func GetUser") {
		t.Fatal("expected function definition in source")
	}
	if strings.Contains(src, "func CreateUser") {
		t.Fatal("should not include next function")
	}
}

func TestReadHandler_MethodReceiver(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	content := `package main

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
`
	os.WriteFile(f, []byte(content), 0o644)

	src, startLine, _, err := ReadHandler(f, "ServeHTTP")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startLine != 3 {
		t.Fatalf("expected start line 3, got %d", startLine)
	}
	if !strings.Contains(src, "ServeHTTP") {
		t.Fatal("expected function name in source")
	}
}

func TestReadHandler_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte(`package main

func OtherFunc() {}
`), 0o644)

	_, _, _, err := ReadHandler(f, "GetUser")
	if err == nil {
		t.Fatal("expected error for missing handler")
	}
}

func TestReadHandler_FileNotFound(t *testing.T) {
	_, _, _, err := ReadHandler("/nonexistent/file.go", "Handler")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadHandler_TrimsTrailingBlanks(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	content := `package main

func Handler(c interface{}) {
	c.String(200, "ok")
}

// Next function comment
func NextHandler() {}
`
	os.WriteFile(f, []byte(content), 0o644)

	src, _, endLine, err := ReadHandler(f, "Handler")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not include trailing blank lines or the comment for next function
	lines := strings.Split(src, "\n")
	lastLine := lines[len(lines)-1]
	if strings.TrimSpace(lastLine) == "" || strings.HasPrefix(strings.TrimSpace(lastLine), "//") {
		t.Fatal("should have trimmed trailing blanks/comments")
	}
	if endLine != 5 {
		t.Fatalf("expected end line 5, got %d", endLine)
	}
}

func TestReadHandler_ScannerError(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	// Create a file with a very long line to trigger scanner buffer error
	longLine := make([]byte, 1024*1024) // 1MB line
	for i := range longLine {
		longLine[i] = 'x'
	}
	content := append([]byte("package main\n\nfunc Handler() {\n"), longLine...)
	content = append(content, []byte("\n}\n")...)
	os.WriteFile(f, content, 0o644)

	_, _, _, err := ReadHandler(f, "Handler")
	if err == nil {
		t.Fatal("expected scanner error for very long line")
	}
}

func TestReadHandler_LastFunction(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	content := `package main

func Handler(c interface{}) {
	c.String(200, "ok")
}
`
	os.WriteFile(f, []byte(content), 0o644)

	src, startLine, endLine, err := ReadHandler(f, "Handler")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startLine != 3 {
		t.Fatalf("expected start 3, got %d", startLine)
	}
	if endLine != 5 {
		t.Fatalf("expected end 5, got %d", endLine)
	}
	if !strings.Contains(src, "func Handler") {
		t.Fatal("expected function in source")
	}
}
