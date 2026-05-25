//ff:func feature=scan type=parser control=iteration dimension=1
//ff:what Parses JSON or YAML endpoint input and generates IDs for each
package scanner

import "fmt"

func ParseEndpoints(data []byte) ([]Endpoint, error) {
	raw, err := unmarshalEndpoints(data)
	if err != nil {
		return nil, fmt.Errorf("parse endpoints: %w", err)
	}

	endpoints := make([]Endpoint, 0, len(raw))
	for _, r := range raw {
		ep := parseRawEndpoint(r)
		if ep != nil {
			endpoints = append(endpoints, *ep)
		}
	}

	return endpoints, nil
}
