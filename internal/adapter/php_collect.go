//ff:func feature=adapter type=adapter control=sequence
//ff:what Returns nil coverage for PHP adapter (coverage not supported)
package adapter

// Collect always returns nil for PhpAdapter — PHP/Laravel
// coverage instrumentation is not supported.
func (p *PhpAdapter) Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	return nil, nil
}
