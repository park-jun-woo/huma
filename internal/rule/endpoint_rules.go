//ff:func feature=rule type=helper control=sequence
//ff:what Returns all endpoint validation rules (E-01 to E-09)
package rule

var (
	E01 = Rule{ID: "E-01", Level: "ERROR", Description: "No OpenAPI file found and --from not specified"}
	E02 = Rule{ID: "E-02", Level: "ERROR", Description: "Input file not readable"}
	E03 = Rule{ID: "E-03", Level: "ERROR", Description: "Input is not valid JSON/YAML"}
	E04 = Rule{ID: "E-04", Level: "ERROR", Description: "Endpoint missing method field"}
	E05 = Rule{ID: "E-05", Level: "ERROR", Description: "Endpoint missing path field"}
	E06 = Rule{ID: "E-06", Level: "WARNING", Description: "Endpoint missing handler field"}
	E07 = Rule{ID: "E-07", Level: "WARNING", Description: "Endpoint missing file field"}
	E08 = Rule{ID: "E-08", Level: "WARNING", Description: "Duplicate endpoint"}
	E09 = Rule{ID: "E-09", Level: "WARNING", Description: "OpenAPI auto-detect failed, falling back to endpoint list parser"}
)

// EndpointRules returns all endpoint validation rules.
func EndpointRules() []Rule {
	return []Rule{E01, E02, E03, E04, E05, E06, E07, E08, E09}
}
