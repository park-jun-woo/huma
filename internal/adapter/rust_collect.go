//ff:func feature=adapter type=adapter control=sequence
//ff:what Returns nil coverage for Rust adapter (coverage not supported)
package adapter

// Collect always returns nil for RustAdapter — LLVM source-based coverage
// instrumentation for running servers is not supported.
func (a *RustAdapter) Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	return nil, nil
}
