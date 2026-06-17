package humaquest

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/runner"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// ---------------------------------------------------------------------------
// selectAdapter
// ---------------------------------------------------------------------------

func TestSelectAdapter(t *testing.T) {
	tests := []struct {
		lang string
		want any
	}{
		{"go", &adapter.GoAdapter{}},
		{"fiber", &adapter.GoAdapter{}},
		{"echo", &adapter.GoAdapter{}},
		{"python", &adapter.PythonAdapter{}},
		{"django", &adapter.PythonAdapter{}},
		{"node", &adapter.NodeAdapter{}},
		{"typescript", &adapter.NodeAdapter{}},
		{"deno", &adapter.DenoAdapter{}},
		{"java", &adapter.JavaAdapter{}},
		{"spring", &adapter.JavaAdapter{}},
		{"dotnet", &adapter.DotnetAdapter{}},
		{"csharp", &adapter.DotnetAdapter{}},
		{"php", &adapter.PhpAdapter{}},
		{"laravel", &adapter.PhpAdapter{}},
		{"rust", &adapter.RustAdapter{}},
		{"actix", &adapter.RustAdapter{}},
		{"", &adapter.GoAdapter{}},      // unset → Go fallback
		{"cobol", &adapter.GoAdapter{}}, // unknown → Go fallback
	}
	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			cfg := &config.Config{Scan: config.ScanConfig{Lang: tt.lang}}
			got := selectAdapter(cfg)
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("lang %q → %T, want %T", tt.lang, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// newLiveProbe
// ---------------------------------------------------------------------------

func TestNewLiveProbe(t *testing.T) {
	cfg := &config.Config{Scan: config.ScanConfig{Lang: "python"}}
	p := newLiveProbe(cfg)
	if p == nil {
		t.Fatal("nil probe")
	}
	if p.cfg != cfg {
		t.Error("cfg not wired")
	}
	if _, ok := p.adapter.(*adapter.PythonAdapter); !ok {
		t.Errorf("adapter = %T, want *PythonAdapter (selectAdapter wiring)", p.adapter)
	}
}

// ---------------------------------------------------------------------------
// liveProbe.Up  (adapter.Build — both branches reachable: `true`/`false`)
// ---------------------------------------------------------------------------

func TestLiveProbe_Up_BuildSuccess(t *testing.T) {
	chdir(t, t.TempDir())
	// `true` is a real, fast no-op binary: Build succeeds without any server.
	cfg := &config.Config{
		Scan:   config.ScanConfig{Lang: "go"},
		Server: config.ServerConfig{Build: "true"},
	}
	p := newLiveProbe(cfg)
	if err := p.Up(); err != nil {
		t.Fatalf("Up: %v", err)
	}
}

func TestLiveProbe_Up_BuildFailure(t *testing.T) {
	chdir(t, t.TempDir())
	// Empty build command → adapter.Build returns an error, wrapped by Up.
	cfg := &config.Config{
		Scan:   config.ScanConfig{Lang: "go"},
		Server: config.ServerConfig{Build: ""},
	}
	p := newLiveProbe(cfg)
	if err := p.Up(); err == nil {
		t.Fatal("want build error from Up")
	}
}

// ---------------------------------------------------------------------------
// liveProbe.Reset  (no-op)
// ---------------------------------------------------------------------------

func TestLiveProbe_Reset_Noop(t *testing.T) {
	p := newLiveProbe(&config.Config{Scan: config.ScanConfig{Lang: "go"}})
	if err := p.Reset(); err != nil {
		t.Errorf("Reset = %v, want nil", err)
	}
}

// ---------------------------------------------------------------------------
// liveProbe.Down  (adapter.Stop — proc nil → no-op cleanup)
// ---------------------------------------------------------------------------

func TestLiveProbe_Down_NoProcess(t *testing.T) {
	p := newLiveProbe(&config.Config{Scan: config.ScanConfig{Lang: "go"}})
	if err := p.Down(); err != nil {
		t.Errorf("Down = %v, want nil (no process to stop)", err)
	}
}

// ---------------------------------------------------------------------------
// liveProbe.Measure
//   - missing .hurl  → (nil, nil, nil)  [reachable]
//   - hurl present but server start fails → error            [reachable]
//   - genuine server-success path (cov returned) requires a real instrumented
//     server and is NOT unit-reachable (see DONE justification).
// ---------------------------------------------------------------------------

func measureEP() scanner.Endpoint {
	return scanner.Endpoint{ID: "GET_/api/v1/users", Method: "GET", Path: "/api/v1/users"}
}

func TestLiveProbe_Measure_MissingHurl(t *testing.T) {
	chdir(t, t.TempDir()) // no hurl files anywhere
	cfg := &config.Config{Scan: config.ScanConfig{Lang: "go"}, HurlDir: "hurl"}
	p := newLiveProbe(cfg)

	cov, res, err := p.Measure(measureEP())
	if err != nil {
		t.Fatalf("Measure: %v", err)
	}
	if cov != nil || res != nil {
		t.Errorf("missing hurl → want nil/nil, got cov=%v res=%v", cov, res)
	}
}

func TestLiveProbe_Measure_StartFailsPropagates(t *testing.T) {
	dir := t.TempDir()
	chdir(t, dir)
	ep := measureEP()
	// Drop a .hurl at the conventional path so FindHurlFile resolves it.
	hurlPath := runner.HurlFileName(&ep, "hurl")
	if err := os.MkdirAll(filepath.Dir(hurlPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(hurlPath, []byte("GET {{base_url}}/api/v1/users\nHTTP 200\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Empty Start command → RunWithCoverage fails at Start (no real server spawned).
	cfg := &config.Config{Scan: config.ScanConfig{Lang: "go"}, HurlDir: "hurl",
		Server: config.ServerConfig{Start: ""}}
	p := newLiveProbe(cfg)

	_, _, err := p.Measure(ep)
	if err == nil {
		t.Fatal("want error when the server cannot start")
	}
}
