package adapter

import (
	"os/exec"
	"testing"
	"time"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestGoStop_NilProc(t *testing.T) {
	a := &GoAdapter{proc: nil}
	err := a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGoStop_ProcessExitsGracefully(t *testing.T) {
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

func TestGoStop_ForceKillOnTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 10s timeout test")
	}
	cmd := exec.Command("bash", "-c", "trap '' INT TERM; while true; do sleep 1; done")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	a := &GoAdapter{proc: cmd}
	err := a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.proc != nil {
		t.Fatal("expected proc to be nil after Stop")
	}
}

func TestGoStop_AlreadyExited(t *testing.T) {
	a := &GoAdapter{
		cfg:      &config.ServerConfig{Start: "true"},
		coverDir: t.TempDir(),
	}
	if err := a.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	a.proc.Wait()
	err := a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
