//ff:type feature=config type=model
//ff:what MetadataConfig holds project metadata such as the project name
package config

type MetadataConfig struct {
	Name string `yaml:"name"`
}
