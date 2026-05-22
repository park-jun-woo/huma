//ff:func feature=scan type=parser control=iteration dimension=1
//ff:what Parses a JSON array of endpoints and generates IDs for each
package scanner

import (
	"encoding/json"
	"fmt"
)

func ParseEndpoints(data []byte) ([]Endpoint, error) {
	var raw []rawEndpoint
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse endpoints JSON: %w", err)
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
