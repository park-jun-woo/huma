//ff:func feature=scan type=helper control=iteration dimension=1
//ff:what Normalizes a symbol name to a case- and separator-insensitive key for handler matching
package scanner

import "strings"

// normalizeSymbol lowercases name and drops every non-alphanumeric rune so that
// camelCase, PascalCase, and snake_case spellings collapse to one key
// (createSubscriber, CreateSubscriber, create_subscriber -> "createsubscriber").
func normalizeSymbol(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r - 'A' + 'a')
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		}
	}
	return b.String()
}
