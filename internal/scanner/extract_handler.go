//ff:func feature=scan type=parser control=sequence
//ff:what Extracts the handler function name from a gin route registration line
package scanner

import "strings"

func extractHandler(line string) string {
	parts := strings.Split(line, ",")
	if len(parts) < 2 {
		return ""
	}
	h := strings.TrimSpace(parts[len(parts)-1])
	h = strings.TrimSuffix(h, ")")
	return strings.TrimSpace(h)
}
