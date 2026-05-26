//ff:func feature=scan type=helper control=sequence
//ff:what Converts OpenAPI path parameters from {param} to :param format
package scanner

import "regexp"

var pathParamRe = regexp.MustCompile(`\{([^}]+)\}`)

func convertPathParams(path string) string {
	return pathParamRe.ReplaceAllString(path, ":$1")
}
