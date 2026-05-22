//ff:func feature=cli type=command control=sequence
//ff:what CLI entrypoint that delegates to cobra command execution
package main

import "github.com/park-jun-woo/hurlfill/cmd"

func main() {
	cmd.Execute()
}
