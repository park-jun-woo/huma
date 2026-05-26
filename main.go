//ff:func feature=cli type=command control=sequence
//ff:what CLI entrypoint that delegates to cobra command execution
package main

import "github.com/park-jun-woo/huma/cmd"

var Version = "dev"

func main() {
	cmd.Execute()
}
