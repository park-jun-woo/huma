//ff:type feature=config type=model
//ff:what Config holds all hurlfill configuration including base URL, hurl directory, and server settings
package config

type Config struct {
	BaseURL string       `yaml:"base_url"`
	HurlDir string       `yaml:"hurl_dir"`
	Scan    ScanConfig   `yaml:"scan"`
	Server  ServerConfig `yaml:"server"`
}
