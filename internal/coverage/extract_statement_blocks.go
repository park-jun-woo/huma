//ff:func feature=coverage type=helper control=iteration dimension=1
//ff:what Converts istanbul statement map entries into coverage blocks
package coverage

// extractStatementBlocks converts each statement in an istanbul file coverage entry
// into a coverage Block with the statement's line range and execution count.
func extractStatementBlocks(fc istanbulFileCoverage) []Block {
	var blocks []Block
	for key, rng := range fc.StatementMap {
		count, ok := fc.S[key]
		if !ok {
			count = 0
		}
		blocks = append(blocks, Block{
			File:      fc.Path,
			StartLine: rng.Start.Line,
			EndLine:   rng.End.Line,
			Count:     count,
		})
	}
	return blocks
}
