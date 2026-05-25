//ff:func feature=adapter type=adapter control=sequence
//ff:what Creates a new PythonAdapter with config-based settings
package adapter

import (
	"github.com/park-jun-woo/huma/internal/config"
)

// NewPythonAdapter creates a new Python coverage adapter.
func NewPythonAdapter(cfg *config.Config) *PythonAdapter {
	return &PythonAdapter{
		cfg:     &cfg.Server,
		baseURL: cfg.BaseURL,
	}
}
