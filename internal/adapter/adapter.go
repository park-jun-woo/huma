//ff:type feature=adapter type=adapter
//ff:what Adapter interface defines server lifecycle operations for coverage collection
package adapter

// Adapter manages a server process for coverage collection.
type Adapter interface {
	Build() error
	Start() error
	WaitReady() error
	Stop() error
	Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error)
	Reset() error
}
