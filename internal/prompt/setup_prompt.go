//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the agent instruction text for creating manifest.yaml when it does not exist
package prompt

import "strings"

func SetupPrompt() string {
	var b strings.Builder

	b.WriteString("# SETUP — manifest.yaml not found\n")
	b.WriteString("\n")
	b.WriteString("Create manifest.yaml in the project root:\n")
	b.WriteString("\n")
	b.WriteString("  apiVersion: yongol/v1\n")
	b.WriteString("  kind: Project\n")
	b.WriteString("  metadata:\n")
	b.WriteString("    name: my-project\n")
	b.WriteString("  backend:\n")
	b.WriteString("    lang: go\n")
	b.WriteString("    framework: gin\n")
	b.WriteString("    module: github.com/org/project\n")
	b.WriteString("  testing:\n")
	b.WriteString("    base_url: \"http://localhost:8080\"\n")
	b.WriteString("    hurl_dir: \"hurl\"\n")
	b.WriteString("    hurl_variables:\n")
	b.WriteString("      host: \"http://localhost:8080\"\n")
	b.WriteString("    server:\n")
	b.WriteString("      build: \"go build -o ./server.test ./cmd/server\"\n")
	b.WriteString("      start: \"./server.test\"\n")
	b.WriteString("      ready: \"/api/health\"\n")
	b.WriteString("\n")
	b.WriteString("Then run `huma next` again.\n")

	return b.String()
}
