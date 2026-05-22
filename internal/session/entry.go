//ff:type feature=session type=model
//ff:what Entry combines an endpoint with its test status, coverage, and improvement tracking
package session

import (
	"github.com/park-jun-woo/hurlfill/internal/scanner"
)

type Entry struct {
	scanner.Endpoint
	Status       Status  `json:"status"`
	Coverage     float64 `json:"coverage,omitempty"`
	ImproveCount int     `json:"improve_count,omitempty"`
	PrevCoverage float64 `json:"prev_coverage,omitempty"`
}
