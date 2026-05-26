package cmd

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/analyzer"
)

func TestFilterClientBranches_Nil(t *testing.T) {
	got := filterClientBranches(nil)
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestFilterClientBranches_Empty(t *testing.T) {
	got := filterClientBranches([]analyzer.ResponseBranch{})
	if len(got) != 0 {
		t.Fatalf("expected empty, got %d", len(got))
	}
}

func TestFilterClientBranches_RemovesServerErrors(t *testing.T) {
	branches := []analyzer.ResponseBranch{
		{Status: 200},
		{Status: 400},
		{Status: 500},
		{Status: 503},
	}
	got := filterClientBranches(branches)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	if got[0].Status != 200 {
		t.Fatalf("expected 200, got %d", got[0].Status)
	}
	if got[1].Status != 400 {
		t.Fatalf("expected 400, got %d", got[1].Status)
	}
}

func TestFilterClientBranches_Keeps4xx(t *testing.T) {
	branches := []analyzer.ResponseBranch{
		{Status: 400},
		{Status: 401},
		{Status: 403},
		{Status: 404},
		{Status: 409},
		{Status: 422},
	}
	got := filterClientBranches(branches)
	if len(got) != 6 {
		t.Fatalf("expected 6, got %d", len(got))
	}
}

func TestFilterClientBranches_AllServerErrors(t *testing.T) {
	branches := []analyzer.ResponseBranch{
		{Status: 500},
		{Status: 502},
		{Status: 503},
	}
	got := filterClientBranches(branches)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}
