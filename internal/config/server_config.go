//ff:type feature=config type=model
//ff:what ServerConfig holds server build, start, readiness, and environment settings
package config

type ServerConfig struct {
	Build string            `yaml:"build"`
	Start string            `yaml:"start"`
	Ready string            `yaml:"ready"`
	Env   map[string]string `yaml:"env"`
}
