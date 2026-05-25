package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestPythonStart_EmptyCommand(t *testing.T) {
	a := &PythonAdapter{cfg: &config.ServerConfig{Start: ""}}
	err := a.Start()
	if err == nil {
		t.Fatal("expected error for empty start command")
	}
	if err.Error() != "empty start command" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPythonStart_Success(t *testing.T) {
	a := &PythonAdapter{
		cfg: &config.ServerConfig{Start: "sleep 30", Env: map[string]string{"TEST_VAR": "val"}},
	}
	err := a.Start()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.proc == nil {
		t.Fatal("expected proc to be set")
	}
	a.proc.Process.Kill()
	a.proc.Wait()
}

func TestPythonStart_InvalidCommand(t *testing.T) {
	a := &PythonAdapter{
		cfg: &config.ServerConfig{Start: "nonexistent_binary_xyz_99999"},
	}
	err := a.Start()
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
}
