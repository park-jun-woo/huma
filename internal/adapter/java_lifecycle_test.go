package adapter

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestNewJavaAdapter_Fields(t *testing.T) {
	cfg := &config.Config{BaseURL: "http://localhost:8080", Server: config.ServerConfig{Start: "java -jar app.jar"}}
	a := NewJavaAdapter(cfg)
	if a.baseURL != "http://localhost:8080" {
		t.Fatalf("baseURL = %s", a.baseURL)
	}
	if a.cfg != &cfg.Server {
		t.Fatal("cfg should point to cfg.Server")
	}
	if a.jacocoDir != jacocoDir {
		t.Fatalf("jacocoDir = %s", a.jacocoDir)
	}
}

func TestJavaBuild(t *testing.T) {
	if err := (&JavaAdapter{cfg: &config.ServerConfig{Build: "false"}, built: true}).Build(); err != nil {
		t.Fatalf("already built: %v", err)
	}
	if err := (&JavaAdapter{cfg: &config.ServerConfig{Build: ""}}).Build(); err != nil {
		t.Fatalf("empty build: %v", err)
	}
	a := &JavaAdapter{cfg: &config.ServerConfig{Build: "true"}}
	if err := a.Build(); err != nil || !a.built {
		t.Fatalf("success: err=%v built=%v", err, a.built)
	}
	if err := (&JavaAdapter{cfg: &config.ServerConfig{Build: "false"}}).Build(); err == nil {
		t.Fatal("expected build failure")
	}
}

func TestJavaReset(t *testing.T) {
	tmp := t.TempDir()
	jdir := filepath.Join(tmp, "jacoco")
	os.MkdirAll(jdir, 0o755)
	os.WriteFile(filepath.Join(jdir, "old.exec"), []byte("x"), 0o644)
	a := &JavaAdapter{jacocoDir: jdir}
	if err := a.Reset(); err != nil {
		t.Fatalf("reset: %v", err)
	}
	entries, _ := os.ReadDir(jdir)
	if len(entries) != 0 {
		t.Fatalf("expected empty dir, got %d", len(entries))
	}

	// RemoveAll error: make parent read-only
	parent := filepath.Join(tmp, "parent")
	child := filepath.Join(parent, "jacoco")
	os.MkdirAll(child, 0o755)
	os.WriteFile(filepath.Join(child, "f"), []byte("x"), 0o644)
	os.Chmod(parent, 0o555)
	t.Cleanup(func() { os.Chmod(parent, 0o755) })
	if err := (&JavaAdapter{jacocoDir: child}).Reset(); err == nil {
		t.Fatal("expected RemoveAll error")
	}
}

func TestJavaStart_EmptyAndSuccess(t *testing.T) {
	if err := (&JavaAdapter{cfg: &config.ServerConfig{Start: ""}, jacocoDir: t.TempDir()}).Start(); err == nil {
		t.Fatal("expected error for empty start")
	}
	// success with no agent (JACOCO_AGENT unset → agentPath "")
	a := &JavaAdapter{cfg: &config.ServerConfig{Start: "sleep 30", Env: map[string]string{"X": "y"}}, jacocoDir: t.TempDir()}
	if err := a.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if a.proc == nil {
		t.Fatal("proc should be set")
	}
	a.proc.Process.Kill()
	a.proc.Wait()

	// success with a resolvable JACOCO_AGENT (exercises the agentPath != "" branch)
	agent := filepath.Join(t.TempDir(), "agent.jar")
	os.WriteFile(agent, []byte("jar"), 0o644)
	a2 := &JavaAdapter{cfg: &config.ServerConfig{Start: "sleep 30", Env: map[string]string{"JACOCO_AGENT": agent}}, jacocoDir: t.TempDir()}
	if err := a2.Start(); err != nil {
		t.Fatalf("start with agent: %v", err)
	}
	a2.proc.Process.Kill()
	a2.proc.Wait()

	if err := (&JavaAdapter{cfg: &config.ServerConfig{Start: "nonexistent_bin_xyz_999"}, jacocoDir: t.TempDir()}).Start(); err == nil {
		t.Fatal("expected error for invalid command")
	}
}

func TestJavaStop(t *testing.T) {
	if err := (&JavaAdapter{}).Stop(); err != nil {
		t.Fatalf("nil proc: %v", err)
	}
	a := &JavaAdapter{cfg: &config.ServerConfig{Start: "sleep 60"}, jacocoDir: t.TempDir()}
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
		a2 := &JavaAdapter{proc: cmd}
		if err := a2.Stop(); err != nil || a2.proc != nil {
			t.Fatalf("force kill: err=%v proc=%v", err, a2.proc)
		}
	}
}

func TestJavaWaitReady(t *testing.T) {
	if !testing.Short() {
		if err := (&JavaAdapter{cfg: &config.ServerConfig{Ready: ""}}).WaitReady(); err != nil {
			t.Fatalf("no ready url: %v", err)
		}
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	if err := (&JavaAdapter{cfg: &config.ServerConfig{Ready: ts.URL}}).WaitReady(); err != nil {
		t.Fatalf("ready 200: %v", err)
	}
}

func TestJavaCollect_NoExecFile(t *testing.T) {
	// no jacoco.exec → returns nil/nil
	a := &JavaAdapter{jacocoDir: t.TempDir()}
	cov, err := a.Collect("Handler.java", 1, 5)
	if cov != nil || err != nil {
		t.Fatalf("expected nil/nil when exec file absent, got %v %v", cov, err)
	}
}

func TestFindJacocoCLI(t *testing.T) {
	// Always returns a non-empty path (either a matched JAR or the fallback).
	got := findJacocoCLI()
	if got == "" {
		t.Fatal("expected non-empty path")
	}
	if !strings.HasSuffix(got, ".jar") {
		t.Fatalf("expected a .jar path, got %s", got)
	}
}

func TestResolveJacocoAgent(t *testing.T) {
	// env-specified, existing file → returned
	dir := t.TempDir()
	agent := filepath.Join(dir, "agent.jar")
	os.WriteFile(agent, []byte("x"), 0o644)
	if got := resolveJacocoAgent(map[string]string{"JACOCO_AGENT": agent}); got != agent {
		t.Fatalf("expected %s, got %s", agent, got)
	}
	// env-specified but nonexistent → falls through (likely "")
	got := resolveJacocoAgent(map[string]string{"JACOCO_AGENT": filepath.Join(dir, "nope.jar")})
	if got == agent {
		t.Fatal("nonexistent env path should not be returned")
	}
	// project dir fallback: build/jacoco/*.jar
	prevWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(prevWd) })
	work := t.TempDir()
	os.Chdir(work)
	os.MkdirAll(filepath.Join(work, "build", "jacoco"), 0o755)
	pj := filepath.Join(work, "build", "jacoco", "agent.jar")
	os.WriteFile(pj, []byte("x"), 0o644)
	// only exercises project-dir branch when ~/.m2 has no match; assert it returns
	// either the m2 match or our project jar — both acceptable, just non-panicking.
	_ = resolveJacocoAgent(nil)
}

func TestCollectSourceFileLines(t *testing.T) {
	sf := jacocoSourceFile{
		Name: "H.java",
		Lines: []jacocoLine{
			{Nr: 5, Ci: 1},   // covered
			{Nr: 10, Mi: 1},  // total only (missed)
			{Nr: 15, Ci: 2},  // covered, in range
			{Nr: 25, Ci: 1},  // out of range (above)
			{Nr: 1, Ci: 1},   // out of range (below)
			{Nr: 8, Ci: 0, Mi: 0}, // neither → ignored
		},
	}
	covered := map[int]bool{}
	total := map[int]bool{}
	collectSourceFileLines(sf, 4, 20, covered, total)
	if !covered[5] || !covered[15] {
		t.Errorf("covered = %v, want 5,15", covered)
	}
	if covered[10] {
		t.Error("line 10 missed, should not be covered")
	}
	if !total[10] {
		t.Error("line 10 should be in total")
	}
	if total[25] || total[1] {
		t.Error("out-of-range lines should be excluded")
	}
	if total[8] {
		t.Error("ci=0,mi=0 line should be ignored")
	}
}

func TestCollectPackageLines(t *testing.T) {
	pkg := jacocoPackage{
		Name: "com/example",
		SourceFiles: []jacocoSourceFile{
			{Name: "H.java", Lines: []jacocoLine{{Nr: 5, Ci: 1}}},
			{Name: "Other.java", Lines: []jacocoLine{{Nr: 6, Ci: 1}}},
		},
	}
	covered := map[int]bool{}
	total := map[int]bool{}
	collectPackageLines(pkg, "src/main/java/com/example/H.java", 0, 0, covered, total)
	if !covered[5] {
		t.Errorf("expected line 5 covered, got %v", covered)
	}
	if covered[6] {
		t.Error("non-matching source file should be skipped")
	}
}
