//ff:func feature=coverage type=parser control=sequence
//ff:what Extracts the line number from a line.col format string
package coverage

import (
	"strconv"
	"strings"
)

// parseLineNum extracts the line number from "line.col" format.
func parseLineNum(s string) (int, error) {
	dotIdx := strings.Index(s, ".")
	if dotIdx < 0 {
		return strconv.Atoi(s)
	}
	return strconv.Atoi(s[:dotIdx])
}
