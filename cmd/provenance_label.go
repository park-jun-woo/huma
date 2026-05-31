//ff:func feature=session type=helper control=sequence
//ff:what Returns the provenance label or "n/a" when provenance is unknown
package cmd

// provenanceLabel returns the denominator provenance label, defaulting to "n/a".
func provenanceLabel(prov string) string {
	if prov == "" {
		return "n/a"
	}
	return prov
}
