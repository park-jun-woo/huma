//ff:func feature=prompt type=builder control=selection
//ff:what Generates an example Hurl template based on HTTP method, path, and URL variable name
package prompt

import (
	"fmt"
)

func hurlExample(method, path, urlVar string) string {
	examplePath := replaceParams(path)

	switch method {
	case "POST":
		return fmt.Sprintf(`POST {{%s}}%s
Content-Type: application/json
{"field": "value"}

HTTP 201
[Asserts]
jsonpath "$.id" exists`, urlVar, examplePath)

	case "PUT", "PATCH":
		return fmt.Sprintf(`%s {{%s}}%s
Content-Type: application/json
{"field": "new_value"}

HTTP 200
[Asserts]
jsonpath "$.id" exists`, method, urlVar, examplePath)

	case "DELETE":
		return fmt.Sprintf(`DELETE {{%s}}%s
HTTP 204`, urlVar, examplePath)

	case "GET":
		return getExample(examplePath, path, urlVar)

	default:
		return fmt.Sprintf(`%s {{%s}}%s
HTTP 200`, method, urlVar, examplePath)
	}
}
