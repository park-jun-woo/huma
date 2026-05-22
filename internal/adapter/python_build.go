//ff:func feature=adapter type=engine control=sequence
//ff:what No-op build for Python servers which require no compilation step
package adapter

// Build is a no-op for Python servers (no compilation needed).
func (a *PythonAdapter) Build() error {
	return nil
}
