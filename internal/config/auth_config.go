//ff:type feature=config type=model
//ff:what AuthConfig holds authentication settings including JWT type, secret env, and claims mapping
package config

type AuthConfig struct {
	Type      string            `yaml:"type"`
	SecretEnv string            `yaml:"secret_env"`
	UserTable string            `yaml:"user_table"`
	Claims    map[string]string `yaml:"claims"`
}
