//ff:type feature=hurlcheck type=model
//ff:what entryShape is the comparable signature of a hurl entry used to detect copy-paste reuse
package hurlcheck

// entryShape is the comparable signature of an entry: same method, url, and
// assertion count indicate a likely copy-paste across branches.
type entryShape struct {
	method  string
	url     string
	asserts int
}
