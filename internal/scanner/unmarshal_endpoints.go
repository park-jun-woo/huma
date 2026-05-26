//ff:func feature=scan type=parser control=sequence
//ff:what Tries JSON then YAML formats to parse raw endpoint list data
package scanner

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

func parseEndpointList(data []byte) ([]rawEndpoint, error) {
	var raw []rawEndpoint
	if err := json.Unmarshal(data, &raw); err == nil {
		return raw, nil
	}

	var wrapped struct {
		Endpoints []rawEndpoint `yaml:"endpoints" json:"endpoints"`
	}
	if err := yaml.Unmarshal(data, &wrapped); err == nil && wrapped.Endpoints != nil {
		return wrapped.Endpoints, nil
	}

	if err := yaml.Unmarshal(data, &raw); err == nil {
		return raw, nil
	}

	return nil, fmt.Errorf("input is not valid JSON array, YAML array, or YAML with 'endpoints' key")
}
