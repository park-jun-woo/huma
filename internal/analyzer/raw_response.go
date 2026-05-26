//ff:type feature=analyzer type=model
//ff:what Represents a single response entry from endpoint input JSON with status code and location
package analyzer

type rawResponse struct {
	Status int    `json:"status"`
	Line   int    `json:"line"`
	Code   string `json:"code"`
}
