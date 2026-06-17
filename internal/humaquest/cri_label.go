//ff:func feature=verify type=helper control=selection
//ff:what Maps a CRI tier (0..3) to its display label (UNVERIFIED/SCAFFOLDED/SMOKE/COVERED) for the transparency output.

package humaquest

// criLabel returns the display label for a CRI tier (§4/§5). Used in the PASS
// transparency line and verdict feedback.
func criLabel(tier int) string {
	switch tier {
	case 3:
		return "COVERED"
	case 2:
		return "SMOKE"
	case 1:
		return "SCAFFOLDED"
	default:
		return "UNVERIFIED"
	}
}
