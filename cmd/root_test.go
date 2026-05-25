package cmd

import (
	"os"
	"os/exec"
	"testing"
)

func TestExecute_Success(t *testing.T) {
	// Execute() with --help succeeds without os.Exit.
	rootCmd.SetArgs([]string{"--help"})
	defer rootCmd.SetArgs(nil)

	// Call Execute() directly — it should not panic or exit for --help.
	Execute()
}

func TestExecute_Error(t *testing.T) {
	// Run this test in a subprocess to capture os.Exit(1).
	if os.Getenv("TEST_EXECUTE_ERROR") == "1" {
		rootCmd.SetArgs([]string{"nonexistent-command"})
		Execute()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExecute_Error")
	cmd.Env = append(os.Environ(), "TEST_EXECUTE_ERROR=1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected exit error, got nil")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected *exec.ExitError, got %T", err)
	}
	if exitErr.ExitCode() == 0 {
		t.Fatal("expected non-zero exit code")
	}
}
