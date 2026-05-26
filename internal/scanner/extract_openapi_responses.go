//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Builds a JSON response array from OpenAPI operation responses map for Endpoint.Responses
package scanner

import (
	"encoding/json"
	"sort"
)

func extractOpenAPIResponses(op map[string]interface{}) json.RawMessage {
	responsesRaw, ok := op["responses"]
	if !ok {
		return nil
	}

	codes := collectResponseCodes(responsesRaw)
	if len(codes) == 0 {
		return nil
	}

	sort.Ints(codes)

	type entry struct {
		Status int `json:"status"`
	}

	entries := make([]entry, len(codes))
	for i, code := range codes {
		entries[i] = entry{Status: code}
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return nil
	}

	return data
}
