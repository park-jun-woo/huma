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
	// Setup is the test-variable harness step run once before the cover loop
	// (its [Captures] become injected hurl variables). Phase 009 / 2-A.
	Setup SetupConfig `yaml:"setup"`
	// Auth wires testing.auth for the JWT-mint path (Phase 009 / 2-B): huma
	// signs a token directly from a secret instead of capturing via login.
	Auth AuthConfig `yaml:"auth"`
	// RequireCRI is the explicit minimum cheese-resistance index gate.
	// nil means unset (auto-require mode maximum).
	RequireCRI *int `yaml:"require_cri"`
}
