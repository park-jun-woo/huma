//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Extracts response fields from an OpenAPI operation by finding 2xx response schemas and resolving refs
package scanner

import "sort"

func extractFieldsFromOp(op map[string]interface{}, components map[string]interface{}) []ResponseField {
	responsesRaw, ok := op["responses"]
	if !ok {
		return nil
	}

	successCodes := []int{200, 201}
	var fields []ResponseField

	for _, code := range successCodes {
		respObj, ok := lookupResponse(responsesRaw, code)
		if !ok {
			continue
		}

		content, ok := respObj["content"].(map[string]interface{})
		if !ok {
			continue
		}

		jsonContent, ok := content["application/json"].(map[string]interface{})
		if !ok {
			continue
		}

		schema, ok := jsonContent["schema"].(map[string]interface{})
		if !ok {
			continue
		}

		extracted := extractResponseFields(schema, components)
		fields = append(fields, extracted...)
	}

	if len(fields) == 0 {
		return nil
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Path < fields[j].Path
	})

	return fields
}
