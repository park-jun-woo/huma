//ff:type feature=adapter type=model
//ff:what JaCoCo XML line element with instruction coverage counts
package adapter

type jacocoLine struct {
	Nr int `xml:"nr,attr"`
	Mi int `xml:"mi,attr"`
	Ci int `xml:"ci,attr"`
}
