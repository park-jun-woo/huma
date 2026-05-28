//ff:func feature=adapter type=adapter control=sequence
//ff:what Creates a new RustAdapter with config-based settings
package adapter

import "github.com/park-jun-woo/huma/internal/config"

// NewRustAdapter creates a new Rust adapter for Actix-web servers.
func NewRustAdapter(cfg *config.Config) *RustAdapter {
	return &RustAdapter{
		cfg:     &cfg.Server,
		baseURL: cfg.BaseURL,
	}
}
