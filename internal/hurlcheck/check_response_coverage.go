//ff:func feature=hurlcheck type=engine control=iteration dimension=1
//ff:what Compares analyzer response branches against hurl status codes to find missing test coverage
package hurlcheck

import "github.com/park-jun-woo/huma/internal/analyzer"

// CheckResponseCoverage compares expected response branches with tested statuses.
func CheckResponseCoverage(branches []analyzer.ResponseBranch, hurlStatuses []int) *ResponseCoverageResult {
	if len(branches) == 0 {
		return &ResponseCoverageResult{
			Covered: 0,
			Total:   0,
			Percent: 0,
		}
	}

	tested := make(map[int]bool)
	for _, s := range hurlStatuses {
		tested[s] = true
	}

	var missing []analyzer.ResponseBranch
	covered := 0

	// Deduplicate branches by status code (keep first occurrence)
	seen := make(map[int]bool)
	var unique []analyzer.ResponseBranch
	for _, b := range branches {
		if !seen[b.Status] {
			seen[b.Status] = true
			unique = append(unique, b)
		}
	}

	for _, b := range unique {
		if tested[b.Status] {
			covered++
		} else {
			missing = append(missing, b)
		}
	}

	total := len(unique)
	pct := float64(0)
	if total > 0 {
		pct = float64(covered) / float64(total) * 100
	}

	return &ResponseCoverageResult{
		Covered: covered,
		Total:   total,
		Percent: pct,
		Missing: missing,
	}
}
