//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the agent instruction text for configuring testing.server in manifest.yaml
package prompt

import "strings"

func ConfigPrompt() string {
	var b strings.Builder

	b.WriteString("# CONFIG — testing.server not configured\n")
	b.WriteString("\n")
	b.WriteString("Add server settings to manifest.yaml:\n")
	b.WriteString("\n")
	b.WriteString("  testing:\n")
	b.WriteString("    server:\n")
	b.WriteString("      build: \"go build -o ./server.test ./cmd/server\"\n")
	b.WriteString("      start: \"./server.test\"\n")
	b.WriteString("      ready: \"/api/health\"\n")
	b.WriteString("\n")
	b.WriteString("Then run `huma next` again.\n")

	return b.String()
}
