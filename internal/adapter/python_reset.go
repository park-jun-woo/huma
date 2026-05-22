//ff:func feature=adapter type=engine control=sequence
//ff:what Removes the .coverage file to prepare for a fresh coverage run
package adapter

import (
	"os"
)

// Reset removes the .coverage file so the next coverage run starts fresh.
func (a *PythonAdapter) Reset() error {
	err := os.Remove(coverageFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
