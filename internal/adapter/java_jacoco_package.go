//ff:type feature=adapter type=model
//ff:what JaCoCo XML package element containing source file coverage data
package adapter

type jacocoPackage struct {
	Name        string             `xml:"name,attr"`
	SourceFiles []jacocoSourceFile `xml:"sourcefile"`
}
