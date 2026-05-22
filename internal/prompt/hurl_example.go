//ff:func feature=prompt type=builder control=selection
//ff:what Generates an example Hurl template based on HTTP method and path
package prompt

import (
	"fmt"
)

func hurlExample(method, path string) string {
	examplePath := replaceParams(path)

	switch method {
	case "POST":
		return fmt.Sprintf(`POST {{base_url}}%s
Content-Type: application/json
{"field": "value"}

HTTP 201
[Asserts]
jsonpath "$.id" exists`, examplePath)

	case "PUT", "PATCH":
		return fmt.Sprintf(`%s {{base_url}}%s
Content-Type: application/json
{"field": "new_value"}

HTTP 200
[Asserts]
jsonpath "$.id" exists`, method, examplePath)

	case "DELETE":
		return fmt.Sprintf(`DELETE {{base_url}}%s
HTTP 204`, examplePath)

	case "GET":
		return getExample(examplePath, path)

	default:
		return fmt.Sprintf(`%s {{base_url}}%s
HTTP 200`, method, examplePath)
	}
}
