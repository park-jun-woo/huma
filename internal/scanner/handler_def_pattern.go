//ff:type feature=scan type=model
//ff:what Language-specific regexes capturing identifiers in handler-definition positions
package scanner

import "regexp"

// identPattern matches one identifier in capture group 1 inside the
// language-specific definition shapes below. The captured identifier is later
// normalized (normalizeSymbol) and compared to the target handler name, which
// absorbs camelCase <-> PascalCase differences (§2.3).
const identPattern = `([A-Za-z_][A-Za-z0-9_]*)`

// These match only function/method DEFINITION shapes. Object-literal keys
// (`name:`) are deliberately NOT matched for function-style languages so JS
// client maps no longer false-match Go handlers (BUG-002, §2.2).
var (
	// Go: `func Name(` or `func (recv *T) Name(`.
	goDefRe = regexp.MustCompile(`\bfunc\s+(?:\([^)]*\)\s*)?` + identPattern + `\s*\(`)
	// Python: `def Name(` (and async def).
	pyDefRe = regexp.MustCompile(`\bdef\s+` + identPattern + `\s*\(`)
	// Rust: `fn Name(`.
	rsDefRe = regexp.MustCompile(`\bfn\s+` + identPattern + `\s*\(`)
	// Java / C#: a method signature `... Name(` preceded by an access/type
	// token (avoids bare calls). Requires at least one preceding word + space.
	javaDefRe = regexp.MustCompile(`\b[A-Za-z_<>\[\].]+\s+` + identPattern + `\s*\(`)
	// PHP: `function Name(`.
	phpDefRe = regexp.MustCompile(`\bfunction\s+` + identPattern + `\s*\(`)
)

// jsDefRes are tried in order for .js/.ts. Object-literal keys (`Name:`) are
// last because in a JS API client they denote references, not definitions.
var jsDefRes = []*regexp.Regexp{
	regexp.MustCompile(`\bfunction\s+` + identPattern + `\s*\(`),       // function Name(
	regexp.MustCompile(`\b` + identPattern + `\s*=\s*(?:async\s*)?\(`), // Name = (...) =>
	regexp.MustCompile(`\b` + identPattern + `\s*\(`),                  // Name( method
	regexp.MustCompile(`\b` + identPattern + `\s*:`),                   // Name: (object literal, last resort)
}
