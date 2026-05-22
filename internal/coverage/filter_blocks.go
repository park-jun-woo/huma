//ff:func feature=coverage type=parser control=iteration dimension=1
//ff:what Filters coverage blocks by file suffix match and line range overlap
package coverage

// FilterBlocks returns blocks whose file suffix matches the given file
// and whose line range overlaps [startLine, endLine].
func FilterBlocks(blocks []Block, file string, startLine, endLine int) []Block {
	var filtered []Block
	for _, b := range blocks {
		if !matchFile(b.File, file) {
			continue
		}
		if b.EndLine < startLine || b.StartLine > endLine {
			continue
		}
		filtered = append(filtered, b)
	}
	return filtered
}
