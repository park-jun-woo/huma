package hurlcheck

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/analyzer"
)

func TestCheckResponseCoverage_AllCovered(t *testing.T) {
	branches := []analyzer.ResponseBranch{
		{Status: 201, File: "h.go", Line: 10},
		{Status: 400, File: "h.go", Line: 15},
	}
	statuses := []int{201, 400}

	result := CheckResponseCoverage(branches, statuses)
	if result.Covered != 2 {
		t.Fatalf("expected 2 covered, got %d", result.Covered)
	}
	if result.Total != 2 {
		t.Fatalf("expected 2 total, got %d", result.Total)
	}
	if result.Percent != 100 {
		t.Fatalf("expected 100%%, got %.0f%%", result.Percent)
	}
	if len(result.Missing) != 0 {
		t.Fatalf("expected 0 missing, got %d", len(result.Missing))
	}
}

func TestCheckResponseCoverage_SomeMissing(t *testing.T) {
	branches := []analyzer.ResponseBranch{
		{Status: 201, File: "h.go", Line: 10},
		{Status: 400, File: "h.go", Line: 15},
		{Status: 409, File: "h.go", Line: 20},
		{Status: 500, File: "h.go", Line: 25},
	}
	statuses := []int{201, 400}

	result := CheckResponseCoverage(branches, statuses)
	if result.Covered != 2 {
		t.Fatalf("expected 2 covered, got %d", result.Covered)
	}
	if result.Total != 4 {
		t.Fatalf("expected 4 total, got %d", result.Total)
	}
	if result.Percent != 50 {
		t.Fatalf("expected 50%%, got %.0f%%", result.Percent)
	}
	if len(result.Missing) != 2 {
		t.Fatalf("expected 2 missing, got %d", len(result.Missing))
	}
	if result.Missing[0].Status != 409 {
		t.Fatalf("expected 409, got %d", result.Missing[0].Status)
	}
	if result.Missing[1].Status != 500 {
		t.Fatalf("expected 500, got %d", result.Missing[1].Status)
	}
}

func TestCheckResponseCoverage_NoBranches(t *testing.T) {
	result := CheckResponseCoverage(nil, []int{200})
	if result.Total != 0 {
		t.Fatalf("expected 0 total, got %d", result.Total)
	}
}

func TestCheckResponseCoverage_NoStatuses(t *testing.T) {
	branches := []analyzer.ResponseBranch{
		{Status: 200, File: "h.go", Line: 10},
	}

	result := CheckResponseCoverage(branches, nil)
	if result.Covered != 0 {
		t.Fatalf("expected 0 covered, got %d", result.Covered)
	}
	if result.Total != 1 {
		t.Fatalf("expected 1 total, got %d", result.Total)
	}
	if len(result.Missing) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(result.Missing))
	}
}

func TestCheckResponseCoverage_DuplicateBranches(t *testing.T) {
	branches := []analyzer.ResponseBranch{
		{Status: 200, File: "h.go", Line: 10},
		{Status: 200, File: "h.go", Line: 15},
		{Status: 400, File: "h.go", Line: 20},
	}
	statuses := []int{200}

	result := CheckResponseCoverage(branches, statuses)
	if result.Total != 2 {
		t.Fatalf("expected 2 total (dedup), got %d", result.Total)
	}
	if result.Covered != 1 {
		t.Fatalf("expected 1 covered, got %d", result.Covered)
	}
	if len(result.Missing) != 1 {
		t.Fatalf("expected 1 missing, got %d", len(result.Missing))
	}
}
