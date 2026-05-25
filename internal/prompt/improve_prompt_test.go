package prompt

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/adapter"
	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestImprovePrompt2(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/users",
		Handler: "GetUsers", Source: "handler.go", Line: 1,
	}
	covResult := &adapter.CoverageResult{
		Covered: 7,
		Total:   10,
		Percent: 70,
		Uncovered: []adapter.UncoveredLine{
			{File: "handler.go", Line: 15, Code: "if err != nil {"},
		},
	}

	result := ImprovePrompt(ep, "hurl/get_users.hurl", covResult)
	if !strings.Contains(result, "# IMPROVE  GET /users") {
		t.Fatal("expected IMPROVE header")
	}
	if !strings.Contains(result, "Coverage: 70%") {
		t.Fatal("expected coverage percentage")
	}
	if !strings.Contains(result, "UNCOVERED") {
		t.Fatal("expected uncovered section")
	}
	if !strings.Contains(result, "handler.go:15") {
		t.Fatal("expected uncovered line reference")
	}
}

func TestImprovePrompt2_NoUncovered(t *testing.T) {
	ep := &scanner.Endpoint{
		ID: "ep1", Method: "GET", Path: "/test",
	}
	covResult := &adapter.CoverageResult{
		Covered: 10,
		Total:   10,
		Percent: 100,
	}

	result := ImprovePrompt(ep, "test.hurl", covResult)
	if strings.Contains(result, "UNCOVERED") {
		t.Fatal("should not have uncovered section")
	}
}
