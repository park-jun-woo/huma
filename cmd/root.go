//ff:func feature=cli type=command control=sequence
//ff:what Defines the root cobra command and Execute entrypoint for hurlfill CLI
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hurlfill",
	Short: "Auto-generate wall-to-wall Hurl tests for legacy SaaS APIs",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
