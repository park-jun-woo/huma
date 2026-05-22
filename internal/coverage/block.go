//ff:type feature=coverage type=model
//ff:what Block represents a single coverage block from a Go coverage.out file
package coverage

// Block represents a single coverage block from a Go coverage.out file.
type Block struct {
	File      string
	StartLine int
	EndLine   int
	Count     int
}
