//ff:type feature=adapter type=model
//ff:what JaCoCo XML source file element containing line-level coverage
package adapter

type jacocoSourceFile struct {
	Name  string       `xml:"name,attr"`
	Lines []jacocoLine `xml:"line"`
}
