//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Extracts keywords from a hurl filename by splitting on underscores and removing the method prefix and common param tokens
package cmd

import (
	"strings"
)

func extractFilenameKeywords(name string) []string {
	name = strings.TrimSuffix(name, ".hurl")
	parts := strings.Split(strings.ToLower(name), "_")
	if len(parts) < 2 {
		return nil
	}
	parts = parts[1:]
	var keywords []string
	for _, p := range parts {
		if p == "" || p == "id" {
			continue
		}
		keywords = append(keywords, p)
	}
	return keywords
}
