//ff:func feature=analyzer type=helper control=selection
//ff:what Selects the appropriate Analyzer implementation based on language name
package analyzer

// NewAnalyzer returns the Analyzer for the given language, or nil if unsupported.
func NewAnalyzer(lang string) Analyzer {
	switch lang {
	case "go", "fiber", "echo":
		return &GoAnalyzer{}
	case "python", "flask", "django", "drf":
		return &PythonAnalyzer{}
	case "node", "javascript", "typescript", "nestjs", "express", "fastify", "hono":
		return &NodeAnalyzer{}
	case "deno", "edge-functions":
		return &DenoAnalyzer{}
	case "java", "spring", "quarkus":
		return &JavaAnalyzer{}
	case "dotnet", "aspnet", "csharp":
		return &DotnetAnalyzer{}
	case "php", "laravel":
		return &PhpAnalyzer{}
	case "rust", "actix":
		return &RustAnalyzer{}
	default:
		return nil
	}
}
