//ff:type feature=analyzer type=adapter
//ff:what Analyzer interface defines the contract for static response code extraction from handler source
package analyzer

// Analyzer extracts response branches from a handler source file.
type Analyzer interface {
	Analyze(file string, handlerName string, startLine, endLine int) ([]ResponseBranch, error)
}
