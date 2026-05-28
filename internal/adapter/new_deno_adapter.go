//ff:func feature=adapter type=adapter control=sequence
//ff:what Creates a new DenoAdapter with config-based settings
package adapter

import "github.com/park-jun-woo/huma/internal/config"

// NewDenoAdapter creates a new Deno adapter for Supabase Edge Functions.
func NewDenoAdapter(cfg *config.Config) *DenoAdapter {
	return &DenoAdapter{
		cfg:     &cfg.Server,
		baseURL: cfg.BaseURL,
	}
}
