//ff:type feature=gate type=model
//ff:what Item.Payload 왕복 타입. Seed가 쓴 scanner.Endpoint를 임베드하고 IMPROVE 단조성용 PrevCoverage/ImproveCount를 곁들인다. Endpoint 필드를 top-level로 승격해 기존 bare-Endpoint 디코드와 호환된다.

package humaquest

import "github.com/park-jun-woo/huma/internal/scanner"

// payloadState is the Item.Payload round-trip carrier for the CRI verdict. It
// embeds the scanner.Endpoint that Seed wrote, so its JSON keeps the endpoint
// fields at the top level — a bare `scanner.Endpoint` decode (Prepare/Render)
// still round-trips, and decoding an old bare-Endpoint payload here leaves
// PrevCoverage/ImproveCount zero. The two extra fields carry IMPROVE
// monotonicity across retries (huma's Entry.PrevCoverage/ImproveCount, §2):
// Evaluate reads them (read-only) to compare coverage vs the prior attempt;
// Phase 006's cover command writes the updated values before Save (Evaluate
// never writes Payload).
type payloadState struct {
	scanner.Endpoint
	// PrevCoverage is the line-coverage percent measured on the previous attempt.
	PrevCoverage float64 `json:"prev_coverage,omitempty"`
	// ImproveCount is how many IMPROVE retries have already run for this item.
	ImproveCount int `json:"improve_count,omitempty"`
}
