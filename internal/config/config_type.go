//ff:type feature=config type=model
//ff:what Config holds all huma configuration including base URL, hurl directory, and server settings
package config

type Config struct {
	BaseURL       string            `yaml:"base_url"`
	HurlDir       string            `yaml:"hurl_dir"`
	HurlVariables map[string]string `yaml:"hurl_variables"`
	Scan          ScanConfig        `yaml:"scan"`
	Server        ServerConfig      `yaml:"server"`
	Deps          DepsConfig        `yaml:"deps"`
	// RequireCRI is the explicit minimum cheese-resistance index gate.
	// nil means unset (auto-require mode maximum).
	RequireCRI *int `yaml:"require_cri"`
}
