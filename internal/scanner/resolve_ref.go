//ff:func feature=scan type=helper control=sequence
//ff:what Resolves a $ref reference in an OpenAPI schema by replacing it with the referenced component schema
package scanner

import "strings"

func resolveRef(schema map[string]interface{}, components map[string]interface{}, depth int) map[string]interface{} {
	if depth <= 0 {
		return schema
	}

	ref, ok := schema["$ref"].(string)
	if !ok {
		return schema
	}

	const prefix = "#/components/schemas/"
	if !strings.HasPrefix(ref, prefix) {
		return schema
	}

	name := ref[len(prefix):]
	resolved, ok := components[name].(map[string]interface{})
	if !ok {
		return schema
	}

	return resolveRef(resolved, components, depth-1)
}
