//ff:func feature=adapter type=adapter control=sequence
//ff:what Creates a new DotnetAdapter with config-based settings
package adapter

import "github.com/park-jun-woo/huma/internal/config"

// NewDotnetAdapter creates a new Dotnet adapter for ASP.NET Core applications.
func NewDotnetAdapter(cfg *config.Config) *DotnetAdapter {
	return &DotnetAdapter{
		cfg:     &cfg.Server,
		baseURL: cfg.BaseURL,
	}
}
