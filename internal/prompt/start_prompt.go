//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the agent instruction text for starting the server when healthcheck fails
package prompt

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/huma/internal/config"
)

func StartPrompt(cfg *config.Config) string {
	var b strings.Builder

	readyURL := cfg.BaseURL + cfg.Server.Ready
	b.WriteString("# START — Server not responding\n")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("GET %s → connection refused\n", readyURL))
	b.WriteString("\n")
	b.WriteString("Start the server:\n")
	if cfg.Deps.Up != "" {
		b.WriteString(fmt.Sprintf("  %s\n", cfg.Deps.Up))
	}
	b.WriteString(fmt.Sprintf("  %s\n", cfg.Server.Start))
	b.WriteString("\n")
	b.WriteString("Then run `huma next` again.\n")

	return b.String()
}
