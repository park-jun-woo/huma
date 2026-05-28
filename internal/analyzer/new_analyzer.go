//ff:func feature=analyzer type=helper control=selection
//ff:what Selects the appropriate Analyzer implementation based on language name
package analyzer

// NewAnalyzer returns the Analyzer for the given language, or nil if unsupported.
func NewAnalyzer(lang string) Analyzer {
	switch lang {
	case "go":
		return &GoAnalyzer{}
	case "python":
		return &PythonAnalyzer{}
	case "node", "javascript", "typescript", "nestjs":
		return &NodeAnalyzer{}
	default:
		return nil
	}
}
