//ff:func feature=gate type=helper control=sequence
//ff:what scanEndpoints — bak/cmd/scan.go 오케스트레이션 이식. 위치인자(openapiPath [linkSourceRoot])로 FindOpenAPIFile→ParseEndpoints/ParseEdgeFunctions→LinkSource를 호출해 []scanner.Endpoint를 반환한다. 세션 쓰기·리포팅 출력은 reins가 담당하므로 버린다.

package humaquest

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/park-jun-woo/huma/internal/config"
	"github.com/park-jun-woo/huma/internal/scanner"
)

// scanEndpoints ports bak/cmd/scan.go's orchestration. reins scan passes only
// positional args, so the OpenAPI source and the optional --link-source root are
// positionally encoded:
//
//	huma scan [openapi] [link-source-root]
//
// openapi defaults to the first file FindOpenAPIFile discovers. A directory source
// is treated as a Supabase Edge Functions tree; "-" reads OpenAPI/endpoint JSON
// from stdin. When a link-source root is given, handlers are linked to file:line
// (in place) using the backend lang from manifest.yaml. Session persistence and
// link-result reporting from the old cmd glue are dropped — reins owns those.
func scanEndpoints(args ...string) ([]scanner.Endpoint, error) {
	from := ""
	if len(args) >= 1 {
		from = args[0]
	}
	linkRoot := ""
	if len(args) >= 2 {
		linkRoot = args[1]
	}

	if from == "" {
		from = scanner.FindOpenAPIFile()
		if from == "" {
			return nil, fmt.Errorf("[E-01] no OpenAPI file found; pass a path: huma scan <openapi> [link-source-root]")
		}
	}

	cfg, err := config.Load()
	if err != nil && !errors.Is(err, config.ErrNoManifest) {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var endpoints []scanner.Endpoint

	info, statErr := os.Stat(from)
	if from != "-" && statErr == nil && info.IsDir() {
		endpoints, err = scanner.ParseEdgeFunctions(from)
		if err != nil {
			return nil, fmt.Errorf("scan edge functions: %w", err)
		}
	} else {
		var data []byte
		if from == "-" {
			data, err = io.ReadAll(os.Stdin)
		} else {
			data, err = os.ReadFile(from)
		}
		if err != nil {
			return nil, fmt.Errorf("read input: %w", err)
		}
		endpoints, err = scanner.ParseEndpoints(data)
		if err != nil {
			return nil, fmt.Errorf("parse endpoints: %w", err)
		}
	}

	if linkRoot != "" {
		lang := ""
		if cfg != nil {
			lang = cfg.Scan.Lang
		}
		// LinkSource mutates endpoints in place (Source/Line); the returned
		// LinkResult is reporting-only and now owned by reins, so drop it.
		_ = scanner.LinkSource(endpoints, linkRoot, lang)
	}

	return endpoints, nil
}
