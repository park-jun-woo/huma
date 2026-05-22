//ff:type feature=config type=model
//ff:what ScanConfig holds language settings for source code scanning
package config

type ScanConfig struct {
	Lang string `yaml:"lang"`
}
