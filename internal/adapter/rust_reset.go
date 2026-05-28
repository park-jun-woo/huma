//ff:func feature=adapter type=adapter control=sequence
//ff:what No-op reset for Rust adapter (no coverage data to clear)
package adapter

// Reset is a no-op for RustAdapter — there is no coverage data to clear.
func (a *RustAdapter) Reset() error {
	return nil
}
