package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type ScanConfig struct {
	Lang string `yaml:"lang"`
}

type ServerConfig struct {
	Build string            `yaml:"build"`
	Start string            `yaml:"start"`
	Ready string            `yaml:"ready"`
	Env   map[string]string `yaml:"env"`
}

type Config struct {
	BaseURL string       `yaml:"base_url"`
	HurlDir string       `yaml:"hurl_dir"`
	Scan    ScanConfig   `yaml:"scan"`
	Server  ServerConfig `yaml:"server"`
}

func Default() *Config {
	return &Config{
		BaseURL: "http://localhost:8080",
		HurlDir: "hurl",
		Scan: ScanConfig{
			Lang: "go",
		},
	}
}

func Load() (*Config, error) {
	data, err := os.ReadFile("hurlfill.yaml")
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
