//ff:type feature=scan type=model
//ff:what ResponseField represents a single field from an OpenAPI response schema with its JSON path and type
package scanner

type ResponseField struct {
	Path string `json:"path"`
	Type string `json:"type"`
}
