//ff:func feature=source type=helper control=iteration dimension=1
//ff:what Removes trailing blank lines and comment lines from a slice of strings
package source

import "strings"

func trimTrailing(lines []string) []string {
	for len(lines) > 0 {
		trimmed := strings.TrimSpace(lines[len(lines)-1])
		if trimmed == "" || strings.HasPrefix(trimmed, "//") {
			lines = lines[:len(lines)-1]
			continue
		}
		break
	}
	return lines
}
