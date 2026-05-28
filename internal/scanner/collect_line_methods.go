//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Collects HTTP methods from a single line using a regex pattern and appends unseen allowed methods
package scanner

import (
	"regexp"
	"strings"
)

func collectLineMethods(line string, re *regexp.Regexp, seen map[string]bool, methods *[]string) {
	for _, m := range re.FindAllStringSubmatch(line, -1) {
		method := strings.ToUpper(m[1])
		if !allowedMethods[method] || seen[method] {
			continue
		}
		seen[method] = true
		*methods = append(*methods, method)
	}
}
