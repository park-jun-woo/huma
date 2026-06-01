//ff:type feature=scan type=model
//ff:what Source file extension maps: per-language handler extensions and the full fallback set
package scanner

// sourceExts is the union of source file extensions used as the fallback when
// the backend language is unknown/empty (preserves pre-Phase046 behavior).
var sourceExts = map[string]bool{
	".go": true, ".py": true, ".js": true, ".ts": true,
	".java": true, ".cs": true, ".php": true, ".rs": true,
}

// langExts maps a backend.lang (or framework alias) to the file extensions in
// which that language defines handlers. lang=go limits candidates to .go so
// polyglot frontend .js/.ts files are excluded outright (BUG-002 root fix).
var langExts = map[string][]string{
	"go":      {".go"},
	"python":  {".py"},
	"django":  {".py"},
	"fastapi": {".py"},
	"flask":   {".py"},
	"node":    {".js", ".ts"},
	"nestjs":  {".ts"},
	"express": {".js", ".ts"},
	"fastify": {".js", ".ts"},
	"hono":    {".ts"},
	"deno":    {".ts"},
	"java":    {".java"},
	"spring":  {".java"},
	"quarkus": {".java"},
	"dotnet":  {".cs"},
	"php":     {".php"},
	"laravel": {".php"},
	"rust":    {".rs"},
}
