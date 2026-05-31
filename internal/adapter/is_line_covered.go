//ff:func feature=adapter type=helper control=sequence
//ff:what Reports whether a source line is in the runtime-covered line set
package adapter

// IsLineCovered reports whether the given source line was hit at runtime.
func (r *CoverageResult) IsLineCovered(line int) bool {
	if r == nil || r.CoveredLines == nil {
		return false
	}
	return r.CoveredLines[line]
}
