//ff:func feature=config type=reader control=sequence
//ff:what Loads configuration from huma.yaml, falling back to defaults if missing
package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

func Load() (*Config, error) {
	data, err := os.ReadFile("huma.yaml")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Default(), nil
		}
		return nil, err
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
