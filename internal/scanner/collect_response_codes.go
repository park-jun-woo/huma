//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Extracts HTTP status code integers from an OpenAPI responses map handling both string and interface key types
package scanner

func collectResponseCodes(responsesRaw interface{}) []int {
	keys := responseMapKeys(responsesRaw)

	var codes []int
	for _, key := range keys {
		code := toStatusCode(key)
		if code > 0 {
			codes = append(codes, code)
		}
	}

	return codes
}
