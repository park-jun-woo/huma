//ff:func feature=hurlcheck type=helper control=sequence
//ff:what Extracts the numeric status code from a hurl HTTP status line
package hurlcheck

import "strconv"

// parseStatusLine extracts the numeric status from an "HTTP NNN" line.
func parseStatusLine(line string) int {
	m := httpStatusRe.FindStringSubmatch(line)
	code, _ := strconv.Atoi(m[1])
	return code
}
