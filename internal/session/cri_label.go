//ff:func feature=session type=helper control=selection
//ff:what Maps a CRI integer (0..3) to its display label (UNVERIFIED/SCAFFOLDED/SMOKE/COVERED)
package session

// CRILabel returns the human display label derived from a CRI value.
func CRILabel(cri int) string {
	switch {
	case cri >= 3:
		return "COVERED"
	case cri == 2:
		return "SMOKE"
	case cri == 1:
		return "SCAFFOLDED"
	default:
		return "UNVERIFIED"
	}
}
