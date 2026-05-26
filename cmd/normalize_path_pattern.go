//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Normalizes a URL path by replacing colon-prefixed parameter segments with a wildcard placeholder
package cmd

import "strings"

func normalizePathPattern(path string) string {
	segments := strings.Split(path, "/")
	for i, seg := range segments {
		if strings.HasPrefix(seg, ":") {
			segments[i] = ":_"
		}
	}
	return strings.Join(segments, "/")
}
