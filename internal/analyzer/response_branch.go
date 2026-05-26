//ff:type feature=analyzer type=model
//ff:what ResponseBranch represents a single response status code found in handler source code
package analyzer

// ResponseBranch represents a detected response status code in handler source.
type ResponseBranch struct {
	Status int
	File   string
	Line   int
	Code   string
}
