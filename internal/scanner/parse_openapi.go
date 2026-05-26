//ff:func feature=scan type=parser control=iteration dimension=2
//ff:what Parses OpenAPI paths object and extracts endpoints with method, path, handler, source, and line
package scanner

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var httpMethods = map[string]bool{
	"get":    true,
	"post":   true,
	"put":    true,
	"delete": true,
	"patch":  true,
}

func parseOpenAPI(data []byte) ([]Endpoint, error) {
	var doc map[string]interface{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("unmarshal openapi: %w", err)
	}

	pathsRaw, ok := doc["paths"]
	if !ok {
		return nil, fmt.Errorf("openapi document missing 'paths' key")
	}

	paths, ok := pathsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("openapi 'paths' is not a map")
	}

	var endpoints []Endpoint
	for path, methodsRaw := range paths {
		path = strings.TrimRight(path, "/")
		methods, ok := methodsRaw.(map[string]interface{})
		if !ok {
			continue
		}
		for method, opRaw := range methods {
			methodLower := strings.ToLower(method)
			if !httpMethods[methodLower] {
				continue
			}

			op, ok := opRaw.(map[string]interface{})
			if !ok {
				continue
			}

			handler := stringField(op, "operationId")
			if handler == "" {
				handler = generateOperationID(methodLower, path)
			}

			source := stringField(op, "x-source-file")
			line := intField(op, "x-source-line")

			humaPath := convertPathParams(path)
			ep := Endpoint{
				ID:      makeID(strings.ToUpper(methodLower), humaPath),
				Method:  strings.ToUpper(methodLower),
				Path:    humaPath,
				Handler: handler,
				Source:  source,
				Line:    line,
			}
			if resp := extractOpenAPIResponses(op); resp != nil {
				ep.Responses = resp
			}
			endpoints = append(endpoints, ep)
		}
	}

	return endpoints, nil
}
