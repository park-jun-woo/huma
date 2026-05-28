//ff:func feature=adapter type=adapter control=sequence
//ff:what No-op reset for PHP adapter (no coverage data to clear)
package adapter

// Reset is a no-op for PhpAdapter — there is no coverage data to clear.
func (p *PhpAdapter) Reset() error {
	return nil
}
