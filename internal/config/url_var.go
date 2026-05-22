//ff:func feature=config type=helper control=sequence
//ff:what Returns the hurl variable name used for the base URL in templates
package config

func (c *Config) URLVar() string {
	for k, v := range c.HurlVariables {
		if v == c.BaseURL {
			return k
		}
	}
	return "base_url"
}
