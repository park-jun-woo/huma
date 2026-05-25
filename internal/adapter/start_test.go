package adapter

import (
	"os"
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestGoStart_EmptyCommand(t *testing.T) {
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

func TestGoStart_Success(t *testing.T) {
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
	a.proc.Process.Kill()
	a.proc.Wait()
}

func TestGoStart_InvalidCommand(t *testing.T) {
	a := &GoAdapter{
		cfg:      &config.ServerConfig{Start: "nonexistent_binary_xyz_99999"},
		coverDir: t.TempDir(),
	}
	err := a.Start()
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
}

func TestGoStart_AbsDirError(t *testing.T) {
	tmpDir := t.TempDir()
	orig, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(orig) })
	os.Chdir(tmpDir)
	os.RemoveAll(tmpDir)

	a := &GoAdapter{
		cfg:      &config.ServerConfig{Start: "sleep 30"},
		coverDir: "relative/path",
	}
	err := a.Start()
	if err != nil {
		return
	}
	if a.proc != nil {
		a.proc.Process.Kill()
		a.proc.Wait()
	}
}
