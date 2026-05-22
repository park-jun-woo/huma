//ff:type feature=adapter type=model
//ff:what UncoveredLine represents a single source line not covered by tests
package adapter

// UncoveredLine represents a single line of source code that was not covered.
type UncoveredLine struct {
	File string
	Line int
	Code string
}
