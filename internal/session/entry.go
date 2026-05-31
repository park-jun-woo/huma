//ff:type feature=session type=model
//ff:what Entry combines an endpoint with its test status, coverage, CRI tier, and improvement tracking
package session

import (
	"github.com/park-jun-woo/huma/internal/scanner"
)

type Entry struct {
	scanner.Endpoint
	Status       Status  `json:"status"`
	Coverage     float64 `json:"coverage,omitempty"`
	ImproveCount int     `json:"improve_count,omitempty"`
	PrevCoverage float64 `json:"prev_coverage,omitempty"`

	// CRI is the cheese-resistance index (0..3) of this entry's verdict.
	// The display label (SCAFFOLDED/SMOKE/COVERED) is derived from CRI.
	CRI int `json:"cri,omitempty"`
	// AGrade is the measured assertion depth (0..3) of the hurl entries.
	AGrade int `json:"a_grade,omitempty"`
	// Provenance records where the denominator branches came from
	// ("source", "declared", "both", or "" when unknown).
	Provenance string `json:"provenance,omitempty"`
}
