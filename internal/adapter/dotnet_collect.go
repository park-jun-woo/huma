//ff:func feature=adapter type=adapter control=sequence
//ff:what Returns nil coverage for Dotnet adapter (coverage not supported)
package adapter

// Collect always returns nil for DotnetAdapter — ASP.NET Core
// coverage instrumentation is not supported.
func (d *DotnetAdapter) Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	return nil, nil
}
