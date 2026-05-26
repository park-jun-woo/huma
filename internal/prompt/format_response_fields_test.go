package prompt

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/huma/internal/scanner"
)

func TestFormatResponseFields_WithFields(t *testing.T) {
	fields := []scanner.ResponseField{
		{Path: "$.success", Type: "boolean"},
		{Path: "$.data.id", Type: "integer"},
		{Path: "$.data.email", Type: "string"},
	}

	result := formatResponseFields(fields)
	if !strings.Contains(result, "## Response fields") {
		t.Fatal("expected Response fields header")
	}
	if !strings.Contains(result, "$.success — boolean") {
		t.Fatal("expected $.success — boolean")
	}
	if !strings.Contains(result, "$.data.id — integer") {
		t.Fatal("expected $.data.id — integer")
	}
	if !strings.Contains(result, "$.data.email — string") {
		t.Fatal("expected $.data.email — string")
	}
}

func TestFormatResponseFields_Empty(t *testing.T) {
	result := formatResponseFields(nil)
	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

func TestFormatResponseFields_NoType(t *testing.T) {
	fields := []scanner.ResponseField{
		{Path: "$.data", Type: ""},
	}

	result := formatResponseFields(fields)
	if !strings.Contains(result, "  $.data\n") {
		t.Fatalf("expected path without type, got %q", result)
	}
	if strings.Contains(result, "—") {
		t.Fatal("should not have dash when type is empty")
	}
}
