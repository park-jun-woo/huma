//ff:type feature=config type=model
//ff:what TestingConfig holds testing settings including base URL, hurl directory, variables, deps, and server
package config

type TestingConfig struct {
	BaseURL       string            `yaml:"base_url"`
	HurlDir       string            `yaml:"hurl_dir"`
	HurlVariables map[string]string `yaml:"hurl_variables"`
	Deps          DepsConfig        `yaml:"deps"`
	Server        ServerConfig      `yaml:"server"`
	// Setup is the setup .hurl run once before the cover loop; its [Captures]
	// become injected hurl variables (Phase 009 / 2-A).
	Setup SetupConfig `yaml:"setup"`
	// Auth is the JWT-mint config (Phase 009 / 2-B), previously unwired.
	Auth AuthConfig `yaml:"auth"`
	// RequireCRI is the minimum cheese-resistance index (0..3) a verdict must
	// reach to be granted PASS. -1 (unset) means huma auto-requires the maximum
	// CRI reachable in the current mode.
	RequireCRI *int `yaml:"require_cri"`
}
