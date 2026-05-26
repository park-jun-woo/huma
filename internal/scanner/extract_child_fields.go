//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Extracts child-level response fields from a nested property map with a parent path prefix
package scanner

func extractChildFields(parentName string, childProps map[string]interface{}, components map[string]interface{}) []ResponseField {
	var fields []ResponseField
	for childName, childRaw := range childProps {
		child, ok := childRaw.(map[string]interface{})
		if !ok {
			continue
		}
		child = resolveRef(child, components, 10)
		fields = append(fields, ResponseField{
			Path: "$." + parentName + "." + childName,
			Type: stringField(child, "type"),
		})
	}
	return fields
}
