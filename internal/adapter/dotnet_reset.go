//ff:func feature=adapter type=adapter control=sequence
//ff:what No-op reset for Dotnet adapter (no coverage data to clear)
package adapter

// Reset is a no-op for DotnetAdapter — there is no coverage data to clear.
func (d *DotnetAdapter) Reset() error {
	return nil
}
