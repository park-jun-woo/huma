//ff:func feature=prompt type=helper control=iteration dimension=1
//ff:what Replaces path parameter placeholders like :id with the value 1
package prompt

import "strings"

func replaceParams(path string) string {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if strings.HasPrefix(p, ":") {
			parts[i] = "1"
		}
	}
	return strings.Join(parts, "/")
}
