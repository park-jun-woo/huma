//ff:type feature=hurlcheck type=model
//ff:what Regular expressions for parsing hurl entry boundaries, skip options, and body assertions
package hurlcheck

import "regexp"

var (
	entryMethodRe = regexp.MustCompile(`^(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\s+(\S+)`)
	skipOptionRe  = regexp.MustCompile(`(?i)^\s*skip\s*:\s*true\s*$`)
	// body/header assertion predicates that inspect the response payload.
	bodyAssertRe = regexp.MustCompile(`(?i)^\s*(jsonpath|header|xpath|body|bytes|sha256|md5)\b`)
)
