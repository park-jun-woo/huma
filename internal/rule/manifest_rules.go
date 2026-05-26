//ff:func feature=rule type=helper control=sequence
//ff:what Returns all manifest validation rules (M-01 to M-10)
package rule

var (
	M01 = Rule{ID: "M-01", Level: "ERROR", Description: "manifest.yaml not found"}
	M02 = Rule{ID: "M-02", Level: "ERROR", Description: "manifest.yaml parse error (invalid YAML)"}
	M03 = Rule{ID: "M-03", Level: "ERROR", Description: "apiVersion missing or unsupported"}
	M04 = Rule{ID: "M-04", Level: "ERROR", Description: "metadata.name missing"}
	M05 = Rule{ID: "M-05", Level: "ERROR", Description: "backend.lang missing"}
	M06 = Rule{ID: "M-06", Level: "ERROR", Description: "testing.base_url missing"}
	M07 = Rule{ID: "M-07", Level: "ERROR", Description: "testing.hurl_dir missing"}
	M08 = Rule{ID: "M-08", Level: "ERROR", Description: "testing.server.start missing"}
	M09 = Rule{ID: "M-09", Level: "ERROR", Description: "testing.server.ready missing"}
	M10 = Rule{ID: "M-10", Level: "WARNING", Description: "testing.hurl_variables empty"}
)

// ManifestRules returns all manifest validation rules.
func ManifestRules() []Rule {
	return []Rule{M01, M02, M03, M04, M05, M06, M07, M08, M09, M10}
}
