//ff:func feature=scan type=parser control=iteration dimension=1
//ff:what Extracts field paths from a resolved OpenAPI schema up to 2 levels deep
package scanner

func extractResponseFields(schema map[string]interface{}, components map[string]interface{}) []ResponseField {
	resolved := resolveRef(schema, components, 10)

	props, ok := resolved["properties"].(map[string]interface{})
	if !ok {
		return nil
	}

	var fields []ResponseField
	for name, valRaw := range props {
		val, ok := valRaw.(map[string]interface{})
		if !ok {
			continue
		}
		val = resolveRef(val, components, 10)

		childProps, hasChildren := val["properties"].(map[string]interface{})
		if hasChildren {
			fields = append(fields, extractChildFields(name, childProps, components)...)
		} else {
			fields = append(fields, ResponseField{
				Path: "$." + name,
				Type: stringField(val, "type"),
			})
		}
	}

	return fields
}
