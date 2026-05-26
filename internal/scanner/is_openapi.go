//ff:func feature=scan type=parser control=sequence
//ff:what Detects whether byte data is an OpenAPI document by checking for top-level openapi or swagger key
package scanner

import "gopkg.in/yaml.v3"

func isOpenAPI(data []byte) bool {
	var doc map[string]interface{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return false
	}
	if _, ok := doc["openapi"]; ok {
		return true
	}
	if _, ok := doc["swagger"]; ok {
		return true
	}
	return false
}
