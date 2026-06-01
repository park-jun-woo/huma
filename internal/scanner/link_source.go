//ff:func feature=scan type=engine control=iteration dimension=1
//ff:what Links endpoints to backend handler file:line under a root, rejecting ambiguous/ext-mismatch
package scanner

import (
	"path/filepath"
	"strings"
)

// LinkSource scans root for source files restricted to the backend lang and,
// for every endpoint with a Handler name but no Source, links it to the single
// file:line where that handler is defined. Candidates with >1 matching file or
// a language-mismatched extension are rejected and left UNVERIFIED (§0/§2.5).
// Unknown/empty lang falls back to all source extensions (no regression for
// manifest-less repos, §2.1). Returns per-language link distribution and skip
// counts for transparent reporting (§2.6).
func LinkSource(endpoints []Endpoint, root, lang string) LinkResult {
	extSet, langKnown := allowedExts(lang)
	files := collectSourceFiles(root, lang)
	res := LinkResult{LangKnown: langKnown, Lang: lang, ByExt: map[string]int{}}
	for i := range endpoints {
		outcome, msg := linkEndpoint(&endpoints[i], files, root, lang, extSet)
		switch outcome {
		case outcomeLinked:
			res.Linked++
			ext := strings.ToLower(filepath.Ext(endpoints[i].Source))
			res.ByExt[ext]++
		case outcomeAmbiguous:
			res.Skipped++
			res.Ambiguous++
			res.SkipMessages = append(res.SkipMessages, msg)
		case outcomeExtMismatch:
			res.Skipped++
			res.ExtMismatch++
			res.SkipMessages = append(res.SkipMessages, msg)
		}
	}
	return res
}
