//ff:type feature=scan type=model
//ff:what Represents a raw JSON endpoint entry before ID generation and handler parsing
package scanner

type rawEndpoint struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"`
	File    string `json:"file"`
	Line    int    `json:"line"`
}
