//ff:type feature=scan type=model
//ff:what LinkResult aggregates link/skip statistics produced by LinkSource for transparent reporting
package scanner

// LinkResult summarizes a LinkSource run: how many endpoints were linked
// (broken down by the extension/lang of the matched file), and how many were
// skipped with the reason (ext-mismatch vs ambiguous). LangKnown is false when
// the backend lang was empty/unrecognized (full-extension fallback, §2.1).
type LinkResult struct {
	Linked       int
	Skipped      int
	ExtMismatch  int
	Ambiguous    int
	LangKnown    bool
	Lang         string
	ByExt        map[string]int // matched-file extension -> count, e.g. ".go" -> 142
	SkipMessages []string       // user-facing reasons for each skipped endpoint
}
