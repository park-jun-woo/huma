//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Searches default paths for an OpenAPI file and returns the first found
package scanner

import "os"

func FindOpenAPIFile() string {
	candidates := []string{
		"openapi.yaml",
		"api/openapi.yaml",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}
