package humaquest

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/quest"
)

// itemWithLog builds an in-memory Item carrying the given Tries and Attempt log.
// The helper drives the read-only phase derivation without touching disk.
func itemWithLog(tries int, attempts ...quest.Attempt) *quest.Item {
	return &quest.Item{Tries: tries, Log: attempts}
}

func TestLastOutcome(t *testing.T) {
	tests := []struct {
		name string
		it   *quest.Item
		want quest.Outcome
	}{
		{"empty log", itemWithLog(0), ""},
		{"single fail", itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutFail)}), quest.OutFail},
		{"single review", itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutReview)}), quest.OutReview},
		{"single pass", itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutPass)}), quest.OutPass},
		{
			"last of many wins",
			itemWithLog(3,
				quest.Attempt{Try: 1, Outcome: string(quest.OutFail)},
				quest.Attempt{Try: 2, Outcome: string(quest.OutReview)},
				quest.Attempt{Try: 3, Outcome: string(quest.OutPass)},
			),
			quest.OutPass,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lastOutcome(tt.it); got != tt.want {
				t.Errorf("lastOutcome = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLastReason(t *testing.T) {
	tests := []struct {
		name string
		it   *quest.Item
		want string
	}{
		{"empty log", itemWithLog(0), ""},
		{
			"single reason",
			itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutFail), Reason: "R1: status 404 uncovered"}),
			"R1: status 404 uncovered",
		},
		{
			"empty reason on last attempt",
			itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutPass), Reason: ""}),
			"",
		},
		{
			"last of many wins",
			itemWithLog(2,
				quest.Attempt{Try: 1, Outcome: string(quest.OutFail), Reason: "first shortfall"},
				quest.Attempt{Try: 2, Outcome: string(quest.OutFail), Reason: "second shortfall"},
			),
			"second shortfall",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lastReason(tt.it); got != tt.want {
				t.Errorf("lastReason = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDerivePhase(t *testing.T) {
	tests := []struct {
		name       string
		it         *quest.Item
		staticMode bool
		want       phase
	}{
		// Fresh item (no log) → manifest mode decides.
		{"fresh static", itemWithLog(0), true, phaseStatic},
		{"fresh live", itemWithLog(0), false, phaseTodo},

		// REVIEW remaps to UNVERIFIED regardless of mode.
		{
			"review static",
			itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutReview)}),
			true, phaseUnverified,
		},
		{
			"review live",
			itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutReview)}),
			false, phaseUnverified,
		},

		// FAIL → improve regardless of mode (FAIL implies Tries > 0).
		{
			"fail static",
			itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutFail)}),
			true, phaseImprove,
		},
		{
			"fail live",
			itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutFail)}),
			false, phaseImprove,
		},

		// A non-FAIL/REVIEW terminal outcome (e.g. PASS) falls back to mode.
		{
			"pass static falls back to static",
			itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutPass)}),
			true, phaseStatic,
		},
		{
			"pass live falls back to todo",
			itemWithLog(1, quest.Attempt{Try: 1, Outcome: string(quest.OutPass)}),
			false, phaseTodo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := derivePhase(tt.it, tt.staticMode); got != tt.want {
				t.Errorf("derivePhase = %d, want %d", got, tt.want)
			}
		})
	}
}
