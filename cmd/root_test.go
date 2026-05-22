package cmd

import (
	"testing"
)

func TestExecute_Success(t *testing.T) {
	// Execute() with --help succeeds without os.Exit.
	rootCmd.SetArgs([]string{"--help"})
	defer rootCmd.SetArgs(nil)

	// Call Execute() directly — it should not panic or exit for --help.
	Execute()
}
