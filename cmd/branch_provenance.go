//ff:type feature=ratchet type=model
//ff:what BranchProvenance records whether the response-branch denominator came from source, declarations, or both
package cmd

// BranchProvenance records the origin of the response-branch denominator.
// It feeds the O/E axis judgement (§3.2) and the transparency output (§5).
type BranchProvenance struct {
	HasSource   bool // source AST analysis contributed branches (ground-truth floor)
	HasDeclared bool // OpenAPI declarations contributed branches
}
