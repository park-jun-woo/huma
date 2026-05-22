package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan_FindsEndpoints(t *testing.T) {
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "routes.go")
	content := `package main

func setupRoutes(r *gin.Engine) {
	r.GET("/users", listUsers)
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)
}
`
	os.WriteFile(goFile, []byte(content), 0o644)

	endpoints, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 5 {
		t.Fatalf("expected 5 endpoints, got %d", len(endpoints))
	}
}

func TestScan_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "routes_test.go")
	content := `package main

func TestRoutes(t *testing.T) {
	r.GET("/test", testHandler)
}
`
	os.WriteFile(testFile, []byte(content), 0o644)

	endpoints, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 0 {
		t.Fatalf("expected 0 endpoints from test file, got %d", len(endpoints))
	}
}

func TestScan_SkipsVendor(t *testing.T) {
	tmpDir := t.TempDir()
	vendorDir := filepath.Join(tmpDir, "vendor", "lib")
	os.MkdirAll(vendorDir, 0o755)
	os.WriteFile(filepath.Join(vendorDir, "routes.go"), []byte(`package lib
func setup(r *gin.Engine) {
	r.GET("/vendor", handler)
}
`), 0o644)

	endpoints, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 0 {
		t.Fatalf("expected 0 endpoints from vendor, got %d", len(endpoints))
	}
}

func TestScan_SkipsNonGoFiles(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "routes.py"), []byte(`r.GET("/test", handler)`), 0o644)

	endpoints, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 0 {
		t.Fatalf("expected 0 endpoints from non-go file, got %d", len(endpoints))
	}
}

func TestScan_WalkError(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a subdirectory that is not readable
	subDir := filepath.Join(tmpDir, "sub")
	os.MkdirAll(subDir, 0o755)
	os.WriteFile(filepath.Join(subDir, "routes.go"), []byte(`package main
func setup(r *gin.Engine) { r.GET("/test", h) }
`), 0o644)
	// Make subdirectory unreadable to cause walk error
	os.Chmod(subDir, 0o000)
	t.Cleanup(func() { os.Chmod(subDir, 0o755) })

	// Should not fail, just skip the unreadable dir
	endpoints, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 0 {
		t.Fatalf("expected 0 endpoints from unreadable dir, got %d", len(endpoints))
	}
}

func TestScan_ScanFileError(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "routes.go")
	os.WriteFile(f, []byte(`package main
func setup(r *gin.Engine) { r.GET("/test", h) }
`), 0o644)
	// Make file unreadable after walk discovers it
	os.Chmod(f, 0o000)
	t.Cleanup(func() { os.Chmod(f, 0o644) })

	endpoints, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// scanFile error is silently ignored
	if len(endpoints) != 0 {
		t.Fatalf("expected 0, got %d", len(endpoints))
	}
}

func TestScan_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	endpoints, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 0 {
		t.Fatalf("expected 0 endpoints, got %d", len(endpoints))
	}
}

func TestScanFile_MultipleRoutes(t *testing.T) {
	tmpDir := t.TempDir()
	f := filepath.Join(tmpDir, "routes.go")
	content := `package main

func routes(r *gin.Engine) {
	r.GET("/health", healthCheck)
	r.POST("/api/v1/items", createItem)
	r.PATCH("/api/v1/items/:id", patchItem)
}
`
	os.WriteFile(f, []byte(content), 0o644)

	endpoints, err := scanFile(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(endpoints) != 3 {
		t.Fatalf("expected 3, got %d", len(endpoints))
	}
	if endpoints[0].Method != "GET" || endpoints[0].Path != "/health" {
		t.Fatalf("unexpected first endpoint: %+v", endpoints[0])
	}
	if endpoints[2].Method != "PATCH" {
		t.Fatalf("expected PATCH, got %s", endpoints[2].Method)
	}
}

func TestExtractHandler(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{`r.GET("/users", listUsers)`, "listUsers"},
		{`r.POST("/users", createUser)`, "createUser"},
		{`r.GET("/health")`, ""},
	}
	for _, tt := range tests {
		got := extractHandler(tt.line)
		if got != tt.want {
			t.Errorf("extractHandler(%q) = %q, want %q", tt.line, got, tt.want)
		}
	}
}

func TestMakeID(t *testing.T) {
	id1 := makeID("GET", "/users")
	id2 := makeID("POST", "/users")
	id3 := makeID("GET", "/users")

	if id1 == id2 {
		t.Fatal("different method+path should produce different IDs")
	}
	if id1 != id3 {
		t.Fatal("same method+path should produce same ID")
	}
	if len(id1) != 16 { // 8 bytes = 16 hex chars
		t.Fatalf("expected 16 hex chars, got %d", len(id1))
	}
}
