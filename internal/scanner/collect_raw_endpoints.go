//ff:func feature=scan type=parser control=iteration dimension=1
//ff:what Converts a slice of rawEndpoint into a slice of Endpoint by parsing each entry
package scanner

func collectRawEndpoints(raw []rawEndpoint) []Endpoint {
	endpoints := make([]Endpoint, 0, len(raw))
	for _, r := range raw {
		ep := parseRawEndpoint(r)
		if ep != nil {
			endpoints = append(endpoints, *ep)
		}
	}
	return endpoints
}
