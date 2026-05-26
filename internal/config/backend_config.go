//ff:type feature=config type=model
//ff:what BackendConfig holds backend language, framework, module, and auth settings
package config

type BackendConfig struct {
	Lang      string     `yaml:"lang"`
	Framework string     `yaml:"framework"`
	Module    string     `yaml:"module"`
	Auth      AuthConfig `yaml:"auth"`
}
