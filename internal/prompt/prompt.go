package prompt

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/hurlfill/internal/adapter"
	"github.com/park-jun-woo/hurlfill/internal/runner"
	"github.com/park-jun-woo/hurlfill/internal/scanner"
	"github.com/park-jun-woo/hurlfill/internal/source"
)

// TodoPrompt builds the agent instruction for a TODO endpoint.
func TodoPrompt(ep *scanner.Endpoint, hurlDir string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# TODO  %s %s\n", ep.Method, ep.Path))
	b.WriteString(fmt.Sprintf("# Source: %s:%d\n", ep.Source, ep.Line))
	b.WriteString(fmt.Sprintf("# Handler: %s\n", ep.Handler))

	// Handler source
	src, _, _, err := source.ReadHandler(ep.Source, ep.Handler)
	if err == nil && src != "" {
		b.WriteString("\n## Handler source\n\n")
		b.WriteString(src)
		b.WriteString("\n")
	}

	// Hurl example template
	b.WriteString("\n## Hurl example\n\n")
	b.WriteString(hurlExample(ep.Method, ep.Path))
	b.WriteString("\n")

	// Instructions
	hurlFile := runner.HurlFileName(ep, hurlDir)
	b.WriteString("\n## Instructions\n\n")
	b.WriteString(fmt.Sprintf("1. Write %s\n", hurlFile))
	b.WriteString("2. Run `hurlfill next`\n")

	return b.String()
}

// FailPrompt builds the agent instruction for a FAIL endpoint.
func FailPrompt(ep *scanner.Endpoint, hurlFile string, feedback string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# FAIL  %s %s\n", ep.Method, ep.Path))
	b.WriteString(fmt.Sprintf("# File: %s\n", hurlFile))
	b.WriteString("\n")
	b.WriteString(feedback)
	b.WriteString("\n")
	b.WriteString("\n## Instructions\n\n")
	b.WriteString(fmt.Sprintf("1. Edit %s\n", hurlFile))
	b.WriteString("2. Fix the failing assertion\n")
	b.WriteString("3. Run `hurlfill next`\n")

	return b.String()
}

// PassPrompt builds the agent instruction for a PASS then next TODO.
func PassPrompt(ep *scanner.Endpoint) string {
	return fmt.Sprintf("# PASS  %s %s\n", ep.Method, ep.Path)
}

// AllComplete builds the final completion message.
func AllComplete(pass, total int) string {
	var b strings.Builder
	b.WriteString("All endpoints complete!\n\n")
	b.WriteString(fmt.Sprintf("PASS: %d (%d%%)\n", pass, percent(pass, total)))
	b.WriteString(fmt.Sprintf("TODO: %d (%d%%)\n", 0, 0))
	return b.String()
}

// ImprovePrompt builds the agent instruction for an endpoint that passes
// but has uncovered lines.
func ImprovePrompt(ep *scanner.Endpoint, hurlFile string, result *adapter.CoverageResult) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# IMPROVE  %s %s\n", ep.Method, ep.Path))
	b.WriteString(fmt.Sprintf("# Coverage: %.0f%% (%d/%d)\n", result.Percent, result.Covered, result.Total))

	if len(result.Uncovered) > 0 {
		b.WriteString("# UNCOVERED:\n")
		for _, u := range result.Uncovered {
			base := filepath.Base(u.File)
			b.WriteString(fmt.Sprintf("#   %s:%d  %s\n", base, u.Line, u.Code))
		}
	}

	b.WriteString("\n## Instructions\n\n")
	b.WriteString(fmt.Sprintf("1. Edit %s\n", hurlFile))
	b.WriteString("2. Add test entries for the uncovered branches above\n")
	b.WriteString("3. Run `hurlfill next`\n")

	return b.String()
}

func percent(n, total int) int {
	if total == 0 {
		return 0
	}
	return n * 100 / total
}

func hurlExample(method, path string) string {
	// Replace path params like :id with 1
	examplePath := replaceParams(path)

	switch method {
	case "POST":
		return fmt.Sprintf(`POST {{base_url}}%s
Content-Type: application/json
{"field": "value"}

HTTP 201
[Asserts]
jsonpath "$.id" exists`, examplePath)

	case "PUT", "PATCH":
		return fmt.Sprintf(`%s {{base_url}}%s
Content-Type: application/json
{"field": "new_value"}

HTTP 200
[Asserts]
jsonpath "$.id" exists`, method, examplePath)

	case "DELETE":
		return fmt.Sprintf(`DELETE {{base_url}}%s
HTTP 204`, examplePath)

	case "GET":
		if hasParam(path) {
			return fmt.Sprintf(`GET {{base_url}}%s
HTTP 200
[Asserts]
jsonpath "$.id" exists`, examplePath)
		}
		return fmt.Sprintf(`GET {{base_url}}%s
HTTP 200
[Asserts]
jsonpath "$" count > 0`, examplePath)

	default:
		return fmt.Sprintf(`%s {{base_url}}%s
HTTP 200`, method, examplePath)
	}
}

func hasParam(path string) bool {
	return strings.Contains(path, ":")
}

func replaceParams(path string) string {
	parts := strings.Split(path, "/")
	for i, p := range parts {
		if strings.HasPrefix(p, ":") {
			parts[i] = "1"
		}
	}
	return strings.Join(parts, "/")
}
