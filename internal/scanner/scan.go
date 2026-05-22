//ff:func feature=scan type=engine control=sequence
//ff:what Walks a directory tree to find all API route registrations in Go source files
package scanner

import (
	"os"
	"path/filepath"
)

func Scan(dir string) ([]Endpoint, error) {
	var endpoints []Endpoint

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return skipDir(path)
		}
		if !isGoSource(path) {
			return nil
		}

		found, err := scanFile(path)
		if err != nil {
			return nil
		}
		endpoints = append(endpoints, found...)
		return nil
	})

	return endpoints, err
}
