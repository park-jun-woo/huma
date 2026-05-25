//ff:func feature=adapter type=adapter control=sequence
//ff:what Creates a new GoAdapter with config-based settings
package adapter

import (
	"github.com/park-jun-woo/huma/internal/config"
)

// NewGoAdapter creates a new Go coverage adapter.
func NewGoAdapter(cfg *config.Config) *GoAdapter {
	return &GoAdapter{
		cfg:      &cfg.Server,
		baseURL:  cfg.BaseURL,
		coverDir: coverDir,
	}
}
