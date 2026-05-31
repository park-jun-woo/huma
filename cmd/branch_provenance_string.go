//ff:func feature=ratchet type=helper control=selection
//ff:what Renders branch provenance as a short label (source/declared/both/none)
package cmd

// String renders the provenance as a short label.
func (p BranchProvenance) String() string {
	switch {
	case p.HasSource && p.HasDeclared:
		return "both"
	case p.HasSource:
		return "source"
	case p.HasDeclared:
		return "declared"
	default:
		return "none"
	}
}
