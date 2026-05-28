//ff:func feature=adapter type=adapter control=sequence
//ff:what Returns nil coverage for Deno adapter (coverage not supported)
package adapter

// Collect always returns nil for DenoAdapter — Deno/Edge Functions
// coverage instrumentation is not supported.
func (d *DenoAdapter) Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	return nil, nil
}
