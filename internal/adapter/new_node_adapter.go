//ff:func feature=adapter type=adapter control=sequence
//ff:what Creates a new NodeAdapter with config-based settings
package adapter

import (
	"github.com/park-jun-woo/huma/internal/config"
)

// NewNodeAdapter creates a new Node.js coverage adapter.
func NewNodeAdapter(cfg *config.Config) *NodeAdapter {
	return &NodeAdapter{
		cfg:      &cfg.Server,
		baseURL:  cfg.BaseURL,
		coverDir: nodeCoverDir,
	}
}
