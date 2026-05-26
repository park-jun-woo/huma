//ff:type feature=config type=model
//ff:what DepsConfig holds dependency lifecycle commands for up, down, and readiness check
package config

type DepsConfig struct {
	Up    string `yaml:"up"`
	Down  string `yaml:"down"`
	Ready string `yaml:"ready"`
}
