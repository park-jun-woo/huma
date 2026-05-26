//ff:func feature=scan type=parser control=sequence
//ff:what Parses endpoint input by detecting OpenAPI format or falling back to endpoint list format
package scanner

import "fmt"

func ParseEndpoints(data []byte) ([]Endpoint, error) {
	if isOpenAPI(data) {
		return parseOpenAPI(data)
	}

	raw, err := parseEndpointList(data)
	if err != nil {
		return nil, fmt.Errorf("parse endpoints: %w", err)
	}

	return collectRawEndpoints(raw), nil
}
