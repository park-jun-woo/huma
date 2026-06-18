//ff:type feature=runner type=model
//ff:what hurlJSONReport is the subset of `hurl --json` output the setup harness consumes: per-file success plus its entries' [Captures]
package runner

// hurlJSONReport is the subset of `hurl --json` output we consume. Hurl emits one
// JSON object per source file (JSON-lines when several files run); each object has
// a top-level success flag and an entries array whose every entry carries the
// variables it captured in a [Captures] section.
type hurlJSONReport struct {
	Filename string          `json:"filename"`
	Success  bool            `json:"success"`
	Entries  []hurlJSONEntry `json:"entries"`
}
