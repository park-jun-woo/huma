//ff:func feature=analyzer type=parser control=iteration dimension=1
//ff:what Parses the responses field from endpoint input JSON into ResponseBranch slice
package analyzer

import "encoding/json"

// ParseResponses converts a JSON-encoded responses array into ResponseBranch slice.
func ParseResponses(data json.RawMessage, file string) []ResponseBranch {
	if len(data) == 0 {
		return nil
	}

	var raw []rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil
	}

	branches := make([]ResponseBranch, 0, len(raw))
	for _, r := range raw {
		if r.Status <= 0 {
			continue
		}
		branches = append(branches, ResponseBranch{
			Status: r.Status,
			File:   file,
			Line:   r.Line,
			Code:   r.Code,
		})
	}

	return branches
}
