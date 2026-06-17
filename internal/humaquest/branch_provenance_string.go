//ff:func feature=ratchet type=helper control=selection
//ff:what Renders branch provenance as a short label (source/declared/both/none) for the CRI transparency line.

package humaquest

// String renders the provenance as a short label. Ported from
// bak/cmd/branch_provenance_string.go.
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
