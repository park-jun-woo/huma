//ff:type feature=scan type=model
//ff:what Represents a raw JSON endpoint entry before ID generation and handler parsing
package scanner

import "encoding/json"

type rawEndpoint struct {
	Method    string          `json:"method" yaml:"method"`
	Path      string          `json:"path" yaml:"path"`
	Handler   string          `json:"handler" yaml:"handler"`
	File      string          `json:"file" yaml:"file"`
	Line      int             `json:"line" yaml:"line"`
	Responses json.RawMessage `json:"responses" yaml:"responses"`
}
