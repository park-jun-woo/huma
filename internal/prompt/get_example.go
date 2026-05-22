//ff:func feature=prompt type=builder control=sequence
//ff:what Generates a GET-specific Hurl example, choosing list or detail format based on path params
package prompt

import "fmt"

func getExample(examplePath, path string) string {
	if hasParam(path) {
		return fmt.Sprintf(`GET {{base_url}}%s
HTTP 200
[Asserts]
jsonpath "$.id" exists`, examplePath)
	}
	return fmt.Sprintf(`GET {{base_url}}%s
HTTP 200
[Asserts]
jsonpath "$" count > 0`, examplePath)
}
