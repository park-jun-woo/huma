//ff:type feature=adapter type=model
//ff:what JaCoCo XML report root element for coverage parsing
package adapter

import "encoding/xml"

type jacocoReport struct {
	XMLName  xml.Name        `xml:"report"`
	Packages []jacocoPackage `xml:"package"`
}
