//ff:func feature=config type=helper control=sequence
//ff:what Returns a Config with default values for base URL, hurl dir, and scan language
package config

func Default() *Config {
	return &Config{
		BaseURL: "http://localhost:8080",
		HurlDir: "hurl",
		Scan: ScanConfig{
			Lang: "go",
		},
	}
}
