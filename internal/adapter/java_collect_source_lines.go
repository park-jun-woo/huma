//ff:func feature=adapter type=helper control=iteration dimension=1
//ff:what Collects covered and total line maps from a JaCoCo source file within a line range
package adapter

// collectSourceFileLines populates covered/total maps from a JaCoCo source file,
// filtering by the given line range.
func collectSourceFileLines(sf jacocoSourceFile, startLine, endLine int, covered, total map[int]bool) {
	for _, line := range sf.Lines {
		if startLine > 0 && line.Nr < startLine {
			continue
		}
		if endLine > 0 && line.Nr > endLine {
			continue
		}
		if line.Ci > 0 || line.Mi > 0 {
			total[line.Nr] = true
		}
		if line.Ci > 0 {
			covered[line.Nr] = true
		}
	}
}
