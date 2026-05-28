//ff:type feature=scan type=model
//ff:what Regex patterns and allowed method set for Edge Function HTTP method extraction
package scanner

import "regexp"

var edgeFuncMethodPositive = regexp.MustCompile(
	`(?:req|request)\.method\s*===?\s*['"](\w+)['"]`,
)

var edgeFuncMethodCase = regexp.MustCompile(
	`case\s+['"](\w+)['"]\s*:`,
)

var allowedMethods = map[string]bool{
	"GET": true, "POST": true, "PUT": true,
	"PATCH": true, "DELETE": true,
}
