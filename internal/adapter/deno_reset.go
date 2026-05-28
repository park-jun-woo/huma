//ff:func feature=adapter type=adapter control=sequence
//ff:what No-op reset for Deno adapter (no coverage data to clear)
package adapter

// Reset is a no-op for DenoAdapter — there is no coverage data to clear.
func (d *DenoAdapter) Reset() error {
	return nil
}
