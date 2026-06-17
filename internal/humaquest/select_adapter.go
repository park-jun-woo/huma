//ff:func feature=adapter type=helper control=selection
//ff:what cfg.Scan.Lang에 맞는 adapter.Adapter 구현체를 고른다(go/python/node/deno/java/dotnet/php/rust). 미지정/미지원 언어는 GoAdapter로 폴백. bak/cmd의 newAdapterFn 이식.

package humaquest

import (
	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/config"
)

// selectAdapter returns the coverage adapter for the configured backend language,
// porting bak/cmd's newAdapterFn factory. Unknown or unset languages fall back to
// the Go adapter so cover never hard-fails on an unrecognized manifest.
func selectAdapter(cfg *config.Config) adapter.Adapter {
	switch cfg.Scan.Lang {
	case "python", "flask", "django", "drf":
		return adapter.NewPythonAdapter(cfg)
	case "node", "javascript", "typescript", "nestjs", "express", "fastify", "hono":
		return adapter.NewNodeAdapter(cfg)
	case "deno", "edge-functions":
		return adapter.NewDenoAdapter(cfg)
	case "java", "spring", "quarkus":
		return adapter.NewJavaAdapter(cfg)
	case "dotnet", "aspnet", "csharp":
		return adapter.NewDotnetAdapter(cfg)
	case "php", "laravel":
		return adapter.NewPhpAdapter(cfg)
	case "rust", "actix":
		return adapter.NewRustAdapter(cfg)
	case "go", "fiber", "echo":
		return adapter.NewGoAdapter(cfg)
	default:
		return adapter.NewGoAdapter(cfg)
	}
}
