//ff:func feature=prompt type=helper control=sequence
//ff:what Checks if a URL path contains colon-prefixed parameter placeholders
package prompt

import "strings"

func hasParam(path string) bool {
	return strings.Contains(path, ":")
}
