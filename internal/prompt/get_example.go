//ff:func feature=prompt type=builder control=sequence
//ff:what Generates a GET-specific Hurl example, choosing list or detail format based on path params
package prompt

import "fmt"

func getExample(examplePath, path, urlVar string) string {
	if hasParam(path) {
		return fmt.Sprintf(`GET {{%s}}%s
HTTP 200
[Asserts]
jsonpath "$.id" exists`, urlVar, examplePath)
	}
	return fmt.Sprintf(`GET {{%s}}%s
HTTP 200
[Asserts]
jsonpath "$" count > 0`, urlVar, examplePath)
}
