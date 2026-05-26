//ff:func feature=scan type=parser control=iteration dimension=1
//ff:what Checks if a trimmed line starts with an HTTP method and returns the method and raw URL if found
package cmd

import "strings"

func matchLineMethod(line string) (method, url string, ok bool) {
	for _, m := range httpMethods {
		if strings.HasPrefix(line, m+" ") {
			return m, strings.TrimSpace(line[len(m)+1:]), true
		}
	}
	return "", "", false
}
