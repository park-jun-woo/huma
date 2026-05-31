//ff:func feature=config type=reader control=sequence
//ff:what Loads configuration from manifest.yaml, returning ErrNoManifest if missing
package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var ErrNoManifest = errors.New("manifest.yaml not found")

func Load() (*Config, error) {
	data, err := os.ReadFile("manifest.yaml")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoManifest
		}
		return nil, err
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	cfg := &Config{
		BaseURL:       m.Testing.BaseURL,
		HurlDir:       m.Testing.HurlDir,
		HurlVariables: m.Testing.HurlVariables,
		Scan: ScanConfig{
			Lang: m.Backend.Lang,
		},
		Server:     m.Testing.Server,
		Deps:       m.Testing.Deps,
		RequireCRI: m.Testing.RequireCRI,
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "http://localhost:8080"
	}
	if cfg.HurlDir == "" {
		cfg.HurlDir = "hurl"
	}
	if cfg.Scan.Lang == "" {
		cfg.Scan.Lang = "go"
	}

	return cfg, nil
}
