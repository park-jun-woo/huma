//ff:type feature=runner type=model
//ff:what hurlJSONEntry is one request/response entry in a hurl --json report with its captured variables
package runner

// hurlJSONEntry is one request/response entry with its captured variables.
type hurlJSONEntry struct {
	Captures []hurlJSONCapture `json:"captures"`
}
