package adapter

import (
	"os/exec"
	"testing"
	"time"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestDenoStop_NilProc(t *testing.T) {
	a := &DenoAdapter{proc: nil}
	err := a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDenoStop_ProcessExitsGracefully(t *testing.T) {
	a := &DenoAdapter{
		cfg: &config.ServerConfig{Start: "sleep 60"},
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

func TestDenoStop_ForceKillOnTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping 10s timeout test")
	}
	cmd := exec.Command("bash", "-c", "trap '' INT TERM; while true; do sleep 1; done")
	if err := cmd.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	a := &DenoAdapter{proc: cmd}
	err := a.Stop()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.proc != nil {
		t.Fatal("expected proc to be nil after Stop")
	}
}

func TestDenoStop_AlreadyExited(t *testing.T) {
	a := &DenoAdapter{
		cfg: &config.ServerConfig{Start: "true"},
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
