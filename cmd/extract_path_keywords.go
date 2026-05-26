//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Extracts keywords from an endpoint path by splitting on slashes and removing path parameter segments
package cmd

import (
	"strings"
)

func extractPathKeywords(path string) []string {
	segments := strings.Split(strings.ToLower(path), "/")
	var keywords []string
	for _, seg := range segments {
		if seg == "" || strings.HasPrefix(seg, ":") {
			continue
		}
		keywords = append(keywords, seg)
	}
	return keywords
}
