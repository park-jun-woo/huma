//ff:func feature=adapter type=adapter control=sequence
//ff:what Creates a new PhpAdapter with config-based settings
package adapter

import "github.com/park-jun-woo/huma/internal/config"

// NewPhpAdapter creates a new PHP adapter for Laravel applications.
func NewPhpAdapter(cfg *config.Config) *PhpAdapter {
	return &PhpAdapter{
		cfg:     &cfg.Server,
		baseURL: cfg.BaseURL,
	}
}
