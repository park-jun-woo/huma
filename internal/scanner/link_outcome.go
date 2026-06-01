//ff:type feature=scan type=model
//ff:what linkOutcome enumerates how a single endpoint fared during source linking
package scanner

// linkOutcome reports how a single endpoint fared during linking.
type linkOutcome int

const (
	outcomeNoop        linkOutcome = iota // already linked or no handler name; not counted
	outcomeLinked                         // Source/Line set
	outcomeNotFound                       // no candidate matched; left UNVERIFIED (silent)
	outcomeAmbiguous                      // >1 candidate file; rejected per §2.5
	outcomeExtMismatch                    // matched file ext not in lang's langExts; rejected per §2.5
)
