package adapter

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"
	"time"

	"github.com/park-jun-woo/huma/internal/config"
)

// --- DotnetAdapter ---

func TestNewDotnetAdapter_Fields(t *testing.T) {
	cfg := &config.Config{BaseURL: "http://localhost:5000", Server: config.ServerConfig{Start: "dotnet run"}}
	a := NewDotnetAdapter(cfg)
	if a.baseURL != "http://localhost:5000" {
		t.Fatalf("baseURL = %s", a.baseURL)
	}
	if a.cfg != &cfg.Server {
		t.Fatal("cfg should point to cfg.Server")
	}
}

func TestDotnetBuild(t *testing.T) {
	// already built → no-op
	if err := (&DotnetAdapter{cfg: &config.ServerConfig{Build: "false"}, built: true}).Build(); err != nil {
		t.Fatalf("already built: %v", err)
	}
	// empty build → no-op success
	a := &DotnetAdapter{cfg: &config.ServerConfig{Build: ""}}
	if err := a.Build(); err != nil {
		t.Fatalf("empty build: %v", err)
	}
	// success
	a = &DotnetAdapter{cfg: &config.ServerConfig{Build: "true"}}
	if err := a.Build(); err != nil || !a.built {
		t.Fatalf("success build: err=%v built=%v", err, a.built)
	}
	// failure
	if err := (&DotnetAdapter{cfg: &config.ServerConfig{Build: "false"}}).Build(); err == nil {
		t.Fatal("expected build failure")
	}
}

func TestDotnetCollectAndReset(t *testing.T) {
	a := &DotnetAdapter{cfg: &config.ServerConfig{}}
	cov, err := a.Collect("h.cs", 1, 5)
	if cov != nil || err != nil {
		t.Fatalf("Collect should be nil/nil, got %v %v", cov, err)
	}
	if err := a.Reset(); err != nil {
		t.Fatalf("Reset should be nil, got %v", err)
	}
}

func TestDotnetStart_EmptyAndSuccess(t *testing.T) {
	if err := (&DotnetAdapter{cfg: &config.ServerConfig{Start: ""}}).Start(); err == nil {
		t.Fatal("expected error for empty start")
	}
	a := &DotnetAdapter{cfg: &config.ServerConfig{Start: "sleep 30", Env: map[string]string{"X": "y"}}}
	if err := a.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if a.proc == nil {
		t.Fatal("proc should be set")
	}
	a.proc.Process.Kill()
	a.proc.Wait()

	if err := (&DotnetAdapter{cfg: &config.ServerConfig{Start: "nonexistent_bin_xyz_999"}}).Start(); err == nil {
		t.Fatal("expected error for invalid command")
	}
}

func TestDotnetStop(t *testing.T) {
	// nil proc
	if err := (&DotnetAdapter{}).Stop(); err != nil {
		t.Fatalf("nil proc: %v", err)
	}
	// graceful
	a := &DotnetAdapter{cfg: &config.ServerConfig{Start: "sleep 60"}}
	if err := a.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := a.Stop(); err != nil || a.proc != nil {
		t.Fatalf("stop: err=%v proc=%v", err, a.proc)
	}
	// force kill on timeout
	if !testing.Short() {
		cmd := exec.Command("bash", "-c", "trap '' INT TERM; while true; do sleep 1; done")
		cmd.Start()
		time.Sleep(200 * time.Millisecond)
		a2 := &DotnetAdapter{proc: cmd}
		if err := a2.Stop(); err != nil || a2.proc != nil {
			t.Fatalf("force kill: err=%v proc=%v", err, a2.proc)
		}
	}
}

func TestDotnetWaitReady(t *testing.T) {
	if !testing.Short() {
		if err := (&DotnetAdapter{cfg: &config.ServerConfig{Ready: ""}}).WaitReady(); err != nil {
			t.Fatalf("no ready url: %v", err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	if err := (&DotnetAdapter{cfg: &config.ServerConfig{Ready: ts.URL}}).WaitReady(); err != nil {
		t.Fatalf("ready 200: %v", err)
	}
}

// --- PhpAdapter ---

func TestNewPhpAdapter_Fields(t *testing.T) {
	cfg := &config.Config{BaseURL: "http://localhost:8000", Server: config.ServerConfig{Start: "php artisan serve"}}
	a := NewPhpAdapter(cfg)
	if a.baseURL != "http://localhost:8000" {
		t.Fatalf("baseURL = %s", a.baseURL)
	}
	if a.cfg != &cfg.Server {
		t.Fatal("cfg should point to cfg.Server")
	}
}

func TestPhpBuild(t *testing.T) {
	if err := (&PhpAdapter{cfg: &config.ServerConfig{Build: "false"}, built: true}).Build(); err != nil {
		t.Fatalf("already built: %v", err)
	}
	if err := (&PhpAdapter{cfg: &config.ServerConfig{Build: ""}}).Build(); err != nil {
		t.Fatalf("empty build: %v", err)
	}
	a := &PhpAdapter{cfg: &config.ServerConfig{Build: "true"}}
	if err := a.Build(); err != nil || !a.built {
		t.Fatalf("success build: err=%v built=%v", err, a.built)
	}
	if err := (&PhpAdapter{cfg: &config.ServerConfig{Build: "false"}}).Build(); err == nil {
		t.Fatal("expected build failure")
	}
}

func TestPhpCollectAndReset(t *testing.T) {
	a := &PhpAdapter{cfg: &config.ServerConfig{}}
	cov, err := a.Collect("h.php", 1, 5)
	if cov != nil || err != nil {
		t.Fatalf("Collect should be nil/nil, got %v %v", cov, err)
	}
	if err := a.Reset(); err != nil {
		t.Fatalf("Reset should be nil, got %v", err)
	}
}

func TestPhpStart_EmptyAndSuccess(t *testing.T) {
	if err := (&PhpAdapter{cfg: &config.ServerConfig{Start: ""}}).Start(); err == nil {
		t.Fatal("expected error for empty start")
	}
	a := &PhpAdapter{cfg: &config.ServerConfig{Start: "sleep 30", Env: map[string]string{"X": "y"}}}
	if err := a.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	a.proc.Process.Kill()
	a.proc.Wait()
	if err := (&PhpAdapter{cfg: &config.ServerConfig{Start: "nonexistent_bin_xyz_999"}}).Start(); err == nil {
		t.Fatal("expected error for invalid command")
	}
}

func TestPhpStop(t *testing.T) {
	if err := (&PhpAdapter{}).Stop(); err != nil {
		t.Fatalf("nil proc: %v", err)
	}
	a := &PhpAdapter{cfg: &config.ServerConfig{Start: "sleep 60"}}
	if err := a.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if err := a.Stop(); err != nil || a.proc != nil {
		t.Fatalf("stop: err=%v proc=%v", err, a.proc)
	}
	if !testing.Short() {
		cmd := exec.Command("bash", "-c", "trap '' INT TERM; while true; do sleep 1; done")
		cmd.Start()
		time.Sleep(200 * time.Millisecond)
		a2 := &PhpAdapter{proc: cmd}
		if err := a2.Stop(); err != nil || a2.proc != nil {
			t.Fatalf("force kill: err=%v proc=%v", err, a2.proc)
		}
	}
}

func TestPhpWaitReady(t *testing.T) {
	if !testing.Short() {
		if err := (&PhpAdapter{cfg: &config.ServerConfig{Ready: ""}}).WaitReady(); err != nil {
			t.Fatalf("no ready url: %v", err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	if err := (&PhpAdapter{cfg: &config.ServerConfig{Ready: ts.URL}}).WaitReady(); err != nil {
		t.Fatalf("ready 200: %v", err)
	}
}
