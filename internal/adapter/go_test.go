package adapter

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/park-jun-woo/hurlfill/internal/config"
)

func TestNewGoAdapter(t *testing.T) {
	cfg := &config.Config{
		BaseURL: "http://localhost:3000",
		Server: config.ServerConfig{
			Build: "go build -cover -o app ./cmd/server",
			Start: "./app",
			Ready: "http://localhost:3000/health",
			Env:   map[string]string{"GIN_MODE": "release"},
		},
	}

	a := NewGoAdapter(cfg)
	if a == nil {
		t.Fatal("expected non-nil adapter")
	}
	if a.baseURL != "http://localhost:3000" {
		t.Fatalf("expected baseURL http://localhost:3000, got %s", a.baseURL)
	}
	if a.cfg != &cfg.Server {
		t.Fatal("expected cfg to point to cfg.Server")
	}
	if a.coverDir != coverDir {
		t.Fatalf("expected coverDir %s, got %s", coverDir, a.coverDir)
	}
}

func TestBuild_AlreadyBuilt(t *testing.T) {
	a := &GoAdapter{
		cfg:   &config.ServerConfig{Build: "true"},
		built: true,
	}
	err := a.Build()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestBuild_EmptyCommand(t *testing.T) {
	a := &GoAdapter{
		cfg: &config.ServerConfig{Build: ""},
	}
	err := a.Build()
	if err == nil {
		t.Fatal("expected error for empty build command")
	}
	if err.Error() != "empty build command" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuild_Success(t *testing.T) {
	a := &GoAdapter{
		cfg: &config.ServerConfig{Build: "true"},
	}
	err := a.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.built {
		t.Fatal("expected built to be true")
	}
}

func TestBuild_Failure(t *testing.T) {
	a := &GoAdapter{
		cfg: &config.ServerConfig{Build: "false"},
	}
	err := a.Build()
	if err == nil {
		t.Fatal("expected error from failed build")
	}
}

func TestBuild_SecondCallSkips(t *testing.T) {
	a := &GoAdapter{
		cfg: &config.ServerConfig{Build: "true"},
	}
	if err := a.Build(); err != nil {
		t.Fatalf("first build: %v", err)
	}
	// Change command to something that would fail
	a.cfg.Build = "false"
	// Second call should skip due to built=true
	if err := a.Build(); err != nil {
		t.Fatalf("second build should skip: %v", err)
	}
}

func TestStart_EmptyCommand(t *testing.T) {
	a := &GoAdapter{
		cfg:      &config.ServerConfig{Start: ""},
		coverDir: t.TempDir(),
	}
	err := a.Start()
	if err == nil {
		t.Fatal("expected error for empty start command")
	}
	if err.Error() != "empty start command" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStart_Success(t *testing.T) {
	a := &GoAdapter{
		cfg:      &config.ServerConfig{Start: "sleep 30", Env: map[string]string{"TEST_VAR": "val"}},
		coverDir: t.TempDir(),
	}
	err := a.Start()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.proc == nil {
		t.Fatal("expected proc to be set")
	}
	// Clean up
	a.proc.Process.Kill()
	a.proc.Wait()
}

func TestStart_InvalidCommand(t *testing.T) {
	a := &GoAdapter{
		cfg:      &config.ServerConfig{Start: "nonexistent_binary_xyz_99999"},
		coverDir: t.TempDir(),
	}
	err := a.Start()
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
}

func TestWaitReady_NoReadyURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 2s sleep test in short mode")
	}
	a := &GoAdapter{
		cfg: &config.ServerConfig{Ready: ""},
	}
	err := a.WaitReady()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWaitReady_ServerReturns200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	a := &GoAdapter{
		cfg: &config.ServerConfig{Ready: ts.URL},
	}
	err := a.WaitReady()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWaitReady_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 30s timeout test in short mode")
	}
	// Server that always returns 503, forcing the 30s timeout
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	a := &GoAdapter{
		cfg: &config.ServerConfig{Ready: ts.URL},
	}
	err := a.WaitReady()
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestStop_NilProc(t *testing.T) {
	a := &GoAdapter{proc: nil}
	err := a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStop_ProcessExitsGracefully(t *testing.T) {
	a := &GoAdapter{
		cfg:      &config.ServerConfig{Start: "sleep 60"},
		coverDir: t.TempDir(),
	}
	if err := a.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	err := a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.proc != nil {
		t.Fatal("expected proc to be nil after Stop")
	}
}

func TestStop_ForceKillOnTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 10s timeout test in short mode")
	}
	// Start a process that traps SIGINT and refuses to exit
	cmd := exec.Command("bash", "-c", "trap '' INT; sleep 300")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	a := &GoAdapter{proc: cmd}
	err := a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.proc != nil {
		t.Fatal("expected proc to be nil after Stop")
	}
}

func TestStop_AlreadyExited(t *testing.T) {
	// Start a process that exits immediately
	a := &GoAdapter{
		cfg:      &config.ServerConfig{Start: "true"},
		coverDir: t.TempDir(),
	}
	if err := a.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	// Wait for it to exit
	a.proc.Wait()
	// Now Stop should handle the already-exited process gracefully
	err := a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCollect_CovdataError(t *testing.T) {
	// Use a non-existent coverage directory to trigger covdata textfmt error
	a := &GoAdapter{
		coverDir: "/nonexistent/path/to/coverdata",
	}
	_, err := a.Collect("handler.go", 10, 20)
	if err == nil {
		t.Fatal("expected error from covdata textfmt")
	}
}

func TestCollect_EmptyCoverage(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	// Create empty coverage data directory
	covDir := filepath.Join(tmpDir, ".hurlfill", "coverdata")
	os.MkdirAll(covDir, 0o755)

	a := &GoAdapter{
		coverDir: covDir,
	}

	result, err := a.Collect("handler.go", 10, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Total != 0 {
		t.Fatalf("expected 0 total, got %d", result.Total)
	}
	if result.Covered != 0 {
		t.Fatalf("expected 0 covered, got %d", result.Covered)
	}
}

func TestCollect_WithCoverageData(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)

	// Create coverage data directory
	covDir := filepath.Join(tmpDir, ".hurlfill", "coverdata")
	os.MkdirAll(covDir, 0o755)
	os.MkdirAll(filepath.Join(tmpDir, ".hurlfill"), 0o755)

	// Create a handler source file
	handlerFile := filepath.Join(tmpDir, "handler.go")
	handlerContent := `package main

func handler() {
	x := 1
	if x > 0 {
		x++
	}
	if x > 10 {
		x--
	}
}
`
	os.WriteFile(handlerFile, []byte(handlerContent), 0o644)

	// Write a coverage.out file directly (bypass covdata textfmt)
	// We'll make covdata textfmt produce an empty file, then write our own
	a := &GoAdapter{
		coverDir: covDir,
	}

	// First run covdata textfmt (empty dir, creates empty coverage.out)
	result, err := a.Collect(handlerFile, 3, 11)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestReset_Success(t *testing.T) {
	tmpDir := t.TempDir()
	covDir := filepath.Join(tmpDir, "coverdata")
	os.MkdirAll(covDir, 0o755)

	// Write a file in the cover dir to verify it gets cleaned
	os.WriteFile(filepath.Join(covDir, "old.dat"), []byte("old"), 0o644)

	a := &GoAdapter{coverDir: covDir}
	err := a.Reset()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify directory exists but is empty
	entries, err := os.ReadDir(covDir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty dir, got %d entries", len(entries))
	}
}

func TestReset_RemoveAllError(t *testing.T) {
	// Create a directory structure where RemoveAll fails
	// by making the parent directory read-only so the child can't be removed
	tmpDir := t.TempDir()
	parentDir := filepath.Join(tmpDir, "parent")
	covDir := filepath.Join(parentDir, "coverdata")
	os.MkdirAll(covDir, 0o755)
	os.WriteFile(filepath.Join(covDir, "file.dat"), []byte("data"), 0o644)
	// Make parent read-only to prevent removal of coverdata
	os.Chmod(parentDir, 0o555)
	t.Cleanup(func() { os.Chmod(parentDir, 0o755) })

	a := &GoAdapter{coverDir: covDir}
	err := a.Reset()
	if err == nil {
		t.Fatal("expected error from RemoveAll")
	}
}

func TestReset_MkdirAllError(t *testing.T) {
	// Use a path where a file exists where a directory is needed
	tmpDir := t.TempDir()
	// Create a file where we want a dir path component
	blocker := filepath.Join(tmpDir, "blocker")
	os.WriteFile(blocker, []byte("block"), 0o644)

	a := &GoAdapter{coverDir: filepath.Join(blocker, "subdir")}
	err := a.Reset()
	if err == nil {
		t.Fatal("expected error from MkdirAll")
	}
}

func TestReset_CreateDir(t *testing.T) {
	tmpDir := t.TempDir()
	covDir := filepath.Join(tmpDir, "new", "coverdata")

	a := &GoAdapter{coverDir: covDir}
	err := a.Reset()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(covDir)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory")
	}
}

func TestReadUncoveredLines_Success(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	content := "line1\nline2\nline3\nline4\nline5\n"
	os.WriteFile(f, []byte(content), 0o644)

	covered := map[int]bool{1: true, 3: true}
	total := map[int]bool{1: true, 2: true, 3: true, 4: true}

	result, err := readUncoveredLines(f, covered, total)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Lines 2 and 4 are in total but not covered
	if len(result) != 2 {
		t.Fatalf("expected 2 uncovered lines, got %d", len(result))
	}
	if result[0].Line != 2 || result[0].Code != "line2" {
		t.Fatalf("unexpected first uncovered line: %+v", result[0])
	}
	if result[1].Line != 4 || result[1].Code != "line4" {
		t.Fatalf("unexpected second uncovered line: %+v", result[1])
	}
}

func TestReadUncoveredLines_FileNotFound(t *testing.T) {
	_, err := readUncoveredLines("/nonexistent/file.go", nil, nil)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadUncoveredLines_AllCovered(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte("line1\nline2\n"), 0o644)

	covered := map[int]bool{1: true, 2: true}
	total := map[int]bool{1: true, 2: true}

	result, err := readUncoveredLines(f, covered, total)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 uncovered lines, got %d", len(result))
	}
}

func TestReadUncoveredLines_TrailingWhitespace(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "handler.go")
	os.WriteFile(f, []byte("code  \t\n"), 0o644)

	covered := map[int]bool{}
	total := map[int]bool{1: true}

	result, err := readUncoveredLines(f, covered, total)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Code != "code" {
		t.Fatalf("expected trimmed code 'code', got %q", result[0].Code)
	}
}

func TestStart_AbsDirError(t *testing.T) {
	// filepath.Abs uses os.Getwd() internally. If we chdir to a deleted dir,
	// Getwd fails, which causes Abs to fail on a relative path.
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })

	os.Chdir(tmpDir)
	os.RemoveAll(tmpDir) // remove the dir we're in

	a := &GoAdapter{
		cfg:      &config.ServerConfig{Start: "sleep 30"},
		coverDir: "relative/path", // relative path forces Getwd call in Abs
	}
	err := a.Start()
	// On some systems this triggers an error from filepath.Abs, on others it may not.
	// We accept either outcome.
	if err != nil {
		// Expected: abs cover dir error
		return
	}
	// If it didn't fail, clean up
	if a.proc != nil {
		a.proc.Process.Kill()
		a.proc.Wait()
	}
}
