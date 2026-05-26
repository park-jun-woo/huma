//ff:type feature=config type=model
//ff:what Manifest holds the yongol-compatible project manifest with backend and testing configuration
package config

type Manifest struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   MetadataConfig `yaml:"metadata"`
	Backend    BackendConfig  `yaml:"backend"`
	Testing    TestingConfig  `yaml:"testing"`
}
