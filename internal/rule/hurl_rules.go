//ff:func feature=rule type=helper control=sequence
//ff:what Returns all hurl file validation rules (H-01 to H-05)
package rule

var (
	H01 = Rule{ID: "H-01", Level: "ERROR", Description: "Hurl file not found at expected path"}
	H02 = Rule{ID: "H-02", Level: "ERROR", Description: "Hurl execution failed"}
	H03 = Rule{ID: "H-03", Level: "ERROR", Description: "Hurl test failed"}
	H04 = Rule{ID: "H-04", Level: "WARNING", Description: "Existing hurl file name doesn't match naming convention"}
	H05 = Rule{ID: "H-05", Level: "WARNING", Description: "Hurl file missing {{host}} variable"}
)

// HurlRules returns all hurl file validation rules.
func HurlRules() []Rule {
	return []Rule{H01, H02, H03, H04, H05}
}
