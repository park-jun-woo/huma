//ff:type feature=config type=model
//ff:what SetupConfig holds the test-variable setup step: a user-authored .hurl run once before the cover loop whose [Captures] become injected hurl variables (e.g. {{token}})
package config

// SetupConfig describes the setup step run once before the cover loop. The
// referenced .hurl is executed via `hurl --json` and any variables it declares
// in a [Captures] section are merged into every endpoint's hurl run (capture
// wins over static hurl_variables). Token-only is the MVP, but the mechanism
// captures ANY variable the setup hurl declares — fixtures come free.
type SetupConfig struct {
	// Hurl is the path to the user-authored setup .hurl (login -> [Captures]).
	// Empty means no setup step is configured.
	Hurl string `yaml:"hurl"`
}
