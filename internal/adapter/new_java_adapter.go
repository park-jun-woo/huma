//ff:func feature=adapter type=adapter control=sequence
//ff:what Creates a new JavaAdapter with config-based settings
package adapter

import (
	"github.com/park-jun-woo/huma/internal/config"
)

// NewJavaAdapter creates a new Java coverage adapter using JaCoCo.
func NewJavaAdapter(cfg *config.Config) *JavaAdapter {
	return &JavaAdapter{
		cfg:       &cfg.Server,
		baseURL:   cfg.BaseURL,
		jacocoDir: jacocoDir,
	}
}
