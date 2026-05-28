//ff:func feature=scan type=parser control=iteration dimension=1
//ff:what Scans a Supabase Edge Functions directory to discover endpoints from file structure and method patterns
package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

func ParseEdgeFunctions(dir string) ([]Endpoint, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var endpoints []Endpoint

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		funcName := entry.Name()
		if strings.HasPrefix(funcName, "_") {
			continue
		}

		entryFile := findEntryFile(filepath.Join(dir, funcName))
		if entryFile == "" {
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, funcName, entryFile))
		if err != nil {
			continue
		}

		methods := extractMethods(string(content))
		path := "/functions/v1/" + funcName

		for _, method := range methods {
			ep := Endpoint{
				ID:      makeID(method, path),
				Method:  method,
				Path:    path,
				Handler: funcName,
				Source:  filepath.Join(dir, funcName, entryFile),
				Line:    0,
			}
			endpoints = append(endpoints, ep)
		}
	}

	return endpoints, nil
}
